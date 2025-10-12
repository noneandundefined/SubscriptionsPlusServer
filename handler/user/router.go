package user

import (
	"net/http"
	"subscriptionplus/server/middleware"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/users").Subrouter()

	userRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	// access: все
	userRouter.Handle("/me", httpx.ErrorHandler(h.GetMeHandler)).Methods(http.MethodGet)
}
