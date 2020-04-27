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

func ParseURL(url *url.URL) (params URLParams, err error) {
	const (
		serviceIndex  = 1
		widthIndex    = 2
		heightIndex   = 3
		urlStartIndex = 4
	)
	var allowServices = []string{"fill", "resize", "fit"}

	path := url.Path
	params = URLParams{}
	slicedParams := strings.Split(path, "/")

	if len(slicedParams) <= 5 {
		return params, fmt.Errorf("not enough params in path: %s", path)
	}

	//check and parse for allow services
	slicedParams[serviceIndex] = strings.ToLower(slicedParams[serviceIndex])
	err = fmt.Errorf("invalid service type: %s. Allow types: %s", slicedParams[serviceIndex], strings.Join(allowServices, ", "))
	for _, t := range allowServices {
		if t == slicedParams[1] {
			params.Service = slicedParams[serviceIndex]
			err = nil
			break
		}
	}
	if err != nil {
		return params, err
	}

	//check and parse width
	params.Width, err = strconv.Atoi(slicedParams[widthIndex])
	if err != nil || params.Width <= 0 {
		return params, fmt.Errorf("invalid Width: %s", slicedParams[widthIndex])
	}

	//check and parse height
	params.Height, err = strconv.Atoi(slicedParams[heightIndex])
	if err != nil || params.Height <= 0 {
		return params, fmt.Errorf("invalid Height: %s", slicedParams[heightIndex])
	}

	//check and parse required remote url
	params.RequestURL = strings.Join(slicedParams[urlStartIndex:], "/")
	if slicedParams[urlStartIndex] == "" || slicedParams[urlStartIndex+1] == "" {
		return params, fmt.Errorf("invalid requst Url: %s", params.RequestURL)
	}

	return params, nil
}
