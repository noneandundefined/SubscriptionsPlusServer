package auth

type AuthCreatePayload struct {
	Email string `json:"email" validate:"required,email"`
}
