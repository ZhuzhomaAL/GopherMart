package service

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/user"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/gofrs/uuid"
	"time"
)

type UserService interface {
	Register(ctx context.Context, login, password string) (user.User, error)
	Login(ctx context.Context, login, password string) (user.User, error)
}

type OrderService interface {
	LoadOrderByNumber(ctx context.Context, number string, userID uuid.UUID) error
	GetUserOrders(ctx context.Context, userID uuid.UUID) ([]OrderInfo, error)
	UpdateOrdersAndBalance(ctx context.Context, info []clients.OrderLoyaltyInfo) []error
	RemoveOrder(ctx context.Context, number string) error
}

type NewOrderProcessor interface {
	ProcessNewOrder(ctx context.Context, number string) error
}

type BalanceService interface {
	GetUserBalance(ctx context.Context, userID uuid.UUID) (int, error)
	GetUserWithdrawSum(ctx context.Context, userID uuid.UUID) (int, error)
	GetUserWithdraws(ctx context.Context, userID uuid.UUID) ([]transaction.Transaction, error)
	Withdraw(ctx context.Context, sum int, orderNumber string, userID uuid.UUID) error
}

type OrderInfo struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type Balance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}
