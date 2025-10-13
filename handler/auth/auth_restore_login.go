package auth

import (
	"net/http"
	"strings"
	"subscriptionplus/server/infra/encryption"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"
	"time"
)

func (h *Handler) AuthRestoreLoginHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	token := r.URL.Query().Get("token")
	if token == "" {
		return httperr.Conflict("")
	}

	decrypted, err := encryption.Decrypt(token)
	if err != nil {
		h.Logger.Error("failed to decrypt token: %v", err)
		return httperr.Conflict("")
	}

	parts := strings.Split(decrypted, ";")
	if len(parts) != 3 {
		return httperr.Conflict("")
	}

	createdAtStr := parts[0]
	email := parts[1]
	endedAtStr := parts[2]

	createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return httperr.Conflict("")
	}

	expiredAt, err := time.Parse("2006-01-02 15:04:05", endedAtStr)
	if err != nil {
		return httperr.Conflict("")
	}

	now := time.Now()

	if now.After(expiredAt) {
		return httperr.Conflict("")
	}

	if now.Sub(createdAt) > time.Hour {
		return httperr.Conflict("")
	}

	user, err := h.Store.Users.Get_UserCoreByEmail(ctx, email)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Conflict("")
	}

	if user == nil {
		return httperr.Conflict("")
	}

	if user.RefreshToken == nil {
		httpx.HttpResponse(w, r, http.StatusNoContent, nil)
		return nil
	}

	if *user.RefreshToken == token {
		httpx.HttpResponse(w, r, http.StatusOK, user.AccessToken)
		return nil
	}

	httpx.HttpResponse(w, r, http.StatusNoContent, nil)
	return nil
}
