package transaction

import (
	"net/http"
	"subscriptionplus/server/infra/store/postgres/models"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"
	"subscriptionplus/server/util"

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

	httpx.HttpResponse(w, r, http.StatusCreated, xtoken)
	return nil
}
