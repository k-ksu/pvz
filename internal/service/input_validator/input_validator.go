package input_validator

import (
	"HomeWork_1/internal/model/errs"
	"errors"
	"fmt"
	"strings"
	"time"

	"HomeWork_1/internal/model"
)

type Validator struct{}

func NewValidator() Validator {
	return Validator{}
}

func (v Validator) ValidateGetOrderFromCourier(orderID, clientID, date, pack string, price, weight int) (*model.Order, error) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println(err)
	}

	dateTime, err := time.ParseInLocation("02.01.2006-15:04", date, location)
	if err != nil {
		return nil, errs.ErrDateTimeIncorrectType
	}

	if len(date) == 0 {
		return nil, errs.ErrEmptyDate
	}

	if len(orderID) == 0 {
		return nil, errs.ErrEmptyOrderID
	}

	if len(clientID) == 0 {
		return nil, errs.ErrEmptyClientID
	}

	if price < 1 {
		return nil, errs.ErrSmallPrice
	}

	if weight < 1 {
		return nil, errs.ErrSmallWeight
	}

	if len(pack) == 0 {
		order := model.Order{
			OrderID:   model.OrderID(orderID),
			ClientID:  model.ClientID(clientID),
			ArrivedAt: dateTime,
			Package:   model.WithoutPackage,
		}
		return &order, errs.ErrPackageDoesNotSet
	}

	order := model.Order{
		OrderID:   model.OrderID(orderID),
		ClientID:  model.ClientID(clientID),
		ArrivedAt: dateTime,
		Package:   model.PackageType(pack),
	}
	return &order, nil
}

func (v Validator) ValidateGiveOrder(orders, pack string, loadedPackages []model.Package) ([]model.OrderID, *model.PackageType, error) {
	if len(orders) == 0 {
		return nil, nil, errs.ErrEmptyListOrders
	}

	listOfOrdersStrings := strings.Split(orders, ",")

	listOfOrders := make([]model.OrderID, 0, len(listOfOrdersStrings))
	for _, id := range listOfOrdersStrings {
		listOfOrders = append(listOfOrders, model.OrderID(id))
	}

	if len(pack) == 0 {
		return listOfOrders, nil, errs.ErrPackageDoesNotSet
	}

	readyPack := model.PackageType(pack)
	for _, packi := range loadedPackages {
		if packi.Package == readyPack {
			return listOfOrders, &readyPack, nil
		}
	}
	return listOfOrders, nil, errs.ErrIncorrectPackageType
}

func (v Validator) ValidateReturnOrder(orderID string) (model.OrderID, error) {
	if len(orderID) == 0 {
		return "", errs.ErrEmptyOrderID
	}
	return model.OrderID(orderID), nil
}

func (v Validator) ValidateListOrders(clientID, action string) (model.ClientID, error) {
	if len(clientID) == 0 {
		return "", errs.ErrEmptyClientID
	}
	if !(action == model.AllOrdersStr || action == model.ActualOrdersStr) {
		return "", errs.ErrIncorrectAction
	}

	return model.ClientID(clientID), nil
}

func (v Validator) ValidateReturnFromClient(orderID, clientID string) (model.OrderID, model.ClientID, error) {
	if len(orderID) == 0 {
		return "", "", errs.ErrEmptyOrderID
	}

	if len(clientID) == 0 {
		return "", "", errs.ErrEmptyClientID
	}
	return model.OrderID(orderID), model.ClientID(clientID), nil
}

func (v Validator) ValidateListReturns(pageSize, pageNumber int) (int, int, error) {
	if pageSize < 1 {
		return 0, 0, errs.ErrNegativePageSize
	}

	if pageNumber < 1 {
		return 0, 0, errs.ErrNegativePageNumber
	}

	return pageSize, pageNumber, nil
}

func (v Validator) ValidatePackage(weight int, pack model.PackageType, loadedPackages []model.Package) error {
	for _, packi := range loadedPackages {
		if packi.Package != pack {
			continue
		}
		if weight >= packi.PackageMaxWeight && packi.PackageMaxWeight != -1 {
			return errors.New("вес товара превыщает допустимые параметры для упаковки типа " + string(packi.Package))
		}
		return nil
	}

	return errs.ErrIncorrectPackageType
}
