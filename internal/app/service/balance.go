package service

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/repository"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/storage/postgres"
	"github.com/gofrs/uuid"
	"math"
	"time"
)

type BalanceService struct {
	repo     repository.TransactionRepository
	txHelper *postgres.TransactionHelper
}

func NewBalanceService(repo repository.TransactionRepository, txHelper *postgres.TransactionHelper) *BalanceService {
	return &BalanceService{repo: repo, txHelper: txHelper}
}

func (bs BalanceService) GetUserBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	return bs.repo.GetBalanceByUser(ctx, userID, nil)
}

func (bs BalanceService) GetUserWithdrawalSum(ctx context.Context, userID uuid.UUID) (float64, error) {
	return bs.repo.GetWithdrawalSumByUser(ctx, userID)
}

func (bs BalanceService) GetUserWithdraws(ctx context.Context, userID uuid.UUID) ([]transaction.Transaction, error) {
	withdraws, err := bs.repo.GetWithdrawalsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(withdraws) == 0 {
		return nil, &service.NoData{}
	}
	for i := range withdraws {
		withdraws[i].Sum = math.Abs(withdraws[i].Sum)
	}

	return withdraws, nil
}

func (bs BalanceService) Withdraw(ctx context.Context, sum float64, orderNumber string, userID uuid.UUID) error {
	if !order.ValidateOrderFormat(orderNumber) {
		return &order.InvalidFormat{OrderNumber: orderNumber}
	}
	tx, err := bs.txHelper.GetTransaction(ctx)
	if err != nil {
		return err
	}
	balance, err := bs.repo.GetBalanceByUser(ctx, userID, tx)
	if err != nil {
		return err
	}
	if balance < sum {
		return &transaction.NotEnoughMoney{}
	}
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	if err := bs.repo.CreateTransaction(ctx, transaction.Transaction{
		ID:          id,
		UserID:      userID,
		OrderNumber: orderNumber,
		Sum:         -sum,
		ProcessedAt: time.Now(),
		Type:        transaction.TypeWithdraw,
	}, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
