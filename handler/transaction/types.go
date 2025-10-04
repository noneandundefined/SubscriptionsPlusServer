package transaction

type SubscriptionPayPayload struct {
	PlanId int `json:"plan_id" validate:"required"`
}
