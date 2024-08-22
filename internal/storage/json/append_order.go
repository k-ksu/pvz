package json

import (
	"encoding/json"
	"os"

	"HomeWork_1/internal/model"
)

func (s *Storage) AppendOrder(order model.Order, hash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkFileExistence(s.fileOrders); err != nil {
		return err
	}

	b, err := os.ReadFile(s.fileOrders)
	if err != nil {
		return nil
	}

	var orderBook model.OrderBook

	if err := json.Unmarshal(b, &orderBook); err != nil {
		return err
	}

	orderBook.Orders = append(orderBook.Orders, order)
	orderBook.Hash = hash

	bWrite, err := json.MarshalIndent(orderBook, "  ", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.fileOrders, bWrite, RWPermissio)
}
