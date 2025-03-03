package login_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"kratos/ui/node"

	"github.com/ory/kratos-client-go"

	"github.com/gobuffalo/httptest"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/x/assertx"
	"github.com/ory/x/urlx"

	"github.com/ory/herodot"

	"kratos/internal"
	"kratos/internal/testhelpers"
	"kratos/schema"
	"kratos/selfservice/flow"
	"kratos/selfservice/flow/login"
	"kratos/text"
	"kratos/x"
)

func TestHandleError(t *testing.T) {
	conf, reg := internal.NewFastRegistryWithMocks(t)
	public, admin := testhelpers.NewKratosServer(t, reg)

	router := httprouter.New()
	ts := httptest.NewServer(router)
	t.Cleanup(ts.Close)

	testhelpers.NewLoginUIFlowEchoServer(t, reg)
	testhelpers.NewErrorTestServer(t, reg)

	h := reg.LoginFlowErrorHandler()
	sdk := testhelpers.NewSDKClient(admin)

	var loginFlow *login.Flow
	var flowError error
	var ct node.Group
	router.GET("/error", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.WriteFlowError(w, r, loginFlow, ct, flowError)
	})

	reset := func() {
		loginFlow = nil
		flowError = nil
		ct = ""
	}

	newFlow := func(t *testing.T, ttl time.Duration, ft flow.Type) *login.Flow {
		req := &http.Request{URL: urlx.ParseOrPanic("/")}
		f := login.NewFlow(conf, ttl, "csrf_token", req, ft)
		for _, s := range reg.LoginStrategies(context.Background()) {
			require.NoError(t, s.PopulateLoginMethod(req, f))
		}

		require.NoError(t, reg.LoginFlowPersister().CreateLoginFlow(context.Background(), f))
		return f
	}

	expectErrorUI := func(t *testing.T) ([]map[string]interface{}, *http.Response) {
		res, err := ts.Client().Get(ts.URL + "/error")
		require.NoError(t, err)
		defer res.Body.Close()
		require.Contains(t, res.Request.URL.String(), conf.SelfServiceFlowErrorURL().String()+"?error=")

		sse, _, err := sdk.PublicApi.GetSelfServiceError(context.Background()).Error_(res.Request.URL.Query().Get("error")).Execute()
		require.NoError(t, err)

		return sse.Errors, nil
	}

	anHourAgo := time.Now().Add(-time.Hour)

	t.Run("case=error with nil flow defaults to error ui redirect", func(t *testing.T) {
		t.Cleanup(reset)

		flowError = herodot.ErrInternalServerError.WithReason("system error")
		ct = node.PasswordGroup

		sse, _ := expectErrorUI(t)
		assertx.EqualAsJSON(t, []interface{}{flowError}, sse)
	})

	t.Run("case=error with nil flow detects application/json", func(t *testing.T) {
		t.Cleanup(reset)

		flowError = herodot.ErrInternalServerError.WithReason("system error")
		ct = node.PasswordGroup

		res, err := ts.Client().Do(testhelpers.NewHTTPGetJSONRequest(t, ts.URL+"/error"))
		require.NoError(t, err)
		defer res.Body.Close()
		assert.Contains(t, res.Header.Get("Content-Type"), "application/json")
		assert.NotContains(t, res.Request.URL.String(), conf.SelfServiceFlowErrorURL().String()+"?error=")

		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Contains(t, string(body), "system error")
	})

	t.Run("flow=api", func(t *testing.T) {
		t.Run("case=expired error", func(t *testing.T) {
			t.Cleanup(reset)

			loginFlow = newFlow(t, time.Minute, flow.TypeAPI)
			flowError = login.NewFlowExpiredError(anHourAgo)
			ct = node.PasswordGroup

			res, err := ts.Client().Do(testhelpers.NewHTTPGetJSONRequest(t, ts.URL+"/error"))
			require.NoError(t, err)
			defer res.Body.Close()
			require.Contains(t, res.Request.URL.String(), public.URL+login.RouteGetFlow)
			require.Equal(t, http.StatusOK, res.StatusCode)

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, int(text.ErrorValidationLoginFlowExpired), int(gjson.GetBytes(body, "ui.messages.0.id").Int()))
			assert.NotEqual(t, loginFlow.ID.String(), gjson.GetBytes(body, "id").String())
		})

		t.Run("case=validation error", func(t *testing.T) {
			t.Cleanup(reset)

			loginFlow = newFlow(t, time.Minute, flow.TypeAPI)
			flowError = schema.NewInvalidCredentialsError()
			ct = node.PasswordGroup

			res, err := ts.Client().Do(testhelpers.NewHTTPGetJSONRequest(t, ts.URL+"/error"))
			require.NoError(t, err)
			defer res.Body.Close()
			require.Equal(t, http.StatusBadRequest, res.StatusCode)

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, int(text.ErrorValidationInvalidCredentials), int(gjson.GetBytes(body, "ui.messages.0.id").Int()), "%s", body)
			assert.Equal(t, loginFlow.ID.String(), gjson.GetBytes(body, "id").String())
		})

		t.Run("case=generic error", func(t *testing.T) {
			t.Cleanup(reset)

			loginFlow = newFlow(t, time.Minute, flow.TypeAPI)
			flowError = herodot.ErrInternalServerError.WithReason("system error")
			ct = node.PasswordGroup

			res, err := ts.Client().Do(testhelpers.NewHTTPGetJSONRequest(t, ts.URL+"/error"))
			require.NoError(t, err)
			defer res.Body.Close()
			require.Equal(t, http.StatusInternalServerError, res.StatusCode)

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			assert.JSONEq(t, x.MustEncodeJSON(t, flowError), gjson.GetBytes(body, "error").Raw)
		})
	})

	t.Run("flow=browser", func(t *testing.T) {
		expectLoginUI := func(t *testing.T) (*kratos.LoginFlow, *http.Response) {
			res, err := ts.Client().Get(ts.URL + "/error")
			require.NoError(t, err)
			defer res.Body.Close()
			assert.Contains(t, res.Request.URL.String(), conf.SelfServiceFlowLoginUI().String()+"?flow=")

			lf, _, err := sdk.PublicApi.GetSelfServiceLoginFlow(context.Background()).Id(res.Request.URL.Query().Get("flow")).Execute()
			require.NoError(t, err)
			return lf, res
		}

		t.Run("case=expired error", func(t *testing.T) {
			t.Cleanup(reset)

			loginFlow = &login.Flow{Type: flow.TypeBrowser}
			flowError = login.NewFlowExpiredError(anHourAgo)
			ct = node.PasswordGroup

			lf, _ := expectLoginUI(t)
			require.Len(t, lf.Ui.Messages, 1)
			assert.Equal(t, int(text.ErrorValidationLoginFlowExpired), int(lf.Ui.Messages[0].Id))
		})

		t.Run("case=validation error", func(t *testing.T) {
			t.Cleanup(reset)

			loginFlow = newFlow(t, time.Minute, flow.TypeBrowser)
			flowError = schema.NewInvalidCredentialsError()
			ct = node.PasswordGroup

			lf, _ := expectLoginUI(t)
			require.NotEmpty(t, lf.Ui.Nodes, x.MustEncodeJSON(t, lf))
			require.Len(t, lf.Ui.Messages, 1, x.MustEncodeJSON(t, lf))
			assert.Equal(t, int(text.ErrorValidationInvalidCredentials), int(lf.Ui.Messages[0].Id), x.MustEncodeJSON(t, lf))
		})

		t.Run("case=generic error", func(t *testing.T) {
			t.Cleanup(reset)

			loginFlow = newFlow(t, time.Minute, flow.TypeBrowser)
			flowError = herodot.ErrInternalServerError.WithReason("system error")
			ct = node.PasswordGroup

			sse, _ := expectErrorUI(t)
			assertx.EqualAsJSON(t, []interface{}{flowError}, sse)
		})
	})
}
