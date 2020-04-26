package inmemory

import (
	"fmt"
)

type InMemory struct {
	storage map[string][]byte
}

func New() InMemory {
	return InMemory{storage: map[string][]byte{}}
}

func (s *InMemory) Add(id string, content []byte) error {
	_, exist := s.storage[id]
	if exist {
		return fmt.Errorf("Fail on Add, content with id: %s already exist", id)
	}
	s.storage[id] = content
	return nil
}

func (s *InMemory) Del(id string) error {
	_, exist := s.storage[id]
	if !exist {
		return fmt.Errorf("Fail on Del, content with id: %s not exist", id)
	}
	delete(s.storage, id)
	return nil
}

func (s *InMemory) Get(id string) ([]byte, error) {
	_, exist := s.storage[id]
	if !exist {
		return nil, fmt.Errorf("Fail on Get, content with id: %s not exist", id)
	}
	return s.storage[id], nil
}

func (s *InMemory) GetListSize() map[string]int {
	usage := make(map[string]int, len(s.storage))
	for id, c := range s.storage {
		usage[id] = len(c)
	}
	return usage
}
