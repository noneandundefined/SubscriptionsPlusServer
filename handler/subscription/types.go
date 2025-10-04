package subscription

import "time"

type SubscriptionAddPayload struct {
	Name            string     `json:"name" validate:"required,min=3"`
	Price           float64    `json:"price" validate:"required"`
	DatePay         time.Time  `json:"date_pay" validate:"required"`
	DateNotifyOne   *time.Time `json:"date_notify_one,omitempty"`
	DateNotifyTwo   *time.Time `json:"date_notify_two,omitempty"`
	DateNotifyThree *time.Time `json:"date_notify_three,omitempty"`
}

type SubscriptionsEditPayload struct {
	Name            string     `json:"name" validate:"required,min=3"`
	Price           float64    `json:"price" validate:"required"`
	DatePay         time.Time  `json:"date_pay" validate:"required"`
	DateNotifyOne   *time.Time `json:"date_notify_one,omitempty"`
	DateNotifyTwo   *time.Time `json:"date_notify_two,omitempty"`
	DateNotifyThree *time.Time `json:"date_notify_three,omitempty"`
}
