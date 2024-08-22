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

func (r *Repository) LoadPackByType(ctx context.Context, packageType model.PackageType) (*model.Package, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.loadPackByType")
	defer span.Finish()

	var packInfo model.Package
	query := `SELECT package, max_weight, surcharge FROM packages_info WHERE package = $1;`
	err := r.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		err := pgxscan.Get(ctx, tx, &packInfo, query, packageType)
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

	return &packInfo, nil
}
