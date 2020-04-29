package http

import (
	"context"
	"io"
	"net/http"
)

func makeRequest(ctx context.Context, url string, header http.Header, body io.Reader) (*http.Response, error) {
	httpsResponse, err := makeRequestByProto(ctx, "https://", url, header, body)
	if err != nil {
		httpResponse, err := makeRequestByProto(ctx, "http://", url, header, body)
		if err != nil {
			return nil, err
		}
		return httpResponse, nil
	}
	return httpsResponse, nil
}

func makeRequestByProto(ctx context.Context, proto string, url string, header http.Header, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("GET", proto+url, body)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	request.Header = header

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
