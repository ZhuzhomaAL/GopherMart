package event

import (
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/handlers"
)

func Subscribe(fetchHandler handlers.OrderFetchInfoHandler, updateHandler handlers.OrderUpdateHandler) {
	go fetchHandler.FetchOrderStatus()
	go updateHandler.UpdateStatusAndBalance()
}
