package api

import (
	"context"
	"strconv"

	"HomeWork_1/internal/model"
	"HomeWork_1/pkg/api/pvz"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Handler) GetOrders(ctx context.Context, req *pvz.GetOrdersRequest) (*pvz.GetOrdersResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Orders")
	defer span.Finish()

	if orders, ok := s.OrdersInMem.Get(model.ClientID(req.ClientID)); ok {
		return &pvz.GetOrdersResponse{Status: "success", Orders: model.ListOrderToProto(orders)}, nil
	}

	readyClientID, err := s.inputValidator.ValidateListOrders(req.ClientID, req.Action)
	if err != nil {
		log.Errorf(ctx, "GetOrders finished with error: %v", err)
		return &pvz.GetOrdersResponse{Status: "failed"}, status.Errorf(codes.Internal, "Orders finished with error: "+unexpectedError)
	}

	actionInt, err := strconv.ParseInt(req.Action, 10, 64)
	if err != nil {
		log.Errorf(ctx, "GetOrders finished with error: %v", err)
		return &pvz.GetOrdersResponse{Status: "failed"}, status.Errorf(codes.Internal, "Orders finished with error: "+unexpectedError)
	}

	listOfOrders, _, err := s.module.ListOrders(ctx, readyClientID, int(actionInt))
	if err != nil {
		log.Errorf(ctx, "GetOrders finished with error: %v", err)
		return &pvz.GetOrdersResponse{Status: "failed"}, status.Errorf(codes.Internal, "Orders finished with error: "+unexpectedError)
	}

	s.OrdersInMem.Put(model.ClientID(req.ClientID), listOfOrders)

	return &pvz.GetOrdersResponse{Status: "success", Orders: model.ListOrderToProto(listOfOrders)}, nil
}
