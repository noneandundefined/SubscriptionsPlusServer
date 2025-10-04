package transaction

import (
	"net/http"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"
)

func (h *Handler) SubscriptionGetPending(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	transactionPending, err := h.Store.Transactions.Get_TransactionPendingByUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, transactionPending)
	return nil
}
