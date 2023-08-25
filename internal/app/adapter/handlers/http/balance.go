package http

import (
	"encoding/json"
	"errors"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/auth"
	"github.com/gofrs/uuid"
	"io"
	"net/http"
)

type BalanceHandler struct {
	bs service.BalanceService
}

func NewBalanceHandler(bs service.BalanceService) *BalanceHandler {
	return &BalanceHandler{bs: bs}
}

type withdrawRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

func (b BalanceHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.ContextUserID).(uuid.UUID)
	if !ok {
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	current, err := b.bs.GetUserBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	withdrawn, err := b.bs.GetUserWithdrawSum(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	balance := service.Balance{
		Current:   current,
		Withdrawn: withdrawn,
	}

	resp, err := json.Marshal(balance)
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

func (b BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	withdraw := withdrawRequest{}
	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		if err == io.EOF {
			http.Error(w, "request is empty, expected not empty", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	userID, ok := r.Context().Value(auth.ContextUserID).(uuid.UUID)
	if !ok {
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	if err := b.bs.Withdraw(r.Context(), withdraw.Sum, withdraw.Order, userID); err != nil {
		if errors.Is(err, order.InvalidFormat{}) {
			http.Error(w, "Invalid order format", http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, transaction.NotEnoughMoney{}) {
			http.Error(w, "Not enough money", http.StatusPaymentRequired)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (b BalanceHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.ContextUserID).(uuid.UUID)
	if !ok {
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	withdrawals, err := b.bs.GetUserWithdraws(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.NoData{}) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(withdrawals)
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