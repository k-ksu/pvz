package api

import (
	"HomeWork_1/internal/model/errs"
	"context"
	"errors"

	"HomeWork_1/internal/model"
	"HomeWork_1/pkg/api/pvz"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AcceptOrder ...
func (s *Handler) AcceptOrder(ctx context.Context, req *pvz.AcceptOrderRequest) (*pvz.AcceptOrderResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AcceptOrder")
	defer span.Finish()

	readyOrder, err := s.inputValidator.ValidateGetOrderFromCourier(req.OrderID, req.ClientID, req.Date, req.Package, int(req.Price), int(req.Weight))

	if err != nil && !errors.Is(err, errs.ErrPackageDoesNotSet) {
		log.Errorf(ctx, "AcceptOrder finished with error: %v", err)
		return &pvz.AcceptOrderResponse{Status: "failed"}, status.Errorf(codes.Internal, unexpectedError)
	}

	if readyOrder.Package != model.WithoutPackage {
		loadedPackages, err := s.module.LoadPackagesToCheck(ctx)
		if err != nil {
			log.Errorf(ctx, "AcceptOrder finished with error: %v", err)
			return &pvz.AcceptOrderResponse{Status: "failed"}, status.Errorf(codes.Internal, unexpectedError)
		}

		err = s.inputValidator.ValidatePackage(int(req.Weight), readyOrder.Package, loadedPackages)
		if err != nil {
			log.Errorf(ctx, "AcceptOrder finished with error: %v", err)
			return &pvz.AcceptOrderResponse{Status: "failed"}, status.Errorf(codes.Internal, unexpectedError)
		}
	}

	readyOrder.MaxWeight = int(req.Weight)
	readyOrder.Price = int(req.Price)
	err = s.module.GetOrderFromCourier(ctx, readyOrder)
	if err != nil {
		log.Errorf(ctx, "AcceptOrder finished with error: %v", err)
		return &pvz.AcceptOrderResponse{Status: "failed"}, status.Errorf(codes.Internal, "AcceptOrder finished with error: "+unexpectedError)
	}

	return &pvz.AcceptOrderResponse{Status: "success"}, nil
}
