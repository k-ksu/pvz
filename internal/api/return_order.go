package api

import (
	"context"

	"HomeWork_1/pkg/api/pvz"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Handler) ReturnOrder(ctx context.Context, req *pvz.ReturnOrderRequest) (*pvz.ReturnOrderResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReturnOrder")
	defer span.Finish()

	readyOrderID, err := s.inputValidator.ValidateReturnOrder(req.OrderID)
	if err != nil {
		log.Errorf(ctx, "ReturnOrder finished with error: %v", err)
		return &pvz.ReturnOrderResponse{Status: "failed"}, status.Errorf(codes.Internal, unexpectedError)

	}

	err = s.module.ReturnOrder(ctx, readyOrderID)
	if err != nil {
		log.Errorf(ctx, "ReturnOrder finished with error: %v", err)
		return &pvz.ReturnOrderResponse{Status: "failed"}, status.Errorf(codes.Internal, "ReturnOrder finished with error: "+unexpectedError)
	}

	return &pvz.ReturnOrderResponse{Status: "success"}, nil
}
