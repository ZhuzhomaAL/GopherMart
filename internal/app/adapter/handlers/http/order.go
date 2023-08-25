package http

import (
	"encoding/json"
	"errors"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/auth"
	"github.com/gofrs/uuid"
	"io"
	"net/http"
)

type OrderHandler struct {
	os service.OrderService
}

func NewOrderHandler(os service.OrderService) *OrderHandler {
	return &OrderHandler{os: os}
}

func (oh OrderHandler) LoadOrder(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		http.Error(w, "request is empty, expected not empty", http.StatusBadRequest)
		return
	}

	request, err := io.ReadAll(r.Body)
	if err != nil {
		//TODO log
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	userID, ok := r.Context().Value(auth.ContextUserID).(uuid.UUID)
	if !ok {
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	if err := oh.os.LoadOrderByNumber(r.Context(), string(request), userID); err != nil {
		var errAlreadyLoaded *order.AlreadyLoaded
		var errInvalidFormat *order.InvalidFormat
		if errors.As(err, &errAlreadyLoaded) {
			if errAlreadyLoaded.UserId == userID {
				w.WriteHeader(http.StatusOK)
				return
			}
			http.Error(w, "already loaded by other user", http.StatusConflict)
			return
		}
		if errors.As(err, &errInvalidFormat) {
			http.Error(w, "invalid format of order number", http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}

func (oh OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.ContextUserID).(uuid.UUID)
	if !ok {
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	orders, err := oh.os.GetUserOrders(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.NoData{}) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
}
