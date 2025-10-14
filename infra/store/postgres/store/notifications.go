package store

import (
	"context"
	"database/sql"
	"errors"
	"subscriptionplus/server/infra/logger"
	"subscriptionplus/server/infra/store/postgres/models"
	"subscriptionplus/server/pkg/httpx/httperr"
	"time"
)

type NotificationStore struct {
	db     *sql.DB
	logger *logger.Logger
}

func (s *NotificationStore) Create_NotificationToken(ctx context.Context, notify *models.NotificationToken) error {
	query := `
		INSERT INTO notification_tokens (user_uuid, token)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, notify.UserUUID, notify.Token)
	if err != nil {
		return err
	}

	return nil
}

func (s *NotificationStore) Get_NotificationTokenByUuid(ctx context.Context, uuid string) (*models.NotificationToken, error) {
	notify := models.NotificationToken{}

	query := `
		SELECT * FROM notification_tokens WHERE user_uuid = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, uuid)

	if err := row.Scan(
		&notify.ID,
		&notify.CreatedAt,
		&notify.UpdatedAt,
		&notify.UserUUID,
		&notify.Token,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &notify, nil
}

func (s *NotificationStore) Update_NotificationTokenByUuid(ctx context.Context, notify *models.NotificationToken) error {
	query := `
		UPDATE notification_tokens
		SET
		    token = $1
		WHERE user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	upd, err := s.db.ExecContext(ctx, query, notify.Token, notify.UserUUID)
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
