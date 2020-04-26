// +build integration

package previewer_test

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"sync"

	"testing"
)

func TestStressAndRace(t *testing.T) {
	count := 100
	wg := sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(clientId int) {
			resp, err := http.Get("http://previewer:8013/resize/300/200/nginx/test_image.jpg")
			if err != nil {
				t.Error(fmt.Sprintf("Fail on client: %d get remote image, err: %s", clientId, err))
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Error(fmt.Sprintf("On resize existing image, Service return status code: %d, but expected code: %d",
					resp.StatusCode, http.StatusOK))
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}
