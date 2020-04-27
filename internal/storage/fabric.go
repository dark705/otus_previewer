package storage

import (
	"fmt"

	"github.com/dark705/otus_previewer/internal/storage/disk"
	"github.com/dark705/otus_previewer/internal/storage/inmemory"
	"github.com/sirupsen/logrus"
)

func Create(storageType, storagePath string, l *logrus.Logger) Storage {
	switch storageType {
	case "inmemory":
		l.Info("use operative memory as cache storage")
		s := inmemory.New()
		return &s
	case "disk":
		fallthrough
	default:
		l.Info(fmt.Sprintf("use disk folder %s as cache storage", storagePath))
		s := disk.New(storagePath)
		return &s
	}
}
