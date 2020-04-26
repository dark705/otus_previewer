// +build integration

package previewer_test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRemoteServerDoNotExist(t *testing.T) {
	resp, err := http.Get("http://previewer:8013/resize/300/200/some_fail_server.com/some_image.jpg")
	if err != nil {
		t.Error("Fail on client get remote image", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadGateway {
		t.Error(fmt.Sprintf("On a non-existing server, Service return status code: %d, but expected code: %d",
			resp.StatusCode, http.StatusBadGateway))
	}
}

func TestImageNotExist(t *testing.T) {
	resp, err := http.Get("http://previewer:8013/resize/300/200/nginx/image_not_exist.jpg")
	if err != nil {
		t.Error("Fail on client get remote image", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Error(fmt.Sprintf("On a non-existing image, Service return status code: %d, but expected code: %d",
			resp.StatusCode, http.StatusNotFound))
	}
}
