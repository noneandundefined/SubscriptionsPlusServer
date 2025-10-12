package store

import (
	"context"
	"database/sql"
	"errors"
	"subscriptionplus/server/infra/logger"
	"subscriptionplus/server/infra/store/postgres/models"
	"time"
)

type PlanStore struct {
	db     *sql.DB
	logger *logger.Logger
}

func (s *PlanStore) Get_Plans(ctx context.Context) (*[]models.Plan, error) {
	plans := []models.Plan{}

	query := `
		SELECT
			id,
			name,
			price,
			currency,
			auto_renewal_subscriptions,
			email_notification_subscriptions,
			max_total_subscriptions,
			auto_find_subscriptions
		FROM plans
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var plan models.Plan

		err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Price,
			&plan.Currency,
			&plan.AutoRenewalSubscriptions,
			&plan.EmailNotificationSubscriptions,
			&plan.MaxTotalSubscriptions,
			&plan.AutoFindSubscriptions,
		)

		if err != nil {
			return nil, err
		}

		plans = append(plans, plan)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &plans, nil
}

func (s *PlanStore) Get_PlanById(ctx context.Context, id uint64) (*models.Plan, error) {
	plan := models.Plan{}

	query := `
		SELECT
			id,
			name,
			price,
			currency,
			billing_period,
			auto_renewal_subscriptions,
			email_notification_subscriptions,
			max_total_subscriptions,
			auto_find_subscriptions
		FROM plans
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, id)

	if err := row.Scan(
		&plan.ID,
		&plan.Name,
		&plan.Price,
		&plan.Currency,
		&plan.BillingPeriod,
		&plan.AutoRenewalSubscriptions,
		&plan.EmailNotificationSubscriptions,
		&plan.MaxTotalSubscriptions,
		&plan.AutoFindSubscriptions,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, err
	}

	return &plan, nil
}

func (s *PlanStore) Get_PlanByUserSubUuid(ctx context.Context, uuid string) (*models.Plan, error) {
	plan := models.Plan{}

	query := `
		SELECT
			plans.id,
			plans.name,
			plans.price,
			plans.currency,
			plans.billing_period,
			plans.auto_renewal_subscriptions,
			plans.email_notification_subscriptions,
			plans.max_total_subscriptions,
			plans.auto_find_subscriptions
		FROM plans
		JOIN user_subscriptions ON plans.id = user_subscriptions.plan_id
		WHERE user_subscriptions.user_uuid = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, uuid)

	if err := row.Scan(
		&plan.ID,
		&plan.Name,
		&plan.Price,
		&plan.Currency,
		&plan.BillingPeriod,
		&plan.AutoRenewalSubscriptions,
		&plan.EmailNotificationSubscriptions,
		&plan.MaxTotalSubscriptions,
		&plan.AutoFindSubscriptions,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, err
	}

	return &plan, nil
}
