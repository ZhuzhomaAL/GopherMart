package postgres

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/storage/postgres"
	"github.com/gofrs/uuid"
)

type TransactionRepository struct {
	client *postgres.Client
}

func NewTransactionRepository(client *postgres.Client) *TransactionRepository {
	return &TransactionRepository{client: client}
}

func (tr TransactionRepository) CreateTransaction(ctx context.Context, transaction transaction.Transaction) error {
	_, err := tr.client.NewInsert().Model(&transaction).Exec(ctx)
	return err
}

func (tr TransactionRepository) GetBalanceByUser(ctx context.Context, userId uuid.UUID) (int, error) {
	var balance int
	err := tr.client.NewRaw(
		"SELECT SUM(sum) FROM transactions WHERE user_id = ? GROUP BY user_id",
		userId.String(),
	).Scan(ctx, &balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (tr TransactionRepository) GetWithdrawSumByUser(ctx context.Context, userId uuid.UUID) (int, error) {
	var sum int
	err := tr.client.NewRaw(
		"SELECT SUM(sum) FROM transactions WHERE user_id = ? AND type = ? GROUP BY user_id",
		userId.String(), transaction.TypeWithdraw,
	).Scan(ctx, &sum)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func (tr TransactionRepository) GetWithdrawsByUser(ctx context.Context, userId uuid.UUID) ([]transaction.Transaction, error) {
	transactions := make([]transaction.Transaction, 0)
	err := tr.client.NewSelect().Model(&transactions).
		Where("user_id = ?", userId.String()).
		Where("type = ?", transaction.TypeWithdraw).
		Scan(ctx)

	return transactions, err
}
