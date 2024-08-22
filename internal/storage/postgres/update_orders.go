package postgres

import (
	"context"
	"github.com/opentracing/opentracing-go"

	"HomeWork_1/internal/model"

	"github.com/jackc/pgx/v4"
)

func (r *Repository) UpdateOrders(ctx context.Context, orders []model.Order, hash string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.UpdateOrders")
	defer span.Finish()

	query := `UPDATE orders SET client_id=$1, condition=$2, arrived_at=$3, received_at=$4, package=$5, price=$6, max_weight=$7 WHERE id=$8;`
	return r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		for _, order := range orders {
			_, err := tx.Exec(ctx, query, order.ClientID, order.Condition, order.ArrivedAt, order.ReceivedAt, order.Package, order.Price, order.MaxWeight, order.OrderID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
