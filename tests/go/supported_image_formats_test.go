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

func TestFormatJpeg(t *testing.T) {
	resp, err := http.Get("http://previewer:8013/resize/300/200/nginx/test_image.jpg")
	if err != nil {
		t.Error("fail on client get remote image", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Error(fmt.Sprintf("on resize existing image, Service return status code: %d, but expected code: %d",
			resp.StatusCode, http.StatusOK))
	}

	_, format, err := image.DecodeConfig(resp.Body)
	if err != nil {
		t.Error("fail on decode received image", err)
	}

	if format != "jpeg" {
		t.Error(fmt.Sprintf("on resize image format: jpeg, Service return format: %s", format))
	}
}

func TestFormatPng(t *testing.T) {
	resp, err := http.Get("http://previewer:8013/resize/300/200/nginx/test_image.png")
	if err != nil {
		t.Error("fail on client get remote image", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Error(fmt.Sprintf("on resize existing image, Service return status code: %d, but expected code: %d",
			resp.StatusCode, http.StatusOK))
	}

	_, format, err := image.DecodeConfig(resp.Body)
	if err != nil {
		t.Error("fail on decode received image", err)
	}

	if format != "png" {
		t.Error(fmt.Sprintf("on resize image format: png, Service return format: %s", format))
	}
}

func TestFormatGif(t *testing.T) {
	resp, err := http.Get("http://previewer:8013/resize/300/200/nginx/test_image.gif")
	if err != nil {
		t.Error("fail on client get remote image", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Error(fmt.Sprintf("on resize existing image, Service return status code: %d, but expected code: %d",
			resp.StatusCode, http.StatusOK))
	}

	_, format, err := image.DecodeConfig(resp.Body)
	if err != nil {
		t.Error("fail on decode received image", err)
	}

	if format != "gif" {
		t.Error(fmt.Sprintf("on resize image format: gif, Service return format: %s", format))
	}
}
