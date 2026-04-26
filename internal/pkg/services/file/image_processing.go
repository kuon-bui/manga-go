package fileservice

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"

	"github.com/HugoSmits86/nativewebp"
	xdraw "golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

func decodeImage(raw []byte) (image.Image, error) {
	if len(raw) == 0 {
		return nil, errors.New("empty image payload")
	}

	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	b := img.Bounds()
	if b.Dx() <= 0 || b.Dy() <= 0 {
		return nil, errors.New("invalid image dimensions")
	}

	return img, nil
}

func encodeWebPVariant(src image.Image, preset imageVariantPreset) ([]byte, int, error) {
	bounds := src.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	targetWidth := srcWidth
	if preset.Width > 0 && preset.Width < srcWidth {
		targetWidth = preset.Width
	}

	targetImage := src
	if targetWidth != srcWidth {
		targetHeight := int(math.Round(float64(srcHeight) * float64(targetWidth) / float64(srcWidth)))
		if targetHeight < 1 {
			targetHeight = 1
		}

		resized := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
		xdraw.CatmullRom.Scale(resized, resized.Bounds(), src, bounds, xdraw.Over, nil)
		targetImage = resized
	}

	buf := bytes.NewBuffer(nil)
	if err := nativewebp.Encode(buf, targetImage, &nativewebp.Options{UseExtendedFormat: false}); err != nil {
		return nil, 0, err
	}

	return buf.Bytes(), targetWidth, nil
}
