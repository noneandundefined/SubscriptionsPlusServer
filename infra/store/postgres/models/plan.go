package models

import "time"

type Plan struct {
	ID                             uint64    `json:"id" db:"id"`
	CreatedAt                      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                      time.Time `json:"updated_at" db:"updated_at"`
	Name                           string    `json:"name" db:"name"`
	Description                    string    `json:"description" db:"description"`
	Price                          float64   `json:"price" db:"price"`
	Currency                       string    `json:"currency" db:"currency"`
	BillingPeriod                  string    `json:"billing_period" db:"billing_period"`
	AutoRenewalSubscriptions       bool      `json:"auto_renewal_subscriptions" db:"auto_renewal_subscriptions"`
	EmailNotificationSubscriptions bool      `json:"email_notification_subscriptions" db:"email_notification_subscriptions"`
	MaxTotalSubscriptions          *int      `json:"max_total_subscriptions,omitempty" db:"max_total_subscriptions"`
	AutoFindSubscriptions          bool      `json:"auto_find_subscriptions" db:"auto_find_subscriptions"`
}
