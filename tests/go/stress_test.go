// +build integration

package previewer_test

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"sync"

	"testing"
)

func TestStressAndRace(t *testing.T) {
	count := 1000
	wg := sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(clientId int) {
			response, err := http.Get("http://previewer:8013/resize/300/200/nginx/test_image.jpg")
			if err != nil {
				t.Errorf("fail on client: %d get remote image, err: %s", clientId, err)
			}
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				t.Errorf("on resize existing image, Service return status code: %d, but expected code: %d",
					response.StatusCode, http.StatusOK)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}
