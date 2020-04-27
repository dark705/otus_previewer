package http

import (
	"io"
	"net/http"
)

func makeRequest(url string, header http.Header, body io.Reader) (*http.Response, error) {
	httpsResponse, err := makeRequestByProto("https://", url, header, body)
	if err != nil {
		httpResponse, err := makeRequestByProto("http://", url, header, body)
		if err != nil {
			return nil, err
		}
		return httpResponse, nil
	}
	return httpsResponse, nil
}

func makeRequestByProto(proto string, url string, header http.Header, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("GET", proto+url, body)
	if err != nil {
		return nil, err
	}
	request.Header = header

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
