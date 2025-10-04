package subscription

import (
	"net/http"
	"strconv"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"

	"github.com/gorilla/mux"
)

func (h *Handler) GetSubscriptionById(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	idParam := mux.Vars(r)["id"]
	if idParam == "" {
		return httperr.NotFound("subscription not found")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.InternalServerError(err.Error())
	}

	sub, err := h.Store.Subscriptions.Get_SubscriptionById(ctx, uint64(id), authToken.User.UserUUID)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, sub)
	return nil
}
