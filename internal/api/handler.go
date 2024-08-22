package api

import (
	"HomeWork_1/internal/pkg/cache/inmemory"
	"context"
	"log"
	"time"

	"HomeWork_1/internal/model"
	"HomeWork_1/pkg/api/pvz"
)

const unexpectedError = "unexpected error"

type orderModule interface {
	GetOrderFromCourier(ctx context.Context, order *model.Order) error
	GiveOrder(ctx context.Context, orders []model.OrderID) (int, error)
	ReturnOrder(ctx context.Context, orderID model.OrderID) error
	ListOrders(ctx context.Context, clientID model.ClientID, action int) ([]model.Order, map[model.PackageType]model.Package, error)
	ReturnFromClient(ctx context.Context, orderID model.OrderID, clientID model.ClientID) error
	ListReturns(ctx context.Context) ([]model.Order, map[model.PackageType]model.Package, error)
	LoadPackagesToCheck(ctx context.Context) ([]model.Package, error)
	GiveOrderWithNewPackage(ctx context.Context, orders []model.OrderID, pack model.PackageType) (int, error)
}

type inputValidator interface {
	ValidateGetOrderFromCourier(orderID, clientID, date, pack string, price, weight int) (*model.Order, error)
	ValidateGiveOrder(orders, pack string, loadedPackages []model.Package) ([]model.OrderID, *model.PackageType, error)
	ValidateReturnOrder(orderID string) (model.OrderID, error)
	ValidateListOrders(clientID, action string) (model.ClientID, error)
	ValidateReturnFromClient(orderID, clientID string) (model.OrderID, model.ClientID, error)
	ValidateListReturns(pageSize, pageNumber int) (int, int, error)
	ValidatePackage(weight int, pack model.PackageType, loadedPackages []model.Package) error
}

type sender interface {
	SendMessage(event *model.EventMessage) error
}

type Handler struct {
	pvz.UnimplementedPVZServer

	module         orderModule
	inputValidator inputValidator
	sender         sender
	OrdersInMem    *inmemory.InMemoryCache[model.ClientID, []model.Order]
}

func NewHandler(module orderModule, inputValidator inputValidator, sender sender) *Handler {
	ordersInMem, err := inmemory.NewInMemoryCache[model.ClientID, []model.Order](5, 1*time.Minute, inmemory.LfuCache)
	if err != nil {
		log.Fatal(err)
	}

	return &Handler{
		module:         module,
		inputValidator: inputValidator,
		sender:         sender,
		OrdersInMem:    ordersInMem,
	}
}
