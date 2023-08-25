package clients

import "fmt"

type NoOrderError struct {
	Order int
}

func (e NoOrderError) Error() string {
	return fmt.Sprintf("Order %d didn't exist in loyalty system", e.Order)
}

type TooManyRequests struct{}

func (TooManyRequests) Error() string {
	return fmt.Sprintf("Too many requests to loyalty service")
}

type LoyaltyServiceError struct {
	Msg string
}

func (e LoyaltyServiceError) Error() string {
	return fmt.Sprintf("Loyalty service error: %s", e.Msg)
}
