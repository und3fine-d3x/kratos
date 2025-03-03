package identities

import (
	"context"
	"strings"
	"testing"

	"github.com/ory/x/cmdx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kratos/identity"
)

func TestListCmd(t *testing.T) {
	reg := setup(t, ListCmd)
	require.NoError(t, ListCmd.Flags().Set(cmdx.FlagQuiet, "true"))

	var deleteIdentities = func(t *testing.T, is []*identity.Identity) {
		for _, i := range is {
			require.NoError(t, reg.Persister().DeleteIdentity(context.Background(), i.ID))
		}
	}

	t.Run("case=lists all identities with default pagination", func(t *testing.T) {
		is, ids := makeIdentities(t, reg, 5)
		defer deleteIdentities(t, is)

		stdOut := execNoErr(t, ListCmd)

		for _, i := range ids {
			assert.Contains(t, stdOut, i)
		}
	})

	t.Run("case=lists all identities with pagination", func(t *testing.T) {
		is, ids := makeIdentities(t, reg, 6)
		defer deleteIdentities(t, is)

		stdoutP1 := execNoErr(t, ListCmd, "1", "3")
		stdoutP2 := execNoErr(t, ListCmd, "2", "3")

		for _, id := range ids {
			// exactly one of page 1 and 2 should contain the id
			assert.True(t, strings.Contains(stdoutP1, id) != strings.Contains(stdoutP2, id), "%s \n %s", stdoutP1, stdoutP2)
		}
	})
}
