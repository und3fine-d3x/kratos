package settings_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ory/kratos-client-go"

	"kratos/corpx"

	"github.com/ory/x/ioutilx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/x/urlx"

	"kratos/driver/config"
	"kratos/identity"
	"kratos/internal"
	"kratos/internal/testhelpers"
	"kratos/selfservice/flow"
	"kratos/selfservice/flow/login"
	"kratos/selfservice/flow/settings"
	"kratos/x"
)

func init() {
	corpx.RegisterFakes()
}

func TestHandler(t *testing.T) {
	conf, reg := internal.NewFastRegistryWithMocks(t)
	conf.MustSet(config.ViperKeyDefaultIdentitySchemaURL, "file://./stub/identity.schema.json")
	testhelpers.StrategyEnable(t, conf, identity.CredentialsTypePassword.String(), true)
	testhelpers.StrategyEnable(t, conf, settings.StrategyProfile, true)

	_ = testhelpers.NewSettingsUITestServer(t, conf)
	_ = testhelpers.NewErrorTestServer(t, reg)

	conf.MustSet(config.ViperKeySelfServiceSettingsPrivilegedAuthenticationAfter, "1ns")
	primaryIdentity := &identity.Identity{ID: x.NewUUID(), Traits: identity.Traits(`{}`)}
	publicTS, adminTS, clients := testhelpers.NewSettingsAPIServer(t, reg, map[string]*identity.Identity{
		"primary":   primaryIdentity,
		"secondary": {ID: x.NewUUID(), Traits: identity.Traits(`{}`)}})

	primaryUser, otherUser := clients["primary"], clients["secondary"]
	newExpiredFlow := func() *settings.Flow {
		return settings.NewFlow(conf, -time.Minute,
			&http.Request{URL: urlx.ParseOrPanic(publicTS.URL + login.RouteInitBrowserFlow)},
			primaryIdentity, flow.TypeBrowser)
	}

	t.Run("daemon=admin", func(t *testing.T) {
		t.Run("description=fetching a non-existent flow should return a 404 error", func(t *testing.T) {
			_, _, err := testhelpers.NewSDKCustomClient(adminTS, otherUser).PublicApi.GetSelfServiceSettingsFlow(context.Background()).Id("i-do-not-exist").Execute()
			require.Error(t, err)

			require.IsType(t, new(kratos.GenericOpenAPIError), err)
			assert.Equal(t, int64(http.StatusNotFound), gjson.GetBytes(err.(*kratos.GenericOpenAPIError).Body(), "error.code").Int())
		})

		t.Run("description=fetching an expired flow returns 410", func(t *testing.T) {
			pr := newExpiredFlow()
			require.NoError(t, reg.SettingsFlowPersister().CreateSettingsFlow(context.Background(), pr))

			_, _, err := testhelpers.NewSDKCustomClient(adminTS, primaryUser).PublicApi.GetSelfServiceSettingsFlow(context.Background()).Id(pr.ID.String()).Execute()
			require.Error(t, err)

			require.IsType(t, new(kratos.GenericOpenAPIError), err, "%T", err)
			assert.Equal(t, int64(http.StatusGone), gjson.GetBytes(err.(*kratos.GenericOpenAPIError).Body(), "error.code").Int())
		})
	})

	t.Run("daemon=public", func(t *testing.T) {
		t.Run("description=fetching a non-existent flow should return a 403 error", func(t *testing.T) {
			_, _, err := testhelpers.NewSDKCustomClient(publicTS, otherUser).PublicApi.GetSelfServiceSettingsFlow(context.Background()).Id("i-do-not-exist").Execute()
			require.Error(t, err)

			require.IsType(t, new(kratos.GenericOpenAPIError), err, "%T", err)
			assert.Equal(t, int64(http.StatusForbidden), gjson.GetBytes(err.(*kratos.GenericOpenAPIError).Body(), "error.code").Int())
		})

		t.Run("description=fetching an expired flow returns 410", func(t *testing.T) {
			pr := newExpiredFlow()
			require.NoError(t, reg.SettingsFlowPersister().CreateSettingsFlow(context.Background(), pr))

			_, _, err := testhelpers.NewSDKCustomClient(publicTS, primaryUser).PublicApi.GetSelfServiceSettingsFlow(context.Background()).Id(pr.ID.String()).Execute()
			require.Error(t, err)

			require.IsType(t, new(kratos.GenericOpenAPIError), err, "%T", err)
			assert.Equal(t, int64(http.StatusGone), gjson.GetBytes(err.(*kratos.GenericOpenAPIError).Body(), "error.code").Int())
		})

		t.Run("description=should fail to fetch request if identity changed", func(t *testing.T) {
			t.Run("type=api", func(t *testing.T) {
				user1 := testhelpers.NewHTTPClientWithArbitrarySessionToken(t, reg)
				user2 := testhelpers.NewHTTPClientWithArbitrarySessionToken(t, reg)

				t.Logf("%+v", user1.Jar)
				res, err := user1.Get(publicTS.URL + settings.RouteInitAPIFlow)
				require.NoError(t, err)
				defer res.Body.Close()
				t.Logf("%+v", user1.Jar)
				assert.Len(t, res.Header.Get("Set-Cookie"), 0)
				body := ioutilx.MustReadAll(res.Body)
				id := gjson.GetBytes(body, "id")
				require.NotEmpty(t, id)

				res, err = user2.Get(publicTS.URL + settings.RouteGetFlow)
				require.NoError(t, err)
				defer res.Body.Close()
				require.EqualValues(t, res.StatusCode, http.StatusForbidden)
				body = ioutilx.MustReadAll(res.Body)
				assert.Contains(t, gjson.GetBytes(body, "error.reason").String(), "Access privileges are missing", "%s", body)
			})

			t.Run("type=browser", func(t *testing.T) {
				res, err := primaryUser.Get(publicTS.URL + settings.RouteInitBrowserFlow)
				require.NoError(t, err)

				rid := res.Request.URL.Query().Get("flow")
				require.NotEmpty(t, rid)

				testhelpers.NewSDKCustomClient(publicTS, primaryUser)

				_, _, err = testhelpers.NewSDKCustomClient(publicTS, otherUser).PublicApi.GetSelfServiceSettingsFlow(context.Background()).Id(rid).Execute()
				require.Error(t, err)
				require.IsType(t, new(kratos.GenericOpenAPIError), err, "%T", err)
				assert.Equal(t, int64(http.StatusForbidden), gjson.GetBytes(err.(*kratos.GenericOpenAPIError).Body(), "error.code").Int(), "should return a 403 error because the identities from the cookies do not match")
			})
		})
	})
}
