package loyal

import (
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type LoyaltyClient struct {
	baseURL string
}

const (
	getOrderInfoPath = "/api/orders"
)

func NewLoyaltyClient(baseURL string) *LoyaltyClient {
	return &LoyaltyClient{baseURL: baseURL}
}

func (lc LoyaltyClient) GetOrderProcessingInfo(order string) (clients.OrderLoyaltyInfo, error) {
	client := resty.New()
	orderInfo := new(clients.OrderLoyaltyInfo)
	resp, err := client.
		SetRetryCount(3).
		SetRetryWaitTime(60*time.Second).
		R().
		SetResult(orderInfo).
		AddRetryCondition(
			func(response *resty.Response, err error) bool {
				return response.StatusCode() == http.StatusTooManyRequests
			},
		).
		SetPathParams(
			map[string]string{
				"order": order,
			},
		).
		SetHeader("Accept", "application/json").
		Get(lc.baseURL + getOrderInfoPath + "/{order}")
	if err != nil {
		if responseError, ok := err.(*resty.ResponseError); ok {
			return *orderInfo, clients.LoyaltyServiceError{
				OriginError: responseError,
			}
		}
		return *orderInfo, err
	}

	if resp.StatusCode() == http.StatusNoContent {
		return *orderInfo, clients.NoOrderError{Order: order}
	}

	return *orderInfo, nil
}
