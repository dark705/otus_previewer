package web

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dark705/otus_previewer/internal/image"
)

func GetImage(proto string, url string, h http.Header, body io.Reader, limit int) ([]byte, error) {
	req, err := http.NewRequest("GET", proto+url, body)
	if err != nil {
		return nil, err
	}
	req.Header = h

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Wrong status code of remote server: %d", resp.StatusCode))
	}

	c, err := image.ReadImageAsByte(resp.Body, limit)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return c, err
}
