package driver

import (
	"kratos/metrics/prometheus"

	"github.com/ory/x/tracing"

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

	"kratos/driver/configuration"
	"kratos/identity"
	"kratos/selfservice/errorx"
	password2 "kratos/selfservice/strategy/password"
	"kratos/session"
)

type Registry interface {
	dbal.Driver

	Init() error

	WithConfig(c configuration.Provider) Registry
	WithLogger(l *logrusx.Logger) Registry

	BuildVersion() string
	BuildDate() string
	BuildHash() string
	WithBuildInfo(version, hash, date string) Registry

	WithCSRFHandler(c x.CSRFHandler)
	WithCSRFTokenGenerator(cg x.CSRFToken)

	HealthHandler() *healthx.Handler
	CookieManager() sessions.Store
	ContinuityCookieManager() sessions.Store

	RegisterRoutes(public *x.RouterPublic, admin *x.RouterAdmin)
	RegisterPublicRoutes(public *x.RouterPublic)
	RegisterAdminRoutes(admin *x.RouterAdmin)
	PrometheusManager() *prometheus.MetricsManager
	Tracer() *tracing.Tracer

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

func NewRegistry(c configuration.Provider) (Registry, error) {
	dsn := c.DSN()
	driver, err := dbal.GetDriverFor(dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	registry, ok := driver.(Registry)
	if !ok {
		return nil, errors.Errorf("driver of type %T does not implement interface Registry", driver)
	}

	return registry, nil
}
