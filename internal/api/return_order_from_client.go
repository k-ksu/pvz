package api

import (
	"context"

	"HomeWork_1/pkg/api/pvz"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Handler) ReturnOrderFromClient(ctx context.Context, req *pvz.ReturnOrderFromClientRequest) (*pvz.ReturnOrderFromClientResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReturnOrderFromClient")
	defer span.Finish()

	readyOrderID, readyClientID, err := s.inputValidator.ValidateReturnFromClient(req.OrderID, req.ClientID)
	if err != nil {
		log.Errorf(ctx, "ReturnOrderFromClient finished with error: %v", err)
		return &pvz.ReturnOrderFromClientResponse{Status: "failed"}, status.Errorf(codes.Internal, "ReturnOrderFromClient finished with error: "+unexpectedError)
	}

	err = s.module.ReturnFromClient(ctx, readyOrderID, readyClientID)
	if err != nil {
		log.Errorf(ctx, "ReturnOrderFromClient finished with error: %v", err)
		return &pvz.ReturnOrderFromClientResponse{Status: "failed"}, status.Errorf(codes.Internal, "ReturnOrderFromClient finished with error: "+unexpectedError)
	}

	return &pvz.ReturnOrderFromClientResponse{Status: "success"}, nil
}
