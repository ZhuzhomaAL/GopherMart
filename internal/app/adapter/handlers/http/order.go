package http

import (
	"encoding/json"
	"errors"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/auth"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type OrderHandler struct {
	os  service.OrderService
	log logger.MyLogger
}

func NewOrderHandler(os service.OrderService, log logger.MyLogger) *OrderHandler {
	return &OrderHandler{os: os, log: log}
}

func (oh OrderHandler) LoadOrder(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		oh.log.L.Error("empty request")
		http.Error(w, "request is empty, expected not empty", http.StatusBadRequest)
		return
	}

	request, err := io.ReadAll(r.Body)
	if err != nil {
		oh.log.L.Error("failed to decode request", zap.Error(err))
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	userID, ok := auth.GetUserID(r)
	if !ok {
		oh.log.L.Error("failed to get user")
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	if err := oh.os.LoadOrderByNumber(r.Context(), string(request), userID); err != nil {
		oh.log.L.Error("failed to load order", zap.Error(err))
		var errAlreadyLoaded *order.AlreadyLoaded
		if errors.As(err, &errAlreadyLoaded) {
			if errAlreadyLoaded.UserID == userID {
				w.WriteHeader(http.StatusOK)
				return
			}
			http.Error(w, "already loaded by other user", http.StatusConflict)
			return
		}
		var errInvalidFormat *order.InvalidFormat
		if errors.As(err, &errInvalidFormat) {
			http.Error(w, "invalid format of order number", http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (oh OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		oh.log.L.Error("failed to get user")
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	orders, err := oh.os.GetUserOrders(r.Context(), userID)
	if err != nil {
		oh.log.L.Error("failed to get user orders", zap.Error(err))
		if _, ok := err.(*service.NoData); ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(orders)
	if err != nil {
		oh.log.L.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		oh.log.L.Error("failed to make response", zap.Error(err))
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
}
