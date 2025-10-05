package transaction

import (
	"net/http"
	"subscriptionplus/server/middleware"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	subscriptionPayRouter := router.PathPrefix("/transactions").Subrouter()

	subscriptionPayRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	// access: все
	subscriptionPayRouter.Handle("/history", httpx.ErrorHandler(h.TransactionsHistory)).Methods(http.MethodGet)

	// access: все
	subscriptionPayRouter.Handle("/subscription/pay", httpx.ErrorHandler(h.SubscriptionPay)).Methods(http.MethodPost)

	// access: все
	subscriptionPayRouter.Handle("/subscription/id/{id:[0-9]+}", httpx.ErrorHandler(h.SubscriptionGetById)).Methods(http.MethodGet)

	// access: все
	subscriptionPayRouter.Handle("/subscription/token/{xtoken}", httpx.ErrorHandler(h.SubscriptionGetByXToken)).Methods(http.MethodGet)

	// access: все
	subscriptionPayRouter.Handle("/subscriptions/pending", httpx.ErrorHandler(h.SubscriptionGetPending)).Methods(http.MethodGet)
}
