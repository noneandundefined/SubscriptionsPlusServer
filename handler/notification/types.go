package notification

type TokenPayload struct {
	Token string `json:"token" validate:"required"`
}
