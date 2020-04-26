package web_test

import (
	"net/url"
	"testing"

	"github.com/dark705/otus_previewer/internal/web"
)

func TestLenUrl(t *testing.T) {
	_, err := web.ParseURL(&url.URL{Path: "/fill/300/200/some_site.com/image.jpeg"})
	if err != nil {
		t.Error(err)
	}

	_, err = web.ParseURL(&url.URL{Path: "/fill/300/200/some_site.com"})
	if err == nil {
		t.Error("No error on wrong url")
	}
}

func TestServiceType(t *testing.T) {
	_, err := web.ParseURL(&url.URL{Path: "/some/300/200/some_site.com/image.jpeg"})
	if err == nil {
		t.Error("No error on wrong url service")
	}

	p, err := web.ParseURL(&url.URL{Path: "/fill/300/200/some_site.com/image.jpeg"})
	if err != nil || p.Service != "fill" {
		t.Error("Error on correct url service fill ")
	}

	p, err = web.ParseURL(&url.URL{Path: "/resize/300/200/some_site.com/image.jpeg"})
	if err != nil || p.Service != "resize" {
		t.Error("Error on correct url service resize")
	}

	p, err = web.ParseURL(&url.URL{Path: "/fit/300/200/some_site.com/image.jpeg"})
	if err != nil || p.Service != "fit" {
		t.Error("Error on correct url service fit")
	}
}

func TestWidthHeight(t *testing.T) {
	p, err := web.ParseURL(&url.URL{Path: "/fill/300/200/some_site.com/image.jpeg"})
	if err != nil || p.Width != 300 || p.Height != 200 {
		t.Error("Incorrect parse Width or Height")
	}
}

func TestRequestedUrl(t *testing.T) {
	p, err := web.ParseURL(&url.URL{Path: "/fill/300/200/some_site.com/image.jpeg"})
	if err != nil || p.RequestURL != "some_site.com/image.jpeg" {
		t.Error("Incorrect parse remote requested url")
	}
}
