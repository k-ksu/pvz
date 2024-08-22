package postgres

import (
	"HomeWork_1/internal/model"
	"HomeWork_1/internal/model/errs"
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

// get client order
func (r *Repository) OrderByClientIDOrderID(orderID model.OrderID, clientID model.ClientID, ctx context.Context, hash string) (*model.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.OrderByClientIDOrderID")
	defer span.Finish()

	var order model.Order
	query := `SELECT id, client_id, condition, arrived_at, received_at, package, price, max_weight FROM orders WHERE id=$1 AND client_id=$2;`
	err := r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		err := pgxscan.Get(ctx, tx, &order, query, orderID, clientID)
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

	return &order, nil
}
