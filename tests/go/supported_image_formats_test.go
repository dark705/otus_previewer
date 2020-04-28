// +build integration

package previewer_test

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"testing"
)

func TestFormatJpeg(t *testing.T) {
	response, err := http.Get("http://previewer:8013/resize/300/200/nginx/test_image.jpg")
	if err != nil {
		t.Error("fail on client get remote image", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Errorf("on resize existing image, Service return status code: %d, but expected code: %d",
			response.StatusCode, http.StatusOK)
	}

	_, format, err := image.DecodeConfig(response.Body)
	if err != nil {
		t.Error("fail on decode received image", err)
	}

	if format != "jpeg" {
		t.Errorf("on resize image format: jpeg, Service return format: %s", format)
	}
}

func TestFormatPng(t *testing.T) {
	response, err := http.Get("http://previewer:8013/resize/300/200/nginx/test_image.png")
	if err != nil {
		t.Error("fail on client get remote image", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Errorf("on resize existing image, Service return status code: %d, but expected code: %d",
			response.StatusCode, http.StatusOK)
	}

	_, format, err := image.DecodeConfig(response.Body)
	if err != nil {
		t.Error("fail on decode received image", err)
	}

	if format != "png" {
		t.Errorf("on resize image format: png, Service return format: %s", format)
	}
}

func TestFormatGif(t *testing.T) {
	response, err := http.Get("http://previewer:8013/resize/300/200/nginx/test_image.gif")
	if err != nil {
		t.Error("fail on client get remote image", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Errorf("on resize existing image, Service return status code: %d, but expected code: %d",
			response.StatusCode, http.StatusOK)
	}

	_, format, err := image.DecodeConfig(response.Body)
	if err != nil {
		t.Error("fail on decode received image", err)
	}

	if format != "gif" {
		t.Errorf("on resize image format: gif, Service return format: %s", format)
	}
}
