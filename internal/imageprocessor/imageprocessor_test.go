package imageprocessor

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

func createTestImage(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 0, 255}) //nolint:gosec
		}
	}

	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, nil)
	return buf.Bytes()
}

func TestResize(t *testing.T) {
	t.Parallel()

	originalWidth, originalHeight := 100, 100
	newWidth, newHeight := 50, 50
	imgData := createTestImage(originalWidth, originalHeight)

	resizedImg, err := Resize(imgData, newWidth, newHeight)
	if err != nil {
		t.Fatalf("failed to resize image: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(resizedImg))
	if err != nil {
		t.Fatalf("failed to decode resized image: %v", err)
	}

	if img.Bounds().Dx() != newWidth || img.Bounds().Dy() != newHeight {
		t.Fatalf("unexpected dimensions: got %dx%d, want %dx%d", img.Bounds().Dx(), img.Bounds().Dy(), newWidth, newHeight)
	}
}

func TestGetImageDimensions(t *testing.T) {
	t.Parallel()

	width, height := 200, 150
	imgData := createTestImage(width, height)

	gotWidth, gotHeight, err := GetImageDimensions(imgData)
	if err != nil {
		t.Fatalf("failed to get image dimensions: %v", err)
	}

	if gotWidth != width || gotHeight != height {
		t.Fatalf("unexpected dimensions: got %dx%d, want %dx%d", gotWidth, gotHeight, width, height)
	}
}
