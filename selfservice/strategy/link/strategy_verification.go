package link

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/x/decoderx"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/sqlxx"
	"github.com/ory/x/urlx"
	"kratos/identity"
	"kratos/schema"
	"kratos/selfservice/flow"
	"kratos/selfservice/flow/verification"
	"kratos/text"
	"kratos/ui/node"
	"kratos/x"
)

func (s *Strategy) VerificationStrategyID() string {
	return verification.StrategyVerificationLinkName
}

func (s *Strategy) RegisterPublicVerificationRoutes(public *x.RouterPublic) {
}

func (s *Strategy) RegisterAdminVerificationRoutes(admin *x.RouterAdmin) {
}

func (s *Strategy) PopulateVerificationMethod(r *http.Request, f *verification.Flow) error {
	f.UI.SetCSRF(s.d.GenerateCSRFToken(r))
	f.UI.GetNodes().Upsert(
		// v0.5: form.Field{Name: "email", Type: "email", Required: true}
		node.NewInputField("email", nil, node.VerificationLinkGroup, node.InputAttributeTypeEmail, node.WithRequiredInputAttribute),
	)
	f.UI.GetNodes().Append(node.NewInputField("method", s.VerificationStrategyID(), node.VerificationLinkGroup, node.InputAttributeTypeSubmit).WithMetaLabel(text.NewInfoNodeLabelSubmit()))
	return nil
}

type verificationSubmitPayload struct {
	Method    string `json:"method" form:"method"`
	Token     string `json:"token" form:"token"`
	CSRFToken string `json:"csrf_token" form:"csrf_token"`
	Flow      string `json:"flow" form:"flow"`
	Email     string `json:"email" form:"email"`
}

func (s *Strategy) decodeVerification(r *http.Request) (*verificationSubmitPayload, error) {
	var body verificationSubmitPayload

	compiler, err := decoderx.HTTPRawJSONSchemaCompiler(verificationMethodSchema)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := s.dx.Decode(r, &body, compiler,
		decoderx.HTTPDecoderUseQueryAndBody(),
		decoderx.HTTPKeepRequestBody(true),
		decoderx.HTTPDecoderAllowedMethods("POST", "GET"),
		decoderx.HTTPDecoderSetValidatePayloads(false),
		decoderx.HTTPDecoderJSONFollowsFormFormat(),
	); err != nil {
		return nil, errors.WithStack(err)
	}

	return &body, nil
}

// handleVerificationError is a convenience function for handling all types of errors that may occur (e.g. validation error).
func (s *Strategy) handleVerificationError(w http.ResponseWriter, r *http.Request, f *verification.Flow, body *verificationSubmitPayload, err error) error {
	if f != nil {
		f.UI.SetCSRF(s.d.GenerateCSRFToken(r))
		f.UI.GetNodes().Upsert(
			// v0.5: form.Field{Name: "email", Type: "email", Required: true, Value: body.Body.Email}
			node.NewInputField("email", body.Email, node.VerificationLinkGroup, node.InputAttributeTypeEmail, node.WithRequiredInputAttribute),
		)
	}

	return err
}

// swagger:parameters submitSelfServiceVerificationFlowWithLinkMethod
// nolint:deadcode,unused
type submitSelfServiceVerificationFlowWithLinkMethodParameters struct {
	// in: body
	Body submitSelfServiceVerificationFlowWithLinkMethod

	// Verification Token
	//
	// The verification token which completes the verification request. If the token
	// is invalid (e.g. expired) an error will be shown to the end-user.
	//
	// in: query
	Token string `json:"token"`

	// The Flow ID
	//
	// format: uuid
	// in: query
	Flow string `json:"flow"`
}

func (m *verificationSubmitPayload) GetFlow() uuid.UUID {
	return x.ParseUUID(m.Flow)
}

// nolint:deadcode,unused
type submitSelfServiceVerificationFlowWithLinkMethod struct {
	// Email to Verify
	//
	// Needs to be set when initiating the flow. If the email is a registered
	// verification email, a verification link will be sent. If the email is not known,
	// a email with details on what happened will be sent instead.
	//
	// format: email
	// in: body
	Email string `json:"email"`

	// Sending the anti-csrf token is only required for browser login flows.
	CSRFToken string `form:"csrf_token" json:"csrf_token"`
}

