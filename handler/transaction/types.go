package transaction

type SubscriptionPayPayload struct {
	PlanId int `json:"plan_id" validate:"required"`
}

type SubscriptionWaitCcheckPayload struct {
	PlanId int    `json:"plan_id" validate:"required"`
	XToken string `json:"x_token" validate:"required"`
}
