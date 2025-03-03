package continuity

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/herodot"
	"github.com/ory/x/sqlcon"

	"kratos/session"
	"kratos/x"
)

var _ Manager = new(ManagerCookie)
var ErrNotResumable = *herodot.ErrBadRequest.WithError("session is not resumable").WithReasonf("No resumable session could be found in the HTTP Header.")

const cookieName = "ory_kratos_continuity"

type (
	managerCookieDependencies interface {
		PersistenceProvider
		x.CookieProvider
		session.ManagementProvider
	}
	ManagerCookie struct {
		d managerCookieDependencies
	}
)

func NewManagerCookie(d managerCookieDependencies) *ManagerCookie {
	return &ManagerCookie{d: d}
}

func (m *ManagerCookie) Pause(ctx context.Context, w http.ResponseWriter, r *http.Request, name string, opts ...ManagerOption) error {
	if len(name) == 0 {
		return errors.Errorf("continuity container name must be set")
	}

	o, err := newManagerOptions(opts)
	if err != nil {
		return err
	}
	c := NewContainer(name, *o)

	if err := x.SessionPersistValues(w, r, m.d.ContinuityCookieManager(ctx), cookieName, map[string]interface{}{
		name: c.ID.String(),
	}); err != nil {
		return err
	}

	if err := m.d.ContinuityPersister().SaveContinuitySession(ctx, c); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *ManagerCookie) Continue(ctx context.Context, w http.ResponseWriter, r *http.Request, name string, opts ...ManagerOption) (*Container, error) {
	container, err := m.container(ctx, r, name)
	if err != nil {
		return nil, err
	}

	o, err := newManagerOptions(opts)
	if err != nil {
		return nil, err
	}

	if err := container.Valid(o.iid); err != nil {
		return nil, err
	}

	if o.payloadRaw != nil && container.Payload != nil {
		if err := json.NewDecoder(bytes.NewBuffer(container.Payload)).Decode(o.payloadRaw); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	if err := m.d.ContinuityPersister().DeleteContinuitySession(ctx, container.ID); err != nil {
		return nil, err
	}

	if err := x.SessionUnsetKey(w, r, m.d.ContinuityCookieManager(ctx), cookieName, name); err != nil {
		return nil, err
	}

	return container, nil
}

func (m *ManagerCookie) sid(ctx context.Context, r *http.Request, name string) (uuid.UUID, error) {
	var sid uuid.UUID
	if s, err := x.SessionGetString(r, m.d.ContinuityCookieManager(ctx), cookieName, name); err != nil {
		return sid, errors.WithStack(ErrNotResumable.WithDebugf("%+v", err))
	} else if sid = x.ParseUUID(s); sid == uuid.Nil {
		return sid, errors.WithStack(ErrNotResumable.WithDebug("session id is not a valid uuid"))
	}

	return sid, nil
}

func (m *ManagerCookie) container(ctx context.Context, r *http.Request, name string) (*Container, error) {
	sid, err := m.sid(ctx, r, name)
	if err != nil {
		return nil, err
	}

	container, err := m.d.ContinuityPersister().GetContinuitySession(ctx, sid)
	if errors.Is(err, sqlcon.ErrNoRows) {
		return nil, errors.WithStack(ErrNotResumable.WithDebugf("Resumable ID from cookie could not be found in the datastore: %+v", err))
	} else if err != nil {
		return nil, err
	}

	return container, err
}

func (m ManagerCookie) Abort(ctx context.Context, w http.ResponseWriter, r *http.Request, name string) error {
	sid, err := m.sid(ctx, r, name)
	if errors.Is(err, &ErrNotResumable) {
		return nil
	} else if err != nil {
		return err
	}

	if err := x.SessionUnsetKey(w, r, m.d.ContinuityCookieManager(ctx), cookieName, name); err != nil {
		return err
	}

	return errors.WithStack(m.d.ContinuityPersister().DeleteContinuitySession(ctx, sid))
}
