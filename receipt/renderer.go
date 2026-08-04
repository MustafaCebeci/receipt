package receipt

import (
	"github.com/phpdave11/gofpdf"
)

// renderer handles all gofpdf communication.
// Only this file should import gofpdf.
type renderer struct {
	pdf           *gofpdf.Fpdf
	currentY      float64
	contentWidth  float64
	leftMargin    float64
	topMargin     float64
	defaultFont   string
	defaultSize   float64
}

// newRenderer creates a new renderer instance.
func newRenderer(filename string, doc *Document) (*renderer, error) {
	pdf := gofpdf.New(string(doc.Orientation), "mm", "", "")

	// Set document settings
	pdf.SetAutoPageBreak(doc.AutoPageBreak, doc.PageBreakMargin)

	// Calculate dimensions
	pageWidth := mmToPoint(doc.PaperWidth)
	pageHeight := mmToPoint(doc.PageHeight)

	if doc.PageHeight > 0 {
		pdf.AddPageFormat(string(doc.Orientation), gofpdf.SizeType{
			Wd: pageWidth,
			Ht: pageHeight,
		})
	} else {
		pdf.AddPage()
	}

	// Set margins
	pdf.SetMargins(doc.Margins.Left, doc.Margins.Top, doc.Margins.Right)

	// Calculate content width
	contentWidth := doc.Margins.ContentWidth()

	// Set default font
	defaultFont := "DejaVu"
	err := registerFontWithGofpdf(pdf, defaultFont)
	if err != nil {
		// Fallback to built-in font
		pdf.SetFont("Helvetica", "", 10)
		defaultFont = "Helvetica"
	} else {
		pdf.SetFont(defaultFont, "", 10)
	}

	r := &renderer{
		pdf:          pdf,
		currentY:     pdf.GetY(),
		contentWidth: contentWidth,
		leftMargin:   doc.Margins.Left,
		topMargin:    doc.Margins.Top,
		defaultFont:  defaultFont,
		defaultSize:  10,
	}

	return r, nil
}

// GetPDF returns the underlying gofpdf instance.
func (r *renderer) GetPDF() *gofpdf.Fpdf {
	return r.pdf
}

// SetFont sets the current font.
func (r *renderer) SetFont(family string, style string, size float64) {
	if family == "" {
		family = r.defaultFont
	}
	r.pdf.SetFont(family, style, size)
}

// SetFontSize sets the current font size.
func (r *renderer) SetFontSize(size float64) {
	r.pdf.SetFontSize(size)
}

// GetStringWidth returns the width of a string in the current font.
func (r *renderer) GetStringWidth(s string) float64 {
	return r.pdf.GetStringWidth(s)
}

// GetY returns the current Y position.
func (r *renderer) GetY() float64 {
	return r.pdf.GetY()
}

// SetY sets the Y position.
func (r *renderer) SetY(y float64) {
	r.pdf.SetY(y)
}

// GetX returns the current X position.
func (r *renderer) GetX() float64 {
	return r.pdf.GetX()
}

// SetX sets the X position.
func (r *renderer) SetX(x float64) {
	r.pdf.SetX(x)
}

// GetContentWidth returns the content width.
func (r *renderer) GetContentWidth() float64 {
	return r.contentWidth
}

// AddPage adds a new page.
func (r *renderer) AddPage() {
	r.pdf.AddPage()
	r.currentY = r.pdf.GetY()
}

// drawText draws text at the specified position with alignment.
func (r *renderer) drawText(text string, x, y, width float64, align HorizontalAlign) {
	if text == "" {
		return
	}

	textWidth := r.pdf.GetStringWidth(text)

	switch align {
	case AlignCenter:
		x = x + (width - textWidth) / 2
	case AlignRight:
		x = x + width - textWidth
	}

	r.pdf.Text(x, y, text)
}

// drawLine draws a horizontal line.
func (r *renderer) drawLine(x1, y, x2 float64, dashed bool) {
	if dashed {
		// Draw dashed line using multiple short segments
		segmentLength := 3.0
		gapLength := 2.0
		currentX := x1
		for currentX < x2 {
			nextX := currentX + segmentLength
			if nextX > x2 {
				nextX = x2
			}
			r.pdf.Line(currentX, y, nextX, y)
			currentX = nextX + gapLength
		}
	} else {
		r.pdf.Line(x1, y, x2, y)
	}
}

// drawDashedLine draws a dashed horizontal line.
func (r *renderer) drawDashedLine(x1, y, x2 float64) {
	r.drawLine(x1, y, x2, true)
}

