package handlers

import (
	"context"
	"net/http"
)

// http

type (
	UserHandler interface {
		Register(w http.ResponseWriter, r *http.Request)
		Login(w http.ResponseWriter, r *http.Request)
	}
	OrderHandler interface {
		LoadOrder(w http.ResponseWriter, r *http.Request)
		GetUserOrders(w http.ResponseWriter, r *http.Request)
	}
	BalanceHandler interface {
		GetUserBalance(w http.ResponseWriter, r *http.Request)
		Withdraw(w http.ResponseWriter, r *http.Request)
		GetWithdrawals(w http.ResponseWriter, r *http.Request)
	}
)

// event

type (
	OrderFetchInfoHandler interface {
		FetchOrderStatus(ctx context.Context)
	}
	OrderUpdateHandler interface {
		UpdateStatusAndBalance(ctx context.Context)
	}
)
