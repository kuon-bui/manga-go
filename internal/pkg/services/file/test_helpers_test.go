package fileservice

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
)

var errFakeNotFound = errors.New("fake not found")

type fakeStorage struct {
	files               map[string][]byte
	getErr              map[string]error
	uploaded            map[string][]byte
	uploadedContentType map[string]string
	getCalls            []string
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{
		files:               map[string][]byte{},
		getErr:              map[string]error{},
		uploaded:            map[string][]byte{},
		uploadedContentType: map[string]string{},
		getCalls:            make([]string, 0, 4),
	}
}

func (f *fakeStorage) CreatePresignedURL(_ context.Context, key string) (string, error) {
	return "https://example.com/" + key, nil
}

func (f *fakeStorage) GetFile(_ context.Context, fileName string) ([]byte, error) {
	f.getCalls = append(f.getCalls, fileName)
	if err, ok := f.getErr[fileName]; ok {
		return nil, err
	}
	if b, ok := f.files[fileName]; ok {
		dup := make([]byte, len(b))
		copy(dup, b)
		return dup, nil
	}
	return nil, errFakeNotFound
}

func (f *fakeStorage) UploadFile(_ context.Context, fileName string, body io.Reader, _ int64, contentType string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	f.uploaded[fileName] = b
	f.uploadedContentType[fileName] = contentType
	return nil
}

func (f *fakeStorage) IsNotFoundError(err error) bool {
	return errors.Is(err, errFakeNotFound)
}

func makePNG(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 255), G: uint8(y % 255), B: uint8((x + y) % 255), A: 255})
		}
	}

	buf := bytes.NewBuffer(nil)
	_ = png.Encode(buf, img)
	return buf.Bytes()
}
