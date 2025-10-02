package main

import (
	"net/http"
	"subscriptionplus/server/handler"
	"subscriptionplus/server/handler/auth"
	"subscriptionplus/server/handler/subscription"
	"subscriptionplus/server/handler/user"
	"subscriptionplus/server/middleware"

	"github.com/gorilla/mux"
)

func (s *httpServer) routes() http.Handler {
	router := mux.NewRouter()

	// middleware for logging API request
	router.Use(middleware.NewLogger(s.logger).LoggerMiddleware)
	// middleware for get exception errors
	router.Use(middleware.RecoveryMiddleware())
	// middleware for security API
	router.Use(middleware.SecurityMiddleware())
	// middleware rate limiter
	router.Use(middleware.RateLimiterMiddleware(6, 10))

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	baseHandler := &handler.BaseHandler{
		Db:     s.db,
		Logger: s.logger,
		Store:  s.store,
	}

	// routes path
	// authenticate
	auth.NewHandler(baseHandler).RegisterRoutes(subrouter)
	// subscription
	subscription.NewHandler(baseHandler).RegisterRoutes(subrouter)
	// user
	user.NewHandler(baseHandler).RegisterRoutes(subrouter)

	// docs
	s.docs(subrouter)

	return s.cors(router)
}
