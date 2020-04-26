package storage

import (
	"fmt"

	"github.com/dark705/otus_previewer/internal/storage/disk"
	"github.com/dark705/otus_previewer/internal/storage/inmemory"
	"github.com/sirupsen/logrus"
)

func Create(t, p string, l *logrus.Logger) Storage {
	switch t {
	case "inmemory":
		l.Info("use operative memory as cache storage")
		s := inmemory.New()
		return &s
	case "disk":
		fallthrough
	default:
		l.Info(fmt.Sprintf("use disk folder %s as cache storage", p))
		s := disk.New(p)
		return &s
	}
}
