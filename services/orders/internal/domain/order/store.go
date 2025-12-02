package order

import (
	"sync"
)

type Store struct {
	mu     sync.Mutex
	orders map[string]Status
}

func NewStore() *Store {
	return &Store{orders: make(map[string]Status)}
}

func (s *Store) SetStatus(orderID string, status Status) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[orderID] = status
}

func (s *Store) GetStatus(orderID string) Status {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.orders[orderID]
}
