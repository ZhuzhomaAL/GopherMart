package main

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/adapter/clients/loyal"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/adapter/handlers/event"
	httpHandlers "github.com/ZhuzhomaAL/GopherMart/internal/app/adapter/handlers/http"
	repo "github.com/ZhuzhomaAL/GopherMart/internal/app/adapter/repository/postgres"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/config"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/storage/postgres"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/service"
	"log"
	"net/http"
	"time"
)

func main() {
	conf := config.MakeConfig()
	dbClient, err := postgres.NewPostgresConnection(context.Background(), conf.DatabaseUri)
	if err != nil {
		log.Fatal(err)
	}
	ordersChannel := make(chan string, 0)
	orderInfosChannel := make(chan clients.OrderLoyaltyInfo, 0)
	//var ordersChannel chan string
	//var orderInfosChannel chan clients.OrderLoyaltyInfo

	orderRepo := repo.NewOrderRepository(dbClient)
	userRepo := repo.NewUserRepository(dbClient)
	transactionRepo := repo.NewTransactionRepository(dbClient)

	dbClient.CreateTables(context.Background())

	loyaltyClient := loyal.NewLoyaltyClient(conf.AccrualSystemAddress)

	balanceService := service.NewBalanceService(transactionRepo)
	orderService := service.NewOrderService(orderRepo, &ordersChannel)
	userService := service.NewUserService(userRepo)

	orderProcessor := service.NewOrderProcessor(ordersChannel, orderInfosChannel, loyaltyClient, orderService)

	fetchHandler := event.NewFetchHandler(&ordersChannel, orderProcessor, time.Second*60)
	updateHandler := event.NewUpdateHandler(orderInfosChannel, orderService, time.Second*10)

	balanceHandler := httpHandlers.NewBalanceHandler(balanceService)
	orderHandler := httpHandlers.NewOrderHandler(orderService)
	userHandler := httpHandlers.NewUserHandler(userService)

	router := httpHandlers.GetRouter(userHandler, orderHandler, balanceHandler)
	event.Subscribe(fetchHandler, updateHandler)

	err = http.ListenAndServe(conf.RunAddress, router)
	if err != nil {
		log.Fatal(err)
	}
}
