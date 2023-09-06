package postgres

import (
	"context"
	"database/sql"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/storage/postgres"
	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
	"math"
)

type TransactionRepository struct {
	client *postgres.Client
}

func NewTransactionRepository(client *postgres.Client) *TransactionRepository {
	return &TransactionRepository{client: client}
}

func (tr TransactionRepository) CreateTransaction(ctx context.Context, transaction transaction.Transaction, tx bun.IDB) error {
	if tx == nil {
		tx = tr.client
	}
	_, err := tx.NewInsert().Model(&transaction).Exec(ctx)
	return err
}

func (tr TransactionRepository) GetBalanceByUser(ctx context.Context, userID uuid.UUID, tx bun.IDB) (float64, error) {
	if tx == nil {
		tx = tr.client
	}
	var balance float64
	err := tx.NewRaw(
		"SELECT SUM(sum) FROM transactions WHERE user_id = ? GROUP BY user_id",
		userID.String(),
	).Scan(ctx, &balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return balance, nil
}

func (tr TransactionRepository) GetWithdrawalSumByUser(ctx context.Context, userID uuid.UUID) (float64, error) {
	var sum float64
	err := tr.client.NewRaw(
		"SELECT SUM(sum) FROM transactions WHERE user_id = ? AND type = ? GROUP BY user_id",
		userID.String(), transaction.TypeWithdraw,
	).Scan(ctx, &sum)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return math.Abs(sum), nil
}

func (tr TransactionRepository) GetWithdrawalsByUser(ctx context.Context, userID uuid.UUID) ([]transaction.Transaction, error) {
	transactions := make([]transaction.Transaction, 0)
	err := tr.client.NewSelect().Model(&transactions).
		Where("user_id = ?", userID.String()).
		Where("type = ?", transaction.TypeWithdraw).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return transactions, nil
		}
		return transactions, err
	}

	return transactions, nil
}
