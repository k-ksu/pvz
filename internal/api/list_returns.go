package api

import (
	"context"
	"math"

	"HomeWork_1/internal/model"
	"HomeWork_1/pkg/api/pvz"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Handler) GetReturns(ctx context.Context, req *pvz.GetReturnsRequest) (*pvz.GetReturnsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ListReturns")
	defer span.Finish()

	pageSize, pageNumber, err := s.inputValidator.ValidateListReturns(int(req.PageSize), int(req.PageNumber))
	if err != nil {
		log.Errorf(ctx, "ListReturns finished with error: %v", err)
		return &pvz.GetReturnsResponse{Status: "failed"}, status.Errorf(codes.Internal, "ListReturns finished with error: "+unexpectedError)
	}

	returns, _, err := s.module.ListReturns(ctx)
	if err != nil {
		log.Errorf(ctx, "ListReturns finished with error: %v", err)
		return &pvz.GetReturnsResponse{Status: "failed"}, status.Errorf(codes.Internal, "ListReturns finished with error: "+unexpectedError)

	}
	maxPage := int(math.Ceil(float64(len(returns)) / float64(pageSize)))

	if pageNumber > maxPage {
		log.Errorf(ctx, "ListReturns finished with error: %v", err)
		return &pvz.GetReturnsResponse{Status: "failed"}, status.Errorf(codes.Internal, "ListReturns finished with error: pageNumber > maxPage")
	}

	return &pvz.GetReturnsResponse{
		Status:     "success",
		Orders:     model.ListOrderToProto(returns),
		MaxPage:    int32(maxPage),
		PageSize:   int32(pageSize),
		PageNumber: int32(pageNumber),
	}, nil
}
