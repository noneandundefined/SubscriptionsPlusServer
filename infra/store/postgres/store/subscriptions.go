package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"subscriptionplus/server/infra/logger"
	"subscriptionplus/server/infra/store/postgres/models"
	"subscriptionplus/server/pkg/httpx/httperr"
	"time"
)

type SubscriptionStore struct {
	db     *sql.DB
	logger *logger.Logger
}

func (s *SubscriptionStore) Get_SubscriptionById(ctx context.Context, id uint64, uuid string) (*models.Subscription, error) {
	subscription := models.Subscription{}

	query := `
		SELECT * FROM subscriptions WHERE id = $1 AND user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, id, uuid)

	if err := row.Scan(
		&subscription.ID,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
		&subscription.UserUUID,
		&subscription.Name,
		&subscription.Price,
		&subscription.DatePay,
		&subscription.DateNotifyOne,
		&subscription.DateNotifyTwo,
		&subscription.DateNotifyThree,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &subscription, nil
}

func (s *SubscriptionStore) Create_Subscription(ctx context.Context, sub *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (user_uuid, name, price, date_pay, date_notify_one, date_notify_two, date_notify_three)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, sub.UserUUID, sub.Name, sub.Price, sub.DatePay, sub.DateNotifyOne, sub.DateNotifyTwo, sub.DateNotifyThree)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionStore) Get_SubscriptionsByUuid(ctx context.Context, uuid, search string) (*[]models.Subscription, error) {
	subscriptions := []models.Subscription{}

	var limit sql.NullInt32
	if err := s.db.QueryRowContext(ctx, `
        SELECT COALESCE(plans.max_total_subscriptions, 1000000) FROM user_subscriptions
        JOIN plans ON user_subscriptions.plan_id = plans.id
        WHERE user_subscriptions.user_uuid = $1
        LIMIT 1
    `, uuid).Scan(&limit); err != nil {
		return nil, err
	}

	maxLimit := int(limit.Int32)
	if maxLimit == 0 {
		maxLimit = 10
	}

	query := `
		SELECT * FROM subscriptions
		WHERE user_uuid = $1
	`

	args := []interface{}{uuid}
	paramIndex := 2

	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", paramIndex)
		args = append(args, "%"+search+"%")
		paramIndex++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", paramIndex)

	args = append(args, maxLimit)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sub models.Subscription

		err := rows.Scan(
			&sub.ID,
			&sub.CreatedAt,
			&sub.UpdatedAt,
			&sub.UserUUID,
			&sub.Name,
			&sub.Price,
			&sub.DatePay,
			&sub.DateNotifyOne,
			&sub.DateNotifyTwo,
			&sub.DateNotifyThree,
		)

		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &subscriptions, nil
}

func (s *SubscriptionStore) Update_SubscriptionById(ctx context.Context, sub *models.Subscription, id int) error {
	query := `
		UPDATE subscriptions SET name = $1, price = $2, date_pay = $3, date_notify_one = $4, date_notify_two = $5, date_notify_three = $6 WHERE id = $7 AND user_uuid = $8
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	upd, err := s.db.ExecContext(ctx, query, sub.Name, sub.Price, sub.DatePay, sub.DateNotifyOne, sub.DateNotifyTwo, sub.DateNotifyThree, id, sub.UserUUID)
	if err != nil {
		return err
	}

	updAffected, err := upd.RowsAffected()
	if err != nil {
		return err
	}

	if updAffected == 0 {
		return httperr.Err_NotUpdated
	}

	return nil
}

func (s *SubscriptionStore) Update_SubscriptionsMounth(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	queries := []string{
		"UPDATE subscriptions SET date_pay = date_pay + INTERVAL '1 month' WHERE date_pay <= CURRENT_DATE",
		"UPDATE subscriptions SET date_notify_one = date_notify_one + INTERVAL '1 month' WHERE date_notify_one <= CURRENT_DATE",
		"UPDATE subscriptions SET date_notify_two = date_notify_two + INTERVAL '1 month' WHERE date_notify_two <= CURRENT_DATE",
		"UPDATE subscriptions SET date_notify_three = date_notify_three + INTERVAL '1 month' WHERE date_notify_three <= CURRENT_DATE",
	}

	for _, q := range queries {
		if _, err := s.db.ExecContext(ctx, q); err != nil {
			return err
		}
	}

	return nil
}

func (s *SubscriptionStore) Delete_SubscriptionById(ctx context.Context, id int, uuid string) error {
	query := `
		DELETE FROM subscriptions WHERE id = $1 AND user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	del, err := s.db.ExecContext(ctx, query, id, uuid)
	if err != nil {
		return err
	}

	delAffected, err := del.RowsAffected()
	if err != nil {
		return err
	}

	if delAffected == 0 {
		return httperr.Err_NotDeleted
	}

	return nil
}
