// +build integration

package previewer_test

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"testing"
)

func TestResize(t *testing.T) {
	width := 300
	height := 200

	response, err := http.Get(fmt.Sprintf("http://previewer:8013/resize/%d/%d/nginx/test_image.jpg", width, height))
	if err != nil {
		t.Error("Fail on client get remote image", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Errorf("On resize existing image, Service return status code: %d, but expected code: %d",
			response.StatusCode, http.StatusOK)
	}

	img, _, err := image.DecodeConfig(response.Body)
	if err != nil {
		t.Error("fail on decode received image", err)
	}

	if img.Width != width || img.Height != height {
		t.Errorf("on resize image to: %dx%d, Service return %dx%d dimensions",
			width, height, img.Width, img.Height)
	}
}

func TestFill(t *testing.T) {
	width := 640
	height := 480

	response, err := http.Get(fmt.Sprintf("http://previewer:8013/fill/%d/%d/nginx/test_image.jpg", width, height))
	if err != nil {
		t.Error("fail on client get remote image", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Errorf("On fill existing image, Service return status code: %d, but expected code: %d",
			response.StatusCode, http.StatusOK)
	}

	img, _, err := image.DecodeConfig(response.Body)
	if err != nil {
		t.Error("fail on decode received image", err)
	}

	if img.Width != width || img.Height != height {
		t.Errorf("on fill image to: %dx%d, Service return %dx%d dimensions",
			width, height, img.Width, img.Height)
	}
}

func TestFit(t *testing.T) {
	width := 640
	height := 200
	fitWidth := 284
	fitHeight := 200

	response, err := http.Get(fmt.Sprintf("http://previewer:8013/fit/%d/%d/nginx/test_image.jpg", width, height))
	if err != nil {
		t.Error("fail on client get remote image", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Errorf("on fit existing image, Service return status code: %d, but expected code: %d",
			response.StatusCode, http.StatusOK)
	}

	image, _, err := image.DecodeConfig(response.Body)
	if err != nil {
		t.Error("fail on decode received image", err)
	}

	if image.Width != fitWidth || image.Height != fitHeight {
		t.Errorf("on fit image to: %dx%d, Service return %dx%d dimensions",
			fitWidth, fitHeight, image.Width, image.Height)
	}
}
