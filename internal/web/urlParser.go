package web

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type UrlParams struct {
	Service    string
	Width      int
	Height     int
	RequestUrl string
}

func ParseUrl(url *url.URL) (p UrlParams, err error) {
	const (
		serviceIndex  = 1
		widthIndex    = 2
		heightIndex   = 3
		urlStartIndex = 4
	)
	var allowServices = []string{"fill", "resize", "fit"}

	path := url.Path
	p = UrlParams{}
	ps := strings.Split(path, "/")

	if len(ps) <= 5 {
		return p, errors.New(fmt.Sprintf("Not enough params in path: %s", path))
	}

	//check and parse for allow services
	ps[serviceIndex] = strings.ToLower(ps[serviceIndex])
	err = errors.New(fmt.Sprintf("Invalid service type: %s. Allow types: %s", ps[serviceIndex], strings.Join(allowServices, ", ")))
	for _, t := range allowServices {
		if t == ps[1] {
			p.Service = ps[serviceIndex]
			err = nil
			break
		}
	}
	if err != nil {
		return p, err
	}

	//check and parse width
	p.Width, err = strconv.Atoi(ps[widthIndex])
	if err != nil || p.Width <= 0 {
		return p, errors.New(fmt.Sprintf("Invalid Width: %s", ps[widthIndex]))
	}

	//check and parse height
	p.Height, err = strconv.Atoi(ps[heightIndex])
	if err != nil || p.Height <= 0 {
		return p, errors.New(fmt.Sprintf("Invalid Height: %s", ps[heightIndex]))
	}

	//check and parse required remote url
	p.RequestUrl = strings.Join(ps[urlStartIndex:], "/")
	if ps[urlStartIndex] == "" || ps[urlStartIndex+1] == "" {
		return p, errors.New(fmt.Sprintf("Invalid requst Url: %s", p.RequestUrl))
	}

	return p, nil
}
