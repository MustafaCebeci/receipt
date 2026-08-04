package receipt

import (
	"bytes"

	"github.com/phpdave11/gofpdf"
	"github.com/skip2/go-qrcode"
)

// QRCode generates a QR code from the given content and embeds it in the PDF.
func QRCode(pdf *gofpdf.Fpdf, content string, x, y, size float64) error {
	if content == "" {
		return nil
	}

	// Generate QR code as PNG
	qr, err := qrcode.Encode(content, qrcode.Medium, int(size*3))
	if err != nil {
		return ErrQRCodeGeneration
	}

	// Create a temporary buffer
	buf := new(bytes.Buffer)

	// Write PNG to buffer
	err = pngEncode(buf, qr)
	if err != nil {
		return ErrQRCodeGeneration
	}

	// Get current page
	pdf.AddPage()

	// Register the image
	info := pdf.RegisterImageOptionsReader("qr", gofpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(qr))
	if info == nil {
		return ErrQRCodeGeneration
	}

	// Calculate position to center the QR code
	pageWidth, _ := pdf.GetPageSize()
	qrWidth := info.Width()
	qrHeight := info.Height()

	// Center horizontally
	centerX := (pageWidth - qrWidth) / 2

	// Draw the image
	pdf.Image("qr", centerX, y, qrWidth, qrHeight, false, "PNG", 0, "")

	return nil
}

// pngEncode is a helper to encode data as PNG.
// Since go-qrcode already returns PNG data, this is a pass-through.
func pngEncode(w *bytes.Buffer, data []byte) error {
	w.Write(data)
	return nil
}

// QRCodeSize calculates the appropriate QR code size based on content length.
func QRCodeSize(content string) float64 {
	// Approximate size calculation
	// QR codes have minimum versions based on character count
	length := len(content)
	switch {
	case length <= 25:
		return 25 // QR Version 1
	case length <= 47:
		return 30 // QR Version 2
	case length <= 77:
		return 35 // QR Version 3
	default:
		return 40 // QR Version 4+
	}
}
