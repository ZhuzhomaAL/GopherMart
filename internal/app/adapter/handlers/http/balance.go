package http

import (
	"encoding/json"
	"errors"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/auth"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type BalanceHandler struct {
	bs  service.BalanceService
	log logger.MyLogger
}

func NewBalanceHandler(bs service.BalanceService, log logger.MyLogger) *BalanceHandler {
	return &BalanceHandler{bs: bs, log: log}
}

type withdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (b BalanceHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		b.log.L.Error("failed to get user")
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	current, err := b.bs.GetUserBalance(r.Context(), userID)
	if err != nil {
		b.log.L.Error("failed to get balance", zap.Error(err))
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	withdrawn, err := b.bs.GetUserWithdrawSum(r.Context(), userID)
	if err != nil {
		b.log.L.Error("failed to get withdrawals", zap.Error(err))
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	balance := service.Balance{
		Current:   current,
		Withdrawn: withdrawn,
	}

	resp, err := json.Marshal(balance)
	if err != nil {
		b.log.L.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		b.log.L.Error("failed to make response", zap.Error(err))
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
}

func (b BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	withdraw := withdrawRequest{}
	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		b.log.L.Error("failed to decode request", zap.Error(err))
		if err == io.EOF {
			http.Error(w, "request is empty, expected not empty", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	userID, ok := auth.GetUserID(r)
	if !ok {
		b.log.L.Error("failed to get user")
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	if err := b.bs.Withdraw(r.Context(), withdraw.Sum, withdraw.Order, userID); err != nil {
		b.log.L.Error("failed to process withdrawal", zap.Error(err))
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
}

func (b BalanceHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		b.log.L.Error("failed to get user")
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	withdrawals, err := b.bs.GetUserWithdraws(r.Context(), userID)
	if err != nil {
		b.log.L.Error("failed to get withdrawals", zap.Error(err))
		if _, ok := err.(*service.NoData); ok {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(withdrawals)
	if err != nil {
		b.log.L.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		b.log.L.Error("failed to make response", zap.Error(err))
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}
}
