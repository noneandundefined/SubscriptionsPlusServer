package plan

import (
	"net/http"
	"subscriptionplus/server/middleware"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	planRouter := router.PathPrefix("/plans").Subrouter()

	planRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	// access: все
	planRouter.Handle("", httpx.ErrorHandler(h.GetPlansHandler)).Methods(http.MethodGet)
}
