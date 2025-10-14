package notification

import (
	"net/http"
	"subscriptionplus/server/middleware"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	notifyRouter := router.PathPrefix("/notify").Subrouter()

	notifyRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	// access: все
	notifyRouter.Handle("/token", httpx.ErrorHandler(h.NotificationTokenHandler)).Methods(http.MethodPost)
}
