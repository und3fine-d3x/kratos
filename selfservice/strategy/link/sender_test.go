package link_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/urlx"
	"kratos/driver/config"
	"kratos/identity"
	"kratos/internal"
	"kratos/selfservice/flow"
	"kratos/selfservice/flow/recovery"
	"kratos/selfservice/flow/verification"
	"kratos/selfservice/strategy/link"
)

func TestManager(t *testing.T) {
	conf, reg := internal.NewFastRegistryWithMocks(t)
	conf.MustSet(config.ViperKeyDefaultIdentitySchemaURL, "file://./stub/default.schema.json")
	conf.MustSet(config.ViperKeyPublicBaseURL, "https://www.ory.sh/")
	conf.MustSet(config.ViperKeyCourierSMTPURL, "smtp://foo@bar@dev.null/")

	u := &http.Request{URL: urlx.ParseOrPanic("https://www.ory.sh/")}

	i := identity.NewIdentity(config.DefaultIdentityTraitsSchemaID)
	i.Traits = identity.Traits(`{"email": "tracked@ory.sh"}`)
	require.NoError(t, reg.IdentityManager().Create(context.Background(), i))

	t.Run("method=SendRecoveryLink", func(t *testing.T) {
		f, err := recovery.NewFlow(conf, time.Hour, "", u, reg.RecoveryStrategies(context.Background()), flow.TypeBrowser)
		require.NoError(t, err)

		require.NoError(t, reg.RecoveryFlowPersister().CreateRecoveryFlow(context.Background(), f))

		require.NoError(t, reg.LinkSender().SendRecoveryLink(context.Background(), nil, f, "email", "tracked@ory.sh"))
		require.EqualError(t, reg.LinkSender().SendRecoveryLink(context.Background(), nil, f, "email", "not-tracked@ory.sh"), link.ErrUnknownAddress.Error())

		messages, err := reg.CourierPersister().NextMessages(context.Background(), 12)
		require.NoError(t, err)
		require.Len(t, messages, 2)

		assert.EqualValues(t, "tracked@ory.sh", messages[0].Recipient)
		assert.Contains(t, messages[0].Subject, "Recover access to your account")
		assert.Contains(t, messages[0].Body, urlx.AppendPaths(conf.SelfPublicURL(nil), recovery.RouteSubmitFlow).String()+"?")
		assert.Contains(t, messages[0].Body, "token=")
		assert.Contains(t, messages[0].Body, "flow=")

		assert.EqualValues(t, "not-tracked@ory.sh", messages[1].Recipient)
		assert.Contains(t, messages[1].Subject, "Account access attempted")
		assert.NotContains(t, messages[1].Body, urlx.AppendPaths(conf.SelfPublicURL(nil), recovery.RouteSubmitFlow).String()+"?")
		assert.NotContains(t, messages[1].Body, "token=")
		assert.NotContains(t, messages[1].Body, "flow=")
	})

	t.Run("method=SendVerificationLink", func(t *testing.T) {
		f, err := verification.NewFlow(conf, time.Hour, "", u, reg.VerificationStrategies(context.Background()), flow.TypeBrowser)
		require.NoError(t, err)

		require.NoError(t, reg.VerificationFlowPersister().CreateVerificationFlow(context.Background(), f))

		require.NoError(t, reg.LinkSender().SendVerificationLink(context.Background(), f, "email", "tracked@ory.sh"))
		require.EqualError(t, reg.LinkSender().SendVerificationLink(context.Background(), f, "email", "not-tracked@ory.sh"), link.ErrUnknownAddress.Error())

		messages, err := reg.CourierPersister().NextMessages(context.Background(), 12)
		require.NoError(t, err)
		require.Len(t, messages, 2)

		assert.EqualValues(t, "tracked@ory.sh", messages[0].Recipient)
		assert.Contains(t, messages[0].Subject, "Please verify")
		assert.Contains(t, messages[0].Body, urlx.AppendPaths(conf.SelfPublicURL(nil), verification.RouteSubmitFlow).String()+"?")
		assert.Contains(t, messages[0].Body, "token=")
		assert.Contains(t, messages[0].Body, "flow=")

		assert.EqualValues(t, "not-tracked@ory.sh", messages[1].Recipient)
		assert.Contains(t, messages[1].Subject, "tried to verify")
		assert.NotContains(t, messages[1].Body, urlx.AppendPaths(conf.SelfPublicURL(nil), verification.RouteSubmitFlow).String()+"?")
	})
}
