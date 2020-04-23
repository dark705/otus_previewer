package web

import (
	"net/url"
	"testing"
)

func TestLenUrl(t *testing.T) {
	_, err := ParseUrl(&url.URL{Path: "/fill/300/200/some_site.com/image.jpeg"})
	if err != nil {
		t.Error(err)
	}

	_, err = ParseUrl(&url.URL{Path: "/fill/300/200/some_site.com"})
	if err == nil {
		t.Error("No error on wrong url")
	}
}

func TestServiceType(t *testing.T) {
	_, err := ParseUrl(&url.URL{Path: "/some/300/200/some_site.com/image.jpeg"})
	if err == nil {
		t.Error("No error on wrong url service")
	}

	p1, err := ParseUrl(&url.URL{Path: "/fill/300/200/some_site.com/image.jpeg"})
	p2, err := ParseUrl(&url.URL{Path: "/resize/300/200/some_site.com/image.jpeg"})
	p3, err := ParseUrl(&url.URL{Path: "/fit/300/200/some_site.com/image.jpeg"})
	if err != nil ||
		p1.Service != "fill" ||
		p2.Service != "resize" ||
		p3.Service != "fit" {
		t.Error("Error on correct url service")
	}
}

func TestWidthHeight(t *testing.T) {
	p, err := ParseUrl(&url.URL{Path: "/fill/300/200/some_site.com/image.jpeg"})
	if err != nil || p.Width != 300 || p.Height != 200 {
		t.Error("Incorrect parse Width or Height")
	}
}

func TestRequestedUrl(t *testing.T) {
	p, err := ParseUrl(&url.URL{Path: "/fill/300/200/some_site.com/image.jpeg"})
	if err != nil || p.RequestUrl != "some_site.com/image.jpeg" {
		t.Error("Incorrect parse remote requested url")
	}
}
