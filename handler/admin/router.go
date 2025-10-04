package admin

import (
	"net/http"
	"subscriptionplus/server/middleware"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/admin").Subrouter()

	adminRouter.Use(middleware.AdminTokenMiddleware(h.BaseHandler))

	// access: все
	adminRouter.Handle("/transactions/subscriptions/pending", httpx.ErrorHandler(h.SubscriptionsPending)).Methods(http.MethodGet)

	// access: все
	adminRouter.Handle("/transactions/subscriptions/{id:[0-9]+}", httpx.ErrorHandler(h.SubscriptionsPatchById)).Methods(http.MethodPatch)
}
