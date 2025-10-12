package subscription

import (
	"net/http"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"
)

func (h *Handler) GetSubscriptionsHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	search := r.URL.Query().Get("search")

	subs, err := h.Store.Subscriptions.Get_SubscriptionsByUuid(ctx, authToken.User.UserUUID, search)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, subs)
	return nil
}
