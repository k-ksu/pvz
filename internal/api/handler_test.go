package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"HomeWork_1/internal/api/mocks"
	"HomeWork_1/internal/model"
	"HomeWork_1/internal/service/input_validator"
	"HomeWork_1/pkg/api/pvz"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService_GetOrderFromCourier(t *testing.T) {
	t.Parallel()

	var (
		mc             = minimock.NewController(t)
		ctx            = context.Background()
		inputValidator = input_validator.NewValidator()
		someError      = errors.New("unexpected error")
		packages       = []model.Package{
			{
				Package:          "box",
				PackageSurcharge: 20,
				PackageMaxWeight: 30000,
			},
			{
				Package:          "noPackage",
				PackageSurcharge: 0,
				PackageMaxWeight: -1,
			},
			{
				Package:          "plasticBag",
				PackageSurcharge: 5,
				PackageMaxWeight: 10000,
			},
			{
				Package:          "film",
				PackageSurcharge: 1,
				PackageMaxWeight: -1,
			},
		}
	)

	location, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)

	tests := []struct {
		desc    string
		handler func() Handler
		req     *pvz.GetOrderFromCourierRequest
		wantRsp *pvz.GetOrderFromCourierResponse
		wantErr error
	}{
		{
			desc: "Test case 1: correct order",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				dateTime, err := time.ParseInLocation("02.01.2006-15:04", "02.01.2006-15:04", location)
				require.NoError(t, err)

				module.LoadPackagesToCheckMock.When(ctx).Then(packages, nil)

				module.GetOrderFromCourierMock.When(
					ctx,
					&model.Order{
						OrderID:   "1",
						ClientID:  "1",
						ArrivedAt: dateTime,
						Price:     10,
						Package:   "film",
						MaxWeight: 10,
					},
				).Then(nil)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.GetOrderFromCourierRequest{
				OrderId:  "1",
				ClientId: "1",
				Date:     "02.01.2006-15:04",
				Price:    10,
				Weight:   10,
				Package:  "film",
			},
			wantRsp: &pvz.GetOrderFromCourierResponse{Status: "success"},
			wantErr: nil,
		},
		{
			desc: "Test case 2: package not correct",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				_, err := time.ParseInLocation("02.01.2006-15:04", "02.01.2006-15:04", location)
				require.NoError(t, err)

				module.LoadPackagesToCheckMock.When(ctx).Then(packages, nil)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.GetOrderFromCourierRequest{
				OrderId:  "1",
				ClientId: "1",
				Date:     "02.01.2006-15:04",
				Price:    10,
				Weight:   10,
				Package:  "kkk",
			},
			wantRsp: &pvz.GetOrderFromCourierResponse{Status: "failed"},
			wantErr: status.Errorf(codes.Internal, "unexpected error"),
		},
		{
			desc: "Test case 3: package not correct",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				_, err := time.ParseInLocation("02.01.2006-15:04", "02.01.2006-15:04", location)
				require.NoError(t, err)

				module.LoadPackagesToCheckMock.When(ctx).Then(nil, someError)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.GetOrderFromCourierRequest{
				OrderId:  "1",
				ClientId: "1",
				Date:     "02.01.2006-15:04",
				Price:    10,
				Weight:   10,
				Package:  "film",
			},
			wantRsp: &pvz.GetOrderFromCourierResponse{Status: "failed"},
			wantErr: status.Errorf(codes.Internal, someError.Error()),
		},
		{
			desc: "Test case 4: data in order does not correct",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				_, err := time.ParseInLocation("02.01.2006-15:04", "02.01.2006-15:04", location)
				require.NoError(t, err)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.GetOrderFromCourierRequest{
				OrderId:  "1",
				ClientId: "1",
				Date:     "02.01.2006-15:04",
				Price:    10,
				Weight:   -1,
				Package:  "film",
			},
			wantRsp: &pvz.GetOrderFromCourierResponse{Status: "failed"},
			wantErr: status.Errorf(codes.Internal, "unexpected error"),
		},
		{
			desc: "Test case 5: some err in GetOrderFromCourier",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				dateTime, err := time.ParseInLocation("02.01.2006-15:04", "02.01.2006-15:04", location)
				require.NoError(t, err)

				module.LoadPackagesToCheckMock.When(ctx).Then(packages, nil)

				module.GetOrderFromCourierMock.When(
					ctx,
					&model.Order{
						OrderID:   "1",
						ClientID:  "1",
						ArrivedAt: dateTime,
						Price:     10,
						Package:   "film",
						MaxWeight: 10,
					},
				).Then(someError)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.GetOrderFromCourierRequest{
				OrderId:  "1",
				ClientId: "1",
				Date:     "02.01.2006-15:04",
				Price:    10,
				Weight:   10,
				Package:  "film",
			},
			wantRsp: &pvz.GetOrderFromCourierResponse{Status: "failed"},
			wantErr: status.Errorf(codes.Internal, "GetOrderFromCourier finished with error: "+someError.Error()),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			handler := tc.handler()
			rsp, err := handler.GetOrderFromCourier(ctx, tc.req)
			require.Equal(t, tc.wantRsp, rsp)
			require.Equal(t, tc.wantErr, err)
		})
	}
}

