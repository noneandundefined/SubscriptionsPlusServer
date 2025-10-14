package notification

import (
	"net/http"
	"subscriptionplus/server/infra/store/postgres/models"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"

	"github.com/go-playground/validator"
)

func (h *Handler) NotificationTokenHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *TokenPayload

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

	notifyToken, err := h.Store.Notifications.Get_NotificationTokenByUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	if notifyToken != nil && notifyToken.Token == payload.Token {
		httpx.HttpResponse(w, r, http.StatusNoContent, nil)
		return nil
	}

	if err := h.Store.Notifications.Create_NotificationToken(ctx, &models.NotificationToken{
		UserUUID: authToken.User.UserUUID,
		Token:    payload.Token,
	}); err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusNoContent, nil)
	return nil
}
