package imageprocessor

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

func Resize(img []byte, width int, height int) ([]byte, error) {
	var buf bytes.Buffer

	imgDecoded, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	imgFilled := imaging.Fill(imgDecoded, width, height, imaging.Center, imaging.Lanczos)

	err = jpeg.Encode(&buf, imgFilled, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
