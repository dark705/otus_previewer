package dispatcher

import (
	"errors"
	"fmt"
	"github.com/dark705/otus_previewer/internal/storage"
	"github.com/sirupsen/logrus"
	"sort"
	"time"
)

type StorageDispatcher struct {
	storage storage.Storage
	lastUse map[string]time.Time
	usage   int
	limit   int
	logger  *logrus.Logger
}

func New(st storage.Storage, lim int, l *logrus.Logger) StorageDispatcher {
	lu := make(map[string]time.Time)
	for _, id := range st.GetUniqId() {
		lu[id] = time.Now()
	}

	return StorageDispatcher{
		storage: st,
		lastUse: lu,
		usage:   st.Usage(),
		limit:   lim,
		logger:  l,
	}
}

func (sd *StorageDispatcher) Usage() int {
	return sd.usage
}

func (sd *StorageDispatcher) Exist(id string) bool {
	_, exist := sd.lastUse[id]
	return exist
}

func (sd *StorageDispatcher) Get(id string) ([]byte, error) {
	_, exist := sd.lastUse[id]
	if !exist {
		return nil, errors.New(fmt.Sprintf("Fail on update lastUse, content with id: %s not exist", id))
	}
	sd.lastUse[id] = time.Now()
	return sd.storage.Get(id)
}

func (sd *StorageDispatcher) Add(id string, content []byte) error {
	//storage not full
	if sd.usage+len(content) < sd.limit {
		return sd.addAvailable(id, content)
	}
	//storage is full, need to clean,
	sd.logger.Debugln(fmt.Sprintf("Storage is full, usage: %d need to clean", sd.Usage()))
	err := sd.cleanOldUseOn(len(content))
	if err != nil {
		return err
	}
	//now we can add
	return sd.addAvailable(id, content)
}

func (sd *StorageDispatcher) addAvailable(id string, content []byte) error {
	sd.usage += len(content)
	sd.lastUse[id] = time.Now()
	sd.logger.Debugln(fmt.Sprintf("Storage not full, add content with id: %s, storage usage: %d", id, sd.Usage()))
	return sd.storage.Add(id, content)
}

func (sd *StorageDispatcher) cleanOldUseOn(newContentLen int) error {
	//make sort, smaller index is more old at time
	list := make([]struct {
		i string
		t time.Time
	}, 0, len(sd.lastUse))

	for id, tm := range sd.lastUse {
		list = append(list, struct {
			i string
			t time.Time
		}{i: id, t: tm})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].t.Before(list[j].t)
	})

	for _, val := range list {
		con, err := sd.storage.Get(val.i)
		if err != nil {
			return err
		}
		err = sd.storage.Del(val.i)
		if err != nil {
			return err
		}
		delete(sd.lastUse, val.i)
		sd.usage -= len(con)
		sd.logger.Debugln(fmt.Sprintf("Deleted content with id: %s, storage usage: %d", val.i, sd.usage))
		if sd.usage+newContentLen < sd.limit {
			break
		}
	}
	return nil
}
