package module

import (
	"HomeWork_1/internal/model/errs"
	"context"
	"github.com/opentracing/opentracing-go"
	"time"

	"HomeWork_1/internal/model"
)

func (m *Module) GiveOrder(ctx context.Context, orders []model.OrderID) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "module.GiveOrder")
	defer span.Finish()

	h := m.hashGenerator.Generate()

	// получить заказы по списку айди
	allOrders, err := m.storage.OrdersByGivenOrderIDList(orders, ctx, h)
	if err != nil {
		return 0, err
	}

	ordersToGive := []model.Order{}
	IDsToCount := []model.OrderID{}

	for i := 0; i < len(allOrders); i++ {
		if !(allOrders[i].Condition == model.ConditionAccepted &&
			checkTerm(allOrders[i].ArrivedAt)) {
			continue
		}
		allOrders[i].ReceivedAt = time.Now()
		allOrders[i].Condition = model.ConditionGiven
		IDsToCount = append(IDsToCount, allOrders[i].OrderID)
		ordersToGive = append(ordersToGive, allOrders[i])
	}

	if !m.sameClient(ordersToGive) {
		return 0, errs.ErrOrdersNotOfOneClient
	}

	resultOfUpdate := m.storage.UpdateOrders(ctx, ordersToGive, h)

	OrdersToCount, err := m.storage.OrdersByGivenOrderIDList(IDsToCount, ctx, h)
	if err != nil {
		return 0, err
	}

	finalPrice, err := m.countFinalCost(ctx, OrdersToCount)
	if err != nil {
		return 0, err
	}

	return finalPrice, resultOfUpdate
}
