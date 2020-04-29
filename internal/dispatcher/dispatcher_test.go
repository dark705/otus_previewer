package dispatcher_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"testing"
	"time"

	"github.com/dark705/otus_previewer/internal/dispatcher"
	"github.com/dark705/otus_previewer/internal/storage/inmemory"
	"github.com/sirupsen/logrus"
)

func TestAddGetSame(t *testing.T) {
	storage := inmemory.New()
	imgDispatcher := dispatcher.New(&storage, 10, &logrus.Logger{})
	image := []byte("GIF89a,...") //len 10 bytes
	uniqID := genUniqID()

	err := imgDispatcher.Add(uniqID, image)
	if err != nil {
		t.Error("fail on add image", err)
	}

	imageGet, err := imgDispatcher.Get(uniqID)
	if err != nil {
		t.Error("fail on get image", err)
	}

	if imageGet == nil {
		t.Error("fail, added image do not exist")
	}

	if !bytes.Equal(imageGet, image) {
		t.Error("source and Get image not same")
	}
}

func TestTotalImagesSize(t *testing.T) {
	image := []byte("GIF89a,...") //len 10 bytes
	countImages := 10
	cacheSizeLimit := 100000000000000
	storage := inmemory.New()
	imageDispatcher := dispatcher.New(&storage, cacheSizeLimit, &logrus.Logger{})

	for i := 0; i < countImages; i++ {
		err := imageDispatcher.Add(genUniqID(), image)
		if err != nil {
			t.Error("fail on add image", err)
		}
	}

	if imageDispatcher.TotalImagesSize() != countImages*len(image) {
		t.Error("incorrect total images size")
	}
}

func TestTotalImagesSizeNotBiggerThenLimit(t *testing.T) {
	image := []byte("GIF89a,...") //len 10 bytes
	countImages := 1000
	cacheSizeLimit := 400
	storage := inmemory.New()
	imageDispatcher := dispatcher.New(&storage, cacheSizeLimit, &logrus.Logger{})

	for i := 0; i < countImages; i++ {
		err := imageDispatcher.Add(genUniqID(), image)
		if err != nil {
			t.Error("fail on add image", err)
		}
	}

	if imageDispatcher.TotalImagesSize() > countImages*len(image) {
		t.Errorf("incorrect total images size, limit: %d, TotalImagesSize: %d",
			cacheSizeLimit, imageDispatcher.TotalImagesSize())
	}
}

func TestLeastRecentUsed(t *testing.T) {
	image := []byte("GIF89a,...") //len 10 bytes
	cacheSizeLimit := 30

	storage := inmemory.New()
	imageDispatcher := dispatcher.New(&storage, cacheSizeLimit, &logrus.Logger{})

	//add image1
	uniqID1 := genUniqID()
	_ = imageDispatcher.Add(uniqID1, image)

	//add image2
	uniqID2 := genUniqID()
	_ = imageDispatcher.Add(uniqID2, image)

	//add image3
	uniqID3 := genUniqID()
	_ = imageDispatcher.Add(uniqID3, image)

	//now storage is full
	//get image2, image1, last Recent use updated
	_, _ = imageDispatcher.Get(uniqID2)
	_, _ = imageDispatcher.Get(uniqID1)

	//add image4, it must replace image3
	uniqID4 := genUniqID()
	_ = imageDispatcher.Add(uniqID4, image)

	image3, _ := imageDispatcher.Get(uniqID3)
	if image3 != nil {
		t.Error("image 3 exists but should not")
	}

	//add image5, it must replace image2
	uniqID5 := genUniqID()
	_ = imageDispatcher.Add(uniqID5, image)

	image2, _ := imageDispatcher.Get(uniqID2)
	if image2 != nil {
		t.Error("image 2 exists but should not")
	}

	image5, _ := imageDispatcher.Get(uniqID5)
	image4, _ := imageDispatcher.Get(uniqID4)
	image1, _ := imageDispatcher.Get(uniqID1)

	//now 3 in storage images: 5,4,1
	if image5 == nil || image4 == nil || image1 == nil {
		t.Error("expected images not exist")
	}
}

func TestRace(t *testing.T) {
	image := []byte("GIF89a,...") //len 10 bytes
	cacheSizeLimit := 80
	routines := 500

	storage := inmemory.New()
	imageDispatcher := dispatcher.New(&storage, cacheSizeLimit, &logrus.Logger{})
	wgWriters := sync.WaitGroup{}
	wgReaders := sync.WaitGroup{}

	wgWriters.Add(routines)
	wgReaders.Add(routines)
	for i := 0; i < routines; i++ {
		uniqID := genUniqID()

		go func(uniqId string) {
			err := imageDispatcher.Add(uniqId, image)
			if err != nil {
				t.Error("Cant add image")
			}
			imageDispatcher.TotalImagesSize()
			wgWriters.Done()
		}(uniqID)

		go func(uniqId string) {
			_, err := imageDispatcher.Get(uniqId)
			if err != nil {
				t.Error("Cant add image")
			}
			imageDispatcher.TotalImagesSize()
			wgReaders.Done()
		}(uniqID)
	}
	wgReaders.Wait()
	wgWriters.Wait()
}

func genUniqID() string {
	uniqBytes := sha256.Sum256([]byte(time.Now().String()))
	return hex.EncodeToString(uniqBytes[:])
}
