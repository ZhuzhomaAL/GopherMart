package loyal

import (
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/clients"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http"
	"strconv"
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
	var orderInfo clients.OrderLoyaltyInfo
	resp, err := client.
		R().
		SetResult(&orderInfo).
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
			return orderInfo, clients.LoyaltyServiceError{
				OriginError: responseError,
			}
		}
		return orderInfo, err
	}

	if resp.StatusCode() == http.StatusNoContent {
		return orderInfo, clients.NoOrderError{Order: order}
	}

	if resp.StatusCode() == http.StatusTooManyRequests {
		retryAfter, err := strconv.Atoi(resp.Header().Get("Retry-After"))
		if err != nil {
			return orderInfo, err
		}
		return orderInfo, clients.TooManyRequests{
			Order:      order,
			RetryAfter: retryAfter,
		}
	}

	return orderInfo, nil
}
