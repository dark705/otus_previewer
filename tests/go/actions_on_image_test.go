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

	resp, err := http.Get(fmt.Sprintf("http://previewer:8013/resize/%d/%d/nginx/test_image.jpg", width, height))
	if err != nil {
		t.Error("Fail on client get remote image", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Error(fmt.Sprintf("On resize existing image, Service return status code: %d, but expected code: %d",
			resp.StatusCode, http.StatusOK))
	}

	img, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		t.Error("Fail on decode received image", err)
	}

	if img.Width != width || img.Height != height {
		t.Error(fmt.Sprintf("On resize image to: %dx%d, Service return %dx%d dimensions",
			width, height, img.Width, img.Height))
	}
}

func TestFill(t *testing.T) {
	width := 640
	height := 480

	resp, err := http.Get(fmt.Sprintf("http://previewer:8013/fill/%d/%d/nginx/test_image.jpg", width, height))
	if err != nil {
		t.Error("Fail on client get remote image", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Error(fmt.Sprintf("On fill existing image, Service return status code: %d, but expected code: %d",
			resp.StatusCode, http.StatusOK))
	}

	img, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		t.Error("Fail on decode received image", err)
	}

	if img.Width != width || img.Height != height {
		t.Error(fmt.Sprintf("On fill image to: %dx%d, Service return %dx%d dimensions",
			width, height, img.Width, img.Height))
	}
}

func TestFit(t *testing.T) {
	width := 640
	height := 200
	fitWidth := 284
	fitHeight := 200

	resp, err := http.Get(fmt.Sprintf("http://previewer:8013/fit/%d/%d/nginx/test_image.jpg", width, height))
	if err != nil {
		t.Error("Fail on client get remote image", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Error(fmt.Sprintf("On fit existing image, Service return status code: %d, but expected code: %d",
			resp.StatusCode, http.StatusOK))
	}

	img, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		t.Error("Fail on decode received image", err)
	}

	if img.Width != fitWidth || img.Height != fitHeight {
		t.Error(fmt.Sprintf("On fit image to: %dx%d, Service return %dx%d dimensions",
			fitWidth, fitHeight, img.Width, img.Height))
	}
}
