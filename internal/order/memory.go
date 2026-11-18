
package order

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Store interface {
	Create(o *Order) error
	All() ([]Order, error)
	ByID(id string) (Order, error)
}

type MemoryStore struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewMemoryStore() Store {
	return &MemoryStore{
		orders: make(map[string]*Order),
	}
}

func newID() string {
	b := make([]byte, 2)
	_, _ = rand.Read(b)
	return fmt.Sprintf("o-%s", hex.EncodeToString(b))
}

func (s *MemoryStore) Create(o *Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	o.ID = newID()
	o.PlacedAt = time.Now().UTC()

	s.orders[o.ID] = o
	return nil
}

func (s *MemoryStore) All() ([]Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]Order, 0, len(s.orders))

	for _, o := range s.orders {
		orders = append(orders, *o)
	}

	return orders, nil
}

func (s *MemoryStore) ByID(id string) (Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	o, exists := s.orders[id]
	if !exists {
		return Order{}, ErrOrderNotFound
	}

	return *o, nil
}
