package hook

import (
	"net/http"

	"kratos/driver/config"
	"kratos/identity"
	"kratos/selfservice/flow"
	"kratos/selfservice/flow/registration"
	"kratos/selfservice/flow/settings"
	"kratos/selfservice/flow/verification"
	"kratos/selfservice/strategy/link"
	"kratos/session"
	"kratos/x"
)

var _ registration.PostHookPostPersistExecutor = new(Verifier)
var _ settings.PostHookPostPersistExecutor = new(Verifier)

type (
	verifierDependencies interface {
		link.SenderProvider
		link.VerificationTokenPersistenceProvider
		config.Provider
		x.CSRFTokenGeneratorProvider
		verification.StrategyProvider
		verification.FlowPersistenceProvider
	}
	Verifier struct {
		r verifierDependencies
	}
)

func NewVerifier(r verifierDependencies) *Verifier {
	return &Verifier{r: r}
}

func (e *Verifier) ExecutePostRegistrationPostPersistHook(_ http.ResponseWriter, r *http.Request, f *registration.Flow, s *session.Session) error {
	return e.do(r, s.Identity, f)
}

func (e *Verifier) ExecuteSettingsPostPersistHook(w http.ResponseWriter, r *http.Request, a *settings.Flow, i *identity.Identity) error {
	return e.do(r, i, a)
}

func (e *Verifier) do(r *http.Request, i *identity.Identity, f flow.Flow) error {
	// Ths is called after the identity has been created so we can safely assume that all addresses are available
	// already.

	for k := range i.VerifiableAddresses {
		address := &i.VerifiableAddresses[k]
		if address.Verified {
			continue
		}

		verificationFlow, err := verification.NewPostHookFlow(e.r.Config(r.Context()),
			e.r.Config(r.Context()).SelfServiceFlowVerificationRequestLifespan(),
			e.r.GenerateCSRFToken(r), r, e.r.VerificationStrategies(r.Context()), f)
		if err != nil {
			return err
		}

		if err := e.r.VerificationFlowPersister().CreateVerificationFlow(r.Context(), verificationFlow); err != nil {
			return err
		}

		token := link.NewSelfServiceVerificationToken(address, verificationFlow)
		if err := e.r.VerificationTokenPersister().CreateVerificationToken(r.Context(), token); err != nil {
			return err
		}

		if err := e.r.LinkSender().SendVerificationTokenTo(r.Context(), verificationFlow, address, token); err != nil {
			return err
		}
	}

	return nil
}
