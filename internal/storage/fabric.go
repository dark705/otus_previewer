package storage

import (
	"fmt"
	"github.com/dark705/otus_previewer/internal/storage/disk"
	"github.com/dark705/otus_previewer/internal/storage/inmemory"
	"github.com/sirupsen/logrus"
)

func CreateStorage(t, p string, l *logrus.Logger) Storage {
	switch t {
	case "inmemory":
		l.Info("Use operative memory as cache storage")
		s := inmemory.New()
		return &s
	case "disk":
		fallthrough
	default:
		l.Info(fmt.Sprintf("Use disk folder %s as cache storage", p))
		s := disk.New(p)
		return &s
	}
}
