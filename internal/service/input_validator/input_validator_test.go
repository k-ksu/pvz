package input_validator

import (
	"HomeWork_1/internal/model/errs"
	"errors"
	"testing"
	"time"

	"HomeWork_1/internal/model"

	"github.com/stretchr/testify/require"
)

func TestValidator_ValidateGetOrderFromCourier(t *testing.T) {
	t.Parallel()

	validator := Validator{}
	loc, _ := time.LoadLocation("Europe/Moscow")

	tests := []struct {
		desc         string
		orderID      string
		clientID     string
		date         string
		pack         string
		price        int
		weight       int
		wantOrderID  model.OrderID
		wantClientID model.ClientID
		wantTime     time.Time
		wantPack     model.PackageType
		wantError    error
	}{
		{
			desc:         "Test case 1: everything is correct",
			orderID:      "1",
			clientID:     "1",
			date:         "25.06.2025-14:48",
			pack:         "film",
			price:        100,
			weight:       100,
			wantOrderID:  model.OrderID("1"),
			wantClientID: model.ClientID("1"),
			wantTime:     time.Date(2025, 6, 25, 14, 48, 0, 0, loc),
			wantPack:     model.PackageType("film"),
			wantError:    nil,
		},
		{
			desc:         "Test case 2: date of order is empty",
			orderID:      "1",
			clientID:     "1",
			date:         "01.01",
			pack:         "film",
			price:        100,
			weight:       100,
			wantOrderID:  "",
			wantClientID: "",
			wantTime:     time.Time{},
			wantPack:     "",
			wantError:    errors.New("некорректный тип срока хранения"),
		},
		{
			desc:         "Test case 3: orderID of order is empty",
			orderID:      "",
			clientID:     "1",
			date:         "25.06.2025-14:48",
			pack:         "film",
			price:        100,
			weight:       100,
			wantOrderID:  "",
			wantClientID: "",
			wantTime:     time.Time{},
			wantPack:     "",
			wantError:    errors.New("ID of order is empty"),
		},
		{
			desc:         "Test case 4: ID of client is empty",
			orderID:      "1",
			clientID:     "",
			date:         "25.06.2025-14:48",
			pack:         "film",
			price:        100,
			weight:       100,
			wantOrderID:  "",
			wantClientID: "",
			wantTime:     time.Time{},
			wantPack:     "",
			wantError:    errors.New("ID of client is empty"),
		},
		{
			desc:         "Test case 5: price is too small",
			orderID:      "1",
			clientID:     "1",
			date:         "25.06.2025-14:48",
			pack:         "film",
			price:        0,
			weight:       100,
			wantOrderID:  "",
			wantClientID: "",
			wantTime:     time.Time{},
			wantPack:     "",
			wantError:    errors.New("price is too small"),
		},
		{
			desc:         "Test case 6: weight is too small",
			orderID:      "1",
			clientID:     "1",
			date:         "25.06.2025-14:48",
			pack:         "film",
			price:        100,
			weight:       0,
			wantOrderID:  "",
			wantClientID: "",
			wantTime:     time.Time{},
			wantPack:     "",
			wantError:    errors.New("weight is too small"),
		},
		{
			desc:         "Test case 7: pack does not set",
			orderID:      "1",
			clientID:     "1",
			date:         "25.06.2025-14:48",
			pack:         "",
			price:        100,
			weight:       100,
			wantOrderID:  model.OrderID("1"),
			wantClientID: model.ClientID("1"),
			wantTime:     time.Date(2025, 6, 25, 14, 48, 0, 0, loc),
			wantPack:     model.WithoutPackage,
			wantError:    errs.ErrPackageDoesNotSet,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			o, ErrorResult := validator.ValidateGetOrderFromCourier(tc.orderID, tc.clientID, tc.date, tc.pack, tc.price, tc.weight)
			if o == nil {
				var w *model.Order
				require.Equal(t, o, w)
			} else {
				require.Equal(t, tc.wantOrderID, o.OrderID)
				require.Equal(t, tc.wantClientID, o.ClientID)
				require.Equal(t, tc.wantTime, o.ArrivedAt)
				require.Equal(t, tc.wantPack, o.Package)
			}

			require.Equal(t, tc.wantError, ErrorResult)
		})
	}
}

