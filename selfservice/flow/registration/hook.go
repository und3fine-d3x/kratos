package registration

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/x/sqlcon"

	"kratos/driver/config"
	"kratos/identity"
	"kratos/schema"
	"kratos/selfservice/flow"
	"kratos/session"
	"kratos/x"
)

type (
	PreHookExecutor interface {
		ExecuteRegistrationPreHook(w http.ResponseWriter, r *http.Request, a *Flow) error
	}
	PreHookExecutorFunc func(w http.ResponseWriter, r *http.Request, a *Flow) error

	PostHookPostPersistExecutor interface {
		ExecutePostRegistrationPostPersistHook(w http.ResponseWriter, r *http.Request, a *Flow, s *session.Session) error
	}
	PostHookPostPersistExecutorFunc func(w http.ResponseWriter, r *http.Request, a *Flow, s *session.Session) error

	PostHookPrePersistExecutor interface {
		ExecutePostRegistrationPrePersistHook(w http.ResponseWriter, r *http.Request, a *Flow, i *identity.Identity) error
	}
	PostHookPrePersistExecutorFunc func(w http.ResponseWriter, r *http.Request, a *Flow, i *identity.Identity) error

	HooksProvider interface {
		PreRegistrationHooks(ctx context.Context) []PreHookExecutor
		PostRegistrationPrePersistHooks(ctx context.Context, credentialsType identity.CredentialsType) []PostHookPrePersistExecutor
		PostRegistrationPostPersistHooks(ctx context.Context, credentialsType identity.CredentialsType) []PostHookPostPersistExecutor
	}
)

func PostHookPostPersistExecutorNames(e []PostHookPostPersistExecutor) []string {
	names := make([]string, len(e))
	for k, ee := range e {
		names[k] = fmt.Sprintf("%T", ee)
	}
	return names
}

func (f PreHookExecutorFunc) ExecuteRegistrationPreHook(w http.ResponseWriter, r *http.Request, a *Flow) error {
	return f(w, r, a)
}
func (f PostHookPostPersistExecutorFunc) ExecutePostRegistrationPostPersistHook(w http.ResponseWriter, r *http.Request, a *Flow, s *session.Session) error {
	return f(w, r, a, s)
}
func (f PostHookPrePersistExecutorFunc) ExecutePostRegistrationPrePersistHook(w http.ResponseWriter, r *http.Request, a *Flow, i *identity.Identity) error {
	return f(w, r, a, i)
}

type (
	executorDependencies interface {
		config.Provider
		identity.ManagementProvider
		identity.ValidationProvider
		session.PersistenceProvider
		HooksProvider
		x.LoggingProvider
		x.WriterProvider
	}
	HookExecutor struct {
		d executorDependencies
	}
	HookExecutorProvider interface {
		RegistrationExecutor() *HookExecutor
	}
)

func NewHookExecutor(d executorDependencies) *HookExecutor {
	return &HookExecutor{d: d}
}

func (e *HookExecutor) PostRegistrationHook(w http.ResponseWriter, r *http.Request, ct identity.CredentialsType, a *Flow, i *identity.Identity) error {
	e.d.Logger().
		WithRequest(r).
		WithField("identity_id", i.ID).
		WithField("flow_method", ct).
		Debug("Running PostRegistrationPrePersistHooks.")
	for k, executor := range e.d.PostRegistrationPrePersistHooks(r.Context(), ct) {
		if err := executor.ExecutePostRegistrationPrePersistHook(w, r, a, i); err != nil {
			if errors.Is(err, ErrHookAbortFlow) {
				e.d.Logger().
					WithRequest(r).
					WithField("executor", fmt.Sprintf("%T", executor)).
					WithField("executor_position", k).
					WithField("executors", PostHookPostPersistExecutorNames(e.d.PostRegistrationPostPersistHooks(r.Context(), ct))).
					WithField("identity_id", i.ID).
					WithField("flow_method", ct).
					Debug("A ExecutePostRegistrationPrePersistHook hook aborted early.")
				return nil
			}
			return err
		}

		e.d.Logger().WithRequest(r).
			WithField("executor", fmt.Sprintf("%T", executor)).
			WithField("executor_position", k).
			WithField("executors", PostHookPostPersistExecutorNames(e.d.PostRegistrationPostPersistHooks(r.Context(), ct))).
			WithField("identity_id", i.ID).
			WithField("flow_method", ct).
			Debug("ExecutePostRegistrationPrePersistHook completed successfully.")
	}

	// We need to make sure that the identity has a valid schema before passing it down to the identity pool.
	if err := e.d.IdentityValidator().Validate(r.Context(), i); err != nil {
		return err
		// We're now creating the identity because any of the hooks could trigger a "redirect" or a "session" which
		// would imply that the identity has to exist already.
	} else if err := e.d.IdentityManager().Create(r.Context(), i); err != nil {
		if errors.Is(err, sqlcon.ErrUniqueViolation) {
			return schema.NewDuplicateCredentialsError()
		}
		return err
	}

	e.d.Audit().
		WithRequest(r).
		WithField("identity_id", i.ID).
		Info("A new identity has registered using self-service registration.")

	s := session.NewActiveSession(i, e.d.Config(r.Context()), time.Now().UTC())
	e.d.Logger().
		WithRequest(r).
		WithField("identity_id", i.ID).
		WithField("flow_method", ct).
		Debug("Running PostRegistrationPostPersistHooks.")
	for k, executor := range e.d.PostRegistrationPostPersistHooks(r.Context(), ct) {
		if err := executor.ExecutePostRegistrationPostPersistHook(w, r, a, s); err != nil {
			if errors.Is(err, ErrHookAbortFlow) {
				e.d.Logger().
					WithRequest(r).
					WithField("executor", fmt.Sprintf("%T", executor)).
					WithField("executor_position", k).
					WithField("executors", PostHookPostPersistExecutorNames(e.d.PostRegistrationPostPersistHooks(r.Context(), ct))).
					WithField("identity_id", i.ID).
					WithField("flow_method", ct).
					Debug("A ExecutePostRegistrationPostPersistHook hook aborted early.")
				return nil
			}
			return err
		}

		e.d.Logger().WithRequest(r).
			WithField("executor", fmt.Sprintf("%T", executor)).
			WithField("executor_position", k).
			WithField("executors", PostHookPostPersistExecutorNames(e.d.PostRegistrationPostPersistHooks(r.Context(), ct))).
			WithField("identity_id", i.ID).
			WithField("flow_method", ct).
			Debug("ExecutePostRegistrationPostPersistHook completed successfully.")
	}

	e.d.Logger().
		WithRequest(r).
		WithField("flow_method", ct).
		WithField("identity_id", i.ID).
		Debug("Post registration execution hooks completed successfully.")

	if a.Type == flow.TypeAPI {
		e.d.Writer().Write(w, r, &APIFlowResponse{Identity: i})
		return nil
	}

	return x.SecureContentNegotiationRedirection(w, r, s.Declassify(), a.RequestURL,
		e.d.Writer(), e.d.Config(r.Context()), x.SecureRedirectOverrideDefaultReturnTo(e.d.Config(r.Context()).SelfServiceFlowRegistrationReturnTo(ct.String())))
}

func (e *HookExecutor) PreRegistrationHook(w http.ResponseWriter, r *http.Request, a *Flow) error {
	for _, executor := range e.d.PreRegistrationHooks(r.Context()) {
		if err := executor.ExecuteRegistrationPreHook(w, r, a); err != nil {
			return err
		}
	}

	return nil
}
