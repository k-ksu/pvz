package module

import (
	"HomeWork_1/internal/model/errs"
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"time"

	"HomeWork_1/internal/model"
)

func (m *Module) GetOrderFromCourier(ctx context.Context, order *model.Order) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "module.GetOrderFromCourier")
	defer span.Finish()

	h := m.hashGenerator.Generate()

	if !checkTerm(order.ArrivedAt) {
		return errs.ErrArrivedAtInPast
	}

	order.Condition = model.ConditionAccepted
	order.ReceivedAt = time.Time{}

	orderFound, err := m.storage.OrderByID(ctx, order.OrderID, h)
	if err != nil {
		if !errors.Is(err, errs.ErrObjectNotFound) {
			return err
		}
	}

	if orderFound != nil {
		return errs.ErrOrderWasAlreadyAccepted
	}

	return m.storage.AppendOrder(ctx, *order, h)
}