func (s *Strategy) Verify(w http.ResponseWriter, r *http.Request, f *verification.Flow) (err error) {
	body, err := s.decodeVerification(r)
	if err != nil {
		return s.handleVerificationError(w, r, nil, body, err)
	}

	if len(body.Token) > 0 {
		if err := flow.MethodEnabledAndAllowed(r.Context(), s.VerificationStrategyID(), s.VerificationStrategyID(), s.d); err != nil {
			return s.handleVerificationError(w, r, nil, body, err)
		}

		return s.verificationUseToken(w, r, body)
	}

	if err := flow.MethodEnabledAndAllowed(r.Context(), s.VerificationStrategyID(), body.Method, s.d); err != nil {
		return s.handleVerificationError(w, r, f, body, err)
	}

	if err := f.Valid(); err != nil {
		return s.handleVerificationError(w, r, f, body, err)
	}

	switch f.State {
	case verification.StateChooseMethod:
		fallthrough
	case verification.StateEmailSent:
		// Do nothing (continue with execution after this switch statement)
		return s.verificationHandleFormSubmission(w, r, f)
	case verification.StatePassedChallenge:
		return s.retryVerificationFlowWithMessage(w, r, f.Type, text.NewErrorValidationVerificationRetrySuccess())
	default:
		return s.retryVerificationFlowWithMessage(w, r, f.Type, text.NewErrorValidationVerificationStateFailure())
	}
}

func (s *Strategy) verificationHandleFormSubmission(w http.ResponseWriter, r *http.Request, f *verification.Flow) error {
	var body = new(verificationSubmitPayload)
	body, err := s.decodeVerification(r)
	if err != nil {
		return s.handleVerificationError(w, r, f, body, err)
	}

	if len(body.Email) == 0 {
		return s.handleVerificationError(w, r, f, body, schema.NewRequiredError("#/email", "email"))
	}

	if err := flow.EnsureCSRF(r, f.Type, s.d.Config(r.Context()).DisableAPIFlowEnforcement(), s.d.GenerateCSRFToken, body.CSRFToken); err != nil {
		return s.handleVerificationError(w, r, f, body, err)
	}

	if err := s.d.LinkSender().SendVerificationLink(r.Context(), f, identity.VerifiableAddressTypeEmail, body.Email); err != nil {
		if !errors.Is(err, ErrUnknownAddress) {
			return s.handleVerificationError(w, r, f, body, err)
		}
		// Continue execution
	}

	f.UI.SetCSRF(s.d.GenerateCSRFToken(r))
	f.UI.GetNodes().Upsert(
		// v0.5: form.Field{Name: "email", Type: "email", Required: true, Value: body.Body.Email}
		node.NewInputField("email", body.Email, node.VerificationLinkGroup, node.InputAttributeTypeEmail, node.WithRequiredInputAttribute),
	)

	f.Active = sqlxx.NullString(s.VerificationNodeGroup())
	f.State = verification.StateEmailSent
	f.UI.Messages.Set(text.NewVerificationEmailSent())
	if err := s.d.VerificationFlowPersister().UpdateVerificationFlow(r.Context(), f); err != nil {
		return s.handleVerificationError(w, r, f, body, err)
	}

	return nil
}

// nolint:deadcode,unused
// swagger:parameters selfServiceBrowserVerify
type selfServiceBrowserVerifyParameters struct {
	// required: true
	// in: query
	Token string `json:"token"`
}

