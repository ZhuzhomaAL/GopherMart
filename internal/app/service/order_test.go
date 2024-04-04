package service

import (
	"context"
	"errors"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/repository/mocks"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	storagemocks "github.com/ZhuzhomaAL/GopherMart/internal/app/infra/storage/mocks"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"reflect"
	"testing"
	"time"
)

func TestOrderService_GetUserOrders(t *testing.T) {
	orderNumber := goluhn.Generate(10)
	ctx := context.Background()
	userID, _ := uuid.NewV7()
	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantedErr error
		wantedRes []service.OrderInfo
		mockRes   []service.OrderInfo
		mockErr   error
	}{
		{
			name: "Test_1. Заказы существуют",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			mockRes: []service.OrderInfo{
				{
					Number:     orderNumber,
					Status:     order.StatusNew,
					Accrual:    100,
					UploadedAt: time.Time{},
				},
			},
			wantedRes: []service.OrderInfo{
				{
					Number:     orderNumber,
					Status:     order.StatusNew,
					Accrual:    100,
					UploadedAt: time.Time{},
				},
			},
		},
		{
			name: "Test_2. Нет заказов",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			wantErr:   true,
			mockRes:   []service.OrderInfo{},
			wantedErr: &service.NoData{},
			wantedRes: nil,
		},
		{
			name: "Test_3. Ошибка репозитория",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			wantErr:   true,
			mockRes:   []service.OrderInfo{},
			mockErr:   errors.New("db gone away"),
			wantedErr: errors.New("db gone away"),
			wantedRes: nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rep := mocks.OrderRepository{}
				txHelper := storagemocks.TransactionHelper{}
				os := NewOrderService(&rep, &txHelper)

				rep.On("GetAllByUser", tt.args.ctx, tt.args.userID).Return(tt.mockRes, tt.mockErr)

				expOrder, err := os.GetUserOrders(tt.args.ctx, tt.args.userID)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetUserOrders() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					require.Equal(t, tt.wantedErr, err)
				}
				require.Equal(t, tt.wantedRes, expOrder)
			},
		)
	}
}

func TestOrderService_GetUnprocessedOrders(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	notFinalStatuses := []string{order.StatusNew, order.StatusProcessing}
	ctx := context.Background()
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantedErr error
		mockRes   []order.Order
		mockErr   error
		wantedRes []order.Order
	}{
		{
			name: "Test_1. Заказы есть",
			args: args{
				ctx: ctx,
			},
			mockRes: []order.Order{
				{},
			},
			wantedRes: []order.Order{
				{},
			},
		},
		{
			name: "Test_2. Заказов нет",
			args: args{
				ctx: ctx,
			},
			mockRes:   []order.Order{},
			wantedRes: []order.Order{},
		},
		{
			name: "Test_3. Ошибка репозитория",
			args: args{
				ctx: ctx,
			},
			wantErr:   true,
			mockRes:   nil,
			mockErr:   errors.New("db gone away"),
			wantedRes: nil,
			wantedErr: errors.New("db gone away"),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rep := mocks.OrderRepository{}
				txHelper := storagemocks.TransactionHelper{}
				os := NewOrderService(&rep, &txHelper)
				rep.On("GetAllByStatuses", tt.args.ctx, notFinalStatuses).Return(tt.mockRes, tt.mockErr)
				orders, err := os.GetUnprocessedOrders(tt.args.ctx)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetUnprocessedOrders() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(orders, tt.wantedRes) {
					t.Errorf("GetUnprocessedOrders() got = %v, want %v", orders, tt.wantedRes)
				}
			},
		)
	}
}

