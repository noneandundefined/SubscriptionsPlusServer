package plan

import (
	"net/http"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	planRouter := router.PathPrefix("/plans").Subrouter()

	// access: все
	planRouter.Handle("", httpx.ErrorHandler(h.GetPlansHandler)).Methods(http.MethodGet)
}
