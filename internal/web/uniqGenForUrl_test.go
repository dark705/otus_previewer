package web

import (
	"net/url"
	"testing"
)

func TestGenUniqIdForUrlCorrect(t *testing.T) {
	uniqId := GenUniqIdForUrl(&url.URL{Path: "/fill/300/200/some_site.com/image.jpeg"})
	if uniqId != "d7ddb931c3cd7b83e5b6f1bd9d4717016d57569adb9e74912c2e311bf009813a" {
		t.Error("Incorrect parse uniq url ID")
	}
}
