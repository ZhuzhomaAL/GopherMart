package postgres

import (
	"context"
	"database/sql"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/order"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/transaction"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/storage/postgres"
	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type OrderRepository struct {
	client *postgres.Client
}

func NewOrderRepository(client *postgres.Client) *OrderRepository {
	return &OrderRepository{client: client}
}

func (or OrderRepository) CreateOrder(ctx context.Context, order order.Order, tx bun.IDB) error {
	if tx == nil {
		tx = or.client
	}
	_, err := tx.NewInsert().Model(&order).Exec(ctx)
	return err
}

func (or OrderRepository) GetByNumber(ctx context.Context, number string, tx bun.IDB) (order.Order, error) {
	if tx == nil {
		tx = or.client
	}
	o := new(order.Order)
	err := tx.NewSelect().Model(o).Where("number = ?", number).Scan(ctx)
	return *o, err
}

func (or OrderRepository) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]service.OrderInfo, error) {
	orderInfos := make([]service.OrderInfo, 0)
	err := or.client.NewRaw(
		"SELECT o.number, o.status, t.sum accrual, o.uploaded_at FROM orders as o LEFT JOIN transactions t on t.order = o.number WHERE o."+
			"user_id = ?",
		userID.String(),
	).Scan(ctx, &orderInfos)
	if err != nil {
		return nil, err
	}
	return orderInfos, nil
}

func (or OrderRepository) UpdateOrder(ctx context.Context, order order.Order, tx bun.IDB) error {
	if tx == nil {
		tx = or.client
	}
	_, err := tx.NewUpdate().Model(&order).WherePK().Exec(ctx)
	return err
}

func (or OrderRepository) BatchUpdateOrdersAndBalance(
	ctx context.Context, orders []order.Order, transactions []transaction.Transaction,
) error {
	tx, err := or.client.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if _, err := tx.NewUpdate().Model(&orders).Column("status").Bulk().Exec(ctx); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	if _, err = or.client.NewInsert().Model(&transactions).Exec(ctx); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (or OrderRepository) GetAllByStatuses(ctx context.Context, statuses []string) ([]order.Order, error) {
	orders := make([]order.Order, 0)
	err := or.client.NewSelect().Model(&orders).
		Where("status IN (?)", bun.In(statuses)).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return orders, nil
		}
		return orders, err
	}

	return orders, nil
}

func (or OrderRepository) GetBatchByNumbers(ctx context.Context, orderNumbers []string) ([]order.Order, error) {
	orders := make([]order.Order, 0)
	err := or.client.NewSelect().Model(&orders).
		Where("number IN (?)", bun.In(orderNumbers)).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return orders, nil
		}
		return orders, err
	}

	return orders, nil
}
