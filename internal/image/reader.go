package image

import (
	"errors"
	"fmt"
	"io"
)

func ReadImageAsByte(r io.ReadCloser, limit int) ([]byte, error) {
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
				err := checkBytesIsImage(buf[:read])
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
			err := checkBytesIsImage(buf[:read])
			if err != nil {
				return nil, err
			}
		}
		offset += read
	}
	return content, nil
}