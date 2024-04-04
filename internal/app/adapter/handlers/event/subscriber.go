package event

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/handlers"
)

func Subscribe(ctx context.Context, fetchHandler handlers.OrderFetchInfoHandler, updateHandler handlers.OrderUpdateHandler) {
	go fetchHandler.FetchOrderStatus(ctx)
	go updateHandler.UpdateStatusAndBalance(ctx)
}