func TestValidator_ValidateGiveOrder(t *testing.T) {
	t.Parallel()

	validator := Validator{}

	tests := []struct {
		desc          string
		orders        string
		pack          string
		loadedPackage []model.Package
		wantPack      *model.PackageType
		wantList      []model.OrderID
		wantError     error
	}{
		{
			desc:   "Test case 1: everything is correct",
			orders: "1,2,3",
			pack:   "box",
			loadedPackage: []model.Package{
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
			wantPack:  toPtr(model.Box),
			wantList:  []model.OrderID{model.OrderID("1"), model.OrderID("2"), model.OrderID("3")},
			wantError: nil,
		},
		{
			desc:   "Test case 2: list of orders is empty",
			orders: "",
			pack:   "box",
			loadedPackage: []model.Package{
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
			wantPack:  nil,
			wantList:  nil,
			wantError: errors.New("list of orders is empty"),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			listResult, packResult, errorResult := validator.ValidateGiveOrder(tc.orders, tc.pack, tc.loadedPackage)
			require.Equal(t, tc.wantList, listResult)
			require.Equal(t, tc.wantPack, packResult)
			require.Equal(t, tc.wantError, errorResult)
		})
	}
}

func TestValidator_ValidateReturnOrder(t *testing.T) {
	t.Parallel()

	validator := Validator{}

	tests := []struct {
		desc      string
		order     string
		wantOrder model.OrderID
		wantError error
	}{
		{
			desc:      "Test case 1: everything is correct",
			order:     "1",
			wantOrder: model.OrderID("1"),
			wantError: nil,
		},
		{
			desc:      "Test case 2: ID of order is empty",
			order:     "",
			wantOrder: "",
			wantError: errors.New("ID of order is empty"),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			orderResult, errorResult := validator.ValidateReturnOrder(tc.order)
			require.Equal(t, tc.wantOrder, orderResult)
			require.Equal(t, tc.wantError, errorResult)
		})
	}
}

func TestValidator_ValidateListOrders(t *testing.T) {
	t.Parallel()

	validator := Validator{}

	tests := []struct {
		desc       string
		clientID   string
		action     string
		wantClient model.ClientID
		wantError  error
	}{
		{
			desc:       "Test case 1: everything is correct and action == 1",
			clientID:   "1",
			action:     "1",
			wantClient: model.ClientID("1"),
			wantError:  nil,
		},
		{
			desc:       "Test case 2: everything is correct and action == 2",
			clientID:   "1",
			action:     "2",
			wantClient: model.ClientID("1"),
			wantError:  nil,
		},
		{
			desc:       "Test case 3: clientID is empty",
			clientID:   "",
			action:     "2",
			wantClient: "",
			wantError:  errors.New("clientID is empty"),
		},
		{
			desc:       "Test case 4: action is not correct",
			clientID:   "1",
			action:     "0",
			wantClient: "",
			wantError:  errors.New("action is not correct"),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			clientResult, errorResult := validator.ValidateListOrders(tc.clientID, tc.action)
			require.Equal(t, tc.wantClient, clientResult)
			require.Equal(t, tc.wantError, errorResult)
		})
	}
}

func TestValidator_ValidateReturnFromClient(t *testing.T) {
	t.Parallel()

	validator := Validator{}

	tests := []struct {
		desc         string
		orderID      string
		clientID     string
		wantOrderID  model.OrderID
		wantClientID model.ClientID
		wantError    error
	}{
		{
			desc:         "Test case 1: everything is correct",
			orderID:      "1",
			clientID:     "1",
			wantOrderID:  model.OrderID("1"),
			wantClientID: model.ClientID("1"),
			wantError:    nil,
		},
		{
			desc:         "Test case 2: ID of order is empty",
			orderID:      "",
			clientID:     "1",
			wantOrderID:  "",
			wantClientID: "",
			wantError:    errors.New("ID of order is empty"),
		},
		{
			desc:         "Test case 3: ID of client is empty",
			orderID:      "1",
			clientID:     "",
			wantOrderID:  "",
			wantClientID: "",
			wantError:    errors.New("ID of client is empty"),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			orderIDResult, clientIDResult, errorResult := validator.ValidateReturnFromClient(tc.orderID, tc.clientID)
			require.Equal(t, tc.wantOrderID, orderIDResult)
			require.Equal(t, tc.wantClientID, clientIDResult)
			require.Equal(t, tc.wantError, errorResult)
		})
	}
}