func TestOrderService_InvalidateOrder(t *testing.T) {
	type args struct {
		ctx    context.Context
		number string
	}
	orderNumber := goluhn.Generate(10)
	ctx := context.Background()
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		wantedErr       error
		mockRes         order.Order
		mockGetOrderErr error
		mockUpdateErr   error
	}{
		{
			name: "Test_1. Успешная инвалидация",
			args: args{
				ctx:    ctx,
				number: orderNumber,
			},
			mockRes: order.Order{
				Number: orderNumber,
				Status: order.StatusNew,
			},
			mockGetOrderErr: nil,
			wantedErr:       nil,
		},
		{
			name: "Test_2. Ошибка получения заказа",
			args: args{
				ctx:    ctx,
				number: orderNumber,
			},
			wantErr:         true,
			mockRes:         order.Order{},
			mockGetOrderErr: errors.New("no order"),
			wantedErr:       errors.New("no order"),
		},
		{
			name: "Test_3. Ошибка обновления",
			args: args{
				ctx:    ctx,
				number: orderNumber,
			},
			wantErr: true,
			mockRes: order.Order{
				Number: orderNumber,
				Status: order.StatusNew,
			},
			mockGetOrderErr: nil,
			mockUpdateErr:   errors.New("can not update"),
			wantedErr:       errors.New("can not update"),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rep := mocks.OrderRepository{}
				txHelper := storagemocks.TransactionHelper{}
				os := NewOrderService(&rep, &txHelper)
				tx := storagemocks.Transaction{}
				txHelper.On("StartTransaction", tt.args.ctx).Return(&tx, nil)
				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)
				tx.On("GetTransaction").Return(&bun.Tx{})
				rep.On("GetByNumber", tt.args.ctx, tt.args.number, &bun.Tx{}).Return(tt.mockRes, tt.mockGetOrderErr)
				rep.On(
					"UpdateOrder", tt.args.ctx, order.Order{Number: tt.args.number, Status: order.StatusInvalid}, &bun.Tx{},
				).Return(tt.mockUpdateErr)
				err := os.InvalidateOrder(tt.args.ctx, tt.args.number)
				if (err != nil) != tt.wantErr {
					t.Errorf("InvalidateOrder() error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantErr {
					require.Equal(t, tt.wantedErr, err)
				}
			},
		)
	}
}

func TestOrderService_LoadOrderByNumber(t *testing.T) {
	type args struct {
		ctx    context.Context
		number string
	}
	orderNumber := goluhn.Generate(10)
	userID, _ := uuid.NewV7()
	ctx := context.Background()
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		wantedErr       error
		mockRes         order.Order
		mockGetOrderErr error
		mockCreateErr   error
	}{
		{
			name: "Test_1. Успешное создание",
			args: args{
				ctx:    ctx,
				number: orderNumber,
			},
			mockRes:         order.Order{},
			mockGetOrderErr: errors.New("no order"),
			wantedErr:       nil,
		},
		{
			name: "Test_2. Невалидный номер",
			args: args{
				ctx:    ctx,
				number: "1245",
			},
			mockRes:   order.Order{},
			wantErr:   true,
			wantedErr: &order.InvalidFormat{OrderNumber: "1245"},
		},
		{
			name: "Test_3. Заказ уже есть",
			args: args{
				ctx:    ctx,
				number: orderNumber,
			},
			mockRes: order.Order{
				UserID: uuid.UUID{},
				Number: orderNumber,
			},
			wantErr: true,
			wantedErr: &order.AlreadyLoaded{
				OrderNumber: orderNumber,
				UserID:      uuid.UUID{},
			},
		},
		{
			name: "Test_4. Ошибка создания",
			args: args{
				ctx:    ctx,
				number: orderNumber,
			},
			mockRes:         order.Order{},
			mockGetOrderErr: errors.New("no order"),
			mockCreateErr:   errors.New("can not create order"),
			wantErr:         true,
			wantedErr:       errors.New("can not create order"),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rep := mocks.OrderRepository{}
				txHelper := storagemocks.TransactionHelper{}
				os := NewOrderService(&rep, &txHelper)
				tx := storagemocks.Transaction{}
				txHelper.On("StartTransaction", tt.args.ctx).Return(&tx, nil)
				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)
				tx.On("GetTransaction").Return(&bun.Tx{})
				rep.On("GetByNumber", tt.args.ctx, tt.args.number, &bun.Tx{}).Return(tt.mockRes, tt.mockGetOrderErr)
				rep.On("CreateOrder", tt.args.ctx, mock.AnythingOfType("order.Order"), &bun.Tx{}).Return(tt.mockCreateErr)
				err := os.LoadOrderByNumber(tt.args.ctx, tt.args.number, userID)
				if (err != nil) != tt.wantErr {
					t.Errorf("LoadOrderByNumber() error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantErr {
					require.Equal(t, tt.wantedErr, err)
				}
			},
		)
	}
}

