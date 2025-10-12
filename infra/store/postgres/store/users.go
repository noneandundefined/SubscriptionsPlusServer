package store

import (
	"context"
	"database/sql"
	"errors"
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
		INSERT INTO user_cores (user_uuid, email, access_token)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.UserUUID, user.Email, user.AccessToken)
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
		INSERT INTO user_subscriptions (user_uuid, plan_id, start_date, is_active)
		VALUES ($1, $2, $3, $4)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.UserUUID, user.PlanID, user.StartDate, user.IsActive)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Create_UserUsage(ctx context.Context, tx *sql.Tx, uuid string) error {
	query := `
		INSERT INTO user_usages (user_uuid, auto_renewal_subscriptions, email_notification_subscriptions)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, uuid, false, false)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Get_UserCoreByAccessToken(ctx context.Context, access_token string) (*models.UserCore, error) {
	user := models.UserCore{}

	query := `
		SELECT 
		    id,
		    user_uuid,
		    email,
		    access_token, 
		    refresh_token 
		FROM user_cores 
		WHERE access_token = $1 
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, access_token)

	if err := row.Scan(&user.ID, &user.UserUUID, &user.Email, &user.AccessToken, &user.RefreshToken); err != nil {
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
		SELECT 
		    id, 
		    user_uuid, 
		    email, 
		    access_token, 
		    refresh_token 
		FROM user_cores 
		WHERE user_uuid = $1 
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, uuid)

	if err := row.Scan(&user.ID, &user.UserUUID, &user.Email, &user.AccessToken, &user.RefreshToken); err != nil {
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
		SELECT 
		    id, 
		    user_uuid, 
		    email, 
		    access_token, 
		    refresh_token 
		FROM user_cores 
		WHERE email = $1 
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, email)

	if err := row.Scan(&user.ID, &user.UserUUID, &user.Email, &user.AccessToken, &user.RefreshToken); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Get_UserInfoByUuid(ctx context.Context, uuid string) (*models.UserInfo, error) {
	user := models.UserInfo{}

	query := `
		SELECT
			user_cores.user_uuid,
			user_cores.email,
			user_cores.access_token,
			user_cores.refresh_token,
			plans.name as plan_name,
			plans.price,
			user_subscriptions.end_date,
			user_subscriptions.is_active,
			plans.auto_renewal_subscriptions,
			plans.email_notification_subscriptions,
			user_usages.auto_renewal_subscriptions as auto_renewal_subscriptions_usage,
			user_usages.email_notification_subscriptions as email_notification_subscriptions_usage,
			plans.max_total_subscriptions,
			plans.auto_find_subscriptions
		FROM user_cores
		LEFT JOIN user_subscriptions ON user_cores.user_uuid = user_subscriptions.user_uuid
		LEFT JOIN plans ON user_subscriptions.plan_id = plans.id
		LEFT JOIN user_usages ON user_cores.user_uuid = user_usages.user_uuid
		WHERE user_cores.user_uuid = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, uuid)

	if err := row.Scan(
		&user.UserUUID,
		&user.Email,
		&user.AccessToken,
		&user.RefreshToken,
		&user.PlanName,
		&user.Price,
		&user.EndDate,
		&user.IsActive,
		&user.AutoRenewalSubscriptions,
		&user.EmailNotificationSubscriptions,
		&user.AutoRenewalSubscriptionsUsage,
		&user.EmailNotificationSubscriptionsUsage,
		&user.MaxTotalSubscriptions,
		&user.AutoFindSubscriptions,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Get_UserInfoByAccessToken(ctx context.Context, token string) (*models.UserInfo, error) {
	user := models.UserInfo{}

	query := `
		SELECT
			user_cores.user_uuid,
			user_cores.email,
			user_cores.access_token,
			user_cores.refresh_token,
			plans.name as plan_name,
			plans.price,
			user_subscriptions.end_date,
			user_subscriptions.is_active,
			plans.auto_renewal_subscriptions,
			plans.email_notification_subscriptions,
			user_usages.auto_renewal_subscriptions as auto_renewal_subscriptions_usage,
			user_usages.email_notification_subscriptions as email_notification_subscriptions_usage,
			plans.max_total_subscriptions,
			plans.auto_find_subscriptions
		FROM user_cores
		LEFT JOIN user_subscriptions ON user_cores.user_uuid = user_subscriptions.user_uuid
		LEFT JOIN plans ON user_subscriptions.plan_id = plans.id
		LEFT JOIN user_usages ON user_cores.user_uuid = user_usages.user_uuid
		WHERE user_cores.access_token = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, token)

	if err := row.Scan(
		&user.UserUUID,
		&user.Email,
		&user.AccessToken,
		&user.RefreshToken,
		&user.PlanName,
		&user.Price,
		&user.EndDate,
		&user.IsActive,
		&user.AutoRenewalSubscriptions,
		&user.EmailNotificationSubscriptions,
		&user.AutoRenewalSubscriptionsUsage,
		&user.EmailNotificationSubscriptionsUsage,
		&user.MaxTotalSubscriptions,
		&user.AutoFindSubscriptions,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

func (s *UserStore) Update_UserCoreRefreshTokenByUuid(ctx context.Context, refreshToken, uuid string) error {
	query := `
		UPDATE user_cores SET refresh_token = $1 WHERE user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	upd, err := s.db.ExecContext(ctx, query, refreshToken, uuid)
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

func (s *UserStore) Update_UserUsageByUuid(ctx context.Context, uuid string, auto_renewal_subscriptions, email_notification_subscriptions bool) error {
	query := `
		UPDATE user_usages SET auto_renewal_subscriptions = $1, email_notification_subscriptions = $2 WHERE user_uuid = $3
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	upd, err := s.db.ExecContext(ctx, query, uuid, auto_renewal_subscriptions, email_notification_subscriptions)
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

func (s *UserStore) Update_UserSubscriptionBeforPaySub(tx *sql.Tx, ctx context.Context, user *models.UserSubscription) error {
	query := `
		UPDATE user_subscriptions
		SET
			plan_id = $1,
			start_date = NOW(),
			end_date = NOW() + INTERVAL '1 month',
			is_active = true
		WHERE user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	upd, err := tx.ExecContext(ctx, query, user.PlanID, user.UserUUID)
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

func (s *UserStore) Update_UserSubscriptionBeforEndSub(ctx context.Context) error {
	query := `
		UPDATE user_subscriptions
		SET
			is_active = false,
			plan_id = 2
        WHERE end_day < CURRENT_DATE AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
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