func TestValidator_ValidateListReturns(t *testing.T) {
	t.Parallel()

	validator := Validator{}

	tests := []struct {
		desc           string
		pageSize       int
		pageNumber     int
		wantPageSize   int
		wantPageNumber int
		wantError      error
	}{
		{
			desc:           "Test case 1: everything is correct",
			pageSize:       2,
			pageNumber:     2,
			wantPageSize:   2,
			wantPageNumber: 2,
			wantError:      nil,
		},
		{
			desc:           "Test case 2: pageSize less then 1",
			pageSize:       0,
			pageNumber:     2,
			wantPageSize:   0,
			wantPageNumber: 0,
			wantError:      errors.New("pageSize cannot be negative"),
		},
		{
			desc:           "Test case 3: pageNumber less then 1",
			pageSize:       2,
			pageNumber:     0,
			wantPageSize:   0,
			wantPageNumber: 0,
			wantError:      errors.New("pageNumber cannot be negative"),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			pageSizeResult, pageNumberResult, errorResult := validator.ValidateListReturns(tc.pageSize, tc.pageNumber)
			require.Equal(t, tc.wantPageSize, pageSizeResult)
			require.Equal(t, tc.wantPageNumber, pageNumberResult)
			require.Equal(t, tc.wantError, errorResult)
		})
	}
}

func TestValidator_ValidatePackage(t *testing.T) {
	t.Parallel()

	validator := Validator{}

	tests := []struct {
		desc           string
		weight         int
		pack           model.PackageType
		loadedPackages []model.Package
		wantError      error
	}{
		{
			desc:   "Test case 1: everything is correct, choose box",
			weight: 100,
			pack:   model.Box,
			loadedPackages: []model.Package{
				{
					Package:          model.WithoutPackage,
					PackageSurcharge: 0,
					PackageMaxWeight: -1,
				},
				{
					Package:          model.PlasticBag,
					PackageSurcharge: 5,
					PackageMaxWeight: 10000,
				},
				{
					Package:          model.Box,
					PackageSurcharge: 20,
					PackageMaxWeight: 30000,
				},
				{
					Package:          model.Film,
					PackageSurcharge: 1,
					PackageMaxWeight: -1,
				},
			},
			wantError: nil,
		},
		{
			desc:   "Test case 2: pack is not correct",
			weight: 100,
			pack:   model.PackageType("Something"),
			loadedPackages: []model.Package{
				{
					Package:          model.WithoutPackage,
					PackageSurcharge: 0,
					PackageMaxWeight: -1,
				},
				{
					Package:          model.PlasticBag,
					PackageSurcharge: 5,
					PackageMaxWeight: 10000,
				},
				{
					Package:          model.Box,
					PackageSurcharge: 20,
					PackageMaxWeight: 30000,
				},
				{
					Package:          model.Film,
					PackageSurcharge: 1,
					PackageMaxWeight: -1,
				},
			},
			wantError: errors.New("указан некорректный тип упаковки"),
		},
		{
			desc:   "Test case 3: weight is too big for box",
			weight: 10000000,
			pack:   model.Box,
			loadedPackages: []model.Package{
				{
					Package:          model.WithoutPackage,
					PackageSurcharge: 0,
					PackageMaxWeight: -1,
				},
				{
					Package:          model.PlasticBag,
					PackageSurcharge: 5,
					PackageMaxWeight: 10000,
				},
				{
					Package:          model.Box,
					PackageSurcharge: 20,
					PackageMaxWeight: 30000,
				},
				{
					Package:          model.Film,
					PackageSurcharge: 1,
					PackageMaxWeight: -1,
				},
			},
			wantError: errors.New("вес товара превыщает допустимые параметры для упаковки типа box"),
		},
		{
			desc:   "Test case 4: weight is too big for plasticBag",
			weight: 10000000,
			pack:   model.PlasticBag,
			loadedPackages: []model.Package{
				{
					Package:          model.WithoutPackage,
					PackageSurcharge: 0,
					PackageMaxWeight: -1,
				},
				{
					Package:          model.PlasticBag,
					PackageSurcharge: 5,
					PackageMaxWeight: 10000,
				},
				{
					Package:          model.Box,
					PackageSurcharge: 20,
					PackageMaxWeight: 30000,
				},
				{
					Package:          model.Film,
					PackageSurcharge: 1,
					PackageMaxWeight: -1,
				},
			},
			wantError: errors.New("вес товара превыщает допустимые параметры для упаковки типа plasticBag"),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			errorResult := validator.ValidatePackage(tc.weight, tc.pack, tc.loadedPackages)
			require.Equal(t, tc.wantError, errorResult)
		})
	}
}

func toPtr[T any](t T) *T {
	return &t
}
