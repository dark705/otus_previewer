package http

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type URLParams struct {
	Service    string
	Width      int
	Height     int
	RequestURL string
}

func ParseURL(url *url.URL) (p URLParams, err error) {
	const (
		serviceIndex  = 1
		widthIndex    = 2
		heightIndex   = 3
		urlStartIndex = 4
	)
	var allowServices = []string{"fill", "resize", "fit"}

	path := url.Path
	p = URLParams{}
	ps := strings.Split(path, "/")

	if len(ps) <= 5 {
		return p, fmt.Errorf("not enough params in path: %s", path)
	}

	//check and parse for allow services
	ps[serviceIndex] = strings.ToLower(ps[serviceIndex])
	err = fmt.Errorf("invalid service type: %s. Allow types: %s", ps[serviceIndex], strings.Join(allowServices, ", "))
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
		return p, fmt.Errorf("invalid Width: %s", ps[widthIndex])
	}

	//check and parse height
	p.Height, err = strconv.Atoi(ps[heightIndex])
	if err != nil || p.Height <= 0 {
		return p, fmt.Errorf("invalid Height: %s", ps[heightIndex])
	}

	//check and parse required remote url
	p.RequestURL = strings.Join(ps[urlStartIndex:], "/")
	if ps[urlStartIndex] == "" || ps[urlStartIndex+1] == "" {
		return p, fmt.Errorf("invalid requst Url: %s", p.RequestURL)
	}

	return p, nil
}
