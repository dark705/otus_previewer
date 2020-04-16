package web

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type UrlParams struct {
	Type       string
	Width      int
	Height     int
	RequestUrl string
}

func parseUrl(URL *url.URL) (p UrlParams, err error) {
	const (
		widthIndex    = 2
		heightIndex   = 3
		urlStartIndex = 4
	)
	var allowTypes = []string{"fill"}

	path := URL.Path
	p = UrlParams{}
	ps := strings.Split(path, "/")

	if len(ps) <= 5 {
		return p, errors.New(fmt.Sprintf("Not enough params in path: %s", path))
	}

	for _, t := range allowTypes {
		if t == ps[1] {
			p.Type = ps[1]
			break
		}
		return p, errors.New(fmt.Sprintf("Invalid type: %s. Allow types: %s", ps[0], strings.Join(allowTypes, ", ")))
	}

	p.Width, err = strconv.Atoi(ps[widthIndex])
	if err != nil || p.Width <= 0 {
		return p, errors.New(fmt.Sprintf("Invalid Width: %s", ps[widthIndex]))
	}

	p.Height, err = strconv.Atoi(ps[heightIndex])
	if err != nil || p.Height <= 0 {
		return p, errors.New(fmt.Sprintf("Invalid Height: %s", ps[heightIndex]))
	}

	p.RequestUrl = strings.Join(ps[urlStartIndex:], "/")
	if ps[urlStartIndex] == "" || ps[urlStartIndex+1] == "" {
		return p, errors.New(fmt.Sprintf("Invalid requst Url: %s", p.RequestUrl))
	}

	return p, err
}
