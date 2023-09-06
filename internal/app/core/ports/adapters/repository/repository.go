package repository

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/user"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.1 --name=UserRepository
type UserRepository interface {
	CreateUser(ctx context.Context, user user.User) (user.User, error)
	GetByLogin(ctx context.Context, login string) (user.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.33.1 --name=OrderRepository
type OrderRepository interface {
	CreateOrder(ctx context.Context, order order.Order, tx bun.IDB) error
	GetByNumber(ctx context.Context, number string, tx bun.IDB) (order.Order, error)
	GetAllByUser(ctx context.Context, userID uuid.UUID) ([]service.OrderInfo, error)
	UpdateOrder(ctx context.Context, order order.Order, tx bun.IDB) error
	BatchUpdateOrdersAndBalance(ctx context.Context, orders []order.Order, transactions []transaction.Transaction) error
	GetAllByStatuses(ctx context.Context, statuses []string) ([]order.Order, error)
	GetBatchByNumbers(ctx context.Context, orderNumbers []string) ([]order.Order, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.33.1 --name=TransactionRepository
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction transaction.Transaction, tx bun.IDB) error
	GetBalanceByUser(ctx context.Context, userID uuid.UUID, tx bun.IDB) (float64, error)
	GetWithdrawalSumByUser(ctx context.Context, userID uuid.UUID) (float64, error)
	GetWithdrawalsByUser(ctx context.Context, userID uuid.UUID) ([]transaction.Transaction, error)
}
