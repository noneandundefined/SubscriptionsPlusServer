package auth

import (
	"net/http"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/auth").Subrouter()

	// access: все
	authRouter.Handle("/create", httpx.ErrorHandler(h.AuthCreateUserHandler)).Methods(http.MethodPost)
}
