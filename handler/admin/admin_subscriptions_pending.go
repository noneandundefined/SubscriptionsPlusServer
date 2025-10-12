package admin

import (
	"net/http"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"
)

func (h *Handler) SubscriptionsPendingHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	subs, err := h.Store.Transactions.Get_TransactionsByStatus(ctx, "pending")
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, subs)
	return nil
}