func TestOrderService_UpdateOrdersAndBalance(t *testing.T) {
	type args struct {
		ctx  context.Context
		info map[string]clients.OrderLoyaltyInfo
	}
	orderNumber := goluhn.Generate(10)
	ctx := context.Background()
	tests := []struct {
		name                  string
		args                  args
		wantErr               bool
		wantedErr             []error
		mockGetButchOrders    []order.Order
		mockGetButchOrdersErr error
		mockUpdateErr         error
		orders                []order.Order
		transactions          []transaction.Transaction
	}{
		{
			name: "Test_1. Нормальное создание",
			args: args{
				ctx: ctx,
				info: map[string]clients.OrderLoyaltyInfo{
					orderNumber: {
						Order:   orderNumber,
						Status:  clients.StatusProcessed,
						Accrual: 100,
					},
				},
			},
			mockGetButchOrders: []order.Order{
				{
					Number: orderNumber,
					Status: order.StatusNew,
				},
			},
			orders: []order.Order{
				{
					Number: orderNumber,
					Status: order.StatusProcessed,
				},
			},
		},
		{
			name: "Test_2. Невалидный статус",
			args: args{
				ctx: ctx,
				info: map[string]clients.OrderLoyaltyInfo{
					orderNumber: {
						Order:   orderNumber,
						Status:  "Невалидный статус",
						Accrual: 100,
					},
				},
			},
			mockGetButchOrders: []order.Order{
				{
					Number: orderNumber,
					Status: order.StatusNew,
				},
			},
			orders: []order.Order{
				{
					Number: orderNumber,
					Status: order.StatusNew,
				},
			},
			wantErr: true,
			wantedErr: []error{
				order.InvalidStatus{
					OrderNumber: orderNumber,
					Status:      "Невалидный статус",
				},
			},
		},
		{
			name: "Test_3. Ошибка получения заказов",
			args: args{
				ctx: ctx,
				info: map[string]clients.OrderLoyaltyInfo{
					orderNumber: {
						Order:   orderNumber,
						Status:  clients.StatusProcessed,
						Accrual: 100,
					},
				},
			},
			mockGetButchOrders:    []order.Order{},
			mockGetButchOrdersErr: errors.New("db gone away"),
			orders:                []order.Order{},
			wantErr:               true,
			wantedErr:             []error{errors.New("db gone away")},
		},
		{
			name: "Test_4. Ошибка обновления",
			args: args{
				ctx: ctx,
				info: map[string]clients.OrderLoyaltyInfo{
					orderNumber: {
						Order:   orderNumber,
						Status:  clients.StatusProcessed,
						Accrual: 100,
					},
				},
			},
			mockGetButchOrders: []order.Order{
				{
					Number: orderNumber,
					Status: order.StatusNew,
				},
			},
			orders: []order.Order{
				{
					Number: orderNumber,
					Status: order.StatusProcessed,
				},
			},
			mockUpdateErr: errors.New("can not update"),
			wantErr:       true,
			wantedErr:     []error{errors.New("can not update")},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rep := mocks.OrderRepository{}
				txHelper := storagemocks.TransactionHelper{}
				os := NewOrderService(&rep, &txHelper)
				orderNumbers := make([]string, len(tt.args.info))
				n := 0
				for _, i := range tt.args.info {
					orderNumbers[n] = i.Order
					n++
				}
				rep.On("GetBatchByNumbers", tt.args.ctx, orderNumbers).Return(tt.mockGetButchOrders, tt.mockGetButchOrdersErr)
				rep.On(
					"BatchUpdateOrdersAndBalance", tt.args.ctx, tt.orders, mock.AnythingOfType("[]transaction.Transaction"),
				).Return(tt.mockUpdateErr)
				got := os.UpdateOrdersAndBalance(tt.args.ctx, tt.args.info)
				require.Equal(t, got, tt.wantedErr)
			},
		)
	}
}