func TestService_ReturnOrder(t *testing.T) {
	t.Parallel()

	var (
		mc             = minimock.NewController(t)
		ctx            = context.Background()
		inputValidator = input_validator.NewValidator()
	)

	tests := []struct {
		desc    string
		handler func() Handler
		req     *pvz.ReturnOrderRequest
		wantRsp *pvz.ReturnOrderResponse
		wantErr error
	}{
		{
			desc: "Test case 1: correct order",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				module.ReturnOrderMock.When(
					ctx,
					"1",
				).Then(nil)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.ReturnOrderRequest{
				OrderId: "1",
			},
			wantRsp: &pvz.ReturnOrderResponse{Status: "success"},
			wantErr: nil,
		},
		{
			desc: "Test case 2: err in validator",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.ReturnOrderRequest{
				OrderId: "",
			},
			wantRsp: &pvz.ReturnOrderResponse{Status: "failed"},
			wantErr: status.Errorf(codes.Internal, "unexpected error"),
		},
		{
			desc: "Test case 3: err in ReturnOrder",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				module.ReturnOrderMock.When(
					ctx,
					"1",
				).Then(errors.New("order cannot be returned"))

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.ReturnOrderRequest{
				OrderId: "1",
			},
			wantRsp: &pvz.ReturnOrderResponse{Status: "failed"},
			wantErr: status.Errorf(codes.Internal, "ReturnOrder finished with error: unexpected error"),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			handler := tc.handler()
			rsp, err := handler.ReturnOrder(ctx, tc.req)
			require.Equal(t, tc.wantRsp, rsp)
			require.Equal(t, tc.wantErr, err)
		})
	}
}

func TestService_GiveOrder(t *testing.T) {
	t.Parallel()

	var (
		mc             = minimock.NewController(t)
		ctx            = context.Background()
		inputValidator = input_validator.NewValidator()
		someError      = errors.New("unexpected error")
		packages       = []model.Package{
			{
				Package:          "box",
				PackageSurcharge: 20,
				PackageMaxWeight: 30000,
			},
			{
				Package:          "noPackage",
				PackageSurcharge: 0,
				PackageMaxWeight: -1,
			},
			{
				Package:          "plasticBag",
				PackageSurcharge: 5,
				PackageMaxWeight: 10000,
			},
			{
				Package:          "film",
				PackageSurcharge: 1,
				PackageMaxWeight: -1,
			},
		}
	)

	location, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)

	tests := []struct {
		desc    string
		handler func() Handler
		req     *pvz.GiveOrderRequest
		wantRsp *pvz.GiveOrderResponse
		wantErr error
	}{
		{
			desc: "Test case 1: correct order",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				_, err := time.ParseInLocation("02.01.2006-15:04", "02.01.2006-15:04", location)
				require.NoError(t, err)

				module.LoadPackagesToCheckMock.When(ctx).Then(packages, nil)

				module.GiveOrderMock.When(
					ctx,
					[]model.OrderID{
						"1",
						"2",
						"3",
					},
				).Then(50,
					nil,
				)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.GiveOrderRequest{
				Orders:  "1,2,3",
				Package: "",
			},
			wantRsp: &pvz.GiveOrderResponse{Status: "success", AmountToBePaid: 50},
			wantErr: nil,
		},
		{
			desc: "Test case 2: correct order with new film",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				_, err := time.ParseInLocation("02.01.2006-15:04", "02.01.2006-15:04", location)
				require.NoError(t, err)

				module.LoadPackagesToCheckMock.When(ctx).Then(packages, nil)

				module.GiveOrderWithNewPackageMock.When(
					ctx,
					[]model.OrderID{
						"1",
						"2",
						"3",
					},
					"film",
				).Then(50,
					nil,
				)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.GiveOrderRequest{
				Orders:  "1,2,3",
				Package: "film",
			},
			wantRsp: &pvz.GiveOrderResponse{Status: "success", AmountToBePaid: 50},
			wantErr: nil,
		},
		{
			desc: "Test case 3: not correct package load",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				module.LoadPackagesToCheckMock.When(ctx).Then(nil, someError)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.GiveOrderRequest{
				Orders:  "1,2,3",
				Package: "kkk",
			},
			wantRsp: &pvz.GiveOrderResponse{Status: "failed", AmountToBePaid: 0},
			wantErr: status.Errorf(codes.Internal, "GiveOrder finished with error: "+someError.Error()),
		},
		{
			desc: "Test case 4: err in validator",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				_, err := time.ParseInLocation("02.01.2006-15:04", "02.01.2006-15:04", location)
				require.NoError(t, err)

				module.LoadPackagesToCheckMock.When(ctx).Then(packages, nil)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.GiveOrderRequest{
				Orders:  "",
				Package: "",
			},
			wantRsp: &pvz.GiveOrderResponse{Status: "failed", AmountToBePaid: 0},
			wantErr: status.Errorf(codes.Internal, "GiveOrder finished with error: unexpected error"),
		},
		{
			desc: "Test case 5: err in give order with new package",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				_, err := time.ParseInLocation("02.01.2006-15:04", "02.01.2006-15:04", location)
				require.NoError(t, err)

				module.LoadPackagesToCheckMock.When(ctx).Then(packages, nil)

				module.GiveOrderWithNewPackageMock.When(
					ctx,
					[]model.OrderID{
						"1",
						"2",
						"3",
					},
					"film",
				).Then(0,
					someError,
				)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.GiveOrderRequest{
				Orders:  "1,2,3",
				Package: "film",
			},
			wantRsp: &pvz.GiveOrderResponse{Status: "failed", AmountToBePaid: 0},
			wantErr: status.Errorf(codes.Internal, "GiveOrder finished with error: "+someError.Error()),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			handler := tc.handler()
			rsp, err := handler.GiveOrder(ctx, tc.req)
			require.Equal(t, tc.wantRsp, rsp)
			require.Equal(t, tc.wantErr, err)
		})
	}
}

