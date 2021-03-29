package persistence

import (
	"context"
	"io"

	"github.com/gobuffalo/pop/v5"

	"kratos/continuity"
	"kratos/courier"
	"kratos/identity"
	"kratos/selfservice/errorx"
	"kratos/selfservice/flow/login"
	"kratos/selfservice/flow/recovery"
	"kratos/selfservice/flow/registration"
	"kratos/selfservice/flow/settings"
	"kratos/selfservice/flow/verification"
	"kratos/selfservice/strategy/link"
	"kratos/session"
)

type Provider interface {
	Persister() Persister
}

type Persister interface {
	continuity.Persister
	identity.PrivilegedPool
	registration.FlowPersister
	login.FlowPersister
	settings.FlowPersister
	courier.Persister
	session.Persister
	errorx.Persister
	verification.FlowPersister
	recovery.FlowPersister
	link.RecoveryTokenPersister
	link.VerificationTokenPersister

	Close(context.Context) error
	Ping(context.Context) error
	MigrationStatus(c context.Context, b io.Writer) error
	MigrateDown(c context.Context, steps int) error
	MigrateUp(c context.Context) error
	GetConnection(ctx context.Context) *pop.Connection
	Transaction(ctx context.Context, callback func(ctx context.Context, connection *pop.Connection) error) error
}
