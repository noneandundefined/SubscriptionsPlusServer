package auth

import (
	"net/http"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"
)

func (h *Handler) AuthDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	if err := h.Store.Users.Delete_UserByUuid(ctx, authToken.User.UserUUID); err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, "Account has been deleted")
	return nil
}
