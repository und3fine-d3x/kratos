package oidc_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/viper"

	"kratos/driver/configuration"
	"kratos/identity"
	"kratos/internal"
	"kratos/selfservice/strategy/oidc"
)

func TestConfig(t *testing.T) {
	conf, reg := internal.NewFastRegistryWithMocks(t)

	var c map[string]interface{}
	require.NoError(t, json.NewDecoder(
		bytes.NewBufferString(`{"config":{"providers": [{"provider": "generic"}]}}`)).Decode(&c))
	viper.Set(configuration.ViperKeySelfServiceStrategyConfig+"."+string(identity.CredentialsTypeOIDC), c)

	s := oidc.NewStrategy(reg, conf)
	collection, err := s.Config()
	require.NoError(t, err)

	require.Len(t, collection.Providers, 1)
	assert.Equal(t, "generic", collection.Providers[0].Provider)
}