func TestService_ListOrders(t *testing.T) {
	t.Parallel()

	var (
		mc             = minimock.NewController(t)
		ctx            = context.Background()
		inputValidator = input_validator.NewValidator()
		someError      = errors.New("unexpected error")
		tomorrow       = time.Now().Add(24 * time.Hour)
	)

	tests := []struct {
		desc    string
		handler func() Handler
		req     *pvz.ListOrdersRequest
		wantRsp *pvz.ListOrdersResponse
		wantErr error
	}{
		{
			desc: "Test case 1: correct order",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				module.ListOrdersMock.When(
					ctx,
					"1",
					1,
				).Then(
					[]model.Order{
						{
							OrderID:    "1",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Film,
							MaxWeight:  600,
						},
					},
					map[model.PackageType]model.Package{
						model.Box: {
							Package:          "box",
							PackageSurcharge: 20,
							PackageMaxWeight: 30000,
						},
						model.WithoutPackage: {
							Package:          "noPackage",
							PackageSurcharge: 0,
							PackageMaxWeight: -1,
						},
						model.PlasticBag: {
							Package:          "plasticBag",
							PackageSurcharge: 5,
							PackageMaxWeight: 10000,
						},
						model.Film: {
							Package:          "film",
							PackageSurcharge: 1,
							PackageMaxWeight: -1,
						},
					},
					nil,
				)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.ListOrdersRequest{
				ClientId: "1",
				Action:   "1",
			},
			wantRsp: &pvz.ListOrdersResponse{Status: "success", Orders: []*pvz.Order{
				{
					OrderId:       "1",
					ClientId:      "1",
					Condition:     "accepted",
					DateArrival:   timestamppb.New(tomorrow),
					DateReceiving: timestamppb.New(time.Time{}),
					Price:         700,
					Package:       "film",
					MaxWeight:     600,
				},
			}},
			wantErr: nil,
		},
		{
			desc: "Test case 2: some err in validator",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.ListOrdersRequest{
				ClientId: "",
				Action:   "1",
			},
			wantRsp: &pvz.ListOrdersResponse{Status: "failed", Orders: nil},
			wantErr: status.Errorf(codes.Internal, "ListOrders finished with error: unexpected error"),
		},
		{
			desc: "Test case 3: err in list orders",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				module.ListOrdersMock.When(
					ctx,
					"1",
					1,
				).Then(
					nil,
					nil,
					someError,
				)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.ListOrdersRequest{
				ClientId: "1",
				Action:   "1",
			},
			wantRsp: &pvz.ListOrdersResponse{Status: "failed", Orders: nil},
			wantErr: status.Errorf(codes.Internal, "ListOrders finished with error: "+someError.Error()),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			handler := tc.handler()
			rsp, err := handler.ListOrders(ctx, tc.req)
			require.Equal(t, tc.wantRsp, rsp)
			require.Equal(t, tc.wantErr, err)
		})
	}
}

