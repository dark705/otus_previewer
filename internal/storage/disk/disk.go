package disk

import (
	"fmt"
	"github.com/dark705/otus_previewer/internal/helpers"
	errorsPack "github.com/pkg/errors"
	"io/ioutil"
	"os"
)

type Disk struct {
	path string
}

var ignoreFiles = []string{".gitkeep"}

func New(path string) Disk {
	return Disk{path: path + "/"}
}

func (s *Disk) Add(id string, content []byte) error {
	err := ioutil.WriteFile(s.path+id, content, 0664)
	if err != nil {
		return errorsPack.Wrap(err, fmt.Sprintf("Fail on Add, content with id: %s already exist", id))
	}
	return nil
}

func (s *Disk) Del(id string) error {
	err := os.Remove(s.path + id)
	if err != nil {
		return errorsPack.Wrap(err, fmt.Sprintf("Fail on Del, content with id: %s not exist", id))
	}
	return nil
}

func (s *Disk) Get(id string) ([]byte, error) {
	content, err := ioutil.ReadFile(s.path + id)
	if err != nil {
		return nil, errorsPack.Wrap(err, fmt.Sprintf("Fail on Get, content with id: %s", id))
	}
	return content, nil
}

func (s *Disk) GetListSize() map[string]int {
	files, err := ioutil.ReadDir(s.path)
	helpers.FailOnError(err, "Cant read cache dir")
	usage := make(map[string]int)

	for _, f := range files {
		//skip not content files
		var skip bool
		for _, igF := range ignoreFiles {
			if f.Name() == igF {
				skip = true
				break
			}
		}
		if !skip {
			usage[f.Name()] = int(f.Size())
		}
	}
	return usage
}
