package daemon

import (
	"net/http"
	"strings"
	"sync"

	"github.com/ory/x/stringsx"

	"github.com/rs/cors"

	"github.com/ory/x/corsx"

	"kratos/metrics/prometheus"

	"github.com/ory/analytics-go/v4"

	"github.com/ory/x/flagx"
	"github.com/ory/x/healthx"

	"github.com/gorilla/context"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"

	"github.com/ory/graceful"
	"github.com/ory/x/metricsx"

	"kratos/driver"
	"kratos/identity"
	"kratos/selfservice/errorx"
	"kratos/selfservice/flow/login"
	"kratos/selfservice/flow/logout"
	"kratos/selfservice/flow/registration"
	"kratos/selfservice/flow/settings"
	"kratos/selfservice/flow/verification"
	"kratos/selfservice/strategy/link"
	"kratos/selfservice/strategy/oidc"
	"kratos/selfservice/strategy/password"
	"kratos/selfservice/strategy/profile"
	"kratos/session"
	"kratos/x"
)

func servePublic(d driver.Driver, wg *sync.WaitGroup, cmd *cobra.Command, args []string) {
	defer wg.Done()

	c := d.Configuration()
	l := d.Logger()
	n := negroni.New()
	r := d.Registry()

	router := x.NewRouterPublic()
	csrf := x.NewCSRFHandler(
		router,
		r.Writer(),
		l,
		stringsx.Coalesce(c.SelfPublicURL().Path, "/"),
		c.SelfPublicURL().Hostname(),
		!flagx.MustGetBool(cmd, "dev"),
	)

	n.UseFunc(x.CleanPath) // Prevent double slashes from breaking CSRF.
	r.WithCSRFHandler(csrf)
	n.UseHandler(r.CSRFHandler())

	r.RegisterPublicRoutes(router)
	n.Use(NewNegroniLoggerMiddleware(l, "public#"+c.SelfPublicURL().String()))
	n.Use(sqa(cmd, d))

	if tracer := d.Registry().Tracer(); tracer.IsLoaded() {
		n.Use(tracer)
	}

	var handler http.Handler = n
	if corsx.IsEnabled(l, "serve.public") {
		handler = cors.New(
			corsx.ParseOptions(l, "serve.public")).Handler(handler)
	}

	server := graceful.WithDefaults(&http.Server{
		Addr:    c.PublicListenOn(),
		Handler: context.ClearHandler(handler),
	})

	l.Printf("Starting the public httpd on: %s", server.Addr)
	if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
		l.Fatalln("Failed to gracefully shutdown public httpd")
	}
	l.Println("Public httpd was shutdown gracefully")
}

func serveAdmin(d driver.Driver, wg *sync.WaitGroup, cmd *cobra.Command, args []string) {
	defer wg.Done()

	c := d.Configuration()
	l := d.Logger()
	n := negroni.New()
	r := d.Registry()

	router := x.NewRouterAdmin()
	r.RegisterAdminRoutes(router)
	n.Use(NewNegroniLoggerMiddleware(l, "admin#"+c.SelfAdminURL().String()))
	n.Use(sqa(cmd, d))
	n.Use(d.Registry().PrometheusManager())

	if tracer := d.Registry().Tracer(); tracer.IsLoaded() {
		n.Use(tracer)
	}

	n.UseHandler(router)
	server := graceful.WithDefaults(&http.Server{
		Addr:    c.AdminListenOn(),
		Handler: context.ClearHandler(n),
	})

	l.Printf("Starting the admin httpd on: %s", server.Addr)
	if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
		l.Fatalln("Failed to gracefully shutdown admin httpd")
	}
	l.Println("Admin httpd was shutdown gracefully")
}

func sqa(cmd *cobra.Command, d driver.Driver) *metricsx.Service {
	// Creates only ones
	// instance
	return metricsx.New(
		cmd,
		d.Logger(),
		&metricsx.Options{
			Service: "ory-kratos",
			ClusterID: metricsx.Hash(
				strings.Join([]string{
					d.Configuration().DSN(),
					d.Configuration().SelfPublicURL().String(),
					d.Configuration().SelfAdminURL().String(),
				}, "|"),
			),
			IsDevelopment: flagx.MustGetBool(cmd, "dev"),
			WriteKey:      "qQlI6q8Q4WvkzTjKQSor4sHYOikHIvvi",
			WhitelistedPaths: []string{
				"/",
				healthx.AliveCheckPath,
				healthx.ReadyCheckPath,
				healthx.VersionPath,

				password.RouteRegistration,
				password.RouteLogin,
				password.RouteSettings,

				oidc.RouteBase,

				login.RouteInitBrowserFlow,
				login.RouteInitAPIFlow,
				login.RouteGetFlow,

				logout.RouteBrowser,

				registration.RouteInitBrowserFlow,
				registration.RouteInitAPIFlow,
				registration.RouteGetFlow,

				session.RouteWhoami,
				identity.RouteBase,

				settings.RouteInitBrowserFlow,
				settings.RouteInitAPIFlow,
				settings.RouteGetFlow,

				verification.RouteInitAPIFlow,
				verification.RouteInitBrowserFlow,
				verification.RouteGetFlow,

				profile.RouteSettings,

				link.RouteAdminCreateRecoveryLink,
				link.RouteRecovery,
				link.RouteVerification,

				errorx.RouteGet,
				prometheus.MetricsPrometheusPath,
			},
			BuildVersion: d.Registry().BuildVersion(),
			BuildHash:    d.Registry().BuildHash(),
			BuildTime:    d.Registry().BuildDate(),
			Config: &analytics.Config{
				Endpoint: "https://sqa.ory.sh",
			},
		},
	)
}

func bgTasks(d driver.Driver, wg *sync.WaitGroup, cmd *cobra.Command, args []string) {
	defer wg.Done()

	d.Logger().Println("Courier worker started.")
	if err := graceful.Graceful(d.Registry().Courier().Work, d.Registry().Courier().Shutdown); err != nil {
		d.Logger().WithError(err).Fatalf("Failed to run courier worker.")
	}
	d.Logger().Println("Courier worker was shutdown gracefully.")
}

func ServeAll(d driver.Driver) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		wg.Add(3)
		go servePublic(d, &wg, cmd, args)
		go serveAdmin(d, &wg, cmd, args)
		go bgTasks(d, &wg, cmd, args)
		wg.Wait()
	}
}
