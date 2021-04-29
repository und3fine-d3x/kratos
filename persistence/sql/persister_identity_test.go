package sql_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"kratos/internal/testhelpers"

	"github.com/stretchr/testify/require"

	"github.com/ory/x/sqlxx"
	"github.com/ory/x/urlx"
	"kratos/driver/config"
	"kratos/identity"
	"kratos/schema"
)

// note: it is important that this test runs on clean databases, as the racy behaviour only happens there
func TestPersister_CreateIdentityRacy(t *testing.T) {
	defaultSchema := schema.Schema{
		ID:     config.DefaultIdentityTraitsSchemaID,
		URL:    urlx.ParseOrPanic("file://./stub/identity.schema.json"),
		RawURL: "file://./stub/identity.schema.json",
	}

	ctx := context.Background()

	for name, p := range createCleanDatabases(t) {
		t.Run(fmt.Sprintf("db=%s", name), func(t *testing.T) {
			var wg sync.WaitGroup
			p.Config(context.Background()).MustSet(config.ViperKeyDefaultIdentitySchemaURL, defaultSchema.RawURL)
			_, ps := testhelpers.NewNetwork(t, ctx, p.Persister())

			for i := 0; i < 10; i++ {
				wg.Add(1)
				// capture i
				ii := i
				go func() {
					defer wg.Done()

					id := identity.NewIdentity("")
					id.SetCredentials(identity.CredentialsTypePassword, identity.Credentials{
						Type:        identity.CredentialsTypePassword,
						Identifiers: []string{fmt.Sprintf("racy identity %d", ii)},
						Config:      sqlxx.JSONRawMessage(`{"foo":"bar"}`),
					})
					id.Traits = identity.Traits("{}")

					require.NoError(t, ps.CreateIdentity(context.Background(), id))
				}()
			}

			wg.Wait()
		})

	}
}
