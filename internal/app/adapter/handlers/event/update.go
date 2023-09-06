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

func NewUpdateHandler(
	orderInfosChannel <-chan clients.OrderLoyaltyInfo, os service.OrderService, frequency time.Duration, log logger.MyLogger,
) *UpdateHandler {
	return &UpdateHandler{orderInfosChannel: orderInfosChannel, os: os, frequency: frequency, log: log}
}

func (u UpdateHandler) UpdateStatusAndBalance(ctx context.Context) {
	ticker := time.NewTicker(u.frequency)
	infos := make(map[string]clients.OrderLoyaltyInfo)
	infos["4283279126516590"] = clients.OrderLoyaltyInfo{
		Order:   "4283279126516590",
		Status:  clients.StatusProcessed,
		Accrual: 100,
	}
	for {
		select {
		case <-ticker.C:
			errors := u.os.UpdateOrdersAndBalance(ctx, infos)
			if len(errors) > 0 {
				u.log.L.Error("failed to update orders", zap.Errors("err", errors))
			}
			infos = make(map[string]clients.OrderLoyaltyInfo)
		case info, ok := <-u.orderInfosChannel:
			if !ok {
				break
			}
			infos[info.Order] = info
		}
	}
}
