package identity

import (
	"context"
	"time"

	"github.com/ory/kratos/corp"

	"github.com/gofrs/uuid"
)

const (
	RecoveryAddressTypeEmail RecoveryAddressType = AddressTypeEmail
)

type (
	// RecoveryAddressType must not exceed 16 characters as that is the limitation in the SQL Schema.
	RecoveryAddressType string

	// RecoveryAddressStatus must not exceed 16 characters as that is the limitation in the SQL Schema.
	RecoveryAddressStatus string

	// swagger:model recoveryIdentityAddress
	RecoveryAddress struct {
		// required: true
		ID uuid.UUID `json:"id" db:"id" faker:"-"`

		// required: true
		Value string `json:"value" db:"value"`

		// required: true
		Via RecoveryAddressType `json:"via" db:"via"`

		// IdentityID is a helper struct field for gobuffalo.pop.
		IdentityID uuid.UUID `json:"-" faker:"-" db:"identity_id"`
		// CreatedAt is a helper struct field for gobuffalo.pop.
		CreatedAt time.Time `json:"-" faker:"-" db:"created_at"`
		// UpdatedAt is a helper struct field for gobuffalo.pop.
		UpdatedAt time.Time `json:"-" faker:"-" db:"updated_at"`
		NID       uuid.UUID `json:"-"  faker:"-" db:"nid"`
	}
)

func (v RecoveryAddressType) HTMLFormInputType() string {
	switch v {
	case RecoveryAddressTypeEmail:
		return "email"
	}
	return ""
}

func (a RecoveryAddress) TableName(ctx context.Context) string {
	return corp.ContextualizeTableName(ctx, "identity_recovery_addresses")
}

func NewRecoveryEmailAddress(
	value string,
	identity uuid.UUID,
) *RecoveryAddress {
	return &RecoveryAddress{
		Value:      value,
		Via:        RecoveryAddressTypeEmail,
		IdentityID: identity,
	}
}
