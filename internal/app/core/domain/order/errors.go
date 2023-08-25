package order

import (
	"fmt"
	"github.com/gofrs/uuid"
)

type InvalidStatus struct {
	OrderNumber string
	Status      string
}

func (e InvalidStatus) Error() string {
	return fmt.Sprintf("Invalid status: %s for order number %d", e.Status, e.OrderNumber)
}

type NoSuchOrder struct {
	OrderNumber string
}

func (e NoSuchOrder) Error() string {
	return fmt.Sprintf("Order number %d not found", e.OrderNumber)
}

type AlreadyLoaded struct {
	OrderNumber string
	UserId      uuid.UUID
}

func (e AlreadyLoaded) Error() string {
	return fmt.Sprintf("Order number %d already loaded", e.OrderNumber)
}

type InvalidFormat struct {
	OrderNumber string
}

func (e InvalidFormat) Error() string {
	return fmt.Sprintf("Order number %d has invalid format", e.OrderNumber)
}
