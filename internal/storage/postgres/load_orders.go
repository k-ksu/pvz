package postgres

import (
	"HomeWork_1/internal/model/errs"
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"

	"HomeWork_1/internal/model"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

func (r *Repository) LoadOrders(ctx context.Context, hash string) ([]model.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.LoadOrders")
	defer span.Finish()

	var orders []model.Order
	query := `SELECT id, client_id, condition, arrived_at, received_at, package, price, max_weight FROM orders;`
	err := r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		err := pgxscan.Select(ctx, tx, &orders, query)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrObjectNotFound
		}
		return nil, err
	}

	return orders, nil
}
