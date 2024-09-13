package watermark

import (
	"errors"
	"image"

	"github.com/nfnt/resize"
)

func Resize(imageObject image.Image, width uint, height uint) (
	img image.Image, err error) {

	// validating image object
	if imageObject == nil {
		return nil, errors.New("image object required")
	}
	return resize.Resize(width, height, imageObject, resize.Lanczos3), nil
}
