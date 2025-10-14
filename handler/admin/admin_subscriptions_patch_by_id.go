package admin

import (
	"net/http"
	"strconv"
	"strings"
	"subscriptionplus/server/infra/store/postgres/models"
	"subscriptionplus/server/pkg/httpx"
	"subscriptionplus/server/pkg/httpx/httperr"
	"subscriptionplus/server/pkg/notify"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func (h *Handler) SubscriptionsPatchByIdHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	idParam := mux.Vars(r)["id"]
	if idParam == "" {
		return httperr.NotFound("subscription not found")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.InternalServerError(err.Error())
	}

	var payload *SubscriptionPatchPayload

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

	transaction, err := h.Store.Transactions.Get_TransactionsSubscriptionById(ctx, uint64(id), payload.UserUuid)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	if transaction == nil {
		httpx.HttpResponse(w, r, http.StatusOK, []models.Transaction{})
		return nil
	}

	userNotifyToken, err := h.Store.Notifications.Get_NotificationTokenByUuid(ctx, payload.UserUuid)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, err)
	}

	tx, err := h.Db.BeginTx(ctx, nil)
	if err != nil {
		h.Logger.Error("%v", err)
		return httperr.Db(ctx, httperr.Err_DbNetwork)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	status := strings.ToLower(payload.Status)

	switch status {
	case "failed":
		if errTr := h.Store.Transactions.Update_TransactionStatusById(tx, ctx, "failed", uint64(id)); errTr != nil {
			h.Logger.Error("%v", err)
			return httperr.Db(ctx, errTr)
		}

		if err := tx.Commit(); err != nil {
			h.Logger.Error("%v", err)
			return httperr.Db(ctx, err)
		}

		if err := notify.SendPushNotification(userNotifyToken.Token, "Subscription receipt error", "Unfortunately, we were unable to subscribe. Please try again."); err != nil {
			h.Logger.Error("%v", err)
		}

		httpx.HttpResponse(w, r, http.StatusOK, "Status successfully updated to failed!")
		return nil

	case "success":
		if errTr := h.Store.Transactions.Update_TransactionStatusById(tx, ctx, "success", uint64(id)); errTr != nil {
			h.Logger.Error("%v", err)
			return httperr.Db(ctx, errTr)
		}

		userSub := &models.UserSubscription{
			UserUUID: payload.UserUuid,
			PlanID:   transaction.PlanID,
		}

		if errUTr := h.Store.Users.Update_UserSubscriptionBeforPaySub(tx, ctx, userSub); errUTr != nil {
			h.Logger.Error("%v", err)
			return httperr.Db(ctx, errUTr)
		}

		if err := tx.Commit(); err != nil {
			h.Logger.Error("%v", err)
			return httperr.Db(ctx, err)
		}

		if err := notify.SendPushNotification(userNotifyToken.Token, "Subscription completed", "Your subscription has been activated. Thank you for staying with us."); err != nil {
			h.Logger.Error("%v", err)
		}

		httpx.HttpResponse(w, r, http.StatusOK, "Status successfully updated to success!")
		return nil

	default:
		return httperr.BadRequest("invalid status: must be 'success' or 'failed'")
	}
}
