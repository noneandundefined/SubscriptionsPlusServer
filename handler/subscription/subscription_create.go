package subscription

import (
	"net/http"
	"subscriptionplus/server/infra/store/postgres/models"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"

	"github.com/go-playground/validator"
)

func (h *Handler) AddSubscriptions(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *SubscriptionAddPayload

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

	sub := &models.Subscription{
		UserUUID:        authToken.User.UserUUID,
		Name:            payload.Name,
		Price:           payload.Price,
		DatePay:         payload.DatePay,
		DateNotifyOne:   payload.DateNotifyOne,
		DateNotifyTwo:   payload.DateNotifyTwo,
		DateNotifyThree: payload.DateNotifyThree,
	}

	if err := h.Store.Subscriptions.Create_Subscription(ctx, sub); err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusCreated, "Successful subscription addition!")
	return nil
}
