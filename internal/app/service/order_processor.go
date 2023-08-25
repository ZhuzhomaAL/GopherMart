package service

import (
	"context"
	"errors"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
)

type OrderProcessor struct {
	fetch         chan<- string
	updatech      chan<- clients.OrderLoyaltyInfo
	loyaltyClient clients.LoyalClient
	os            service.OrderService
}

func NewOrderProcessor(
	fetch chan<- string, updatech chan<- clients.OrderLoyaltyInfo, loyaltyClient clients.LoyalClient, os service.OrderService,
) *OrderProcessor {
	return &OrderProcessor{fetch: fetch, updatech: updatech, loyaltyClient: loyaltyClient, os: os}
}

func (op OrderProcessor) ProcessNewOrder(ctx context.Context, number string) error {
	orderInfo, err := op.loyaltyClient.GetOrderProcessingInfo(number)
	if err != nil {
		if errors.Is(err, clients.NoOrderError{}) {
			_ = op.os.RemoveOrder(ctx, number)
			return err
		}
		op.fetch <- number
		return err
	}

	if orderInfo.Status != clients.StatusProcessed && orderInfo.Status != clients.StatusInvalid {
		op.fetch <- number
	}

	op.updatech <- orderInfo

	return nil
}
