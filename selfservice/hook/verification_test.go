package hook_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/sqlxx"
	"github.com/ory/x/urlx"
	"kratos/driver/config"
	"kratos/identity"
	"kratos/internal"
	"kratos/selfservice/flow"
	"kratos/selfservice/flow/registration"
	"kratos/selfservice/flow/settings"
	"kratos/selfservice/flow/verification"
	"kratos/selfservice/hook"
	"kratos/session"
	"kratos/x"
)

func TestVerifier(t *testing.T) {
	u := &http.Request{URL: urlx.ParseOrPanic("https://www.ory.sh/")}
	for k, hf := range map[string]func(*hook.Verifier, *identity.Identity, flow.Flow) error{
		"settings": func(h *hook.Verifier, i *identity.Identity, f flow.Flow) error {
			return h.ExecuteSettingsPostPersistHook(
				httptest.NewRecorder(), u, f.(*settings.Flow), i)
		},
		"register": func(h *hook.Verifier, i *identity.Identity, f flow.Flow) error {
			return h.ExecutePostRegistrationPostPersistHook(
				httptest.NewRecorder(), u, f.(*registration.Flow), &session.Session{ID: x.NewUUID(), Identity: i})
		},
	} {
		t.Run("name="+k, func(t *testing.T) {
			conf, reg := internal.NewFastRegistryWithMocks(t)
			conf.MustSet(config.ViperKeyDefaultIdentitySchemaURL, "file://./stub/verify.schema.json")
			conf.MustSet(config.ViperKeyPublicBaseURL, "https://www.ory.sh/")
			conf.MustSet(config.ViperKeyCourierSMTPURL, "smtp://foo@bar@dev.null/")

			i := identity.NewIdentity(config.DefaultIdentityTraitsSchemaID)
			i.Traits = identity.Traits(`{"emails":["foo@ory.sh","bar@ory.sh","baz@ory.sh"]}`)
			require.NoError(t, reg.IdentityManager().Create(context.Background(), i))

			actual, err := reg.IdentityPool().FindVerifiableAddressByValue(context.Background(), identity.VerifiableAddressTypeEmail, "foo@ory.sh")
			require.NoError(t, err)
			assert.EqualValues(t, "foo@ory.sh", actual.Value)

			actual, err = reg.IdentityPool().FindVerifiableAddressByValue(context.Background(), identity.VerifiableAddressTypeEmail, "bar@ory.sh")
			require.NoError(t, err)
			assert.EqualValues(t, "bar@ory.sh", actual.Value)

			actual, err = reg.IdentityPool().FindVerifiableAddressByValue(context.Background(), identity.VerifiableAddressTypeEmail, "baz@ory.sh")
			require.NoError(t, err)
			assert.EqualValues(t, "baz@ory.sh", actual.Value)
			actual.Status = identity.VerifiableAddressStatusCompleted
			actual.Verified = true
			actual.VerifiedAt = sqlxx.NullTime(time.Now())
			require.NoError(t, reg.PrivilegedIdentityPool().UpdateVerifiableAddress(context.Background(), actual))

			i, err = reg.IdentityPool().GetIdentity(context.Background(), i.ID)
			require.NoError(t, err)

			var originalFlow flow.Flow
			switch k {
			case "settings":
				originalFlow = &settings.Flow{RequestURL: "http://foo.com/settings?after_verification_return_to=verification_callback"}
			case "register":
				originalFlow = &registration.Flow{RequestURL: "http://foo.com/registration?after_verification_return_to=verification_callback"}
			default:
				t.FailNow()
			}

			h := hook.NewVerifier(reg)
			require.NoError(t, hf(h, i, originalFlow))
			expectedVerificationFlow, err := verification.NewPostHookFlow(conf, conf.SelfServiceFlowVerificationRequestLifespan(), "", u, reg.VerificationStrategies(context.Background()), originalFlow)
			require.NoError(t, err)

			var verificationFlow verification.Flow
			require.NoError(t, reg.Persister().GetConnection(context.Background()).First(&verificationFlow))

			assert.Equal(t, expectedVerificationFlow.RequestURL, verificationFlow.RequestURL)

			messages, err := reg.CourierPersister().NextMessages(context.Background(), 12)
			require.NoError(t, err)
			require.Len(t, messages, 2)

			recipients := make([]string, len(messages))
			for k, m := range messages {
				recipients[k] = m.Recipient
			}

			assert.Contains(t, recipients, "foo@ory.sh")
			assert.Contains(t, recipients, "bar@ory.sh")
			assert.NotContains(t, recipients, "baz@ory.sh")
			// Email to baz@ory.sh is skipped because it is verified already.
		})
	}
}
