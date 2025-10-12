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
		Create_UserUsage(ctx context.Context, tx *sql.Tx, uuid string) error

		Get_UserCoreByUuid(ctx context.Context, uuid string) (*models.UserCore, error)
		Get_UserInfoByUuid(ctx context.Context, uuid string) (*models.UserInfo, error)
		Get_UserInfoByAccessToken(ctx context.Context, token string) (*models.UserInfo, error)
		Get_UserCoreByAccessToken(ctx context.Context, access_token string) (*models.UserCore, error)
		Get_UserCoreByEmail(ctx context.Context, email string) (*models.UserCore, error)
		Get_UserSubscriptionAdvancedByUuid(ctx context.Context, uuid string) (*models.UserSubscriptionAdvanced, error)
		Get_UserMe(ctx context.Context, uuid string) (*models.UserMe, error)

		Update_UserCoreRefreshTokenByUuid(ctx context.Context, refreshToken, uuid string) error
		Update_UserSubscriptionBeforPaySub(tx *sql.Tx, ctx context.Context, user *models.UserSubscription) error
		Update_UserSubscriptionBeforEndSub(ctx context.Context) error

		Delete_UserByUuid(ctx context.Context, uuid string) error
	}
	Subscriptions interface {
		Create_Subscription(ctx context.Context, sub *models.Subscription) error

		Get_SubscriptionsByUuid(ctx context.Context, uuid string, search string) (*[]models.Subscription, error)
		Get_SubscriptionById(ctx context.Context, id uint64, uuid string) (*models.Subscription, error)

		Update_SubscriptionsMounth(ctx context.Context) error
		Update_SubscriptionById(ctx context.Context, sub *models.Subscription, id int) error
		Delete_SubscriptionById(ctx context.Context, id int, uuid string) error
	}
	Transactions interface {
		Create_Transaction(ctx context.Context, transaction *models.Transaction) error

		Get_TransactionsByUuid(ctx context.Context, uuid string) (*[]models.Transaction, error)
		Get_TransactionPendingByUuid(ctx context.Context, uuid string) (*models.Transaction, error)
		Get_TransactionsByStatus(ctx context.Context, status string) (*[]models.Transaction, error)
		Get_TransactionsSubscriptionById(ctx context.Context, id uint64, uuid string) (*models.Transaction, error)
		Get_TransactionsSubscriptionByXToken(ctx context.Context, xtoken, uuid string) (*models.Transaction, error)

		Update_TransactionStatusById(tx *sql.Tx, ctx context.Context, status string, id uint64) error

		AutoActivateExpiredTransactions(ctx context.Context) error
	}
	Plans interface {
		Get_Plans(ctx context.Context) (*[]models.Plan, error)
		Get_PlanById(ctx context.Context, id uint64) (*models.Plan, error)
		Get_PlanByUserSubUuid(ctx context.Context, uuid string) (*models.Plan, error)
	}
}

func NewStorage(db *sql.DB, logger *logger.Logger) Storage {
	return Storage{
		Users:         &UserStore{db, logger},
		Subscriptions: &SubscriptionStore{db, logger},
		Transactions:  &TransactionStore{db, logger},
		Plans:         &PlanStore{db, logger},
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
