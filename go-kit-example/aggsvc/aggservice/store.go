package aggservice

import (
	"fmt"

	"github.com/leehaowei/tolling-micro-service/types"
)

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type MemoryStore struct {
	data map[int]float64 //map[OBUID]distance
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(d types.Distance) error {
	m.data[d.OBUID] += d.Value
	return nil
}

func (m *MemoryStore) Get(id int) (float64, error) {
	dist, ok := m.data[id]
	if !ok {
		return 0.0, fmt.Errorf("could not find obu id %d", id)
	}
	return dist, nil
}
