package plan

import (
	"net/http"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"
)

func (h *Handler) GetPlansHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	plans, err := h.Store.Plans.Get_Plans(ctx)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	// httpx.HttpCache(w, 43200) // 12h.
	httpx.HttpResponseWithETag(w, r, http.StatusOK, plans)
	return nil
}
