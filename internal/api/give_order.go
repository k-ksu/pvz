package api

import (
	"HomeWork_1/internal/model/errs"
	"context"
	"errors"

	"HomeWork_1/internal/metric"
	"HomeWork_1/pkg/api/pvz"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Handler) GiveOrder(ctx context.Context, req *pvz.GiveOrderRequest) (*pvz.GiveOrderResponse, error) {
	defer metric.OrdersGiven.Add(1)

	span, ctx := opentracing.StartSpanFromContext(ctx, "GiveOrder")
	defer span.Finish()

	loadedPackages, err := s.module.LoadPackagesToCheck(ctx)
	if err != nil {
		log.Errorf(ctx, "GiveOrder finished with error: %v", err)
		return &pvz.GiveOrderResponse{Status: "failed"}, status.Errorf(codes.Internal, "GiveOrder finished with error: "+unexpectedError)
	}

	readyListOfOrders, readyPack, err := s.inputValidator.ValidateGiveOrder(req.Orders, req.Package, loadedPackages)
	if err != nil && !errors.Is(err, errs.ErrPackageDoesNotSet) {
		log.Errorf(ctx, "GiveOrder finished with error: %v", err)
		return &pvz.GiveOrderResponse{Status: "failed"}, status.Errorf(codes.Internal, "GiveOrder finished with error: "+unexpectedError)
	}

	var summa int
	if errors.Is(err, errs.ErrPackageDoesNotSet) {
		summa, err = s.module.GiveOrder(ctx, readyListOfOrders)
	} else {
		summa, err = s.module.GiveOrderWithNewPackage(ctx, readyListOfOrders, *readyPack)
	}

	if err != nil {
		log.Errorf(ctx, "GiveOrder finished with error: %v", err)
		return &pvz.GiveOrderResponse{Status: "failed"}, status.Errorf(codes.Internal, "GiveOrder finished with error: "+unexpectedError)
	}

	return &pvz.GiveOrderResponse{Status: "success", AmountToBePaid: int32(summa)}, nil
}
