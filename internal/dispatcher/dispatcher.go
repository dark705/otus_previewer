package dispatcher

import (
	"fmt"
	"sort"
	"time"

	"github.com/dark705/otus_previewer/internal/storage"
	"github.com/sirupsen/logrus"
)

type imageInfo struct {
	size    int
	lastUse time.Time
}

type ImageDispatcher struct {
	storage         storage.Storage
	imageState      map[string]imageInfo
	totalImagesSize int
	maxLimit        int
	logger          *logrus.Logger
}

func New(storage storage.Storage, limit int, logger *logrus.Logger) ImageDispatcher {
	state := make(map[string]imageInfo)
	var totalSize int
	for id, size := range storage.GetListSize() {
		state[id] = imageInfo{size: size, lastUse: time.Now()}
		totalSize += size
	}

	return ImageDispatcher{
		storage:         storage,
		imageState:      state,
		totalImagesSize: totalSize,
		maxLimit:        limit,
		logger:          logger,
	}
}

func (imDis *ImageDispatcher) TotalImagesSize() int {
	return imDis.totalImagesSize
}

func (imDis *ImageDispatcher) Exist(id string) bool {
	_, exist := imDis.imageState[id]
	return exist
}

func (imDis *ImageDispatcher) Get(id string) ([]byte, error) {
	_, exist := imDis.imageState[id]
	if !exist {
		return nil, fmt.Errorf("Fail on update lastUse on get, image with id: %s not exist", id)
	}
	imDis.imageState[id] = imageInfo{size: imDis.imageState[id].size, lastUse: time.Now()}
	imDis.logger.Debugln(fmt.Sprintf("Image with id: %s updated, last use time: %s", id, time.Now()))
	return imDis.storage.Get(id)
}

func (imDis *ImageDispatcher) Add(id string, image []byte) error {
	//storage not full
	if imDis.totalImagesSize+len(image) <= imDis.maxLimit {
		return imDis.addAvailable(id, image)
	}
	//storage is full, need to clean,
	imDis.logger.Debugln(fmt.Sprintf("Storage is full, totalImagesSize: %d, make clean", imDis.TotalImagesSize()))
	err := imDis.cleanOldUseImagesOn(len(image))
	if err != nil {
		return err
	}
	//now we can add
	return imDis.addAvailable(id, image)
}

func (imDis *ImageDispatcher) addAvailable(id string, image []byte) error {
	err := imDis.storage.Add(id, image)
	if err != nil {
		return err
	}
	imDis.totalImagesSize += len(image)
	imDis.imageState[id] = imageInfo{size: len(image), lastUse: time.Now()}
	imDis.logger.Debugln(fmt.Sprintf("Storage not full, add image with id: %s, size: %d, now total images size: %d", id, len(image), imDis.TotalImagesSize()))
	return nil
}

func (imDis *ImageDispatcher) cleanOldUseImagesOn(needCleanBytes int) error {
	//make sort images, smaller index is more old at time
	list := make([]struct {
		i string
		t time.Time
	}, 0, len(imDis.imageState))

	for id, is := range imDis.imageState {
		list = append(list, struct {
			i string
			t time.Time
		}{i: id, t: is.lastUse})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].t.Before(list[j].t)
	})

	//delete old used images until free space for new image
	for _, val := range list {
		err := imDis.storage.Del(val.i)
		if err != nil {
			return err
		}
		usedSize := imDis.imageState[val.i].size
		imDis.totalImagesSize -= usedSize
		delete(imDis.imageState, val.i)
		imDis.logger.Debugln(fmt.Sprintf("Deleted image with id: %s, used size: %d", val.i, usedSize))
		if imDis.totalImagesSize+needCleanBytes <= imDis.maxLimit {
			break
		}
	}
	return nil
}
