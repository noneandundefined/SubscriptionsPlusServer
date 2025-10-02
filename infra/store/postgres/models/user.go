package models

import "time"

type UserCore struct {
	ID        uint64    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserUUID  string    `json:"user_uuid" db:"user_uuid"`
	Email     string    `json:"email" db:"email"`
	Token     string    `json:"token" db:"token"`
}

type UserSubscription struct {
	ID        uint64     `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	UserUUID  string     `json:"user_uuid" db:"user_uuid"`
	PlanID    int        `json:"plan_id,omitempty" db:"plan_id"`
	StartDate time.Time  `json:"start_date" db:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" db:"end_date"` // NULL, если активна или бессрочная
	IsActive  bool       `json:"is_active" db:"is_active"`
}

type UserSubscriptionAdvanced struct {
	ID                    uint64     `json:"id" db:"id"`
	UserUUID              string     `json:"user_uuid" db:"user_uuid"`
	Name                  string     `json:"name" db:"name"`
	Price                 float64    `json:"price" db:"price"`
	EndDate               *time.Time `json:"end_date,omitempty" db:"end_date"`
	MaxTotalSubscriptions *int       `json:"max_total_subscriptions,omitempty" db:"max_total_subscriptions"`
	AutoFindSubscriptions bool       `json:"auto_find_subscriptions" db:"auto_find_subscriptions"`
}

type UserMe struct {
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UserUUID  string     `json:"user_uuid" db:"user_uuid"`
	Email     string     `json:"email" db:"email"`
	PlanName  string     `json:"plan_name" db:"plan_name"`
	EndDate   *time.Time `json:"end_date,omitempty" db:"end_date"` // NULL, если активна или бессрочная
	IsActive  bool       `json:"is_active" db:"is_active"`
}
