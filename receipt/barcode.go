package receipt

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/code39"
	"github.com/boombuler/barcode/qr"
	"github.com/phpdave11/gofpdf"
)

// BarcodeType represents the type of barcode.
type BarcodeType string

const (
	// Code128 is a Code 128 barcode.
	Code128 BarcodeType = "CODE128"
	// Code39 is a Code 39 barcode.
	Code39 BarcodeType = "CODE39"
	// QRCodeType is a QR code.
	QRCodeType BarcodeType = "QR"
)

// Barcode generates a barcode from the given content and embeds it in the PDF.
func Barcode(pdf *gofpdf.Fpdf, content string, barcodeType BarcodeType, x, y, width, height float64) error {
	if content == "" {
		return nil
	}

	var bc barcode.Barcode
	var err error

	switch barcodeType {
	case Code128:
		bc, err = code128.Encode(content)
	case Code39:
		bc, err = code39.Encode(content, true, true)
	case QRCodeType:
		bc, err = qr.Encode(content, qr.M, qr.Auto)
	default:
		// Default to Code128
		bc, err = code128.Encode(content)
	}

	if err != nil {
		return ErrBarcodeGeneration
	}

	// Scale the barcode to fit the requested size
	bcWidth, bcHeight := bc.Bounds().Dx(), bc.Bounds().Dy()
	if bcWidth == 0 || bcHeight == 0 {
		return ErrBarcodeGeneration
	}

	scale := min(width/float64(bcWidth), height/float64(bcHeight))
	scaledWidth := float64(bcWidth) * scale
	scaledHeight := float64(bcHeight) * scale

	// Center the barcode if width/height are larger
	drawX := x
	drawY := y
	if width > scaledWidth {
		drawX = x + (width - scaledWidth) / 2
	}
	if height > scaledHeight {
		drawY = y + (height - scaledHeight) / 2
	}

	// Create PNG image from barcode
	pngData, err := barcodeToPNG(bc)
	if err != nil {
		return ErrBarcodeGeneration
	}

	// Register and draw the image
	imgName := "barcode_" + strings.ReplaceAll(content[:min(10, len(content))], "/", "_")
	pdf.RegisterImageOptionsReader(imgName, gofpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(pngData))

	pdf.Image(imgName, drawX, drawY, scaledWidth, scaledHeight, false, "PNG", 0, "")

	return nil
}

// barcodeToPNG converts a barcode to PNG bytes.
func barcodeToPNG(bc barcode.Barcode) ([]byte, error) {
	bounds := bc.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create a grayscale image
	img := image.NewGray(bounds)

	// Draw the barcode
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get the color at this pixel - true/false for barcode
			_, _, _, a := bc.At(x, y).RGBA()
			if a > 32768 {
				img.SetGray(x, y, color.Gray{Y: 0}) // black
			} else {
				img.SetGray(x, y, color.Gray{Y: 255}) // white
			}
		}
	}

	// Encode to PNG
	pngBuf := new(bytes.Buffer)
	err := png.Encode(pngBuf, img)
	if err != nil {
		return nil, err
	}

	return pngBuf.Bytes(), nil
}

// Code128Barcode generates a Code128 barcode.
func Code128Barcode(pdf *gofpdf.Fpdf, content string, x, y, width, height float64) error {
	return Barcode(pdf, content, Code128, x, y, width, height)
}

// Code39Barcode generates a Code39 barcode.
func Code39Barcode(pdf *gofpdf.Fpdf, content string, x, y, width, height float64) error {
	return Barcode(pdf, content, Code39, x, y, width, height)
}
