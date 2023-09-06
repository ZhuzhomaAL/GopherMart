package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/user"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/storage"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type (
	Client struct {
		*OrigClient
	}

	OrigClient = bun.DB
)

func NewPostgresConnection(ctx context.Context, dsn string) (*Client, error) {
	var sqlDB = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	var bunDB = bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	if err := bunDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping connection: %w", err)
	}

	return &Client{bunDB}, nil
}

func (c Client) CreateTables(ctx context.Context) error {
	u := new(user.User)
	o := new(order.Order)
	t := new(transaction.Transaction)
	if _, err := c.NewCreateTable().Model(u).IfNotExists().Exec(ctx); err != nil {
		return err
	}
	if _, err := c.NewCreateTable().Model(o).IfNotExists().Exec(ctx); err != nil {
		return err
	}
	if _, err := c.NewCreateTable().Model(t).IfNotExists().Exec(ctx); err != nil {
		return err
	}
	return nil
}

//go:generate go run github.com/vektra/mockery/v2@v2.33.1 --name=TransactionHelper
type TransactionHelper struct {
	db *Client
}

func NewTransactionHelper(db *Client) *TransactionHelper {
	return &TransactionHelper{db: db}
}

func (t TransactionHelper) StartTransaction(ctx context.Context) (storage.Transaction, error) {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	return PgTransaction{&tx}, nil
}

type PgTransaction struct {
	tx *bun.Tx
}

func (p PgTransaction) Commit() error {
	return p.tx.Commit()
}

func (p PgTransaction) Rollback() error {
	return p.tx.Rollback()
}

func (p PgTransaction) GetTransaction() *bun.Tx {
	return p.tx
}
