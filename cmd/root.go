package cmd

import (
	"errors"
	"fmt"
	"os"

	"kratos/driver/config"

	"kratos/cmd/courier"
	"kratos/cmd/hashers"

	"kratos/cmd/remote"

	"github.com/ory/x/cmdx"
	"kratos/cmd/identities"
	"kratos/cmd/jsonnet"
	"kratos/cmd/migrate"
	"kratos/cmd/serve"

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
			_, _ = fmt.Fprintln(RootCmd.ErrOrStderr(), err)
		}
		os.Exit(1)
	}
}

func init() {
	identities.RegisterCommandRecursive(RootCmd)
	identities.RegisterFlags()

	jsonnet.RegisterCommandRecursive(RootCmd)
	serve.RegisterCommandRecursive(RootCmd)
	migrate.RegisterCommandRecursive(RootCmd)
	remote.RegisterCommandRecursive(RootCmd)
	hashers.RegisterCommandRecursive(RootCmd)
	courier.RegisterCommandRecursive(RootCmd)

	RootCmd.AddCommand(cmdx.Version(&config.Version, &config.Commit, &config.Date))
}
