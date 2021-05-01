package cliclient

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/spf13/cobra"

	"github.com/spf13/pflag"

	"github.com/ory/kratos-client-go"
	"github.com/ory/x/cmdx"
)

const (
	envKeyEndpoint = "KRATOS_ADMIN_URL"
	FlagEndpoint   = "endpoint"
)

type ContextKey int

const (
	ClientContextKey ContextKey = iota + 1
)

func NewClient(cmd *cobra.Command) *kratos.APIClient {
	if f, ok := cmd.Context().Value(ClientContextKey).(func(cmd *cobra.Command) *kratos.APIClient); ok {
		return f(cmd)
	} else if f != nil {
		panic(fmt.Sprintf("ClientContextKey was expected to be *client.OryKratos but it contained an invalid type %T ", f))
	}

	endpoint, err := cmd.Flags().GetString(FlagEndpoint)
	cmdx.Must(err, "flag access error: %s", err)

	if endpoint == "" {
		endpoint = os.Getenv(envKeyEndpoint)
	}

	if endpoint == "" {
		// no endpoint is set
		_, _ = fmt.Fprintln(os.Stderr, "You have to set the remote endpoint, try --help for details.")
		os.Exit(1)
	}

	u, err := url.Parse(endpoint)
	cmdx.Must(err, `Could not parse the endpoint URL "%s".`, endpoint)

	conf := kratos.NewConfiguration()
	conf.HTTPClient = retryablehttp.NewClient().StandardClient()
	conf.HTTPClient.Timeout = time.Second * 10
	conf.Servers = kratos.ServerConfigurations{{URL: u.String()}}
	return kratos.NewAPIClient(conf)
}

func RegisterClientFlags(flags *pflag.FlagSet) {
	flags.StringP(FlagEndpoint, FlagEndpoint[:1], "", fmt.Sprintf("The URL of Ory Kratos' Admin API. Alternatively set using the %s environmental variable.", envKeyEndpoint))
}
