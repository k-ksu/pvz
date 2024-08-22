package module

import (
	"HomeWork_1/internal/model"
	"context"
)

type hashGenerator interface {
	Generate() string
}

type orderStorage interface {
	AppendOrder(ctx context.Context, order model.Order, hash string) error
	OrderByID(ctx context.Context, orderID model.OrderID, hash string) (*model.Order, error)
	UpdateOrders(ctx context.Context, orders []model.Order, hash string) error
	LoadOrders(ctx context.Context, hash string) ([]model.Order, error)
	DeleteOrder(ctx context.Context, orderID model.OrderID, hash string) error
	OrdersWithRefundCondition(ctx context.Context, hash string) ([]model.Order, error)
	OrderByClientIDOrderID(orderID model.OrderID, clientID model.ClientID, ctx context.Context, hash string) (*model.Order, error)
	OrdersWithGivenClientID(clientID model.ClientID, ctx context.Context, hash string) ([]model.Order, error)
	OrdersWithGivenClientIDAndAcceptedCondition(clientID model.ClientID, ctx context.Context, hash string) ([]model.Order, error)
	OrdersByGivenOrderIDList(orderIDs []model.OrderID, ctx context.Context, hash string) ([]model.Order, error)
	LoadPackages(ctx context.Context) ([]model.Package, error)
	LoadPackByType(ctx context.Context, packageType model.PackageType) (*model.Package, error)
}

type Module struct {
	storage       orderStorage
	hashGenerator hashGenerator
}

func NewModule(storage orderStorage, hashGenerator hashGenerator) Module {
	return Module{
		storage:       storage,
		hashGenerator: hashGenerator,
	}
}
