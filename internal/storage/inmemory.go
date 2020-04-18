package storage

import (
	"errors"
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
		return errors.New(fmt.Sprintf("Fail on Add, content with id: %s already exist", id))
	}
	s.storage[id] = content
	return nil
}

func (s *InMemory) Del(id string) error {
	_, exist := s.storage[id]
	if !exist {
		return errors.New(fmt.Sprintf("Fail on Del, content with id: %s not exist", id))
	}
	delete(s.storage, id)
	return nil
}

func (s *InMemory) Get(id string) ([]byte, error) {
	_, exist := s.storage[id]
	if !exist {
		return nil, errors.New(fmt.Sprintf("Fail on Get, content with id: %s not exist", id))
	}
	return s.storage[id], nil
}

func (s *InMemory) Usage() int {
	var usage int
	for _, c := range s.storage {
		usage += len(c)
	}
	return usage
}

func (s *InMemory) GetUniqId() []string {
	uniqId := make([]string, 0, len(s.storage))
	for id := range s.storage {
		uniqId = append(uniqId, id)
	}
	return uniqId
}
