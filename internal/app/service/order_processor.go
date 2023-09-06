package service

import (
	"context"
	"errors"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
)

type OrderProcessor struct {
	orderInfosChannel chan<- clients.OrderLoyaltyInfo
	loyaltyClient     clients.LoyalClient
	os                service.OrderService
}

func NewOrderProcessor(orderInfosChannel chan<- clients.OrderLoyaltyInfo, loyaltyClient clients.LoyalClient, os service.OrderService) *OrderProcessor {
	return &OrderProcessor{orderInfosChannel: orderInfosChannel, loyaltyClient: loyaltyClient, os: os}
}

func (op OrderProcessor) ProcessNewOrder(ctx context.Context, number string) error {
	select {
	case <-ctx.Done():
		return errors.New("context canceled")
	default:
		orderInfo, err := op.loyaltyClient.GetOrderProcessingInfo(number)
		if err != nil {
			if errors.Is(err, clients.NoOrderError{}) {
				_ = op.os.InvalidateOrder(ctx, number)
				return err
			}
			return err
		}

		op.orderInfosChannel <- orderInfo

		return nil
	}
}
