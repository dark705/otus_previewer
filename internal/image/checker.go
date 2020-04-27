package image

import (
	"fmt"
	"net/http"
	"strings"
)

var allowTypes = []string{"image/jpeg", "image/png", "image/gif"}

func checkBytesIsImage(b []byte) error {
	imageType := http.DetectContentType(b)
	for _, allowType := range allowTypes {
		if allowType == imageType {
			return nil
		}
	}
	return fmt.Errorf("invalid image type: %s. Allow types: %s", imageType, strings.Join(allowTypes, ", "))
}

func checkDecodedStringIsImage(decodedImageType string) error {
	for _, allowType := range allowTypes {
		allowType = strings.ReplaceAll(allowType, "image/", "")
		if allowType == decodedImageType {
			return nil
		}
	}
	return fmt.Errorf("fail to decode image, unknown type of source image: %s", decodedImageType)
}
