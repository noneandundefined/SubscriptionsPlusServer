package store

import (
	"context"
	"database/sql"
	"strings"
	"subscriptionplus/server/infra/logger"
	"subscriptionplus/server/infra/store/postgres/models"
	"subscriptionplus/server/pkg/httpx/httperr"
	"time"
)

type UserStore struct {
	db     *sql.DB
	logger *logger.Logger
}

func (s *UserStore) Create_UserCore(ctx context.Context, tx *sql.Tx, user *models.UserCore) error {
	query := `
		INSERT INTO user_cores (user_uuid, email, token)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.UserUUID, user.Email, user.Token)
	if err != nil {
		if strings.Contains(err.Error(), `user_cores_email_key`) {
			return httperr.Err_DuplicateEmail
		}

		return err
	}

	return nil
}

func (s *UserStore) Create_UserSubscription(ctx context.Context, tx *sql.Tx, user *models.UserSubscription) error {
	query := `
		INSERT INTO user_subscriptions (user_uuid, plan_id, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.UserUUID, user.PlanID, user.StartDate, user.EndDate, user.IsActive)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Get_UserCoreByToken(ctx context.Context, token string) (*models.UserCore, error) {
	user := models.UserCore{}

	query := `
		SELECT id, user_uuid, email, token FROM user_cores WHERE token = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, token)

	if err := row.Scan(&user.ID, &user.UserUUID, &user.Email, &user.Token); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Get_UserCoreByUuid(ctx context.Context, uuid string) (*models.UserCore, error) {
	user := models.UserCore{}

	query := `
		SELECT id, user_uuid, email, token FROM user_cores WHERE user_uuid = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, uuid)

	if err := row.Scan(&user.ID, &user.UserUUID, &user.Email, &user.Token); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Get_UserCoreByEmail(ctx context.Context, email string) (*models.UserCore, error) {
	user := models.UserCore{}

	query := `
		SELECT id, user_uuid, email, token FROM user_cores WHERE email = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, email)

	if err := row.Scan(&user.ID, &user.UserUUID, &user.Email, &user.Token); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Get_UserMe(ctx context.Context, uuid string) (*models.UserMe, error) {
	user := models.UserMe{}

	query := `
		SELECT
			user_cores.created_at,
			user_cores.user_uuid,
			user_cores.email,
			plans.name as plan_name,
			user_subscriptions.end_date,
			user_subscriptions.is_active
		FROM user_cores
		LEFT JOIN user_subscriptions ON user_cores.user_uuid = user_subscriptions.user_uuid
		LEFT JOIN plans ON user_subscriptions.plan_id = plans.id
		WHERE user_cores.user_uuid = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, uuid)

	if err := row.Scan(
		&user.CreatedAt,
		&user.UserUUID,
		&user.Email,
		&user.PlanName,
		&user.EndDate,
		&user.IsActive,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Delete_UserByUuid(ctx context.Context, uuid string) error {
	query := `
		DELETE FROM user_cores WHERE user_uuid = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	del, err := s.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	delRows, _ := del.RowsAffected()
	if delRows == 0 {
		return httperr.Err_NotDeleted
	}

	return nil
}

func (s *UserStore) Get_UserSubscriptionAdvancedByUuid(ctx context.Context, uuid string) (*models.UserSubscriptionAdvanced, error) {
	user := models.UserSubscriptionAdvanced{}

	query := `
		SELECT
			user_subscriptions.id,
			user_subscriptions.user_uuid,
			plans.name,
			plans.price,
			user_subscriptions.end_date,
			plans.max_total_subscriptions,
			plans.auto_find_subscriptions
		FROM user_subscriptions
		JOIN plans ON user_subscriptions.plan_id = plans.id
		JOIN user_cores ON user_subscriptions.user_uuid = user_cores.user_uuid
		WHERE user_subscriptions.user_uuid = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, uuid)

	if err := row.Scan(
		&user.ID,
		&user.UserUUID,
		&user.Name,
		&user.Price,
		&user.EndDate,
		&user.MaxTotalSubscriptions,
		&user.AutoFindSubscriptions,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
