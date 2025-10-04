package admin

type SubscriptionPatchPayload struct {
	Status   string `json:"status" validate:"required"`
	UserUuid string `json:"user_uuid" validate:"required"`
}
