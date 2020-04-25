package dispatcher

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/dark705/otus_previewer/internal/storage/inmemory"
	"github.com/sirupsen/logrus"
)

func TestAddGetSame(t *testing.T) {
	storage := inmemory.New()
	imgDispatcher := New(&storage, 10, &logrus.Logger{})
	image := []byte("GIF89a,...") //len 10 bytes
	uniqId := genUniqId()

	err := imgDispatcher.Add(uniqId, image)
	if err != nil {
		t.Error("Fail on add image", err)
	}

	if !imgDispatcher.Exist(uniqId) {
		if err != nil {
			t.Error("Fail, added image do not exist")
		}
	}

	imageGet, err := imgDispatcher.Get(uniqId)
	if err != nil {
		t.Error("Fail on get image", err)
	}

	if !bytes.Equal(imageGet, image) {
		t.Error("Source and Get image not same")
	}
}

func TestTotalImagesSize(t *testing.T) {
	image := []byte("GIF89a,...") //len 10 bytes
	countImages := 10
	cacheSizeLimit := 100000000000000
	storage := inmemory.New()
	imageDispatcher := New(&storage, cacheSizeLimit, &logrus.Logger{})

	for i := 0; i < countImages; i++ {
		err := imageDispatcher.Add(genUniqId(), image)
		if err != nil {
			t.Error("Fail on add image", err)
		}
	}

	if imageDispatcher.TotalImagesSize() != countImages*len(image) {
		t.Error("Incorrect total images size")
	}
}

func TestTotalImagesSizeNotBiggerThenLimit(t *testing.T) {
	image := []byte("GIF89a,...") //len 10 bytes
	countImages := 1000
	cacheSizeLimit := 400
	storage := inmemory.New()
	imageDispatcher := New(&storage, cacheSizeLimit, &logrus.Logger{})

	for i := 0; i < countImages; i++ {
		err := imageDispatcher.Add(genUniqId(), image)
		if err != nil {
			t.Error("Fail on add image", err)
		}
	}

	if imageDispatcher.TotalImagesSize() > countImages*len(image) {
		t.Error(fmt.Sprintf("Incorrect total images size, limit: %d, TotalImagesSize: %d",
			cacheSizeLimit, imageDispatcher.TotalImagesSize()))
	}
}

func TestLeastRecentUsed(t *testing.T) {
	image := []byte("GIF89a,...") //len 10 bytes
	cacheSizeLimit := 30

	storage := inmemory.New()
	imageDispatcher := New(&storage, cacheSizeLimit, &logrus.Logger{})

	//add image1
	uniqId1 := genUniqId()
	_ = imageDispatcher.Add(uniqId1, image)

	//add image2
	uniqId2 := genUniqId()
	_ = imageDispatcher.Add(uniqId2, image)

	//add image3
	uniqId3 := genUniqId()
	_ = imageDispatcher.Add(uniqId3, image)

	//now storage is full
	//get image2, image1, last Recent use updated
	time.Sleep(time.Nanosecond) //windows fix
	_, _ = imageDispatcher.Get(uniqId2)
	time.Sleep(time.Nanosecond) //windows fix
	_, _ = imageDispatcher.Get(uniqId1)

	//add image4, it must replace image3
	uniqId4 := genUniqId()
	_ = imageDispatcher.Add(uniqId4, image)

	if imageDispatcher.Exist(uniqId3) {
		t.Error("Image 3 exists but should not")
	}

	//add image5, it must replace image2
	uniqId5 := genUniqId()
	_ = imageDispatcher.Add(uniqId5, image)

	if imageDispatcher.Exist(uniqId2) {
		t.Error("Image 2 exists but should not")
	}

	//now 3 in storage images: 5,4,1
	if !imageDispatcher.Exist(uniqId5) ||
		!imageDispatcher.Exist(uniqId4) ||
		!imageDispatcher.Exist(uniqId1) {
		t.Error("Expected images not exist")
	}
}

func genUniqId() string {
	time.Sleep(time.Nanosecond) //windows fix
	b := sha256.Sum256([]byte(time.Now().String()))
	return hex.EncodeToString(b[:])
}
