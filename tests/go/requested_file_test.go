// +build integration

package previewer_test

import (
	"net/http"
	"testing"
)

func TestRequestedFileNotImage(t *testing.T) {
	response, err := http.Get("http://previewer:8013/resize/300/200/nginx/not_real_image.jpg")
	if err != nil {
		t.Error("fail on client get remote image", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadGateway {
		t.Errorf("on a non-image file, Service return status code: %d, but expected code: %d",
			response.StatusCode, http.StatusBadGateway)
	}
}
