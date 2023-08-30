package service

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/repository"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/gofrs/uuid"
	"time"
)

type OrderService struct {
	orderRepo     repository.OrderRepository
	ordersChannel chan string
}

func NewOrderService(orderRepo repository.OrderRepository, ordersChannel chan string) *OrderService {
	return &OrderService{orderRepo: orderRepo, ordersChannel: ordersChannel}
}

func (os OrderService) LoadOrderByNumber(ctx context.Context, number string, userID uuid.UUID) error {
	if !order.ValidateOrderFormat(number) {
		return &order.InvalidFormat{OrderNumber: number}
	}
	if o, err := os.orderRepo.GetByNumber(ctx, number); err == nil {
		return &order.AlreadyLoaded{
			OrderNumber: o.Number,
			UserID:      o.UserID,
		}
	}
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	if err := os.orderRepo.CreateOrder(
		ctx, order.Order{
			ID:         id,
			UserID:     userID,
			Number:     number,
			Status:     order.StatusNew,
			UploadedAt: time.Now(),
		},
	); err != nil {
		return err
	}

	//Пишем номер в канал в новом потоке, чтобы не блокировать хэндлер
	go func() {
		os.ordersChannel <- number
	}()
	return nil
}

func (os OrderService) GetUserOrders(ctx context.Context, userId uuid.UUID) ([]service.OrderInfo, error) {
	orders, err := os.orderRepo.GetAllByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, &service.NoData{}
	}

	return orders, nil
}

func (os OrderService) UpdateOrdersAndBalance(ctx context.Context, info []clients.OrderLoyaltyInfo) []error {
	var orders []order.Order
	var transactions []transaction.Transaction
	var errors []error
	errors = append(errors, os.fillInfo(ctx, info, &orders, &transactions)...)

	if err := os.orderRepo.BatchUpdateOrdersAndBalance(ctx, orders, transactions); err != nil {
		errors = append(errors, err)
	}

	return errors
}

func (os OrderService) RemoveOrder(ctx context.Context, number string) error {
	o, err := os.orderRepo.GetByNumber(ctx, number)
	if err != nil {
		return err
	}
	return os.orderRepo.DeleteOrder(ctx, o)
}

func (os OrderService) fillInfo(
	ctx context.Context, info []clients.OrderLoyaltyInfo, orders *[]order.Order, transactions *[]transaction.Transaction,
) []error {
	var errors []error
	for _, i := range info {
		o, err := os.orderRepo.GetByNumber(ctx, i.Order)
		if err != nil {
			errors = append(
				errors, order.NoSuchOrder{
					OrderNumber: i.Order,
				},
			)
			continue
		}
		orderStatus, ok := os.getOrderStatusFromLoyalty(i.Status)
		if !ok {
			errors = append(
				errors, order.InvalidStatus{
					OrderNumber: i.Order,
					Status:      i.Status,
				},
			)
			continue
		}
		o.Status = orderStatus
		*orders = append(*orders, o)
		if i.Accrual > 0 {
			id, _ := uuid.NewV4()
			*transactions = append(
				*transactions, transaction.Transaction{
					ID:          id,
					UserID:      o.UserID,
					OderNumber:  o.Number,
					Sum:         i.Accrual,
					ProcessedAt: time.Now(),
					Type:        transaction.TypeIncome,
				},
			)
		}
	}
	return errors
}

func (os OrderService) getOrderStatusFromLoyalty(loyaltyStatus string) (string, bool) {
	statusMap := map[string]string{
		clients.StatusRegistered: order.StatusNew,
		clients.StatusProcessing: order.StatusProcessing,
		clients.StatusInvalid:    order.StatusInvalid,
		clients.StatusProcessed:  order.StatusProcessed,
	}

	status, ok := statusMap[loyaltyStatus]
	return status, ok
}
