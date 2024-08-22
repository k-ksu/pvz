package postgres

import (
	"HomeWork_1/internal/model/errs"
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"

	"HomeWork_1/internal/model"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
)

// получить заказы по списку айди
func (r *Repository) OrdersByGivenOrderIDList(orderIDs []model.OrderID, ctx context.Context, hash string) ([]model.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.OrdersByGivenOrderIDList")
	defer span.Finish()

	var orders []model.Order
	query := `SELECT id, client_id, condition, arrived_at, received_at, package, price, max_weight FROM orders WHERE id = ANY ($1);`
	err := r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		var orderIDsStr []string
		for _, ord := range orderIDs {
			orderIDsStr = append(orderIDsStr, string(ord))
		}
		err := pgxscan.Select(ctx, tx, &orders, query, pq.Array(orderIDsStr))
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