func (s *Strategy) verificationUseToken(w http.ResponseWriter, r *http.Request, body *verificationSubmitPayload) error {
	token, err := s.d.VerificationTokenPersister().UseVerificationToken(r.Context(), body.Token)
	if err != nil {
		if errors.Is(err, sqlcon.ErrNoRows) {
			return s.retryVerificationFlowWithMessage(w, r, flow.TypeBrowser, text.NewErrorValidationVerificationTokenInvalidOrAlreadyUsed())
		}

		return s.handleVerificationError(w, r, nil, body, err)
	}

	var f *verification.Flow
	if !token.FlowID.Valid {
		f, err = verification.NewFlow(s.d.Config(r.Context()), time.Until(token.ExpiresAt), s.d.GenerateCSRFToken(r), r, s.d.VerificationStrategies(r.Context()), flow.TypeBrowser)
		if err != nil {
			return s.handleVerificationError(w, r, nil, body, err)
		}

		if err := s.d.VerificationFlowPersister().CreateVerificationFlow(r.Context(), f); err != nil {
			return s.handleVerificationError(w, r, nil, body, err)
		}
	} else {
		f, err = s.d.VerificationFlowPersister().GetVerificationFlow(r.Context(), token.FlowID.UUID)
		if err != nil {
			return s.handleVerificationError(w, r, nil, body, err)
		}
	}

	if err := token.Valid(); err != nil {
		return s.handleVerificationError(w, r, f, body, err)
	}

	f.UI.Messages.Clear()
	f.State = verification.StatePassedChallenge
	if err := s.d.VerificationFlowPersister().UpdateVerificationFlow(r.Context(), f); err != nil {
		return s.handleVerificationError(w, r, f, body, err)
	}

	address := token.VerifiableAddress
	address.Verified = true
	address.VerifiedAt = sqlxx.NullTime(time.Now().UTC())
	address.Status = identity.VerifiableAddressStatusCompleted
	if err := s.d.PrivilegedIdentityPool().UpdateVerifiableAddress(r.Context(), address); err != nil {
		return s.handleVerificationError(w, r, f, body, err)
	}

	defaultRedirectURL := s.d.Config(r.Context()).SelfServiceFlowVerificationReturnTo(f.AppendTo(s.d.Config(r.Context()).SelfServiceFlowVerificationUI()))

	verificationRequestURL, err := urlx.Parse(f.GetRequestURL())
	if err != nil {
		s.d.Logger().Debugf("error parsing verification requestURL: %s\n", err)
		http.Redirect(w, r, defaultRedirectURL.String(), http.StatusFound)
		return errors.WithStack(flow.ErrCompletedByStrategy)
	}
	verificationRequest := http.Request{URL: verificationRequestURL}

	returnTo, err := x.SecureRedirectTo(&verificationRequest, defaultRedirectURL,
		x.SecureRedirectAllowSelfServiceURLs(s.d.Config(r.Context()).SelfPublicURL(r)),
		x.SecureRedirectAllowURLs(s.d.Config(r.Context()).SelfServiceBrowserWhitelistedReturnToDomains()),
	)
	if err != nil {
		s.d.Logger().Debugf("error parsing redirectTo from verification: %s\n", err)
		http.Redirect(w, r, defaultRedirectURL.String(), http.StatusFound)
		return errors.WithStack(flow.ErrCompletedByStrategy)
	}

	http.Redirect(w, r, returnTo.String(), http.StatusFound)
	return errors.WithStack(flow.ErrCompletedByStrategy)
}

func (s *Strategy) retryVerificationFlowWithMessage(w http.ResponseWriter, r *http.Request, ft flow.Type, message *text.Message) error {
	s.d.Logger().WithRequest(r).WithField("message", message).Debug("A verification flow is being retried because a validation error occurred.")

	f, err := verification.NewFlow(s.d.Config(r.Context()),
		s.d.Config(r.Context()).SelfServiceFlowVerificationRequestLifespan(), s.d.GenerateCSRFToken(r), r, s.d.VerificationStrategies(r.Context()), ft)
	if err != nil {
		return s.handleVerificationError(w, r, f, nil, err)
	}

	f.UI.Messages.Add(message)
	if err := s.d.VerificationFlowPersister().CreateVerificationFlow(r.Context(), f); err != nil {
		return s.handleVerificationError(w, r, f, nil, err)
	}

	if ft == flow.TypeBrowser {
		http.Redirect(w, r, f.AppendTo(s.d.Config(r.Context()).SelfServiceFlowVerificationUI()).String(), http.StatusFound)
	} else {
		http.Redirect(w, r, urlx.CopyWithQuery(urlx.AppendPaths(s.d.Config(r.Context()).SelfPublicURL(r),
			verification.RouteGetFlow), url.Values{"id": {f.ID.String()}}).String(), http.StatusFound)
	}

	return errors.WithStack(flow.ErrCompletedByStrategy)
}
