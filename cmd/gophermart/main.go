package main

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/adapter/clients/loyal"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/adapter/handlers/event"
	httpHandlers "github.com/ZhuzhomaAL/GopherMart/internal/app/adapter/handlers/http"
	repo "github.com/ZhuzhomaAL/GopherMart/internal/app/adapter/repository/postgres"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/config"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/logger"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/storage/postgres"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/service"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func main() {
	conf := config.MakeConfig()
	mainContext := context.Background()
	l, err := logger.Initialize(conf.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	dbClient, err := postgres.NewPostgresConnection(mainContext, conf.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}
	ordersChannel := make(chan string, 1000)
	orderInfosChannel := make(chan clients.OrderLoyaltyInfo, 1000)

	orderRepo := repo.NewOrderRepository(dbClient)
	userRepo := repo.NewUserRepository(dbClient)
	transactionRepo := repo.NewTransactionRepository(dbClient)

	if err := dbClient.CreateTables(mainContext); err != nil {
		log.Fatal(err)
	}

	loyaltyClient := loyal.NewLoyaltyClient(conf.AccrualSystemAddress, l)

	balanceService := service.NewBalanceService(transactionRepo)
	orderService := service.NewOrderService(orderRepo, ordersChannel)
	userService := service.NewUserService(userRepo)

	orderProcessor := service.NewOrderProcessor(ordersChannel, orderInfosChannel, loyaltyClient, orderService)

	fetchHandler := event.NewFetchHandler(ordersChannel, orderProcessor, time.Second, 20, l)
	updateHandler := event.NewUpdateHandler(orderInfosChannel, orderService, time.Second, l)

	balanceHandler := httpHandlers.NewBalanceHandler(balanceService, l)
	orderHandler := httpHandlers.NewOrderHandler(orderService, l)
	userHandler := httpHandlers.NewUserHandler(userService, l)

	router := httpHandlers.GetRouter(userHandler, orderHandler, balanceHandler)
	event.Subscribe(mainContext, fetchHandler, updateHandler)

	l.L.Info("Running server", zap.String("address", conf.RunAddress))
	err = http.ListenAndServe(conf.RunAddress, router)
	if err != nil {
		log.Fatal(err)
	}
}
