package models

import "time"

type Subscription struct {
	ID              uint64     `json:"id" db:"id"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	UserUUID        string     `json:"user_uuid" db:"user_uuid"`
	Name            string     `json:"name" db:"name"`
	Price           float64    `json:"price" db:"price"`
	DatePay         time.Time  `json:"date_pay" db:"date_pay"`
	DateNotifyOne   *time.Time `json:"date_notify_one,omitempty" db:"date_notify_one"`
	DateNotifyTwo   *time.Time `json:"date_notify_two,omitempty" db:"date_notify_two"`
	DateNotifyThree *time.Time `json:"date_notify_three,omitempty" db:"date_notify_three"`
	AutoRenewal     bool       `json:"auto_renewal" db:"auto_renewal"`
}
