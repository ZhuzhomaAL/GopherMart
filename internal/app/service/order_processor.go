package service

import (
	"context"
	"errors"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
)

type OrderProcessor struct {
	ordersChannel     chan<- string
	orderInfosChannel chan<- clients.OrderLoyaltyInfo
	loyaltyClient     clients.LoyalClient
	os                service.OrderService
}

func NewOrderProcessor(
	ordersChannel chan<- string, orderInfosChannel chan<- clients.OrderLoyaltyInfo, loyaltyClient clients.LoyalClient, os service.OrderService,
) *OrderProcessor {
	return &OrderProcessor{ordersChannel: ordersChannel, orderInfosChannel: orderInfosChannel, loyaltyClient: loyaltyClient, os: os}
}

func (op OrderProcessor) ProcessNewOrder(ctx context.Context, number string) error {
	orderInfo, err := op.loyaltyClient.GetOrderProcessingInfo(number)
	if err != nil {
		if errors.Is(err, clients.NoOrderError{}) {
			_ = op.os.RemoveOrder(ctx, number)
			return err
		}
		op.ordersChannel <- number
		return err
	}

	if orderInfo.Status != clients.StatusProcessed && orderInfo.Status != clients.StatusInvalid {
		op.ordersChannel <- number
	}

	op.orderInfosChannel <- orderInfo

	return nil
}
