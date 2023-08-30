package event

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/logger"
	"go.uber.org/zap"
	"time"
)

type FetchHandler struct {
	ordersChannel chan string
	processor     service.NewOrderProcessor
	frequency     time.Duration
	workersCount  int
	log           logger.MyLogger
}

func NewFetchHandler(
	ordersChannel chan string, processor service.NewOrderProcessor, frequency time.Duration, workersCount int, log logger.MyLogger,
) *FetchHandler {
	return &FetchHandler{ordersChannel: ordersChannel, processor: processor, frequency: frequency, workersCount: workersCount, log: log}
}

func (f *FetchHandler) FetchOrderStatus(ctx context.Context) {
	ticker := time.NewTicker(f.frequency)
	var orders []string
	//Количество потоков запросов к сервису лояльности
	workers := make(chan struct{}, f.workersCount)
Loop:
	for {
		select {
		case <-ticker.C:
			for _, o := range orders {
				orderNumber := o
				//Запускаем получение данных из сервиса лояльности многопоточно
				go func() {
					//"Занимаем" или ожидаем один из потоков
					workers <- struct{}{}
					err := f.processor.ProcessNewOrder(ctx, orderNumber)
					if err != nil {
						f.log.L.Error("failed to get order info", zap.Error(err))
					}
					//"Освобождаем" поток
					<-workers
				}()
			}
			orders = nil
		case o, ok := <-f.ordersChannel:
			if !ok {
				break Loop
			}
			orders = append(orders, o)
		}
	}
}
