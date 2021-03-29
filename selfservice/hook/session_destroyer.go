package hook

import (
	"net/http"

	"kratos/selfservice/flow/login"
	"kratos/session"
)

var _ login.PostHookExecutor = new(SessionDestroyer)

type (
	sessionDestroyerDependencies interface {
		session.ManagementProvider
		session.PersistenceProvider
	}
	SessionDestroyer struct {
		r sessionDestroyerDependencies
	}
)

func NewSessionDestroyer(r sessionDestroyerDependencies) *SessionDestroyer {
	return &SessionDestroyer{r: r}
}

func (e *SessionDestroyer) ExecuteLoginPostHook(_ http.ResponseWriter, r *http.Request, _ *login.Flow, s *session.Session) error {
	if err := e.r.SessionPersister().DeleteSessionsByIdentity(r.Context(), s.Identity.ID); err != nil {
		return err
	}
	return nil
}
