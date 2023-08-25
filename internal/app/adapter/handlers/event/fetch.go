package event

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"time"
)

type FetchHandler struct {
	ochan     *chan string
	lf        service.NewOrderProcessor
	frequency time.Duration
}

func NewFetchHandler(ochan *chan string, lf service.NewOrderProcessor, frequency time.Duration) *FetchHandler {
	return &FetchHandler{ochan: ochan, lf: lf, frequency: frequency}
}

func (f *FetchHandler) FetchOrderStatus() {
	ticker := time.NewTicker(f.frequency)
	var orders []string
Loop:
	for {
		select {
		case o, ok := <-*f.ochan:
			if !ok {
				break Loop
			}
			orders = append(orders, o)
		case <-ticker.C:
			for _, o := range orders {
				err := f.lf.ProcessNewOrder(context.Background(), o)
				if err != nil {
					//TODO log
					//TODO wait
				}
			}
			orders = nil
		}
	}
}
