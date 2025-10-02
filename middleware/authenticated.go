// nolint
package middleware

import (
	"context"
	"net/http"
	"subscriptionplus/server/handler"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/httpx"

	"github.com/gorilla/mux"
)

// IsAuthenticatedMiddleware проверка на аутентифицированного пользователя
func IsAuthenticatedMiddleware(h *handler.BaseHandler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth-token")
			if err != nil || cookie == nil {
				httpx.HttpResponse(w, r, http.StatusUnauthorized, "please connect to account")
				return
			}

			// db: get user core
			user, err := h.Store.Users.Get_UserCoreByToken(r.Context(), cookie.Value)
			if err != nil {
				httpx.HttpResponse(w, r, http.StatusUnauthorized, "please connect to account")
				return
			}

			if user == nil {
				httpx.HttpResponse(w, r, http.StatusUnauthorized, "please connect to account")
				return
			}

			authTokenModel := &types.AuthToken{
				User: *user,
			}

			//nolint
			ctx := context.WithValue(r.Context(), "identity", authTokenModel)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetIdentity(ctx context.Context) *types.AuthToken {
	if value, ok := ctx.Value("identity").(*types.AuthToken); ok {
		return value
	}

	return nil
}
