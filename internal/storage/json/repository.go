package json

import (
	"sync"
)

const RWPermission = 0666

type Storage struct {
	fileOrders string
	mu         sync.Mutex
}

func NewStorage(fileOrders string) Storage {
	return Storage{
		fileOrders: fileOrders,
	}
}
