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

type TransactionStore struct {
	db     *sql.DB
	logger *logger.Logger
}

func (s *TransactionStore) Create_Transaction(ctx context.Context, transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (plan_id, user_uuid, x_token, amount)
		VALUES ($1, $2, $3, $4)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, transaction.PlanID, transaction.UserUUID, transaction.XToken, transaction.Amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionStore) Get_TransactionsByUuid(ctx context.Context, uuid string) (*[]models.Transaction, error) {
	transactions := []models.Transaction{}

	query := `
		SELECT
		    id,
		    created_at,
		    user_uuid,
		    plan_id,
		    status,
		    x_token,
		   	amount,
		    currency
		FROM transactions
		WHERE user_uuid = $1
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.CreatedAt,
			&transaction.UserUUID,
			&transaction.PlanID,
			&transaction.Status,
			&transaction.XToken,
			&transaction.Amount,
			&transaction.Currency,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &transactions, nil
}

func (s *TransactionStore) Get_TransactionPendingByUuid(ctx context.Context, uuid string) (*models.Transaction, error) {
	transaction := models.Transaction{}

	query := `
		SELECT * FROM transactions WHERE status = 'pending' OR status = 'paid' AND user_uuid = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, uuid)

	if err := row.Scan(
		&transaction.ID,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.UserUUID,
		&transaction.PlanID,
		&transaction.Status,
		&transaction.XToken,
		&transaction.Amount,
		&transaction.Currency,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &transaction, nil
}

func (s *TransactionStore) Get_TransactionsByStatus(ctx context.Context, status string) (*[]models.Transaction, error) {
	transactions := []models.Transaction{}

	query := `
		SELECT * FROM transactions WHERE status = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.UserUUID,
			&transaction.PlanID,
			&transaction.Status,
			&transaction.XToken,
			&transaction.Amount,
			&transaction.Currency,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &transactions, nil
}

func (s *TransactionStore) Get_TransactionsSubscriptionById(ctx context.Context, id uint64, uuid string) (*models.Transaction, error) {
	transaction := models.Transaction{}

	query := `
		SELECT
		    id,
		    created_at,
		    plan_id,
		    status,
		    user_uuid,
		    x_token,
		    amount,
		    currency
		FROM transactions
		WHERE id = $1 AND user_uuid = $2
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, id, uuid)

	if err := row.Scan(
		&transaction.ID,
		&transaction.CreatedAt,
		&transaction.PlanID,
		&transaction.Status,
		&transaction.UserUUID,
		&transaction.XToken,
		&transaction.Amount,
		&transaction.Currency,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, err
	}

	return &transaction, nil
}

func (s *TransactionStore) Get_TransactionsSubscriptionByXToken(ctx context.Context, xtoken, uuid string) (*models.Transaction, error) {
	transaction := models.Transaction{}

	query := `
		SELECT
		    id,
		    created_at,
		    plan_id,
		    status,
		    user_uuid,
		    x_token,
		    amount,
		    currency
		FROM transactions
		WHERE x_token = $1 AND user_uuid = $2
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, xtoken, uuid)

	if err := row.Scan(
		&transaction.ID,
		&transaction.CreatedAt,
		&transaction.PlanID,
		&transaction.Status,
		&transaction.UserUUID,
		&transaction.XToken,
		&transaction.Amount,
		&transaction.Currency,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, err
	}

	return &transaction, nil
}

func (s *TransactionStore) Update_TransactionStatusById(tx *sql.Tx, ctx context.Context, status string, id uint64) error {
	query := `
		UPDATE transactions SET status = $1 WHERE id = $2 AND status = 'pending' OR status = 'paid'
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	upd, err := tx.ExecContext(ctx, query, status, id)
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

func (s *TransactionStore) AutoActivateExpiredTransactions(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query1 := `
        UPDATE transactions SET status = 'success', updated_at = NOW()
        WHERE status = 'pending' OR status = 'paid' AND created_at < NOW() - INTERVAL '1 hour'
        RETURNING user_uuid, plan_id
    `

	rows, err := tx.QueryContext(ctx, query1)
	if err != nil {
		return err
	}
	defer rows.Close()

	type Item struct {
		UserUUID string
		PlanID   int
	}
	var items []Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.UserUUID, &it.PlanID); err != nil {
			return err
		}
		items = append(items, it)
	}

	for _, it := range items {
		query2 := `
            UPDATE user_subscriptions
            SET
                plan_id = $1,
                start_date = NOW(),
                end_date = NOW() + INTERVAL '1 month',
                is_active = true
            WHERE user_uuid = $2
        `

		if _, err := tx.ExecContext(ctx, query2, it.PlanID, it.UserUUID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
