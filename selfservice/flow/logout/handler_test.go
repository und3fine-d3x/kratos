package logout_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/gobuffalo/httptest"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/nosurf"
	"kratos/driver/config"
	"kratos/identity"
	"kratos/internal"
	"kratos/internal/testhelpers"
	"kratos/selfservice/flow/logout"
	"kratos/session"
	"kratos/x"
)

func TestLogoutHandler(t *testing.T) {
	conf, reg := internal.NewFastRegistryWithMocks(t)
	handler := reg.LogoutHandler()

	conf.MustSet(config.ViperKeyDefaultIdentitySchemaURL, "file://./stub/registration.schema.json")
	conf.MustSet(config.ViperKeyPublicBaseURL, "http://example.com")

	router := x.NewRouterPublic()
	handler.RegisterPublicRoutes(router)
	reg.WithCSRFHandler(x.NewCSRFHandler(router, reg))
	ts := httptest.NewServer(reg.CSRFHandler())
	defer ts.Close()

	var sess session.Session
	sess.ID = x.NewUUID()
	sess.Identity = new(identity.Identity)
	require.NoError(t, reg.PrivilegedIdentityPool().CreateIdentity(context.Background(), sess.Identity))
	require.NoError(t, reg.SessionPersister().CreateSession(context.Background(), &sess))

	router.GET("/set", testhelpers.MockSetSession(t, reg, conf))

	router.GET("/csrf", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		_, _ = w.Write([]byte(nosurf.Token(r)))
	})

	redirTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer redirTS.Close()

	conf.MustSet(config.ViperKeySelfServiceLogoutBrowserDefaultReturnTo, redirTS.URL)
	conf.MustSet(config.ViperKeyPublicBaseURL, ts.URL)

	client := testhelpers.NewClientWithCookies(t)

	t.Run("case=set initial session", func(t *testing.T) {
		testhelpers.MockHydrateCookieClient(t, client, ts.URL+"/set")
	})

	var token string
	t.Run("case=get csrf token", func(t *testing.T) {
		res, err := ts.Client().Get(ts.URL + "/csrf")
		require.NoError(t, err)
		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		token = string(body)
		require.NotEmpty(t, token)
	})

	t.Run("case=log out", func(t *testing.T) {
		res, err := client.Get(ts.URL + logout.RouteBrowser)
		require.NoError(t, err)

		var found bool
		for _, c := range res.Cookies() {
			if c.Name == config.DefaultSessionCookieName {
				found = true
			}
		}
		require.False(t, found)
		assert.Equal(t, redirTS.URL, res.Request.URL.String())
	})

	t.Run("case=csrf token should be reset", func(t *testing.T) {
		res, err := ts.Client().Get(ts.URL + "/csrf")
		require.NoError(t, err)
		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.NotEmpty(t, body)
		assert.NotEqual(t, token, string(body))
	})

	t.Run("case=respects return_to URI parameter", func(t *testing.T) {
		returnToURL := ts.URL + "/after-logout"
		conf.MustSet(config.ViperKeyURLsWhitelistedReturnToDomains, []string{returnToURL})

		query := url.Values{
			"return_to": {returnToURL},
		}

		res, err := client.Get(ts.URL + logout.RouteBrowser + "?" + query.Encode())
		require.NoError(t, err)
		assert.Equal(t, returnToURL, res.Request.URL.String())
	})
}
