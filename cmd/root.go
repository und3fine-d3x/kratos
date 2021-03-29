package cmd

import (
	"errors"
	"fmt"
	"os"

	"kratos/cmd/hashers"

	"kratos/cmd/remote"

	"kratos/cmd/identities"
	"kratos/cmd/jsonnet"
	"kratos/cmd/migrate"
	"kratos/cmd/serve"
	"kratos/internal/clihelpers"

	"github.com/ory/x/cmdx"

	"github.com/ory/x/viperx"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "kratos",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if !errors.Is(err, cmdx.ErrNoPrintButFail) {
			fmt.Fprintln(RootCmd.ErrOrStderr(), err)
		}
		os.Exit(1)
	}
}

func init() {
	viperx.RegisterConfigFlag(RootCmd, "kratos")

	identities.RegisterCommandRecursive(RootCmd)
	jsonnet.RegisterCommandRecursive(RootCmd)
	serve.RegisterCommandRecursive(RootCmd)
	migrate.RegisterCommandRecursive(RootCmd)
	remote.RegisterCommandRecursive(RootCmd)
	hashers.RegisterCommandRecursive(RootCmd)

	RootCmd.AddCommand(cmdx.Version(&clihelpers.BuildVersion, &clihelpers.BuildGitHash, &clihelpers.BuildTime))
}
