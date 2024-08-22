package json

import (
	"encoding/json"
	"os"

	"HomeWork_1/internal/model"
)

func (s *Storage) OrderByID(orderID model.OrderID, hash string) (*model.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkFileExistence(s.fileOrders); err != nil {
		return nil, err
	}

	b, err := os.ReadFile(s.fileOrders)
	if err != nil {
		return nil, err
	}

	var orderBook model.OrderBook
	if err := json.Unmarshal(b, &orderBook); err != nil {
		return nil, err
	}

	orderBook.Hash = hash

	bWrite, err := json.MarshalIndent(orderBook, "  ", "  ")
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(s.fileOrders, bWrite, 0666)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(orderBook.Orders); i++ {
		if orderID == orderBook.Orders[i].OrderID {
			return &orderBook.Orders[i], nil
		}
	}

	return nil, nil
}
