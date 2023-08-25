package http

import (
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/handlers"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func GetRouter(userHandler handlers.UserHandler, orderHandler handlers.OrderHandler, balanceHandler handlers.BalanceHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(100 * time.Second))

	r.Route(
		"/api/user", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
		},
	)
	r.Route(
		"/api/user/orders", func(r chi.Router) {
			r.Use(auth.Middleware)
			r.Post("/", orderHandler.LoadOrder)
			r.Get("/", orderHandler.GetUserOrders)
		},
	)
	r.Route(
		"/api/user/balance", func(r chi.Router) {
			r.Use(auth.Middleware)
			r.Get("/", balanceHandler.GetUserBalance)
			r.Post("/withdraw", balanceHandler.Withdraw)
		},
	)
	r.Route(
		"/api/user/withdrawals", func(r chi.Router) {
			r.Use(auth.Middleware)
			r.Get("/", balanceHandler.GetWithdrawals)
		},
	)

	return r
}
