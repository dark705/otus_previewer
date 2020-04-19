package storage

type Storage interface {
	Add(id string, content []byte) error
	Del(id string) error
	Get(id string) ([]byte, error)
	GetListSize() map[string]int
}
