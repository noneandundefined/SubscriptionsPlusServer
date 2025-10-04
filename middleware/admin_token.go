package middleware

import (
	"net/http"
	"os"
	"subscriptionplus/server/handler"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

func AdminTokenMiddleware(h *handler.BaseHandler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-API-Token")

			if token != os.Getenv("ADMIN_API_TOKEN") {
				httpx.HttpResponse(w, r, http.StatusForbidden, "invalid API token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
