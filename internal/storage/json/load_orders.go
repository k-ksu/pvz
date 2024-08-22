package json

import (
	"encoding/json"
	"os"

	"HomeWork_1/internal/model"
)

func (s *Storage) LoadOrders(hash string) ([]model.Order, error) {
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

	if err := os.WriteFile(s.fileOrders, bWrite, 0666); err != nil {
		return nil, err
	}

	return orderBook.Orders, nil
}
