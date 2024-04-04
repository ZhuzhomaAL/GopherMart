package storage

import (
	"context"
	"github.com/uptrace/bun"
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.1 --name=TransactionHelper
type TransactionHelper interface {
	StartTransaction(ctx context.Context) (Transaction, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.33.1 --name=Transaction
type Transaction interface {
	Commit() error
	Rollback() error
	GetTransaction() *bun.Tx
}
