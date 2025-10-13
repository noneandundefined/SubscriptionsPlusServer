package transaction

import (
	"fmt"
	"net/http"
	"strconv"
	"subscriptionplus/server/infra/store/postgres/models"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/botx"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"
	"subscriptionplus/server/util"
	"time"

	"github.com/go-playground/validator"
)

func (h *Handler) SubscriptionPayHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *SubscriptionPayPayload

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

	if transactionPending != nil {
		return httperr.Conflict("you already have a pending subscription request")
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

	xtoken := util.GenerateXToken(authToken.User.UserUUID, int(plan.ID))

	transactionModel := &models.Transaction{
		PlanID:   int(plan.ID),
		UserUUID: authToken.User.UserUUID,
		XToken:   xtoken,
		Amount:   plan.Price,
	}

	if err := h.Store.Transactions.Create_Transaction(ctx, transactionModel); err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	// bot
	date := time.Now().Format("2006-01-02 15:04:05")

	amountStr := strconv.FormatFloat(plan.Price, 'f', 2, 64)

	message := fmt.Sprintf(
		"New Transaction in app %s\n\nuuid: %s\namount: %s\nxtoken: %s",
		botx.EscapeMarkdownV2(date),
		botx.EscapeMarkdownV2(authToken.User.UserUUID),
		botx.EscapeMarkdownV2(amountStr),
		botx.EscapeMarkdownV2(xtoken),
	)

	go botx.Send(message)

	httpx.HttpResponse(w, r, http.StatusCreated, xtoken)
	return nil
}
