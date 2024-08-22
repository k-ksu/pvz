package module

import (
	"HomeWork_1/internal/model/errs"
	"context"
	"github.com/opentracing/opentracing-go"
	"time"

	"HomeWork_1/internal/model"
)

func (m *Module) ReturnFromClient(ctx context.Context, orderID model.OrderID, clientID model.ClientID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "module.ReturnFromClient")
	defer span.Finish()

	h := m.hashGenerator.Generate()

	// get client order
	loadOrder, err := m.storage.OrderByClientIDOrderID(orderID, clientID, ctx, h)
	if err != nil {
		return err
	}

	counter := 0
	if loadOrder.Condition != model.ConditionGiven {
		return errs.ErrOrderWasReceived
	}

	if checkTerm(loadOrder.ReceivedAt.Add(48 * time.Hour)) {
		loadOrder.Condition = model.ConditionRefund
		counter = counter + 1
	} else {
		return errs.ErrReturnTimeExpired
	}

	if counter == 0 {
		return errs.ErrOrderDoesNotReceived
	}

	return m.storage.UpdateOrders(ctx, []model.Order{*loadOrder}, h)
}
