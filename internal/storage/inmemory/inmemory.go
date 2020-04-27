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

func (storage *InMemory) Add(id string, content []byte) error {
	_, exist := storage.storage[id]
	if exist {
		return fmt.Errorf("fail on Add, content with id: %storage already exist", id)
	}
	storage.storage[id] = content
	return nil
}

func (storage *InMemory) Del(id string) error {
	_, exist := storage.storage[id]
	if !exist {
		return fmt.Errorf("fail on Del, content with id: %storage not exist", id)
	}
	delete(storage.storage, id)
	return nil
}

func (storage *InMemory) Get(id string) ([]byte, error) {
	_, exist := storage.storage[id]
	if !exist {
		return nil, fmt.Errorf("fail on Get, content with id: %storage not exist", id)
	}
	return storage.storage[id], nil
}

func (storage *InMemory) GetListSize() map[string]int {
	usage := make(map[string]int, len(storage.storage))
	for id, c := range storage.storage {
		usage[id] = len(c)
	}
	return usage
}
