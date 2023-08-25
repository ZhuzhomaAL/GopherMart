package repository

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/user"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/gofrs/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user user.User) (user.User, error)
	GetByLogin(ctx context.Context, login string) (user.User, error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order order.Order) error
	GetByNumber(ctx context.Context, number string) (order.Order, error)
	GetAllByUser(ctx context.Context, userId uuid.UUID) ([]service.OrderInfo, error)
	UpdateOrder(ctx context.Context, order order.Order) error
	DeleteOrder(ctx context.Context, order order.Order) error
	ButchUpdateOrdersAndBalance(ctx context.Context, orders []order.Order, transactions []transaction.Transaction) error
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction transaction.Transaction) error
	GetBalanceByUser(ctx context.Context, userId uuid.UUID) (int, error)
	GetWithdrawSumByUser(ctx context.Context, userId uuid.UUID) (int, error)
	GetWithdrawsByUser(ctx context.Context, userId uuid.UUID) ([]transaction.Transaction, error)
}