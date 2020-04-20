package image

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var allowTypes = []string{"image/jpeg", "image/png"}

func ReadImageContent(r io.ReadCloser, limit int) ([]byte, error) {
	defer r.Close()
	var content []byte
	offset := 0
	buf := make([]byte, 1024)
	for {
		read, err := r.Read(buf)
		content = append(content, buf[:read]...)
		if err == io.EOF {
			// on first step eol, check content is real image
			if offset == 0 {
				err := checkContentIsImage(buf[:read])
				if err != nil {
					return nil, err
				}
			}
			break
		}
		if err != nil {
			return nil, err
		}
		if offset+read > limit {
			return nil, errors.New(fmt.Sprintf("Requested image is bigger limit: %d", limit))
		}
		//check content is real image
		if offset == 0 {
			err := checkContentIsImage(buf[:read])
			if err != nil {
				return nil, err
			}
		}
		offset += read
	}
	return content, nil
}

func checkContentIsImage(b []byte) error {
	imageType := http.DetectContentType(b)
	for _, t := range allowTypes {
		if t == imageType {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid image type: %s. Allow types: %s", imageType, strings.Join(allowTypes, ", ")))
}
