package http

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
)

func GenUniqIDForURL(url *url.URL) string {
	uniqBytes := sha256.Sum256([]byte(url.Path))
	return hex.EncodeToString(uniqBytes[:])
}
