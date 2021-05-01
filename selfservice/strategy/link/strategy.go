package link

import (
	"github.com/ory/x/decoderx"
	"kratos/courier"
	"kratos/driver/config"
	"kratos/identity"
	"kratos/schema"
	"kratos/selfservice/errorx"
	"kratos/selfservice/flow/recovery"
	"kratos/selfservice/flow/settings"
	"kratos/selfservice/flow/verification"
	"kratos/session"
	"kratos/ui/container"
	"kratos/ui/node"
	"kratos/x"
)

var _ recovery.Strategy = new(Strategy)
var _ recovery.AdminHandler = new(Strategy)
var _ recovery.PublicHandler = new(Strategy)

var _ verification.Strategy = new(Strategy)
var _ verification.AdminHandler = new(Strategy)
var _ verification.PublicHandler = new(Strategy)

type (
	// FlowMethod contains the configuration for this selfservice strategy.
	FlowMethod struct {
		*container.Container
	}

	strategyDependencies interface {
		x.CSRFProvider
		x.CSRFTokenGeneratorProvider
		x.WriterProvider
		x.LoggingProvider

		config.Provider

		session.HandlerProvider
		session.ManagementProvider
		settings.HandlerProvider
		settings.FlowPersistenceProvider

		identity.ValidationProvider
		identity.ManagementProvider
		identity.PoolProvider
		identity.PrivilegedPoolProvider

		courier.Provider

		errorx.ManagementProvider

		recovery.ErrorHandlerProvider
		recovery.FlowPersistenceProvider
		recovery.StrategyProvider

		verification.ErrorHandlerProvider
		verification.FlowPersistenceProvider
		verification.StrategyProvider

		RecoveryTokenPersistenceProvider
		VerificationTokenPersistenceProvider
		SenderProvider

		schema.IdentityTraitsProvider
	}

	Strategy struct {
		d  strategyDependencies
		dx *decoderx.HTTP
	}
)

func NewStrategy(d strategyDependencies) *Strategy {
	return &Strategy{d: d, dx: decoderx.NewHTTP()}
}

func (s *Strategy) RecoveryNodeGroup() node.Group {
	return node.RecoveryLinkGroup
}

func (s *Strategy) VerificationNodeGroup() node.Group {
	return node.VerificationLinkGroup
}
