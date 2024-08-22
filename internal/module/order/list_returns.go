package module

import (
	"context"
	"github.com/opentracing/opentracing-go"

	"HomeWork_1/internal/model"
)

func (m *Module) ListReturns(ctx context.Context) ([]model.Order, map[model.PackageType]model.Package, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "module.ListReturns")
	defer span.Finish()

	h := m.hashGenerator.Generate()

	// get refund orders
	orders, err := m.storage.OrdersWithRefundCondition(ctx, h)
	if err != nil {
		return nil, nil, err
	}

	packagesMap, err := m.loadPackagesInMap(ctx)
	if err != nil {
		return nil, nil, err
	}

	return orders, packagesMap, nil
}
