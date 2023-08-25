package event

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"time"
)

type UpdateHandler struct {
	infoChan  <-chan clients.OrderLoyaltyInfo
	os        service.OrderService
	frequency time.Duration
}

func NewUpdateHandler(infoChan <-chan clients.OrderLoyaltyInfo, os service.OrderService, frequency time.Duration) *UpdateHandler {
	return &UpdateHandler{infoChan: infoChan, os: os, frequency: frequency}
}

func (u UpdateHandler) UpdateStatusAndBalance() {
	ticker := time.NewTicker(u.frequency)
	var infos []clients.OrderLoyaltyInfo
Loop:
	for {
		select {
		case info, ok := <-u.infoChan:
			if !ok {
				break Loop
			}
			infos = append(infos, info)
		case <-ticker.C:
			errors := u.os.UpdateOrdersAndBalance(context.Background(), infos)
			if len(errors) > 0 {
				//TODO log
			}
			infos = nil
		}
	}
}
