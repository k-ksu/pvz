package module

import (
	"context"
	"time"

	"HomeWork_1/internal/model"
)

func (m *Module) loadPackagesInMap(ctx context.Context) (map[model.PackageType]model.Package, error) {
	packages, err := m.storage.LoadPackages(ctx)
	if err != nil {
		return nil, err
	}

	packagesMap := make(map[model.PackageType]model.Package)
	for _, pack := range packages {
		packagesMap[pack.Package] = pack
	}
	return packagesMap, nil
}

func (m *Module) countFinalCost(ctx context.Context, orders []model.Order) (int, error) {
	finalCost := 0
	for _, order := range orders {
		payPack, err := m.storage.LoadPackByType(ctx, order.Package)
		if err != nil {
			return 0, err
		}
		finalCost += payPack.PackageSurcharge + order.Price
	}
	return finalCost, nil
}

func (m *Module) sameClient(orders []model.Order) bool {
	for i := 0; i < len(orders)-1; i++ {
		if orders[i].ClientID != orders[i+1].ClientID {
			return false
		}
	}

	return true
}

func checkTerm(ArrivedAt time.Time) bool {
	return ArrivedAt.After(time.Now())
}

func canBeReturned(order *model.Order) bool {
	return order.Condition == model.ConditionAccepted && !checkTerm(order.ReceivedAt)
}

func isKnownListAction(action int) bool {
	return action == model.AllOrders || action == model.ActualOrders
}