func TestService_ReturnFromClient(t *testing.T) {
	t.Parallel()

	var (
		mc             = minimock.NewController(t)
		ctx            = context.Background()
		inputValidator = input_validator.NewValidator()
		someError      = errors.New("unexpected error")
	)

	tests := []struct {
		desc    string
		handler func() Handler
		req     *pvz.ReturnFromClientRequest
		wantRsp *pvz.ReturnFromClientResponse
		wantErr error
	}{
		{
			desc: "Test case 1: correct order",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				module.ReturnFromClientMock.When(
					ctx,
					"1",
					"1",
				).Then(nil)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.ReturnFromClientRequest{
				OrderId:  "1",
				ClientId: "1",
			},
			wantRsp: &pvz.ReturnFromClientResponse{Status: "success"},
			wantErr: nil,
		},
		{
			desc: "Test case 2: err in validator",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.ReturnFromClientRequest{
				OrderId:  "",
				ClientId: "1",
			},
			wantRsp: &pvz.ReturnFromClientResponse{Status: "failed"},
			wantErr: status.Errorf(codes.Internal, "ReturnFromClient finished with error: unexpected error"),
		},
		{
			desc: "Test case 3: err in return from client",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				module.ReturnFromClientMock.When(
					ctx,
					"1",
					"1",
				).Then(someError)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.ReturnFromClientRequest{
				OrderId:  "1",
				ClientId: "1",
			},
			wantRsp: &pvz.ReturnFromClientResponse{Status: "failed"},
			wantErr: status.Errorf(codes.Internal, "ReturnFromClient finished with error: "+someError.Error()),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			handler := tc.handler()
			rsp, err := handler.ReturnFromClient(ctx, tc.req)
			require.Equal(t, tc.wantRsp, rsp)
			require.Equal(t, tc.wantErr, err)
		})
	}
}

func TestService_ListReturns(t *testing.T) {
	t.Parallel()

	var (
		mc             = minimock.NewController(t)
		ctx            = context.Background()
		inputValidator = input_validator.NewValidator()
		someError      = errors.New("unexpected error")
		tomorrow       = time.Now().Add(24 * time.Hour)
	)

	tests := []struct {
		desc    string
		handler func() Handler
		req     *pvz.ListReturnsRequest
		wantRsp *pvz.ListReturnsResponse
		wantErr error
	}{
		{
			desc: "Test case 1: correct order",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				module.ListReturnsMock.When(
					ctx,
				).Then(
					[]model.Order{
						{
							OrderID:    "1",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Film,
							MaxWeight:  600,
						},
					},
					map[model.PackageType]model.Package{
						model.Box: {
							Package:          "box",
							PackageSurcharge: 20,
							PackageMaxWeight: 30000,
						},
						model.WithoutPackage: {
							Package:          "noPackage",
							PackageSurcharge: 0,
							PackageMaxWeight: -1,
						},
						model.PlasticBag: {
							Package:          "plasticBag",
							PackageSurcharge: 5,
							PackageMaxWeight: 10000,
						},
						model.Film: {
							Package:          "film",
							PackageSurcharge: 1,
							PackageMaxWeight: -1,
						},
					},
					nil,
				)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.ListReturnsRequest{
				PageNumber: 1,
				PageSize:   1,
			},
			wantRsp: &pvz.ListReturnsResponse{Status: "success", Orders: []*pvz.Order{
				{
					OrderId:       "1",
					ClientId:      "1",
					Condition:     "accepted",
					DateArrival:   timestamppb.New(tomorrow),
					DateReceiving: timestamppb.New(time.Time{}),
					Price:         700,
					Package:       "film",
					MaxWeight:     600,
				},
			},
				MaxPage:    1,
				PageSize:   1,
				PageNumber: 1,
			},
			wantErr: nil,
		},
		{
			desc: "Test case 2: not positive page number",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.ListReturnsRequest{
				PageNumber: 1,
				PageSize:   0,
			},
			wantRsp: &pvz.ListReturnsResponse{Status: "failed", Orders: nil,
				MaxPage:    0,
				PageSize:   0,
				PageNumber: 0,
			},
			wantErr: status.Errorf(codes.Internal, "ListReturns finished with error: unexpected error"),
		},
		{
			desc: "Test case 3: err in listReturns",
			handler: func() Handler {
				module := mocks.NewModule(mc)

				module.ListReturnsMock.When(
					ctx,
				).Then(
					nil,
					nil,
					someError,
				)

				return Handler{
					module:         module,
					inputValidator: inputValidator,
					sender:         nil}
			},
			req: &pvz.ListReturnsRequest{
				PageNumber: 1,
				PageSize:   1,
			},
			wantRsp: &pvz.ListReturnsResponse{Status: "failed", Orders: nil,
				MaxPage:    0,
				PageSize:   0,
				PageNumber: 0,
			},
			wantErr: status.Errorf(codes.Internal, "ListReturns finished with error: "+someError.Error()),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			handler := tc.handler()
			rsp, err := handler.ListReturns(ctx, tc.req)
			require.Equal(t, tc.wantRsp, rsp)
			require.Equal(t, tc.wantErr, err)
		})
	}
}
