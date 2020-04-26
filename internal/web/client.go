package web

import (
	"io"
	"net/http"
)

func makeRequest(proto string, url string, h http.Header, body io.Reader) (*http.Response, error) {
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
	return resp, nil
}
