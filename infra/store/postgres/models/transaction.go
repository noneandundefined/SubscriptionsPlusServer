package models

import "time"

type Transaction struct {
	ID        uint64    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserUUID  string    `json:"user_uuid" db:"user_uuid"`
	Amount    float64   `json:"amount" db:"amount"`
	Status    string    `json:"status" db:"status"` // pending, success, failed
}
