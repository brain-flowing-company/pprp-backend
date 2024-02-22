package utils

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/nfnt/resize"
)

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

type ImageProcessor struct {
	img    image.Image
	width  int
	height int
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (ip *ImageProcessor) LoadPNG(file io.Reader) error {
	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	ip.img = img
	ip.width = img.Bounds().Max.X
	ip.height = img.Bounds().Max.Y

	return err
}

func (ip *ImageProcessor) LoadJPEG(file io.Reader) error {
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	ip.img = img
	ip.width = img.Bounds().Max.X
	ip.height = img.Bounds().Max.Y

	return err
}

func (ip *ImageProcessor) SquareCropped() error {
	subImg, ok := ip.img.(subImager)
	if !ok {
		return errors.New("could not crop image")
	}

	size := Min(ip.width, ip.height)
	x, y := (ip.width-size)/2, (ip.height-size)/2

	ip.img = subImg.SubImage(image.Rect(x, y, x+size, y+size))

	return nil
}

func (ip *ImageProcessor) Resize(minSize int) error {
	w, h := 0, 0
	if ip.width < ip.height {
		w = Min(ip.width, minSize)
	} else {
		h = Min(ip.height, minSize)
	}

	ip.img = resize.Resize(uint(w), uint(h), ip.img, resize.Bilinear)

	return nil
}

func (ip *ImageProcessor) Save() (io.Reader, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, ip.img, nil)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
