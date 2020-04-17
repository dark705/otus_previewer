package web

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func GetContext(url string, h http.Header, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest("GET", url, body)
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

	c, err := ioutil.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return c, err
}
