package models

import "time"

type UserCore struct {
	ID           uint64    `json:"id" db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	UserUUID     string    `json:"user_uuid" db:"user_uuid"`
	Email        string    `json:"email" db:"email"`
	AccessToken  string    `json:"access_token" db:"access_token"`
	RefreshToken *string   `json:"refresh_token" db:"refresh_token"`
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

type UserUsage struct {
	ID                             uint64    `json:"id" db:"id"`
	CreatedAt                      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                      time.Time `json:"updated_at" db:"updated_at"`
	UserUUID                       string    `json:"user_uuid" db:"user_uuid"`
	AutoRenewalSubscriptions       bool      `json:"auto_renewal_subscriptions" db:"auto_renewal_subscriptions"`
	EmailNotificationSubscriptions bool      `json:"email_notification_subscriptions" db:"email_notification_subscriptions"`
}

type UserInfo struct {
	ID                                  uint64     `json:"id" db:"id"`
	UserUUID                            string     `json:"user_uuid" db:"user_uuid"`
	Email                               string     `json:"email" db:"email"`
	AccessToken                         string     `json:"access_token" db:"access_token"`
	RefreshToken                        *string    `json:"refresh_token" db:"refresh_token"`
	PlanName                            string     `json:"plan_name" db:"plan_name"`
	Price                               float64    `json:"price" db:"price"`
	EndDate                             *time.Time `json:"end_date,omitempty" db:"end_date"` // NULL, если активна или бессрочная
	IsActive                            bool       `json:"is_active" db:"is_active"`
	AutoRenewalSubscriptions            bool       `json:"auto_renewal_subscriptions" db:"auto_renewal_subscriptions"`
	EmailNotificationSubscriptions      bool       `json:"email_notification_subscriptions" db:"email_notification_subscriptions"`
	AutoRenewalSubscriptionsUsage       bool       `json:"auto_renewal_subscriptions_usage" db:"auto_renewal_subscriptions_usage"`
	EmailNotificationSubscriptionsUsage bool       `json:"email_notification_subscriptions_usage" db:"email_notification_subscriptions_usage"`
	MaxTotalSubscriptions               *int       `json:"max_total_subscriptions,omitempty" db:"max_total_subscriptions"`
	AutoFindSubscriptions               bool       `json:"auto_find_subscriptions" db:"auto_find_subscriptions"`
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
