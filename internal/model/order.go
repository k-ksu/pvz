package model

import (
	"time"

	"HomeWork_1/pkg/api/pvz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	OrderID  string
	ClientID string
)

type Order struct {
	OrderID    OrderID `db:"id"`
	ClientID   ClientID
	Condition  Condition
	ArrivedAt  time.Time
	ReceivedAt time.Time
	Price      int
	Package    PackageType
	MaxWeight  int
}

func ListOrderToProto(listOfOrders []Order) (orders []*pvz.Order) {
	for _, order := range listOfOrders {
		ord := &pvz.Order{
			OrderID:    string(order.OrderID),
			ClientID:   string(order.ClientID),
			Condition:  string(order.Condition),
			ArrivedAt:  timestamppb.New(order.ArrivedAt),
			ReceivedAt: timestamppb.New(order.ReceivedAt),
			Price:      int32(order.Price),
			Package:    string(order.Package),
			MaxWeight:  int32(order.MaxWeight),
		}
		orders = append(orders, ord)
	}

	return
}
