package transaction

import (
	"fmt"
	"net/http"
	"strconv"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/botx"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"
	"time"

	"github.com/go-playground/validator"
)

func (h *Handler) SubscriptionWaitCheckHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *SubscriptionWaitCcheckPayload

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

	if payload.PlanId <= 0 {
		return httperr.BadRequest("not all fields are filled")
	}

	transactionPending, err := h.Store.Transactions.Get_TransactionPendingByUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	if transactionPending == nil {
		return httperr.Conflict("you don't have any active subscriptions.")
	}

	tx, err := h.Db.BeginTx(ctx, nil)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, httperr.Err_DbNetwork)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := h.Store.Transactions.Update_TransactionStatusById(tx, ctx, "paid", uint64(transactionPending.ID)); err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	if err := tx.Commit(); err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	plan, err := h.Store.Plans.Get_PlanById(ctx, uint64(payload.PlanId))
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	if plan == nil {
		h.Logger.Error("%v", err)
		return httperr.InternalServerError("we don't have an active plan")
	}

	// bot
	date := time.Now().Format("2006-01-02 15:04:05")

	amountStr := strconv.FormatFloat(plan.Price, 'f', 2, 64)

	message := fmt.Sprintf(
		"New Transaction in app %s\n\nuuid: %s\namount: %s\nxtoken: %s",
		botx.EscapeMarkdownV2(date),
		botx.EscapeMarkdownV2(authToken.User.UserUUID),
		botx.EscapeMarkdownV2(amountStr),
		botx.EscapeMarkdownV2(payload.XToken),
	)

	go botx.Send(message)

	httpx.HttpResponse(w, r, http.StatusOK, "Give us a couple of minutes to verify the payment.")
	return nil
}
