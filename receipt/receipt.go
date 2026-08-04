// Package receipt provides a fluent API for generating 80mm thermal POS receipts.
// It wraps gofpdf completely, hiding all third-party API from the developer.
//
// Example usage:
//
//	pdf := receipt.New("receipt.pdf")
//	pdf.Title("SWIFTY CAFE")
//	pdf.Line()
//	pdf.Table(receipt.Columns(50, 10, 20))
//	pdf.Row(
//	    receipt.Cell("Latte"),
//	    receipt.Cell("2", receipt.Center()),
//	    receipt.Cell(receipt.Money(270), receipt.Right()),
//	)
//	pdf.Line()
//	pdf.KeyValue("Toplam", receipt.Money(270))
//	pdf.QRCode("https://example.com")
//	pdf.Save()
package receipt

import (
	"bytes"
	"fmt"

	"github.com/phpdave11/gofpdf"
	"github.com/skip2/go-qrcode"
)

// Receipt is the main type for generating receipts.
// It provides a fluent API for building receipt content.
type Receipt struct {
	filename    string
	renderer    *renderer
	doc         *Document
	currentTable *Table
	tableRows   []*row
	headerFn    func(*Receipt)
	footerFn    func(*Receipt)
	closed      bool
}

// New creates a new Receipt for the specified filename.
// The filename should have a .pdf extension.
func New(filename string) *Receipt {
	doc := NewDocument()

	r, err := newRenderer(filename, doc)
	if err != nil {
		// If renderer creation fails, we still return a Receipt
		// but Save() will fail with an appropriate error
		return &Receipt{
			filename: filename,
			doc:      doc,
		}
	}

	return &Receipt{
		filename: filename,
		renderer: r,
		doc:      doc,
	}
}

// Paper80 is a constant for 80mm paper (kept for API compatibility).
const Paper80 = 80.0

// Title sets the receipt title (centered, bold, large).
func (r *Receipt) Title(text string) *Receipt {
	if r.renderer == nil {
		return r
	}
	r.renderer.Title(text, 14)
	return r
}

// SubTitle sets the receipt subtitle (centered).
func (r *Receipt) SubTitle(text string) *Receipt {
	if r.renderer == nil {
		return r
	}
	r.renderer.Subtitle(text, 12)
	return r
}

// Line draws a solid horizontal line.
func (r *Receipt) Line() *Receipt {
	if r.renderer == nil {
		return r
	}
	r.renderer.Line(LineSolid)
	return r
}

// DashedLine draws a dashed horizontal line.
func (r *Receipt) DashedLine() *Receipt {
	if r.renderer == nil {
		return r
	}
	r.renderer.Line(LineDashed)
	return r
}

// Separator draws a dashed line with space above and below.
func (r *Receipt) Separator() *Receipt {
	if r.renderer == nil {
		return r
	}
	r.renderer.SpaceMM(2)
	r.renderer.Line(LineDashed)
	r.renderer.SpaceMM(2)
	return r
}

// Space adds vertical space of the specified millimeters.
func (r *Receipt) Space(mm float64) *Receipt {
	if r.renderer == nil {
		return r
	}
	r.renderer.SpaceMM(mm)
	return r
}

// Table starts a new table with the specified columns.
// Columns can be created with Columns(50, 10, 20) for fixed widths in mm.
// Use 0 for auto-width columns.
func (r *Receipt) Table(columns []Column) *Receipt {
	if r.renderer == nil {
		return r
	}

	// If there's an existing table being built, finalize it first
	if r.currentTable != nil {
		r.finishTable()
	}

	// Calculate content width for auto columns
	contentWidth := r.doc.Margins.ContentWidth()
	for _, col := range columns {
		if col.Width > 0 {
			contentWidth -= col.Width
		}
	}

	// Create new table
	r.currentTable = TableNew(columns)

	// Calculate column widths
	r.currentTable.CalculateColumnWidths(contentWidth)

	return r
}

// Row adds a row to the current table.
// Must be called after Table().
func (r *Receipt) Row(cells ...*cell) *Receipt {
	if r.renderer == nil {
		return r
	}

	if r.currentTable == nil {
		// Create a default table if none exists
		colWidths := make([]float64, len(cells))
		contentWidth := r.doc.Margins.ContentWidth()
		colWidth := contentWidth / float64(len(cells))
		for i := range cells {
			colWidths[i] = colWidth
		}
		r.currentTable = TableNew(Columns(colWidths...))
		r.currentTable.CalculateColumnWidths(contentWidth)
	}

	row := &row{Cells: cells}
	r.currentTable.AddRow(row)

	return r
}

// finishTable finalizes the current table and renders it.
func (r *Receipt) finishTable() {
	if r.currentTable == nil || r.renderer == nil {
		return
	}

	// Render the table
	x := r.renderer.GetX()
	y := r.renderer.GetY()

	// Get column widths in points
	colWidths := r.currentTable.GetColumnWidths()

	// Calculate total table width
	totalWidth := 0.0
	for _, w := range colWidths {
		totalWidth += w
	}

	// Center the table
	x = r.renderer.GetX() + (mmToPoint(r.doc.Margins.ContentWidth())-totalWidth)/2

	r.currentTable.Render(r.renderer.GetPDF(), x, y)

	// Clear current table
	r.currentTable = nil
}

