package middleware

import (
	"net/http"
	"subscriptionplus/server/pkg/httpx"
	"sync"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getLimiter(ip string, rps float64, burst int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if l, exists := limiters[ip]; exists {
		return l
	}

	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	limiters[ip] = limiter
	return limiter
}

// RateLimiterMiddleware ограничивает кол-во запросов от клиентов
func RateLimiterMiddleware(rps float64, burst int) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			limiter := getLimiter(ip, rps, burst)

			if !limiter.Allow() {
				httpx.HttpResponse(w, r, http.StatusTooManyRequests, "rate limiter exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
