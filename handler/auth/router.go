package auth

import (
	"net/http"
	"subscriptionplus/server/middleware"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/auth").Subrouter()

	authProtectedRouter := router.PathPrefix("/auth").Subrouter()
	authProtectedRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	// access: все
	authRouter.Handle("/create", httpx.ErrorHandler(h.AuthCreateUserHandler)).Methods(http.MethodPost)

	// access: все
	authProtectedRouter.Handle("/delete", httpx.ErrorHandler(h.AuthDelete)).Methods(http.MethodDelete)
}