// KeyValue adds a key-value line.
// The key is left-aligned, the value is right-aligned.
func (r *Receipt) KeyValue(key interface{}, value interface{}, styles ...Style) *Receipt {
	if r.renderer == nil {
		return r
	}

	// Finalize any pending table
	if r.currentTable != nil {
		r.finishTable()
	}

	ctx := DefaultStyleContext()
	Apply(ctx, styles...)

	keyStr := toString(key)
	valueStr := toString(value)

	r.renderer.KeyValue(keyStr, valueStr, ctx)
	return r
}

// QRCode adds a QR code to the receipt.
// The QR code is centered and sized appropriately.
func (r *Receipt) QRCode(content string) *Receipt {
	if r.renderer == nil || content == "" {
		return r
	}

	// Calculate size
	size := QRCodeSize(content)

	// Generate QR code
	qr, err := qrcode.Encode(content, qrcode.Medium, int(size*3))
	if err != nil {
		return r
	}

	// Get PDF reference
	pdf := r.renderer.GetPDF()
	pdf.AddPage()

	// Get page width for centering
	pageWidth, _ := pdf.GetPageSize()

	// Create a unique image name
	imgName := "qrcode_" + content[:min(10, len(content))]

	// Register image
	info := pdf.RegisterImageOptionsReader(imgName, gofpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(qr))
	if info == nil {
		return r
	}

	// Calculate dimensions to fit on page
	imgWidth := info.Width()
	imgHeight := info.Height()

	// Scale to fit page width with margins
	maxWidth := pageWidth - mmToPoint(20) // 10mm margin each side
	if imgWidth > maxWidth {
		scale := maxWidth / imgWidth
		imgWidth *= scale
		imgHeight *= scale
	}

	// Center horizontally
	x := (pageWidth - imgWidth) / 2
	y := r.renderer.GetY()

	// Draw image
	pdf.Image(imgName, x, y, imgWidth, imgHeight, false, "PNG", 0, "")

	// Update Y position
	r.renderer.SetY(y + imgHeight + 5)

	return r
}

// Barcode adds a barcode to the receipt.
func (r *Receipt) Barcode(content string, barcodeType BarcodeType) *Receipt {
	if r.renderer == nil || content == "" {
		return r
	}

	// Finalize any pending table
	if r.currentTable != nil {
		r.finishTable()
	}

	pdf := r.renderer.GetPDF()
	x := r.renderer.GetX()
	y := r.renderer.GetY()

	width := mmToPoint(r.doc.Margins.ContentWidth() * 0.8)
	height := mmToPoint(15)

	err := Barcode(pdf, content, barcodeType, x, y, width, height)
	if err != nil {
		return r
	}

	r.renderer.SetY(y + height + 5)
	return r
}

// Footer sets a footer callback that is called on each page.
func (r *Receipt) Footer(fn func(*Receipt)) *Receipt {
	r.footerFn = fn
	if r.renderer != nil {
		r.renderer.SetFooter(func(p *gofpdf.Fpdf) {
			if r.footerFn != nil {
				// Create a minimal receipt for footer context
				r.footerFn(r)
			}
		})
	}
	return r
}

// Header sets a header callback that is called on each page.
func (r *Receipt) Header(fn func(*Receipt)) *Receipt {
	r.headerFn = fn
	if r.renderer != nil {
		r.renderer.SetHeader(func(p *gofpdf.Fpdf) {
			if r.headerFn != nil {
				r.headerFn(r)
			}
		})
	}
	return r
}

// Save saves the receipt to the file.
// This should be called after all content has been added.
func (r *Receipt) Save() error {
	if r.renderer == nil {
		return ErrSaveFailed
	}

	// Finalize any pending table
	if r.currentTable != nil {
		r.finishTable()
	}

	return r.renderer.Save(r.filename)
}

// GetPDF returns the underlying gofpdf instance (for advanced use cases).
// This should rarely be needed.
func (r *Receipt) GetPDF() *gofpdf.Fpdf {
	if r.renderer == nil {
		return nil
	}
	return r.renderer.GetPDF()
}

// SetFont sets the default font family.
func (r *Receipt) SetFont(family string, size float64) *Receipt {
	if r.renderer == nil {
		return r
	}
	r.renderer.SetFont(family, "", size)
	return r
}

// GetY returns the current Y position.
func (r *Receipt) GetY() float64 {
	if r.renderer == nil {
		return 0
	}
	return r.renderer.GetY()
}

// SetY sets the Y position.
func (r *Receipt) SetY(y float64) *Receipt {
	if r.renderer == nil {
		return r
	}
	r.renderer.SetY(y)
	return r
}

// toString converts an interface{} to string.
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return s
	case fmt.Stringer:
		return s.String()
	case error:
		return s.Error()
	default:
		return fmt.Sprintf("%v", s)
	}
}
