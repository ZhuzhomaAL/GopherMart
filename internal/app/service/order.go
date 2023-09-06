package service

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/repository"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/storage"
	"github.com/gofrs/uuid"
	"time"
)

type OrderService struct {
	orderRepo repository.OrderRepository
	txHelper  storage.TransactionHelper
}

func NewOrderService(orderRepo repository.OrderRepository, txHelper storage.TransactionHelper) *OrderService {
	return &OrderService{orderRepo: orderRepo, txHelper: txHelper}
}

func (os OrderService) LoadOrderByNumber(ctx context.Context, number string, userID uuid.UUID) error {
	if !order.ValidateOrderFormat(number) {
		return &order.InvalidFormat{OrderNumber: number}
	}
	tx, err := os.txHelper.StartTransaction(ctx)
	if err != nil {
		return err
	}
	if o, err := os.orderRepo.GetByNumber(ctx, number, tx.GetTransaction()); err == nil {
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
		}, tx.GetTransaction(),
	); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (os OrderService) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]service.OrderInfo, error) {
	orders, err := os.orderRepo.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, &service.NoData{}
	}

	return orders, nil
}

func (os OrderService) UpdateOrdersAndBalance(ctx context.Context, info map[string]clients.OrderLoyaltyInfo) []error {
	orders, transactions, errors := os.makeOrdersAndTransactions(ctx, info)

	if err := os.orderRepo.BatchUpdateOrdersAndBalance(ctx, orders, transactions); err != nil {
		errors = append(errors, err)
	}

	return errors
}

func (os OrderService) InvalidateOrder(ctx context.Context, number string) error {
	tx, err := os.txHelper.StartTransaction(ctx)
	if err != nil {
		return err
	}
	o, err := os.orderRepo.GetByNumber(ctx, number, tx.GetTransaction())
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	o.Status = order.StatusInvalid
	if err := os.orderRepo.UpdateOrder(ctx, o, tx.GetTransaction()); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (os OrderService) makeOrdersAndTransactions(
	ctx context.Context, info map[string]clients.OrderLoyaltyInfo,
) ([]order.Order, []transaction.Transaction, []error) {
	var errors []error
	var transactions []transaction.Transaction
	orderNumbers := os.getOrderNumbersByOrderInfos(info)
	orders, err := os.orderRepo.GetBatchByNumbers(ctx, orderNumbers)
	if err != nil {
		errors = append(errors, err)
		return orders, transactions, errors
	}
	if len(orders) == 0 {
		return orders, transactions, errors
	}
	for n, o := range orders {
		i, ok := info[o.Number]
		if !ok {
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
		if o.Status == orderStatus {
			continue
		}
		orders[n].Status = orderStatus
		if i.Accrual > 0 {
			id, _ := uuid.NewV4()
			transactions = append(
				transactions, transaction.Transaction{
					ID:          id,
					UserID:      o.UserID,
					OrderNumber: o.Number,
					Sum:         i.Accrual,
					ProcessedAt: time.Now(),
					Type:        transaction.TypeIncome,
				},
			)
		}
	}
	return orders, transactions, errors
}

func (os OrderService) GetUnprocessedOrders(ctx context.Context) ([]order.Order, error) {
	notFinalStatuses := []string{order.StatusNew, order.StatusProcessing}
	return os.orderRepo.GetAllByStatuses(ctx, notFinalStatuses)
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

func (os OrderService) getOrderNumbersByOrderInfos(info map[string]clients.OrderLoyaltyInfo) []string {
	orderNumbers := make([]string, len(info))
	n := 0
	for _, i := range info {
		orderNumbers[n] = i.Order
		n++
	}

	return orderNumbers
}
