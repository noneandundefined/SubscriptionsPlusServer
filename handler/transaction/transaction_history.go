package transaction

import (
	"net/http"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"
)

func (h *Handler) TransactionsHistory(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	history, err := h.Store.Transactions.Get_TransactionsByUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, history)
	return nil
}
