package logout

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"kratos/driver/config"
	"kratos/selfservice/errorx"
	"kratos/session"
	"kratos/x"
)

const (
	RouteBrowser = "/self-service/browser/flows/logout"
)

type (
	handlerDependencies interface {
		x.CSRFProvider
		session.ManagementProvider
		errorx.ManagementProvider
		config.Provider
	}
	HandlerProvider interface {
		LogoutHandler() *Handler
	}
	Handler struct {
		d handlerDependencies
	}
)

func NewHandler(d handlerDependencies) *Handler {
	return &Handler{d: d}
}

func (h *Handler) RegisterPublicRoutes(router *x.RouterPublic) {
	router.GET(RouteBrowser, h.logout)
}

// swagger:route GET /self-service/browser/flows/logout public initializeSelfServiceBrowserLogoutFlow
//
// Initialize Browser-Based Logout User Flow
//
// This endpoint initializes a logout flow.
//
// > This endpoint is NOT INTENDED for API clients and only works
// with browsers (Chrome, Firefox, ...).
//
// On successful logout, the browser will be redirected (HTTP 302 Found) to the `return_to` parameter of the initial request
// or fall back to `urls.default_return_to`.
//
// More information can be found at [Ory Kratos User Logout Documentation](https://www.ory.sh/docs/next/kratos/self-service/flows/user-logout).
//
//     Schemes: http, https
//
//     Responses:
//       302: emptyResponse
//       500: genericError
func (h *Handler) logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_ = h.d.CSRFHandler().RegenerateToken(w, r)

	if err := h.d.SessionManager().PurgeFromRequest(r.Context(), w, r); err != nil {
		h.d.SelfServiceErrorManager().Forward(r.Context(), w, r, err)
		return
	}

	ret, err := x.SecureRedirectTo(r, h.d.Config(r.Context()).SelfServiceFlowLogoutRedirectURL(),
		x.SecureRedirectUseSourceURL(r.RequestURI),
		x.SecureRedirectAllowURLs(h.d.Config(r.Context()).SelfServiceBrowserWhitelistedReturnToDomains()),
		x.SecureRedirectAllowSelfServiceURLs(h.d.Config(r.Context()).SelfPublicURL(r)),
	)
	if err != nil {
		fmt.Printf("\n%s\n\n", err.Error())
		h.d.SelfServiceErrorManager().Forward(r.Context(), w, r, err)
		return
	}

	http.Redirect(w, r, ret.String(), http.StatusFound)
}
