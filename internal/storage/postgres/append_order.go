package postgres

import (
	"context"
	"github.com/opentracing/opentracing-go"

	"HomeWork_1/internal/model"

	"github.com/jackc/pgx/v4"
)

func (r *Repository) AppendOrder(ctx context.Context, order model.Order, hash string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.AppendOrder")
	defer span.Finish()

	query := `
		INSERT INTO orders(
			id, 
			client_id, 
			condition, 
			arrived_at, 
			received_at, 
			package, 
			price, 
			max_weight
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	return r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(
			ctx,
			query,
			order.OrderID,
			order.ClientID,
			order.Condition,
			order.ArrivedAt,
			order.ReceivedAt,
			order.Package,
			order.Price,
			order.MaxWeight,
		)

		return err
	})
}
