package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"subscriptionplus/server/infra/encryption"
	"subscriptionplus/server/pkg"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"
	"time"

	"github.com/go-playground/validator"
)

func (h *Handler) AuthRequestRestoreAccessHandler(w http.ResponseWriter, r *http.Request) error {
	env := os.Getenv("GO_ENV")

	var payload *AuthReqRestoreAccessPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		h.Logger.Error("%v", err)
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		h.Logger.Error("%v", err)
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(err))
		}

		return httperr.BadRequest("not all fields are filled")
	}

	tokenDto := fmt.Sprintf("%s;%s;%s", time.Now().Format("2006-01-02 15:04:05"), payload.Email, time.Now().Add(time.Hour).Format("2006-01-02 15:04:05"))

	refreshToken, err := encryption.Encrypt(tokenDto)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.InternalServerError(err.Error())
	}

	var link string

	if env == "DEV" {
		link = fmt.Sprintf("%s://%s:8080/api/v1/auth/restore_access?token=%s", os.Getenv("BACKEND_PROTOCOL"), os.Getenv("BACKEND_HOST"), url.QueryEscape(refreshToken))
	} else {
		link = fmt.Sprintf("%s://%s/api/v1/auth/restore_access?token=%s", os.Getenv("BACKEND_PROTOCOL"), os.Getenv("BACKEND_HOST"), url.QueryEscape(refreshToken))
	}

	if errSendEmail := pkg.SendEmail(payload.Email, "Restoring a Subscription+ account", fmt.Sprintf(`
	<body>
		<p>To reset your account access, follow the link below:</p>

		<a href="%s">
			%s
		</a>

		<p>If you did not request password recovery, simply ignore this email..</p>
	</body>
	`, link, link)); errSendEmail != nil {
		h.Logger.Error("%v", errSendEmail)
		return httperr.InternalServerError(errSendEmail.Error())
	}

	h.Logger.Info("Reset link for %s: %s", payload.Email, link)

	httpx.HttpResponse(w, r, http.StatusOK, map[string]interface{}{
		"token": refreshToken,
		"msg":   "A link to restore your access has been sent to your email.",
	})
	return nil
}
