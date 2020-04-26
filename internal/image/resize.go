package image

import (
	"bytes"
	"fmt"
	"image"

	"github.com/disintegration/imaging"

	_ "image/gif"  // Register some gif support
	_ "image/jpeg" // Register some jpeg support
	_ "image/png"  // Register some png support
)

type ResizeConfig struct {
	Action string
	Width  int
	Height int
}

func Resize(srcImageContent []byte, p ResizeConfig) ([]byte, error) {
	//Decode
	srcImage, ds, err := image.Decode(bytes.NewReader(srcImageContent))
	if err != nil {
		return nil, err
	}
	err = checkDecodedStringIsImage(ds)
	if err != nil {
		return nil, err
	}

	//resize
	var destImage *image.NRGBA
	switch p.Action {
	case "fill":
		destImage = imaging.Fill(srcImage, p.Width, p.Height, imaging.Center, imaging.Lanczos)
	case "resize":
		destImage = imaging.Resize(srcImage, p.Width, p.Height, imaging.Lanczos)
	case "fit":
		destImage = imaging.Fit(srcImage, p.Width, p.Height, imaging.Lanczos)
	default:
		return nil, fmt.Errorf("unknown action on image: %s", p.Action)
	}

	//encode
	var buf bytes.Buffer
	switch ds {
	case "png":
		err = imaging.Encode(&buf, destImage, imaging.PNG)
	case "jpeg":
		err = imaging.Encode(&buf, destImage, imaging.JPEG)
	case "gif":
		err = imaging.Encode(&buf, destImage, imaging.GIF)
	default:
		err = fmt.Errorf("fail encode image, type: %s", ds)
	}
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
