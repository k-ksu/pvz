package postgres

import (
	"HomeWork_1/internal/model/errs"
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"

	"HomeWork_1/internal/model"

	"github.com/jackc/pgx/v4"
)

func (r *Repository) DeleteOrder(ctx context.Context, orderID model.OrderID, hash string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.DeleteOrder")
	defer span.Finish()

	query := `DELETE FROM orders WHERE id = ($1);`
	err := r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, query, orderID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrObjectNotFound
		}
		return err
	}

	return nil
}
