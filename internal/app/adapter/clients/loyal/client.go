package loyal

import "github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"

type LoyaltyClient struct {
	baseUrl string
}

func NewLoyaltyClient(baseUrl string) *LoyaltyClient {
	return &LoyaltyClient{baseUrl: baseUrl}
}

func (LoyaltyClient) GetOrderProcessingInfo(order string) (clients.OrderLoyaltyInfo, error) {
	//TODO implement me
	panic("implement me")
}
