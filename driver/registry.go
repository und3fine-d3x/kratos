package driver

import (
	"context"

	"github.com/ory/x/tracing"
	"kratos/metrics/prometheus"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"

	"github.com/ory/x/logrusx"

	"kratos/continuity"
	"kratos/courier"
	"kratos/hash"
	"kratos/schema"
	"kratos/selfservice/flow/recovery"
	"kratos/selfservice/flow/settings"
	"kratos/selfservice/flow/verification"
	"kratos/selfservice/strategy/link"

	"github.com/ory/x/healthx"

	"kratos/persistence"
	"kratos/selfservice/flow/login"
	"kratos/selfservice/flow/logout"
	"kratos/selfservice/flow/registration"

	"kratos/x"

	"github.com/ory/x/dbal"

	"kratos/driver/config"
	"kratos/identity"
	"kratos/selfservice/errorx"
	password2 "kratos/selfservice/strategy/password"
	"kratos/session"
)

type Registry interface {
	dbal.Driver

	Init(ctx context.Context) error

	WithLogger(l *logrusx.Logger) Registry

	WithCSRFHandler(c x.CSRFHandler)
	WithCSRFTokenGenerator(cg x.CSRFToken)

	HealthHandler(ctx context.Context) *healthx.Handler
	CookieManager(ctx context.Context) sessions.Store
	MetricsHandler() *prometheus.Handler
	ContinuityCookieManager(ctx context.Context) sessions.Store

	RegisterRoutes(ctx context.Context, public *x.RouterPublic, admin *x.RouterAdmin)
	RegisterPublicRoutes(ctx context.Context, public *x.RouterPublic)
	RegisterAdminRoutes(ctx context.Context, admin *x.RouterAdmin)
	PrometheusManager() *prometheus.MetricsManager
	Tracer(context.Context) *tracing.Tracer

	config.Provider
	WithConfig(c *config.Config) Registry

	x.CSRFProvider
	x.WriterProvider
	x.LoggingProvider

	continuity.ManagementProvider
	continuity.PersistenceProvider

	courier.Provider

	persistence.Provider

	errorx.ManagementProvider
	errorx.HandlerProvider
	errorx.PersistenceProvider

	hash.HashProvider

	identity.HandlerProvider
	identity.ValidationProvider
	identity.PoolProvider
	identity.PrivilegedPoolProvider
	identity.ManagementProvider
	identity.ActiveCredentialsCounterStrategyProvider

	schema.HandlerProvider

	password2.ValidationProvider

	session.HandlerProvider
	session.ManagementProvider
	session.PersistenceProvider

	settings.HandlerProvider
	settings.ErrorHandlerProvider
	settings.FlowPersistenceProvider
	settings.StrategyProvider

	login.FlowPersistenceProvider
	login.ErrorHandlerProvider
	login.HooksProvider
	login.HookExecutorProvider
	login.HandlerProvider
	login.StrategyProvider

	logout.HandlerProvider

	registration.FlowPersistenceProvider
	registration.ErrorHandlerProvider
	registration.HooksProvider
	registration.HookExecutorProvider
	registration.HandlerProvider
	registration.StrategyProvider

	verification.FlowPersistenceProvider
	verification.ErrorHandlerProvider
	verification.HandlerProvider
	verification.StrategyProvider

	link.SenderProvider
	link.VerificationTokenPersistenceProvider
	link.RecoveryTokenPersistenceProvider

	recovery.FlowPersistenceProvider
	recovery.ErrorHandlerProvider
	recovery.HandlerProvider
	recovery.StrategyProvider

	x.CSRFTokenGeneratorProvider
}

func NewRegistryFromDSN(c *config.Config, l *logrusx.Logger) (Registry, error) {
	driver, err := dbal.GetDriverFor(c.DSN())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	registry, ok := driver.(Registry)
	if !ok {
		return nil, errors.Errorf("driver of type %T does not implement interface Registry", driver)
	}

	return registry.WithLogger(l).WithConfig(c), nil
}
