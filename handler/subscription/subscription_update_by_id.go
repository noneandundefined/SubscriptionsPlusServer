package subscription

import (
	"net/http"
	"strconv"
	"subscriptionplus/server/infra/store/postgres/models"
	"subscriptionplus/server/infra/types"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func (h *Handler) EditSubscriptions(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *SubscriptionsEditPayload

	idParam := mux.Vars(r)["id"]
	if idParam == "" {
		return httperr.NotFound("subscription not found")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.InternalServerError(err.Error())
	}

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

	if err := h.Store.Subscriptions.Update_SubscriptionById(ctx, sub, id); err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, "Subscription changed")
	return nil
}
