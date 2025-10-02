package models

import "time"

type Plan struct {
	ID                    uint64    `json:"id" db:"id"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
	Name                  string    `json:"name" db:"name"`
	Price                 float64   `json:"price" db:"price"`
	MaxTotalSubscriptions *int      `json:"max_total_subscriptions,omitempty" db:"max_total_subscriptions"`
	AutoFindSubscriptions bool      `json:"auto_find_subscriptions" db:"auto_find_subscriptions"`
}
