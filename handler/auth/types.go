package auth

type AuthCreatePayload struct {
	Email string `json:"email" validate:"required,email"`
}

type AuthReqRestoreAccessPayload struct {
	Email string `json:"email" validate:"required,email"`
}