// Space adds vertical space.
func (r *renderer) Space(mm float64) {
	r.currentY += mmToPoint(mm)
	r.pdf.SetY(r.currentY)
}

// Title renders a title text centered.
func (r *renderer) Title(text string, size float64) {
	if size <= 0 {
		size = 14
	}

	r.pdf.SetFont(r.defaultFont, "B", size)
	r.pdf.SetTextColor(0, 0, 0)

	titleWidth := r.pdf.GetStringWidth(text)
	x := r.leftMargin + (r.contentWidth-titleWidth)/2

	r.pdf.Text(x, r.currentY+size, text)
	r.currentY += size * 1.5
	r.pdf.SetY(r.currentY)
}

// Subtitle renders a subtitle text centered.
func (r *renderer) Subtitle(text string, size float64) {
	if size <= 0 {
		size = 12
	}

	r.pdf.SetFont(r.defaultFont, "", size)
	r.pdf.SetTextColor(0, 0, 0)

	subtitleWidth := r.pdf.GetStringWidth(text)
	x := r.leftMargin + (r.contentWidth-subtitleWidth)/2

	r.pdf.Text(x, r.currentY+size, text)
	r.currentY += size * 1.3
	r.pdf.SetY(r.currentY)
}

// Cell renders a single cell of text.
func (r *renderer) Cell(text string, x, y, width, height float64, align HorizontalAlign, valign VerticalAlign, style *StyleContext) {
	if text == "" {
		return
	}

	// Set font
	fontStyle := ""
	if style.Bold {
		fontStyle += "B"
	}
	if style.Italic {
		fontStyle += "I"
	}
	if style.Underline {
		fontStyle += "U"
	}

	fontSize := style.FontSize
	if style.Small {
		fontSize = 8
	} else if style.Header {
		fontSize = 10
	} else if style.Footer {
		fontSize = 8
	}

	r.pdf.SetFont(style.FontFamily, fontStyle, fontSize)
	r.pdf.SetTextColor(0, 0, 0)

	// Calculate text position
	textX := x + mmToPoint(style.PaddingLeft)

	// Horizontal alignment
	switch align {
	case AlignCenter:
		textX = x + width/2 - r.pdf.GetStringWidth(text)/2
	case AlignRight:
		textX = x + width - r.pdf.GetStringWidth(text) - mmToPoint(style.PaddingRight)
	}

	// Vertical alignment
	textY := y + mmToPoint(style.PaddingTop)
	lineHeight := fontSize * 1.2

	switch valign {
	case AlignMiddle:
		textY = y + height/2 - lineHeight/2
	case AlignBottom:
		textY = y + height - lineHeight - mmToPoint(style.PaddingBottom)
	}

	r.pdf.Text(textX, textY, text)
}

// KeyValue renders a key-value pair.
func (r *renderer) KeyValue(key, value string, style *StyleContext) {
	r.pdf.SetFont(style.FontFamily, "", style.FontSize)
	r.pdf.SetTextColor(0, 0, 0)

	valueWidth := r.pdf.GetStringWidth(value)

	keyX := r.leftMargin + mmToPoint(style.MarginLeft)
	valueX := r.leftMargin + r.contentWidth - valueWidth - mmToPoint(style.MarginRight)

	r.pdf.Text(keyX, r.currentY, key)
	r.pdf.Text(valueX, r.currentY, value)

	r.currentY += style.FontSize * 1.5
	r.pdf.SetY(r.currentY)
}

// Line renders a horizontal line across the content width.
func (r *renderer) Line(style LineType) {
	y := r.currentY
	x1 := r.leftMargin
	x2 := r.leftMargin + r.contentWidth

	if style == LineDashed {
		r.drawDashedLine(x1, y, x2)
	} else {
		r.drawLine(x1, y, x2, false)
	}

	r.currentY += 1
	r.pdf.SetY(r.currentY)
}

// SpaceMM adds vertical space of the specified millimeters.
func (r *renderer) SpaceMM(mm float64) {
	r.Space(mm)
}

// Save saves the PDF to a file.
func (r *renderer) Save(filename string) error {
	err := r.pdf.OutputFileAndClose(filename)
	if err != nil {
		return ErrSaveFailed
	}
	return nil
}

// SetFooter sets a footer callback.
func (r *renderer) SetFooter(fn func(*gofpdf.Fpdf)) {
	r.pdf.SetFooterFunc(func() {
		fn(r.pdf)
	})
}

// SetHeader sets a header callback.
func (r *renderer) SetHeader(fn func(*gofpdf.Fpdf)) {
	r.pdf.SetHeaderFunc(func() {
		fn(r.pdf)
	})
}
