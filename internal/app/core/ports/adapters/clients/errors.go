package clients

import "fmt"

type NoOrderError struct {
	Order string
}

func (e NoOrderError) Error() string {
	return fmt.Sprintf("Order %s didn't exist in loyalty system", e.Order)
}

type LoyaltyServiceError struct {
	OriginError error
}

func (e LoyaltyServiceError) Error() string {
	return fmt.Sprintf("Loyalty service error: %s", e.OriginError.Error())
}
