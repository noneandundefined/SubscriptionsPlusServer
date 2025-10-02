package main

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func (s *httpServer) cors(handler http.Handler) http.Handler {
	origins := handlers.AllowedOrigins([]string{"http://localhost:8001", "http://localhost:8002"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With", "Authorization", "X-Captcha-Token"})
	exposed := handlers.ExposedHeaders([]string{"X-Captcha-Required", "X-Captcha-Token"})

	allowCredentials := handlers.AllowCredentials()

	return handlers.CORS(origins, methods, headers, exposed, allowCredentials)(handler)
}
