package loyal

import (
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type LoyaltyClient struct {
	baseURL string
	log     logger.MyLogger
}

func NewLoyaltyClient(baseURL string, log logger.MyLogger) *LoyaltyClient {
	return &LoyaltyClient{baseURL: baseURL, log: log}
}

const (
	getOrderInfoPath = "/api/orders"
)

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
	lc.log.L.Info("request", zap.String("URL", resp.Request.URL))
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
