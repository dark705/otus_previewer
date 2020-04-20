package dispatcher

import (
	"errors"
	"fmt"
	"github.com/dark705/otus_previewer/internal/storage"
	"github.com/sirupsen/logrus"
	"sort"
	"time"
)

type contentInfo struct {
	size    int
	lastUse time.Time
}

type StorageDispatcher struct {
	storage          storage.Storage
	contentState     map[string]contentInfo
	totalContentSize int
	limit            int
	logger           *logrus.Logger
}

func New(storage storage.Storage, limit int, logger *logrus.Logger) StorageDispatcher {
	state := make(map[string]contentInfo)
	var totalSize int
	for id, size := range storage.GetListSize() {
		state[id] = contentInfo{size: size, lastUse: time.Now()}
		totalSize += size
	}

	return StorageDispatcher{
		storage:          storage,
		contentState:     state,
		totalContentSize: totalSize,
		limit:            limit,
		logger:           logger,
	}
}

func (sd *StorageDispatcher) TotalContentSize() int {
	return sd.totalContentSize
}

func (sd *StorageDispatcher) Exist(id string) bool {
	_, exist := sd.contentState[id]
	return exist
}

func (sd *StorageDispatcher) Get(id string) ([]byte, error) {
	_, exist := sd.contentState[id]
	if !exist {
		return nil, errors.New(fmt.Sprintf("Fail on update lastUse on get, content with id: %s not exist", id))
	}
	sd.contentState[id] = contentInfo{size: sd.contentState[id].size, lastUse: time.Now()}
	sd.logger.Debugln(fmt.Sprintf("Content with id: %s updated, last use time: %s", id, time.Now()))
	return sd.storage.Get(id)
}

func (sd *StorageDispatcher) Add(id string, content []byte) error {
	//storage not full
	if sd.totalContentSize+len(content) < sd.limit {
		return sd.addAvailable(id, content)
	}
	//storage is full, need to clean,
	sd.logger.Debugln(fmt.Sprintf("Storage is full, totalContentSize: %d need clean", sd.TotalContentSize()))
	err := sd.cleanOldUseContentOn(len(content))
	if err != nil {
		return err
	}
	//now we can add
	return sd.addAvailable(id, content)
}

func (sd *StorageDispatcher) addAvailable(id string, content []byte) error {
	err := sd.storage.Add(id, content)
	if err != nil {
		return err
	}
	sd.totalContentSize += len(content)
	sd.contentState[id] = contentInfo{size: len(content), lastUse: time.Now()}
	sd.logger.Debugln(fmt.Sprintf("Storage not full, add content with id: %s, size: %d, now total content size: %d", id, len(content), sd.TotalContentSize()))
	return nil
}

func (sd *StorageDispatcher) cleanOldUseContentOn(needCleanBytes int) error {
	//make sort, smaller index is more old at time
	list := make([]struct {
		i string
		t time.Time
	}, 0, len(sd.contentState))

	for id, ci := range sd.contentState {
		list = append(list, struct {
			i string
			t time.Time
		}{i: id, t: ci.lastUse})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].t.Before(list[j].t)
	})
	//delete old used content until free space for new content
	for _, val := range list {
		err := sd.storage.Del(val.i)
		if err != nil {
			return err
		}
		usedSize := sd.contentState[val.i].size
		sd.totalContentSize -= usedSize
		delete(sd.contentState, val.i)
		sd.logger.Debugln(fmt.Sprintf("Deleted content with id: %s, used size: %d", val.i, usedSize))
		if sd.totalContentSize+needCleanBytes < sd.limit {
			break
		}
	}
	return nil
}
