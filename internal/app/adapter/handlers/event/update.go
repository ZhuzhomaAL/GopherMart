package event

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/logger"
	"go.uber.org/zap"
	"time"
)

type UpdateHandler struct {
	orderInfosChannel <-chan clients.OrderLoyaltyInfo
	os                service.OrderService
	frequency         time.Duration
	log               logger.MyLogger
}

func NewUpdateHandler(orderInfosChannel <-chan clients.OrderLoyaltyInfo, os service.OrderService, frequency time.Duration, log logger.MyLogger) *UpdateHandler {
	return &UpdateHandler{orderInfosChannel: orderInfosChannel, os: os, frequency: frequency, log: log}
}

func (u UpdateHandler) UpdateStatusAndBalance(ctx context.Context) {
	ticker := time.NewTicker(u.frequency)
	var infos []clients.OrderLoyaltyInfo
Loop:
	for {
		select {
		case <-ticker.C:
			errors := u.os.UpdateOrdersAndBalance(ctx, infos)
			if len(errors) > 0 {
				u.log.L.Error("failed to update orders", zap.Errors("err", errors))
			}
			infos = nil
		case info, ok := <-u.orderInfosChannel:
			if !ok {
				break Loop
			}
			infos = append(infos, info)
		}
	}
}
