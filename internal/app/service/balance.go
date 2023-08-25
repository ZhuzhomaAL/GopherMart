package service

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/repository"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/gofrs/uuid"
	"time"
)

type BalanceService struct {
	repo repository.TransactionRepository
}

func NewBalanceService(repo repository.TransactionRepository) *BalanceService {
	return &BalanceService{repo: repo}
}

func (bs BalanceService) GetUserBalance(ctx context.Context, userId uuid.UUID) (int, error) {
	return bs.repo.GetBalanceByUser(ctx, userId)
}

func (bs BalanceService) GetUserWithdrawSum(ctx context.Context, userId uuid.UUID) (int, error) {
	return bs.repo.GetWithdrawSumByUser(ctx, userId)
}

func (bs BalanceService) GetUserWithdraws(ctx context.Context, userId uuid.UUID) ([]transaction.Transaction, error) {
	withdraws, err := bs.repo.GetWithdrawsByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	if len(withdraws) == 0 {
		return nil, &service.NoData{}
	}

	return withdraws, nil
}

func (bs BalanceService) Withdraw(ctx context.Context, sum int, orderNumber string, userId uuid.UUID) error {
	if !order.ValidateOrderFormat(orderNumber) {
		return &order.InvalidFormat{OrderNumber: orderNumber}
	}

	balance, err := bs.repo.GetBalanceByUser(ctx, userId)
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
	return bs.repo.CreateTransaction(
		ctx, transaction.Transaction{
			ID:          id,
			UserId:      userId,
			OderNumber:  orderNumber,
			Sum:         -sum,
			ProcessedAt: time.Now(),
			Type:        transaction.TypeWithdraw,
		},
	)
}
