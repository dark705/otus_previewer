package previewer_test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRequestedFileNotImage(t *testing.T) {
	resp, err := http.Get("http://previewer:8013/resize/300/200/nginx/not_real_image.jpg")
	if err != nil {
		t.Error("Fail on client get remote image", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadGateway {
		t.Error(fmt.Sprintf("On a non-image file, Service return status code: %d, but expected code: %d",
			resp.StatusCode, http.StatusBadGateway))
	}
}
