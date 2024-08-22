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

func (r *Repository) LoadPackages(ctx context.Context) ([]model.Package, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.LoadPackages")
	defer span.Finish()

	var packages []model.Package
	query := `SELECT package, max_weight, surcharge FROM packages_info;`
	err := r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		err := pgxscan.Select(ctx, tx, &packages, query)
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

	return packages, nil
}
