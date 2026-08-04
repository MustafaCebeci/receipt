package receipt

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	"github.com/phpdave11/gofpdf"
)

// DejaVu Sans TTF font URL (Google Fonts CDN).
const dejavuSansURL = "https://github.com/dejavu-fonts/dejavu-fonts/raw/master/out/ttf/DejaVuSans.ttf"

var (
	fontCache     sync.Map // map[string]*registeredFont
	fontCacheOnce sync.Map // map[string]chan struct{}
)

// registeredFont holds cached font data.
type registeredFont struct {
	regular []byte
	bold    []byte
	italic  []byte
	boldItalic []byte
}

// getFont retrieves or downloads and caches the DejaVu Sans font.
func getFont() (*registeredFont, error) {
	// Check cache first
	if cached, ok := fontCache.Load("dejavu"); ok {
		return cached.(*registeredFont), nil
	}

	// Prevent concurrent download
	once, _ := fontCacheOnce.LoadOrStore("dejavu", make(chan struct{}))
	ch := once.(chan struct{})
	select {
	case <-ch:
		// Download completed by another goroutine, check cache again
		if cached, ok := fontCache.Load("dejavu"); ok {
			return cached.(*registeredFont), nil
		}
		return nil, ErrFontLoading
	default:
		// We're the first, proceed with download
	}

	defer func() {
		fontCacheOnce.Delete("dejavu")
		close(ch)
	}()

	// Download the font
	resp, err := http.Get(dejavuSansURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrFontLoading
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// For DejaVu Sans, we use the same file for all variants
	// A real implementation would download separate bold/italic files
	font := &registeredFont{
		regular:    data,
		bold:       data,
		italic:     data,
		boldItalic: data,
	}

	fontCache.Store("dejavu", font)
	return font, nil
}

// registerFontWithGofpdf registers the cached font with gofpdf.
func registerFontWithGofpdf(pdf *gofpdf.Fpdf, family string) error {
	font, err := getFont()
	if err != nil {
		return err
	}

	// Register font with all variants using the same file for simplicity
	// AddFontFromBytes(family, style string, file, file2 []byte)
	pdf.AddFontFromBytes(family, "", font.regular, font.bold)
	pdf.AddFontFromBytes(family, "B", font.bold, font.regular)
	pdf.AddFontFromBytes(family, "I", font.italic, font.regular)
	pdf.AddFontFromBytes(family, "BI", font.boldItalic, font.regular)

	return nil
}

// SetDefaultFont sets the default font family for the PDF.
func SetDefaultFont(pdf *gofpdf.Fpdf, family string) error {
	err := registerFontWithGofpdf(pdf, family)
	if err != nil {
		return err
	}
	pdf.SetFont(family, "", 10)
	return nil
}

// StringWidth calculates the width of a string in the current font.
func StringWidth(pdf *gofpdf.Fpdf, s string) float64 {
	return pdf.GetStringWidth(s)
}

// pointToMM converts points to millimeters.
func pointToMM(pt float64) float64 {
	return pt / 72 * 25.4
}

// mmToPoint converts millimeters to points.
func mmToPoint(mm float64) float64 {
	return mm / 25.4 * 72
}

// bytesReadCloser wraps a byte slice as an io.ReadCloser.
type bytesReadCloser struct {
	*bytes.Reader
}

func (b *bytesReadCloser) Close() error {
	return nil
}
