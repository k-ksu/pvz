package module

import (
	"HomeWork_1/internal/model/errs"
	"context"
	"errors"
	"testing"
	"time"

	"HomeWork_1/internal/model"
	"HomeWork_1/internal/module/order/mocks"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestModule_GetOrderFromCourier(t *testing.T) {
	t.Parallel()

	var (
		mc        = minimock.NewController(t)
		ctx       = context.Background()
		tomorrow  = time.Now().Add(24 * time.Hour)
		someError = errors.New("some error")
	)

	hashGenerator := mocks.NewHashGenerator(mc)
	hashGenerator.GenerateMock.Return("someHash")

	tests := []struct {
		desc      string
		module    func() Module
		orderID   model.OrderID
		clientID  model.ClientID
		dateTime  time.Time
		pack      model.PackageType
		weight    int
		price     int
		wantError error
	}{
		{
			desc: "Test case 1: correct order - finish with nil",
			module: func() Module {
				storage := mocks.NewStorage(mc)
				storage.OrderByIDMock.When(
					ctx,
					"1",
					"someHash",
				).Then(
					nil,
					errs.ErrObjectNotFound,
				)

				o := model.Order{
					OrderID:    "1",
					Condition:  model.ConditionAccepted,
					ArrivedAt:  tomorrow,
					ReceivedAt: time.Time{},
				}
				storage.AppendOrderMock.When(
					ctx,
					o,
					"someHash",
				).Then(
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orderID:   "1",
			dateTime:  tomorrow,
			wantError: nil,
		},
		{
			desc: "Test case 2: dateTime is in past",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orderID:   "1",
			dateTime:  time.Now().Add(-24 * time.Hour),
			wantError: errors.New("срок хранения не может быть в прошлом"),
		},
		{
			desc: "Test case 3: object found",
			module: func() Module {
				storage := mocks.NewStorage(mc)
				storage.OrderByIDMock.When(
					ctx,
					"1",
					"someHash",
				).Then(
					&model.Order{},
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orderID:   "1",
			dateTime:  time.Now().Add(24 * time.Hour),
			wantError: errors.New("заказ с таким ID уже был принят на ПВЗ"),
		},
		{
			desc: "Test case 4: some error happened in OrderByOrderID",
			module: func() Module {
				storage := mocks.NewStorage(mc)
				storage.OrderByIDMock.When(
					ctx,
					"1",
					"someHash",
				).Then(
					nil,
					someError,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orderID:   "1",
			dateTime:  time.Now().Add(24 * time.Hour),
			wantError: someError,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			module := tc.module()
			order := model.Order{
				OrderID:   tc.orderID,
				ClientID:  tc.clientID,
				ArrivedAt: tc.dateTime,
				Package:   tc.pack,
				MaxWeight: tc.weight,
				Price:     tc.price,
			}
			result := module.GetOrderFromCourier(ctx, &order)
			require.Equal(t, tc.wantError, result)
		})
	}
}

func TestModule_GiveOrder(t *testing.T) {
	t.Parallel()

	var (
		mc       = minimock.NewController(t)
		ctx      = context.Background()
		tomorrow = time.Now().Add(24 * time.Hour)
		//previousDay = time.Now().Add(-24 * time.Hour)
		someError = errors.New("some error")
	)

	hashGenerator := mocks.NewHashGenerator(mc)
	hashGenerator.GenerateMock.Return("someHash")

	tests := []struct {
		desc      string
		module    func() Module
		orders    []model.OrderID
		wantSum   int
		wantError error
	}{
		{
			desc: "Test case 1: correct order - finish with nil",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11", "22"},
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				storage.LoadPackByTypeMock.When(
					ctx,
					model.Box,
				).Then(
					&model.Package{
						Package:          "Box",
						PackageSurcharge: 20,
						PackageMaxWeight: 30000,
					},
					nil,
				)

				storage.UpdateOrdersMock.Return(nil)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11"},
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orders:    []model.OrderID{"11", "22"},
			wantSum:   720,
			wantError: nil,
		},
		{
			desc: "Test case 2: some error in OrdersByGivenOrderIDList",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11", "22"},
					ctx,
					"someHash",
				).Then(
					nil,
					someError,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orders:    []model.OrderID{"11", "22"},
			wantSum:   0,
			wantError: someError,
		},
		{
			desc: "Test case 3: not same client",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11", "22"},
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
						{
							OrderID:    "22",
							ClientID:   "2",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orders:    []model.OrderID{"11", "22"},
			wantSum:   0,
			wantError: errors.New("не все заказы принадлежат одному клиенту"),
		},
		{
			desc: "Test case 4: someError in OrdersByGivenOrderIDListMock",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11", "22"},
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				storage.UpdateOrdersMock.Return(nil)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11"},
					ctx,
					"someHash",
				).Then(
					nil,
					someError,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orders:    []model.OrderID{"11", "22"},
			wantSum:   0,
			wantError: someError,
		},
		{
			desc: "Test case 5: someError in CountFinalCost",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11", "22"},
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				storage.UpdateOrdersMock.Return(nil)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11"},
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      720,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				storage.LoadPackByTypeMock.When(
					ctx,
					model.Box,
				).Then(
					nil,
					someError,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orders:    []model.OrderID{"11", "22"},
			wantSum:   0,
			wantError: someError,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			module := tc.module()
			resultSum, resultErr := module.GiveOrder(ctx, tc.orders)
			require.Equal(t, tc.wantSum, resultSum)
			require.Equal(t, tc.wantError, resultErr)
		})
	}
}

func TestModule_GiveOrderWithNewPackage(t *testing.T) {
	t.Parallel()

	var (
		mc       = minimock.NewController(t)
		ctx      = context.Background()
		tomorrow = time.Now().Add(24 * time.Hour)
		//previousDay = time.Now().Add(-24 * time.Hour)
		someError = errors.New("some error")
	)

	hashGenerator := mocks.NewHashGenerator(mc)
	hashGenerator.GenerateMock.Return("someHash")

	tests := []struct {
		desc      string
		module    func() Module
		orders    []model.OrderID
		pack      model.PackageType
		wantSum   int
		wantError error
	}{
		{
			desc: "Test case 1: correct order - finish with nil",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11", "22"},
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				storage.LoadPackByTypeMock.When(
					ctx,
					model.Box,
				).Then(
					&model.Package{
						Package:          "Box",
						PackageSurcharge: 20,
						PackageMaxWeight: 30000,
					},
					nil,
				)

				storage.UpdateOrdersMock.Return(nil)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11"},
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orders:    []model.OrderID{"11", "22"},
			pack:      model.Box,
			wantSum:   720,
			wantError: nil,
		},
		{
			desc: "Test case 2: some error in OrdersByGivenOrderIDList",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11", "22"},
					ctx,
					"someHash",
				).Then(
					nil,
					someError,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orders:    []model.OrderID{"11", "22"},
			pack:      model.Box,
			wantSum:   0,
			wantError: someError,
		},
		{
			desc: "Test case 3: not same client",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11", "22"},
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
						{
							OrderID:    "22",
							ClientID:   "2",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				storage.LoadPackByTypeMock.When(
					ctx,
					model.Box,
				).Then(
					&model.Package{
						Package:          "Box",
						PackageSurcharge: 20,
						PackageMaxWeight: 30000,
					},
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orders:    []model.OrderID{"11", "22"},
			pack:      model.Box,
			wantSum:   0,
			wantError: errors.New("не все заказы принадлежат одному клиенту"),
		},
		{
			desc: "Test case 4: someError in OrdersByGivenOrderIDListMock",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11", "22"},
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  tomorrow,
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				storage.UpdateOrdersMock.Return(nil)

				storage.OrdersByGivenOrderIDListMock.When(
					[]model.OrderID{"11"},
					ctx,
					"someHash",
				).Then(
					nil,
					someError,
				)

				storage.LoadPackByTypeMock.When(
					ctx,
					model.Box,
				).Then(
					&model.Package{
						Package:          "Box",
						PackageSurcharge: 20,
						PackageMaxWeight: 30000,
					},
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orders:    []model.OrderID{"11", "22"},
			pack:      model.Box,
			wantSum:   0,
			wantError: someError,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			module := tc.module()
			resultSum, resultErr := module.GiveOrderWithNewPackage(ctx, tc.orders, tc.pack)
			require.Equal(t, tc.wantSum, resultSum)
			require.Equal(t, tc.wantError, resultErr)
		})
	}
}

func TestModule_ReturnOrder(t *testing.T) {
	t.Parallel()

	var (
		mc        = minimock.NewController(t)
		ctx       = context.Background()
		someError = errors.New("some error")
	)

	hashGenerator := mocks.NewHashGenerator(mc)
	hashGenerator.GenerateMock.Return("someHash")

	tests := []struct {
		desc      string
		module    func() Module
		orderID   model.OrderID
		wantError error
	}{
		{
			desc: "Test case 1: correct order - finish with nil",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrderByIDMock.When(
					ctx,
					"11",
					"someHash",
				).Then(
					&model.Order{
						OrderID:    "11",
						ClientID:   "1",
						Condition:  model.ConditionAccepted,
						ArrivedAt:  time.Time{},
						ReceivedAt: time.Time{},
						Price:      700,
						Package:    model.Box,
						MaxWeight:  600,
					},
					nil,
				)

				storage.DeleteOrderMock.When(
					ctx,
					"11",
					"someHash",
				).Then(
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orderID:   "11",
			wantError: nil,
		},
		{
			desc: "Test case 2: OrderByIDMock finish with error",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrderByIDMock.When(
					ctx,
					"11",
					"someHash",
				).Then(
					nil,
					someError,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orderID:   "11",
			wantError: someError,
		},
		{
			desc: "Test case 3: OrderByIDMock finish with order = nil",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrderByIDMock.When(
					ctx,
					"11",
					"someHash",
				).Then(
					nil,
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orderID:   "11",
			wantError: errors.New("no such order"),
		},
		{
			desc: "Test case 4: order cannot be returned",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrderByIDMock.When(
					ctx,
					"11",
					"someHash",
				).Then(
					&model.Order{
						OrderID:    "11",
						ClientID:   "1",
						Condition:  model.ConditionRefund,
						ArrivedAt:  time.Time{},
						ReceivedAt: time.Time{},
						Price:      700,
						Package:    model.Box,
						MaxWeight:  600,
					},
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			orderID:   "11",
			wantError: errors.New("order cannot be returned"),
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			module := tc.module()
			result := module.ReturnOrder(ctx, tc.orderID)
			require.Equal(t, tc.wantError, result)
		})
	}
}

func TestModule_ListOrders(t *testing.T) {
	t.Parallel()

	var (
		mc        = minimock.NewController(t)
		ctx       = context.Background()
		someError = errors.New("some error")
	)

	hashGenerator := mocks.NewHashGenerator(mc)
	hashGenerator.GenerateMock.Return("someHash")

	tests := []struct {
		desc      string
		module    func() Module
		clientID  model.ClientID
		action    int
		wantList  []model.Order
		wantMap   map[model.PackageType]model.Package
		wantError error
	}{
		{
			desc: "Test case 1: everything correct, action == 1",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.LoadPackagesMock.When(
					ctx,
				).Then(
					[]model.Package{
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
					},
					nil,
				)

				storage.OrdersWithGivenClientIDMock.When(
					"1",
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionRefund,
							ArrivedAt:  time.Time{},
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			clientID: "1",
			action:   1,
			wantList: []model.Order{
				{
					OrderID:    "11",
					ClientID:   "1",
					Condition:  model.ConditionRefund,
					ArrivedAt:  time.Time{},
					ReceivedAt: time.Time{},
					Price:      700,
					Package:    model.Box,
					MaxWeight:  600,
				},
			},
			wantMap: map[model.PackageType]model.Package{
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
			wantError: nil,
		},
		{
			desc: "Test case 2: someError in OrdersWithGivenClientIDMock",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.LoadPackagesMock.When(
					ctx,
				).Then(
					[]model.Package{
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
					},
					nil,
				)

				storage.OrdersWithGivenClientIDMock.When(
					"1",
					ctx,
					"someHash",
				).Then(
					nil,
					someError,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			clientID:  "1",
			action:    1,
			wantList:  nil,
			wantMap:   nil,
			wantError: someError,
		},

		{
			desc: "Test case 3: everything correct, action == 2",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.LoadPackagesMock.When(
					ctx,
				).Then(
					[]model.Package{
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
					},
					nil,
				)

				storage.OrdersWithGivenClientIDAndAcceptedConditionMock.When(
					"1",
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionAccepted,
							ArrivedAt:  time.Time{},
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			clientID: "1",
			action:   2,
			wantList: []model.Order{
				{
					OrderID:    "11",
					ClientID:   "1",
					Condition:  model.ConditionAccepted,
					ArrivedAt:  time.Time{},
					ReceivedAt: time.Time{},
					Price:      700,
					Package:    model.Box,
					MaxWeight:  600,
				},
			},
			wantMap: map[model.PackageType]model.Package{
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
			wantError: nil,
		},

		{
			desc: "Test case 4: someError in OrdersWithGivenClientIDAndAcceptedConditionMock",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.LoadPackagesMock.When(
					ctx,
				).Then(
					[]model.Package{
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
					},
					nil,
				)

				storage.OrdersWithGivenClientIDAndAcceptedConditionMock.When(
					"1",
					ctx,
					"someHash",
				).Then(
					nil,
					someError,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			clientID:  "1",
			action:    2,
			wantList:  nil,
			wantMap:   nil,
			wantError: someError,
		},

		{
			desc: "Test case 5: unknown action",
			module: func() Module {
				return Module{
					hashGenerator: hashGenerator,
				}
			},
			clientID:  "1",
			action:    3,
			wantList:  nil,
			wantMap:   nil,
			wantError: errors.New("unknown action"),
		},
		{
			desc: "Test case 6: someError in LoadPackages",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.LoadPackagesMock.When(
					ctx,
				).Then(
					nil,
					someError,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			clientID:  "1",
			action:    2,
			wantList:  nil,
			wantMap:   nil,
			wantError: someError,
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			module := tc.module()
			resultList, resultMap, resultError := module.ListOrders(ctx, tc.clientID, tc.action)
			require.Equal(t, tc.wantList, resultList)
			require.Equal(t, tc.wantMap, resultMap)
			require.Equal(t, tc.wantError, resultError)
		})
	}
}

func TestModule_ReturnFromClient(t *testing.T) {
	t.Parallel()

	var (
		mc                  = minimock.NewController(t)
		ctx                 = context.Background()
		someError           = errors.New("some error")
		previousDay         = time.Now().Add(-24 * time.Hour)
		previousDayForError = time.Now().Add(-56 * time.Hour)
	)

	hashGenerator := mocks.NewHashGenerator(mc)
	hashGenerator.GenerateMock.Return("someHash")

	tests := []struct {
		desc      string
		module    func() Module
		clientID  model.ClientID
		orderID   model.OrderID
		wantError error
	}{
		{
			desc: "Test case 1: correct order - finish with nil",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrderByClientIDOrderIDMock.When(
					"11",
					"1",
					ctx,
					"someHash",
				).Then(
					&model.Order{
						OrderID:    "11",
						ClientID:   "1",
						Condition:  model.ConditionGiven,
						ArrivedAt:  time.Time{},
						ReceivedAt: previousDay,
						Price:      700,
						Package:    model.Box,
						MaxWeight:  600,
					},
					nil,
				)

				storage.UpdateOrdersMock.When(
					ctx,
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionRefund,
							ArrivedAt:  time.Time{},
							ReceivedAt: previousDay,
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					"someHash",
				).Then(
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			clientID:  "1",
			orderID:   "11",
			wantError: nil,
		},

		{
			desc: "Test case 2: someError in OrderByClientIDOrderID",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrderByClientIDOrderIDMock.When(
					"11",
					"1",
					ctx,
					"someHash",
				).Then(
					nil,
					someError,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			clientID:  "1",
			orderID:   "11",
			wantError: someError,
		},

		{
			desc: "Test case 3: заказ не был получен или уже возвращен",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrderByClientIDOrderIDMock.When(
					"11",
					"1",
					ctx,
					"someHash",
				).Then(
					&model.Order{
						OrderID:    "11",
						ClientID:   "1",
						Condition:  model.ConditionAccepted,
						ArrivedAt:  time.Time{},
						ReceivedAt: previousDay,
						Price:      700,
						Package:    model.Box,
						MaxWeight:  600,
					},
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			clientID:  "1",
			orderID:   "11",
			wantError: errors.New("заказ не был получен или уже возвращен"),
		},

		{
			desc: "Test case 4: время для возврата товара истекло",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrderByClientIDOrderIDMock.When(
					"11",
					"1",
					ctx,
					"someHash",
				).Then(
					&model.Order{
						OrderID:    "11",
						ClientID:   "1",
						Condition:  model.ConditionGiven,
						ArrivedAt:  time.Time{},
						ReceivedAt: previousDayForError,
						Price:      700,
						Package:    model.Box,
						MaxWeight:  600,
					},
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			clientID:  "1",
			orderID:   "11",
			wantError: errors.New("время для возврата товара истекло"),
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			module := tc.module()
			resultError := module.ReturnFromClient(ctx, tc.orderID, tc.clientID)
			require.Equal(t, tc.wantError, resultError)
		})
	}
}

func TestModule_ListReturns(t *testing.T) {
	t.Parallel()

	var (
		mc        = minimock.NewController(t)
		ctx       = context.Background()
		someError = errors.New("some error")
	)

	hashGenerator := mocks.NewHashGenerator(mc)
	hashGenerator.GenerateMock.Return("someHash")

	tests := []struct {
		desc      string
		module    func() Module
		clientID  model.ClientID
		orderID   model.OrderID
		wantList  []model.Order
		wantMap   map[model.PackageType]model.Package
		wantError error
	}{
		{
			desc: "Test case 1: correct order - finish with nil",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrdersWithRefundConditionMock.When(
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionRefund,
							ArrivedAt:  time.Time{},
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				storage.LoadPackagesMock.When(
					ctx,
				).Then(
					[]model.Package{
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
					},
					nil,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			clientID: "1",
			orderID:  "11",
			wantList: []model.Order{
				{
					OrderID:    "11",
					ClientID:   "1",
					Condition:  model.ConditionRefund,
					ArrivedAt:  time.Time{},
					ReceivedAt: time.Time{},
					Price:      700,
					Package:    model.Box,
					MaxWeight:  600,
				},
			},
			wantMap: map[model.PackageType]model.Package{
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
			wantError: nil,
		},

		{
			desc: "Test case 2: someError in OrdersWithRefundConditionMock",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrdersWithRefundConditionMock.When(
					ctx,
					"someHash",
				).Then(
					nil,
					someError,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			clientID:  "1",
			orderID:   "11",
			wantList:  nil,
			wantMap:   nil,
			wantError: someError,
		},

		{
			desc: "Test case 3: someError in LoadPackagesMock",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.OrdersWithRefundConditionMock.When(
					ctx,
					"someHash",
				).Then(
					[]model.Order{
						{
							OrderID:    "11",
							ClientID:   "1",
							Condition:  model.ConditionRefund,
							ArrivedAt:  time.Time{},
							ReceivedAt: time.Time{},
							Price:      700,
							Package:    model.Box,
							MaxWeight:  600,
						},
					},
					nil,
				)

				storage.LoadPackagesMock.When(
					ctx,
				).Then(
					nil,
					someError,
				)

				return Module{
					storage:       storage,
					hashGenerator: hashGenerator,
				}
			},
			clientID:  "1",
			orderID:   "11",
			wantList:  nil,
			wantMap:   nil,
			wantError: someError,
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			module := tc.module()
			resultList, resultMap, resultError := module.ListReturns(ctx)
			require.Equal(t, tc.wantList, resultList)
			require.Equal(t, tc.wantMap, resultMap)
			require.Equal(t, tc.wantError, resultError)
		})
	}
}

func TestModule_LoadPackagesToCheck(t *testing.T) {
	t.Parallel()

	var (
		mc        = minimock.NewController(t)
		ctx       = context.Background()
		someError = errors.New("some error")
	)

	tests := []struct {
		desc      string
		module    func() Module
		wantList  []model.Package
		wantError error
	}{
		{
			desc: "Test case 1: everything correct",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.LoadPackagesMock.When(
					ctx,
				).Then(
					[]model.Package{
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
					},
					nil,
				)

				return Module{
					storage: storage,
				}
			},
			wantList: []model.Package{
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
			},
			wantError: nil,
		},
		{
			desc: "Test case 2: someError in ",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.LoadPackagesMock.When(
					ctx,
				).Then(
					nil,
					someError,
				)

				return Module{
					storage: storage,
				}
			},
			wantList:  nil,
			wantError: someError,
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			module := tc.module()
			resultList, resultError := module.LoadPackagesToCheck(ctx)
			require.Equal(t, tc.wantList, resultList)
			require.Equal(t, tc.wantError, resultError)
		})
	}
}

func TestModule_CountFinalCost(t *testing.T) {
	t.Parallel()

	var (
		mc        = minimock.NewController(t)
		ctx       = context.Background()
		someError = errors.New("some error")
	)

	tests := []struct {
		desc      string
		module    func() Module
		orders    []model.Order
		wantSum   int
		wantError error
	}{
		{
			desc: "Test case 1: everything correct",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.LoadPackByTypeMock.When(
					ctx,
					model.Box,
				).Then(
					&model.Package{
						Package:          "Box",
						PackageSurcharge: 20,
						PackageMaxWeight: 30000,
					},
					nil,
				)

				return Module{
					storage: storage,
				}
			},
			orders: []model.Order{
				{
					OrderID:    "11",
					ClientID:   "1",
					Condition:  model.ConditionRefund,
					ArrivedAt:  time.Time{},
					ReceivedAt: time.Time{},
					Price:      700,
					Package:    model.Box,
					MaxWeight:  600,
				},
			},
			wantSum:   720,
			wantError: nil,
		},
		{
			desc: "Test case 1: everything correct",
			module: func() Module {
				storage := mocks.NewStorage(mc)

				storage.LoadPackByTypeMock.When(
					ctx,
					model.Box,
				).Then(
					nil,
					someError,
				)

				return Module{
					storage: storage,
				}
			},
			orders: []model.Order{
				{
					OrderID:    "11",
					ClientID:   "1",
					Condition:  model.ConditionRefund,
					ArrivedAt:  time.Time{},
					ReceivedAt: time.Time{},
					Price:      700,
					Package:    model.Box,
					MaxWeight:  600,
				},
			},
			wantSum:   0,
			wantError: someError,
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			module := tc.module()
			resultSum, resultError := module.countFinalCost(ctx, tc.orders)
			require.Equal(t, tc.wantSum, resultSum)
			require.Equal(t, tc.wantError, resultError)
		})
	}
}

func TestModule_SameClient(t *testing.T) {
	t.Parallel()

	module := Module{}

	tests := []struct {
		desc       string
		orders     []model.Order
		wantResult bool
	}{
		{
			desc: "Test case 1: all orders have same clientID",
			orders: []model.Order{
				{
					ClientID: "1",
				},
				{
					ClientID: "1",
				},
				{
					ClientID: "1",
				},
			},
			wantResult: true,
		},

		{
			desc: "Test case 2: all orders have not same clientID",
			orders: []model.Order{
				{
					ClientID: "2",
				},
				{
					ClientID: "1",
				},
				{
					ClientID: "1",
				},
			},
			wantResult: false,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			result := module.sameClient(tc.orders)
			require.Equal(t, tc.wantResult, result)
		})
	}
}
