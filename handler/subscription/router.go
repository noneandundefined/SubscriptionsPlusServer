package subscription

import (
	"net/http"
	"subscriptionplus/server/middleware"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	subRouter := router.PathPrefix("/subs").Subrouter()
	subFreeRouter := router.PathPrefix("/subs").Subrouter()

	subRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	// access: все
	subRouter.Handle("", httpx.ErrorHandler(h.GetSubscriptionsHandler)).Methods(http.MethodGet)

	// access: все
	subRouter.Handle("", httpx.ErrorHandler(h.AddSubscriptionsHandler)).Methods(http.MethodPost)

	// access: все
	subRouter.Handle("/{id:[0-9]+}", httpx.ErrorHandler(h.GetSubscriptionByIdHandler)).Methods(http.MethodGet)

	// access: все
	subRouter.Handle("/{id:[0-9]+}", httpx.ErrorHandler(h.EditSubscriptionsHandler)).Methods(http.MethodPut)

	// access: все
	subRouter.Handle("/{id:[0-9]+}", httpx.ErrorHandler(h.DeleteSubscriptionsHandler)).Methods(http.MethodDelete)

	// access: все
	subFreeRouter.Handle("/images/w350/{logo}", httpx.ErrorHandler(h.GetSubscriptionLogoHandler)).Methods(http.MethodGet)
}
