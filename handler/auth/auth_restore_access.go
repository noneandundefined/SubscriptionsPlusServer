package auth

import (
	"fmt"
	"net/http"
	"strings"
	"subscriptionplus/server/infra/encryption"
	"time"
)

func renderHtml(title, message string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Subscription+ Restore Access</title>
	<style>
	  body { margin:0; padding:0; font-family: ui-sans-serif,system-ui,sans-serif,"Apple Color Emoji","Segoe UI Emoji",Segoe UI Symbol,"Noto Color Emoji"; background-color:#000; display:flex; justify-content:center; align-items:center; height:100vh; width:100vw; }
	  .container { max-width:80vw; width:90%%; border-radius:6px; box-shadow:0 8px 20px rgba(0,0,0,0.4); padding:40px 30px; text-align:center; justify-content:center; color:#ccc; }
	  .logo { max-width:100px; margin-bottom:30px; border-radius:20%%; box-shadow:0 4px 15px rgba(0,0,0,0.2); }
	  h1 { font-size:28px; font-weight:700; margin-bottom:20px; color:#fff; }
	  p { font-size:15px; color:#bbb; margin:10px 0; line-height:1.6; white-space:pre-wrap; }
	</style>
	</head>
	<body>
	<div class="container">
	  <img class="logo" src="https://github.com/noneandundefined/SubscriptionsPlus/blob/main/assets/images/sub-icon-base.png?raw=true" alt="Subscription+ Logo">
	  <h1>%s</h1>
	  <p>%s</p>
	</div>
	</body>
	</html>`, title, message)
}

func (h *Handler) AuthRestoreAccessHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	renderError := func(msg string) {
		html := renderHtml("Failed restore access", msg)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		renderError("Please open the link from your email to restore access.")
		return nil
	}

	decrypted, err := encryption.Decrypt(token)
	if err != nil {
		h.Logger.Error("failed to decrypt token: %v", err)

		renderError("The token seems incorrect. Request a new password reset link from the app.")
		return nil
	}

	parts := strings.Split(decrypted, ";")
	if len(parts) != 3 {
		renderError("The link may be outdated. Please request a new reset link.")
		return nil
	}

	createdAtStr := parts[0]
	email := parts[1]
	endedAtStr := parts[2]

	createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		renderError("Internal error. Try again or contact support.")
		return nil
	}

	expiredAt, err := time.Parse("2006-01-02 15:04:05", endedAtStr)
	if err != nil {
		renderError("Internal error. Try again or contact support.")
		return nil
	}

	now := time.Now()

	if now.After(expiredAt) {
		renderError("This link has expired. Please request a new password reset link from the app.")
		return nil
	}

	if now.Sub(createdAt) > time.Hour {
		renderError("This link has expired. Please request a new password reset link from the app.")
		return nil
	}

	user, err := h.Store.Users.Get_UserCoreByEmail(ctx, email)
	if err != nil {
		h.Logger.Error("%v", err)

		renderError("Unexpected server error. Try again later.")
		return nil
	}

	if user == nil {
		renderError("We could not find your account. Check your email or contact support.")
		return nil
	}

	if err := h.Store.Users.Update_UserCoreRefreshTokenByUuid(ctx, token, user.UserUUID); err != nil {
		h.Logger.Error("%v", err)

		renderError("Temporary server error. Try again in a few minutes.")
		return nil
	}

	html := renderHtml("Access Restored", "Your access has been successfully restored.\nPlease open your Subscription+ application to continue. Your refresh token is active for 1 hour.")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))

	return nil
}
