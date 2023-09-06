package clients

import "fmt"

type NoOrderError struct {
	Order string
}

func (e NoOrderError) Error() string {
	return fmt.Sprintf("Order %s didn't exist in loyalty system", e.Order)
}

type TooManyRequests struct {
	Order      string
	RetryAfter int
}

func (e TooManyRequests) Error() string {
	return fmt.Sprintf("Too many requests, retry after %d seconds", e.RetryAfter)
}

type LoyaltyServiceError struct {
	OriginError error
}

func (e LoyaltyServiceError) Error() string {
	return fmt.Sprintf("Loyalty service error: %s", e.OriginError.Error())
}
