package subscription

import (
	"net/http"
	"subscriptionplus/server/middleware"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	subRouter := router.PathPrefix("/subs").Subrouter()

	subRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	// access: все
	subRouter.Handle("", httpx.ErrorHandler(h.GetSubscriptions)).Methods(http.MethodGet)

	// access: все
	subRouter.Handle("", httpx.ErrorHandler(h.AddSubscriptions)).Methods(http.MethodPost)

	// access: все
	subRouter.Handle("/{id:[0-9]+}", httpx.ErrorHandler(h.GetSubscriptionById)).Methods(http.MethodGet)

	// access: все
	subRouter.Handle("/{id:[0-9]+}", httpx.ErrorHandler(h.EditSubscriptions)).Methods(http.MethodPut)

	// access: все
	subRouter.Handle("/{id:[0-9]+}", httpx.ErrorHandler(h.DeleteSubscriptions)).Methods(http.MethodDelete)
}
