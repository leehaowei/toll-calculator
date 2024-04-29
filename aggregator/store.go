package main

import "github.com/leehaowei/tolling-micro-service/types"

type MemoryStore struct {
	data map[int]float64  //map[OBUID]distance
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
