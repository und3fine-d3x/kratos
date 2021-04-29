package oidc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kratos/driver/config"
	"kratos/identity"
	"kratos/internal"
	"kratos/selfservice/strategy/oidc"
)

func TestConfig(t *testing.T) {
	conf, reg := internal.NewFastRegistryWithMocks(t)

	var c map[string]interface{}
	require.NoError(t, json.NewDecoder(
		bytes.NewBufferString(`{"config":{"providers": [{"provider": "generic"}]}}`)).Decode(&c))
	conf.MustSet(config.ViperKeySelfServiceStrategyConfig+"."+string(identity.CredentialsTypeOIDC), c)

	s := oidc.NewStrategy(reg)
	collection, err := s.Config(context.Background())
	require.NoError(t, err)

	require.Len(t, collection.Providers, 1)
	assert.Equal(t, "generic", collection.Providers[0].Provider)
}
