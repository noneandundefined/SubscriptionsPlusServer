package middleware

import (
	"net/http"
	"os"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

func ApiTokenMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-App-Key")

			if token != os.Getenv("APPKEY_TOKEN") {
				httpx.HttpResponse(w, r, http.StatusForbidden, "unauthorized client")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
