package disk

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dark705/otus_previewer/internal/helpers"
	errorsPack "github.com/pkg/errors"
)

type Disk struct {
	path string
}

var ignoreFiles = []string{".gitkeep"}

func New(path string) Disk {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0775)
		helpers.FailOnError(err, "fail to create cache directory")
	}

	return Disk{path: path + "/"}
}

func (storage *Disk) Add(id string, content []byte) error {
	err := ioutil.WriteFile(storage.path+id, content, 0664)
	if err != nil {
		return errorsPack.Wrap(err, fmt.Sprintf("fail on Add, content with id: %s already exist", id))
	}
	return nil
}

func (storage *Disk) Del(id string) error {
	err := os.Remove(storage.path + id)
	if err != nil {
		return errorsPack.Wrap(err, fmt.Sprintf("fail on Del, content with id: %s not exist", id))
	}
	return nil
}

func (storage *Disk) Get(id string) ([]byte, error) {
	content, err := ioutil.ReadFile(storage.path + id)
	if err != nil {
		return nil, errorsPack.Wrap(err, fmt.Sprintf("fail on Get, content with id: %s", id))
	}
	return content, nil
}

func (storage *Disk) GetListSize() map[string]int {
	files, err := ioutil.ReadDir(storage.path)
	helpers.FailOnError(err, "cant read cache dir")
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
