package service

import (
	"context"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/repository/mocks"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	storagemocks "github.com/ZhuzhomaAL/GopherMart/internal/app/infra/storage/mocks"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"testing"
	"time"
)

func TestBalanceService_GetUserBalance(t *testing.T) {
	ctx := context.Background()
	userID, _ := uuid.NewV7()
	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		mockRes float64
		mockErr error
		wantErr bool
	}{
		{
			name: "Test_1.Баланс - целое число",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			mockRes: 100,
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "Test_2.Нулевой баланс",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			mockRes: 0,
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "Test_3.Отрицательный баланс",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			mockRes: -1,
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "Test_4.Баланс - дробное число",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			mockRes: 0.11111,
			mockErr: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rep := mocks.TransactionRepository{}
				txHelper := storagemocks.TransactionHelper{}
				bs := NewBalanceService(&rep, &txHelper)
				rep.On("GetBalanceByUser", tt.args.ctx, tt.args.userID, nil).Return(tt.mockRes, tt.mockErr)
				balance, err := bs.GetUserBalance(tt.args.ctx, tt.args.userID)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetUserBalance() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				require.Equal(t, tt.mockRes, balance)
			},
		)
	}
}

func TestBalanceService_GetUserWithdrawalSum(t *testing.T) {
	ctx := context.Background()
	userID, _ := uuid.NewV7()
	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		mockRes float64
	}{
		{
			name: "Test_1.Списание - целое число",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			mockRes: 1,
		},
		{
			name: "Test_2.Нулевое списание",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			mockRes: 0,
		},
		{
			name: "Test_3.Отрицательное списание",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			mockRes: -1,
		},
		{
			name: "Test_4.Списание - дробное число",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			mockRes: 0.1,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rep := mocks.TransactionRepository{}
				txHelper := storagemocks.TransactionHelper{}
				bs := NewBalanceService(&rep, &txHelper)
				rep.On("GetWithdrawalSumByUser", tt.args.ctx, tt.args.userID).Return(tt.mockRes, nil)
				withdrawal, err := bs.GetUserWithdrawalSum(tt.args.ctx, tt.args.userID)
				if err != nil {
					t.Errorf("GetUserWithdrawalSum() error = %v", err)
					return
				}
				require.Equal(t, tt.mockRes, withdrawal)
			},
		)
	}
}

func TestBalanceService_GetUserWithdraws(t *testing.T) {
	ctx := context.Background()
	userID, _ := uuid.NewV7()
	ID, _ := uuid.NewV7()
	orderNumber := goluhn.Generate(10)

	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name        string
		args        args
		transaction []transaction.Transaction
		wantErr     bool
		mockErr     error
	}{
		{
			name: "Test_1.Метод отрабатывает без ошибки. Списание",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			transaction: []transaction.Transaction{
				{
					ID:          userID,
					UserID:      ID,
					OrderNumber: orderNumber,
					Sum:         1,
					ProcessedAt: time.Time{},
					Type:        transaction.TypeWithdraw,
				},
			},
			wantErr: false,
		},
		{
			name: "Test_2.Метод отрабатывает без ошибки. Пополнение",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			transaction: []transaction.Transaction{
				{
					ID:          userID,
					UserID:      ID,
					OrderNumber: orderNumber,
					Sum:         1,
					ProcessedAt: time.Time{},
					Type:        transaction.TypeIncome,
				},
			},
			wantErr: false,
		},
		{
			name: "Test_3.Метод отрабатывает с ошибкой, нулевое списание",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			transaction: nil,
			wantErr:     true,
			mockErr:     &service.NoData{},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rep := mocks.TransactionRepository{}
				txHelper := storagemocks.TransactionHelper{}
				bs := NewBalanceService(&rep, &txHelper)
				rep.On("GetWithdrawalsByUser", tt.args.ctx, tt.args.userID).Return(tt.transaction, nil)
				withdrawal, err := bs.GetUserWithdraws(tt.args.ctx, tt.args.userID)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetUserWithdraws() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					require.Equal(t, tt.mockErr, err)
				}
				require.Equal(t, tt.transaction, withdrawal)
			},
		)
	}
}

func TestBalanceService_Withdraw(t *testing.T) {
	ctx := context.Background()
	userID, _ := uuid.NewV7()
	orderNumber := goluhn.Generate(10)
	type args struct {
		ctx         context.Context
		userID      uuid.UUID
		orderNumber string
		sum         float64
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		mockBalance float64
		mockErr     error
	}{
		{
			name: "Test_1.Есть остаток после списания",
			args: args{
				ctx:         ctx,
				userID:      userID,
				orderNumber: orderNumber,
				sum:         499,
			},
			mockBalance: 500,
			mockErr:     nil,
			wantErr:     false,
		},
		{
			name: "Test_2.Нулевой остаток после списания",
			args: args{
				ctx:         ctx,
				userID:      userID,
				orderNumber: orderNumber,
				sum:         500,
			},
			mockBalance: 500,
			mockErr:     nil,
			wantErr:     false,
		},
		{
			name: "Test_3.Метод возвращает ошибку. Баланс меньше списания",
			args: args{
				ctx:         ctx,
				userID:      userID,
				orderNumber: orderNumber,
				sum:         501,
			},
			mockBalance: 500,
			mockErr:     &transaction.NotEnoughMoney{},
			wantErr:     true,
		},
		{
			name: "Test_4.Метод возвращает ошибку. Невалидный формат номера заказа",
			args: args{
				ctx:         ctx,
				userID:      userID,
				orderNumber: "123",
				sum:         1,
			},
			mockBalance: 500,
			mockErr:     &order.InvalidFormat{OrderNumber: "123"},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rep := mocks.TransactionRepository{}
				txHelper := storagemocks.TransactionHelper{}
				bs := NewBalanceService(&rep, &txHelper)
				tx := storagemocks.Transaction{}
				txHelper.On("StartTransaction", tt.args.ctx).Return(&tx, nil)
				rep.On("GetBalanceByUser", tt.args.ctx, tt.args.userID, &bun.Tx{}).Return(tt.mockBalance, nil)
				rep.On("CreateTransaction", tt.args.ctx, mock.AnythingOfType("transaction.Transaction"), &bun.Tx{}).Return(nil)
				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)
				tx.On("GetTransaction").Return(&bun.Tx{})
				err := bs.Withdraw(tt.args.ctx, tt.args.sum, tt.args.orderNumber, tt.args.userID)
				if (err != nil) != tt.wantErr {
					t.Errorf("Withdraw() error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantErr {
					require.Equal(t, tt.mockErr, err)
				}
			},
		)
	}
}
