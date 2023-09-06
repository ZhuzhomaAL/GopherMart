package event

import (
	"context"
	"errors"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/logger"
	"go.uber.org/zap"
	"time"
)

type FetchHandler struct {
	processor    service.NewOrderProcessor
	orderService service.OrderService
	frequency    time.Duration
	workersCount int
	log          logger.MyLogger
}

func NewFetchHandler(
	processor service.NewOrderProcessor, orderService service.OrderService, frequency time.Duration, workersCount int, log logger.MyLogger,
) *FetchHandler {
	return &FetchHandler{processor: processor, orderService: orderService, frequency: frequency, workersCount: workersCount, log: log}
}

func (f *FetchHandler) FetchOrderStatus(ctx context.Context) {
	ticker := time.NewTicker(f.frequency)
	//Количество потоков запросов к сервису лояльности
	workers := make(chan struct{}, f.workersCount)
	sleepSignal := make(chan int)
	requestCtx, cancelRequests := context.WithCancel(ctx)
	for {
		select {
		case <-ticker.C:
			orders, err := f.orderService.GetUnprocessedOrders(ctx)
			if err != nil {
				f.log.L.Error("failed to get unprocessed orders", zap.Error(err))
			}
			if len(orders) == 0 {
				continue
			}
			for _, o := range orders {
				select {
				case sleepTime := <-sleepSignal:
					time.Sleep(time.Duration(sleepTime) * time.Second)
				default:
					orderNumber := o.Number
					//Запускаем получение данных из сервиса лояльности многопоточно
					//"Занимаем" или ожидаем один из потоков
					workers <- struct{}{}
					go func() {
						err := f.processor.ProcessNewOrder(requestCtx, orderNumber)
						if err != nil {
							var tooManyRequests *clients.TooManyRequests
							if errors.As(err, &tooManyRequests) {
								cancelRequests()
								sleepSignal <- tooManyRequests.RetryAfter
							}
							f.log.L.Error("failed to get order info", zap.Error(err))
						}
						//"Освобождаем" поток
						<-workers
					}()
				}
			}
		}
	}
}
