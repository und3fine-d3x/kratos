package password

import (
	"encoding/json"

	"kratos/ui/node"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/ory/x/decoderx"

	"kratos/continuity"
	"kratos/driver/config"
	"kratos/hash"
	"kratos/identity"
	"kratos/selfservice/errorx"
	"kratos/selfservice/flow/login"
	"kratos/selfservice/flow/registration"
	"kratos/selfservice/flow/settings"
	"kratos/session"
	"kratos/x"
)

var _ login.Strategy = new(Strategy)
var _ registration.Strategy = new(Strategy)
var _ identity.ActiveCredentialsCounter = new(Strategy)

type registrationStrategyDependencies interface {
	x.LoggingProvider
	x.WriterProvider
	x.CSRFTokenGeneratorProvider
	x.CSRFProvider

	config.Provider

	continuity.ManagementProvider

	errorx.ManagementProvider
	ValidationProvider
	hash.HashProvider

	registration.HandlerProvider
	registration.HooksProvider
	registration.ErrorHandlerProvider
	registration.HookExecutorProvider
	registration.FlowPersistenceProvider

	login.HooksProvider
	login.ErrorHandlerProvider
	login.HookExecutorProvider
	login.FlowPersistenceProvider
	login.HandlerProvider

	settings.FlowPersistenceProvider
	settings.HookExecutorProvider
	settings.HooksProvider
	settings.ErrorHandlerProvider

	identity.PrivilegedPoolProvider
	identity.ValidationProvider

	session.HandlerProvider
	session.ManagementProvider
}

type Strategy struct {
	d  registrationStrategyDependencies
	v  *validator.Validate
	hd *decoderx.HTTP
}

func NewStrategy(d registrationStrategyDependencies) *Strategy {
	return &Strategy{
		d:  d,
		v:  validator.New(),
		hd: decoderx.NewHTTP(),
	}
}

func (s *Strategy) CountActiveCredentials(cc map[identity.CredentialsType]identity.Credentials) (count int, err error) {
	for _, c := range cc {
		if c.Type == s.ID() && len(c.Config) > 0 {
			var conf CredentialsConfig
			if err = json.Unmarshal(c.Config, &conf); err != nil {
				return 0, errors.WithStack(err)
			}

			if len(c.Identifiers) > 0 && len(c.Identifiers[0]) > 0 &&
				(hash.IsBcryptHash([]byte(conf.HashedPassword)) || hash.IsArgon2idHash([]byte(conf.HashedPassword))) {
				count++
			}
		}
	}
	return
}

func (s *Strategy) ID() identity.CredentialsType {
	return identity.CredentialsTypePassword
}

func (s *Strategy) NodeGroup() node.Group {
	return node.PasswordGroup
}
