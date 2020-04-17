package web

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
)

func GenUniqIdForUrl(url *url.URL) string {
	b := sha256.Sum256([]byte(url.Path))
	s := hex.EncodeToString(b[:])
	return s
}
