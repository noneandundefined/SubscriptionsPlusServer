package store

import (
	"context"
	"database/sql"
	"subscriptionplus/server/infra/logger"
	"subscriptionplus/server/infra/store/postgres/models"
)

type Storage struct {
	Users interface { //nolint
		Create_UserCore(ctx context.Context, tx *sql.Tx, user *models.UserCore) error
		Create_UserSubscription(ctx context.Context, tx *sql.Tx, user *models.UserSubscription) error

		Get_UserCoreByUuid(ctx context.Context, uuid string) (*models.UserCore, error)
		Get_UserCoreByToken(ctx context.Context, token string) (*models.UserCore, error)
		Get_UserCoreByEmail(ctx context.Context, email string) (*models.UserCore, error)
		Get_UserSubscriptionAdvancedByUuid(ctx context.Context, uuid string) (*models.UserSubscriptionAdvanced, error)
		Get_UserMe(ctx context.Context, uuid string) (*models.UserMe, error)

		Delete_UserByUuid(ctx context.Context, uuid string) error
	}
	Subscriptions interface {
		Create_Subscription(ctx context.Context, sub *models.Subscription) error

		Get_SubscriptionsByUuid(ctx context.Context, uuid string) ([]models.Subscription, error)

		Update_SubscriptionById(ctx context.Context, sub *models.Subscription, id int) error
		Delete_SubscriptionById(ctx context.Context, id int, uuid string) error
	}
}

func NewStorage(db *sql.DB, logger *logger.Logger) Storage {
	return Storage{
		Users:         &UserStore{db, logger},
		Subscriptions: &SubscriptionStore{db, logger},
	}
}

func WithTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
