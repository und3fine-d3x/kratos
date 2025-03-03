package courier

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/herodot"

	gomail "github.com/ory/mail/v3"

	"kratos/driver/config"
	"kratos/x"
)

type (
	smtpDependencies interface {
		PersistenceProvider
		x.LoggingProvider
		config.Provider
	}
	Courier struct {
		Dialer *gomail.Dialer
		d      smtpDependencies
	}
	Provider interface {
		Courier(ctx context.Context) *Courier
	}
)

func NewSMTP(d smtpDependencies, c *config.Config) *Courier {
	uri := c.CourierSMTPURL()
	password, _ := uri.User.Password()
	port, _ := strconv.ParseInt(uri.Port(), 10, 0)

	var ssl bool
	var tlsConfig *tls.Config
	if uri.Scheme == "smtps" {
		ssl = true
		sslSkipVerify, _ := strconv.ParseBool(uri.Query().Get("skip_ssl_verify"))
		// #nosec G402 This is ok (and required!) because it is configurable and disabled by default.
		tlsConfig = &tls.Config{InsecureSkipVerify: sslSkipVerify, ServerName: uri.Hostname()}
	}

	return &Courier{
		d: d,
		Dialer: &gomail.Dialer{
			/* #nosec we need to support SMTP servers without TLS */
			TLSConfig:    tlsConfig,
			Host:         uri.Hostname(),
			Port:         int(port),
			Username:     uri.User.Username(),
			Password:     password,
			SSL:          ssl,
			Timeout:      time.Second * 10,
			RetryFailure: true,
		},
	}
}

func (m *Courier) QueueEmail(ctx context.Context, t EmailTemplate) (uuid.UUID, error) {
	recipient, err := t.EmailRecipient()
	if err != nil {
		return uuid.Nil, err
	}

	subject, err := t.EmailSubject()
	if err != nil {
		return uuid.Nil, err
	}

	bodyPlaintext, err := t.EmailBodyPlaintext()
	if err != nil {
		return uuid.Nil, err
	}

	templateType, err := GetTemplateType(t)
	if err != nil {
		return uuid.Nil, err
	}

	templateData, err := json.Marshal(t)
	if err != nil {
		return uuid.Nil, err
	}

	message := &Message{
		Status:       MessageStatusQueued,
		Type:         MessageTypeEmail,
		Recipient:    recipient,
		Body:         bodyPlaintext,
		Subject:      subject,
		TemplateType: templateType,
		TemplateData: templateData,
	}
	if err := m.d.CourierPersister().AddMessage(ctx, message); err != nil {
		return uuid.Nil, err
	}
	return message.ID, nil
}

func (m *Courier) Work(ctx context.Context) error {
	errChan := make(chan error)
	defer close(errChan)

	go m.watchMessages(ctx, errChan)

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil
		}
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

func (m *Courier) watchMessages(ctx context.Context, errChan chan error) {
	for {
		if err := backoff.Retry(func() error {
			return m.DispatchQueue(ctx)
		}, backoff.NewExponentialBackOff()); err != nil {
			errChan <- err
			return
		}
		time.Sleep(time.Second)
	}
}

func (m *Courier) DispatchMessage(ctx context.Context, msg Message) error {
	switch msg.Type {
	case MessageTypeEmail:
		from := m.d.Config(ctx).CourierSMTPFrom()
		fromName := m.d.Config(ctx).CourierSMTPFromName()
		gm := gomail.NewMessage()
		if fromName == "" {
			gm.SetHeader("From", from)
		} else {
			gm.SetAddressHeader("From", from, fromName)
		}

		gm.SetHeader("To", msg.Recipient)
		gm.SetHeader("Subject", msg.Subject)
		gm.SetBody("text/plain", msg.Body)

		tmpl, err := NewEmailTemplateFromMessage(m.d.Config(ctx), msg)
		if err != nil {
			m.d.Logger().
				WithError(err).
				WithField("message_id", msg.ID).
				Error(`Unable to get email template from message.`)
		} else {
			htmlBody, err := tmpl.EmailBody()
			if err != nil {
				m.d.Logger().
					WithError(err).
					WithField("message_id", msg.ID).
					Error(`Unable to get email body from template.`)
			} else {
				gm.AddAlternative("text/html", htmlBody)
			}
		}

		if err := m.Dialer.DialAndSend(ctx, gm); err != nil {
			m.d.Logger().
				WithError(err).
				WithField("smtp_server", fmt.Sprintf("%s:%d", m.Dialer.Host, m.Dialer.Port)).
				WithField("smtp_ssl_enabled", m.Dialer.SSL).
				// WithField("email_to", msg.Recipient).
				WithField("message_from", from).
				Error("Unable to send email using SMTP connection.")
			return errors.WithStack(err)
		}

		if err := m.d.CourierPersister().SetMessageStatus(ctx, msg.ID, MessageStatusSent); err != nil {
			m.d.Logger().
				WithError(err).
				WithField("message_id", msg.ID).
				Error(`Unable to set the message status to "sent".`)
			return err
		}

		m.d.Logger().
			WithField("message_id", msg.ID).
			WithField("message_type", msg.Type).
			WithField("message_template_type", msg.TemplateType).
			WithField("message_subject", msg.Subject).
			Debug("Courier sent out message.")
		return nil
	}
	return errors.Errorf("received unexpected message type: %d", msg.Type)
}

func (m *Courier) DispatchQueue(ctx context.Context) error {
	if len(m.Dialer.Host) == 0 {
		return errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Courier tried to deliver an email but courier.smtp_url is not set!"))
	}

	messages, err := m.d.CourierPersister().NextMessages(ctx, 10)
	if err != nil {
		if errors.Is(err, ErrQueueEmpty) {
			return nil
		}
		return err
	}

	for k := range messages {
		var msg = messages[k]
		if err := m.DispatchMessage(ctx, msg); err != nil {
			for _, replace := range messages[k:] {
				if err := m.d.CourierPersister().SetMessageStatus(ctx, replace.ID, MessageStatusQueued); err != nil {
					m.d.Logger().
						WithError(err).
						WithField("message_id", replace.ID).
						Error(`Unable to reset the failed message's status to "queued".`)
				}
			}

			return err
		}
	}

	return nil
}
