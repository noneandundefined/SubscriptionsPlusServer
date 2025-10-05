package transaction

import (
	"net/http"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"

	"github.com/gorilla/mux"
)

func (h *Handler) SubscriptionGetByXToken(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	xtoken := mux.Vars(r)["xtoken"]
	if xtoken == "" {
		return httperr.NotFound("subscription not found")
	}

	transactionSub, err := h.Store.Transactions.Get_TransactionsSubscriptionByXToken(ctx, xtoken, authToken.User.UserUUID)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, transactionSub)
	return nil
}
