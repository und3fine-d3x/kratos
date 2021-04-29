package persistence

import (
	"context"

	"github.com/ory/x/networkx"

	"github.com/gofrs/uuid"

	"github.com/gobuffalo/pop/v5"

	"github.com/ory/x/popx"

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
	Ping() error
	MigrationStatus(c context.Context) (popx.MigrationStatuses, error)
	MigrateDown(c context.Context, steps int) error
	MigrateUp(c context.Context) error
	Migrator() *popx.Migrator
	GetConnection(ctx context.Context) *pop.Connection
	Transaction(ctx context.Context, callback func(ctx context.Context, connection *pop.Connection) error) error
	Networker
}

type Networker interface {
	WithNetworkID(sid uuid.UUID) Persister
	NetworkID() uuid.UUID
	DetermineNetwork(ctx context.Context) (*networkx.Network, error)
}
