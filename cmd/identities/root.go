package identities

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"

	"kratos/cmd/cliclient"
)

// identitiesCmd represents the identity command
var identitiesCmd = &cobra.Command{
	Use:   "identities",
	Short: "Tools to interact with remote identities",
}

func RegisterCommandRecursive(parent *cobra.Command) {
	parent.AddCommand(identitiesCmd)

	identitiesCmd.AddCommand(ImportCmd)
	identitiesCmd.AddCommand(ValidateCmd)
	identitiesCmd.AddCommand(ListCmd)
	identitiesCmd.AddCommand(GetCmd)
	identitiesCmd.AddCommand(DeleteCmd)
	identitiesCmd.AddCommand(PatchCmd)
}

func RegisterFlags() {
	cliclient.RegisterClientFlags(identitiesCmd.PersistentFlags())
	cmdx.RegisterFormatFlags(identitiesCmd.PersistentFlags())
}
