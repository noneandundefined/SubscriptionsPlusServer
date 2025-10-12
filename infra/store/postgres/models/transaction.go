package models

import "time"

type Transaction struct {
	ID        uint64    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserUUID  string    `json:"user_uuid" db:"user_uuid"`
	PlanID    int       `json:"plan_id,omitempty" db:"plan_id"`
	Status    string    `json:"status" db:"status"`
	XToken    string    `json:"x_token" db:"x_token"`
	Amount    float64   `json:"amount" db:"amount"`
	Currency  string    `json:"currency" db:"currency"`
}
