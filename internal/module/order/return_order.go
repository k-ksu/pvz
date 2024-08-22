package module

import (
	"HomeWork_1/internal/model/errs"
	"context"
	"github.com/opentracing/opentracing-go"

	"HomeWork_1/internal/model"
)

func (m *Module) ReturnOrder(ctx context.Context, orderID model.OrderID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "module.ReturnOrder")
	defer span.Finish()

	h := m.hashGenerator.Generate()

	order, err := m.storage.OrderByID(ctx, orderID, h)
	if err != nil {
		return err
	}
	if order == nil {
		return errs.ErrUnknownOrder
	}

	if !canBeReturned(order) {
		return errs.ErrOrderCannotBeReturned

	}

	return m.storage.DeleteOrder(ctx, orderID, h)
}
