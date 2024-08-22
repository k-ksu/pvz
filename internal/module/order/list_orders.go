package module

import (
	"HomeWork_1/internal/model/errs"
	"context"
	"github.com/opentracing/opentracing-go"

	"HomeWork_1/internal/model"
)

func (m *Module) ListOrders(ctx context.Context, clientID model.ClientID, action int) ([]model.Order, map[model.PackageType]model.Package, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "module.ListOrders")
	defer span.Finish()

	h := m.hashGenerator.Generate()

	if !isKnownListAction(action) {
		return nil, nil, errs.ErrUnknownAction
	}

	packagesMap, err := m.loadPackagesInMap(ctx)
	if err != nil {
		return nil, nil, err
	}

	// get orders where clientid = clientid
	if action == model.AllOrders {
		orders, err := m.storage.OrdersWithGivenClientID(clientID, ctx, h)
		if err != nil {
			return nil, nil, err
		}
		return orders, packagesMap, nil
	}

	// get client accepted orders where clientid = clienid and condition = accepted
	if action == model.ActualOrders {
		orders, err := m.storage.OrdersWithGivenClientIDAndAcceptedCondition(clientID, ctx, h)
		if err != nil {
			return nil, nil, err
		}
		return orders, packagesMap, nil
	}
	return nil, nil, nil
}
