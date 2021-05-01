package daemon

import (
	"net/http"
	"sync"

	"kratos/selfservice/flow/recovery"

	"github.com/ory/x/reqlog"

	"kratos/cmd/courier"
	"kratos/driver/config"

	"github.com/rs/cors"

	"kratos/metrics/prometheus"

	"github.com/ory/analytics-go/v4"

	"github.com/ory/x/healthx"
	"github.com/ory/x/networkx"

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
	"kratos/session"
	"kratos/x"
)

type options struct {
	mwf []func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

func newOptions(opts []Option) *options {
	o := new(options)
	for _, f := range opts {
		f(o)
	}
	return o
}

type Option func(*options)

func WithRootMiddleware(m func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)) Option {
	return func(o *options) {
		o.mwf = append(o.mwf, m)
	}
}

func ServePublic(r driver.Registry, wg *sync.WaitGroup, cmd *cobra.Command, args []string, opts ...Option) {
	defer wg.Done()
	modifiers := newOptions(opts)

	ctx := cmd.Context()

	c := r.Config(cmd.Context())
	l := r.Logger()
	n := negroni.New()
	for _, mw := range modifiers.mwf {
		n.UseFunc(mw)
	}

	router := x.NewRouterPublic()
	csrf := x.NewCSRFHandler(router, r)

	n.UseFunc(x.CleanPath) // Prevent double slashes from breaking CSRF.
	r.WithCSRFHandler(csrf)
	n.UseHandler(r.CSRFHandler())

	r.RegisterPublicRoutes(ctx, router)
	n.Use(reqlog.NewMiddlewareFromLogger(l, "public#"+c.SelfPublicURL(nil).String()))
	n.Use(sqa(cmd, r))
	n.Use(r.PrometheusManager())

	if tracer := r.Tracer(ctx); tracer.IsLoaded() {
		n.Use(tracer)
	}

	var handler http.Handler = n
	options, enabled := r.Config(ctx).CORS("public")
	if enabled {
		handler = cors.New(options).Handler(handler)
	}

	server := graceful.WithDefaults(&http.Server{Handler: context.ClearHandler(handler)})
	addr := c.PublicListenOn()

	l.Printf("Starting the public httpd on: %s", addr)
	if err := graceful.Graceful(func() error {
		listener, err := networkx.MakeListener(addr, c.PublicSocketPermission())
		if err != nil {
			return err
		}

		return server.Serve(listener)
	}, server.Shutdown); err != nil {
		l.Fatalf("Failed to gracefully shutdown public httpd: %s", err)
	}
	l.Println("Public httpd was shutdown gracefully")
}

func ServeAdmin(r driver.Registry, wg *sync.WaitGroup, cmd *cobra.Command, args []string, opts ...Option) {
	defer wg.Done()
	modifiers := newOptions(opts)
	ctx := cmd.Context()

	c := r.Config(cmd.Context())
	l := r.Logger()
	n := negroni.New()
	for _, mw := range modifiers.mwf {
		n.UseFunc(mw)
	}

	router := x.NewRouterAdmin()
	r.RegisterAdminRoutes(ctx, router)
	n.Use(reqlog.NewMiddlewareFromLogger(l, "admin#"+c.SelfPublicURL(nil).String()))
	n.Use(sqa(cmd, r))
	n.Use(r.PrometheusManager())

	if tracer := r.Tracer(ctx); tracer.IsLoaded() {
		n.Use(tracer)
	}

	n.UseHandler(router)
	server := graceful.WithDefaults(&http.Server{Handler: context.ClearHandler(n)})
	addr := c.AdminListenOn()

	l.Printf("Starting the admin httpd on: %s", addr)
	if err := graceful.Graceful(func() error {
		listener, err := networkx.MakeListener(addr, c.AdminSocketPermission())
		if err != nil {
			return err
		}

		return server.Serve(listener)
	}, server.Shutdown); err != nil {
		l.Fatalf("Failed to gracefully shutdown admin httpd: %s", err)
	}
	l.Println("Admin httpd was shutdown gracefully")
}

func sqa(cmd *cobra.Command, d driver.Registry) *metricsx.Service {
	// Creates only ones
	// instance
	return metricsx.New(
		cmd,
		d.Logger(),
		d.Config(cmd.Context()).Source(),
		&metricsx.Options{
			Service:       "ory-kratos",
			ClusterID:     metricsx.Hash(d.Persister().NetworkID().String()),
			IsDevelopment: d.Config(cmd.Context()).IsInsecureDevMode(),
			WriteKey:      "qQlI6q8Q4WvkzTjKQSor4sHYOikHIvvi",
			WhitelistedPaths: []string{
				"/",
				healthx.AliveCheckPath,
				healthx.ReadyCheckPath,
				healthx.VersionPath,

				oidc.RouteBase,

				login.RouteInitBrowserFlow,
				login.RouteInitAPIFlow,
				login.RouteGetFlow,
				login.RouteSubmitFlow,

				logout.RouteBrowser,

				registration.RouteInitBrowserFlow,
				registration.RouteInitAPIFlow,
				registration.RouteGetFlow,
				registration.RouteSubmitFlow,

				session.RouteWhoami,
				identity.RouteBase,

				settings.RouteInitBrowserFlow,
				settings.RouteInitAPIFlow,
				settings.RouteGetFlow,
				settings.RouteSubmitFlow,

				verification.RouteInitAPIFlow,
				verification.RouteInitBrowserFlow,
				verification.RouteGetFlow,
				verification.RouteSubmitFlow,

				recovery.RouteInitAPIFlow,
				recovery.RouteInitBrowserFlow,
				recovery.RouteGetFlow,
				recovery.RouteSubmitFlow,

				link.RouteAdminCreateRecoveryLink,

				errorx.RouteGet,
				prometheus.MetricsPrometheusPath,
			},
			BuildVersion: config.Version,
			BuildHash:    config.Commit,
			BuildTime:    config.Date,
			Config: &analytics.Config{
				Endpoint: "https://sqa.ory.sh",
			},
		},
	)
}

func bgTasks(d driver.Registry, wg *sync.WaitGroup, cmd *cobra.Command, args []string) {
	defer wg.Done()

	if d.Config(cmd.Context()).IsBackgroundCourierEnabled() {
		go courier.Watch(cmd.Context(), d)
	}
}

func ServeAll(d driver.Registry, opts ...Option) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		wg.Add(3)
		go ServePublic(d, &wg, cmd, args, opts...)
		go ServeAdmin(d, &wg, cmd, args, opts...)
		go bgTasks(d, &wg, cmd, args)
		wg.Wait()
	}
}
