package image

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var allowTypes = []string{"image/jpeg", "image/png", "image/gif"}

func checkBytesIsImage(b []byte) error {
	imageType := http.DetectContentType(b)
	for _, t := range allowTypes {
		if t == imageType {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid image type: %s. Allow types: %s", imageType, strings.Join(allowTypes, ", ")))
}

func checkDecodedStringIsImage(decodedImageType string) error {
	for _, t := range allowTypes {
		t = strings.ReplaceAll(t, "image/", "")
		if t == decodedImageType {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Fail to decode image, unknown type of source image: %s", decodedImageType))
}
