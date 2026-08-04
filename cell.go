package receipt

import (
	"fmt"
	"strings"

	"github.com/phpdave11/gofpdf"
)

// cell represents a table cell with content and styles.
type cell struct {
	// Content is the cell's text content.
	Content interface{}
	// Styles holds the composable styles for this cell.
	Styles []Style
	// colSpan tracks how many columns this cell spans.
	colSpan int
	// rowSpan tracks how many rows this cell spans.
	rowSpan int
	// width is the cell width in mm.
	width float64
}

// Cell creates a new Cell with the given content and styles.
// Example: receipt.Cell("Latte", receipt.Bold())
func Cell(content interface{}, styles ...Style) *cell {
	c := &cell{
		Content: content,
		Styles:  styles,
	}

	// Apply styles to extract span information
	ctx := DefaultStyleContext()
	Apply(ctx, styles...)

	c.colSpan = ctx.ColSpan
	c.rowSpan = ctx.RowSpan

	return c
}

// String returns the cell content as a string.
func (c *cell) String() string {
	if c.Content == nil {
		return ""
	}
	switch v := c.Content.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case error:
		return v.Error()
	default:
		return fmt.Sprint(v)
	}
}

// calculateHeight calculates the required height for the cell content.
func (c *cell) calculateHeight(pdf *gofpdf.Fpdf, availableWidth float64, style *StyleContext) float64 {
	if c.width > 0 {
		availableWidth = c.width
	}

	// Account for padding
	usableWidth := availableWidth - style.PaddingLeft - style.PaddingRight

	// Account for margin
	usableWidth -= style.MarginLeft + style.MarginRight

	if usableWidth <= 0 {
		usableWidth = availableWidth
	}

	// Get text content
	text := c.String()
	if text == "" {
		return style.PaddingTop + style.PaddingBottom
	}

	// Calculate font height
	fontSize := style.FontSize
	if style.Small {
		fontSize = 8
	} else if style.Header || style.Title {
		fontSize = 12
	} else if style.Footer {
		fontSize = 8
	}

	// Set font to calculate string width
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
	if style.Strike {
		fontStyle += "D"
	}

	pdf.SetFont(style.FontFamily, fontStyle, fontSize)

	// If no wrap is requested, return single line height
	if style.NoWrap || !style.Wrap {
		lineHeight := fontSize * 1.2
		return lineHeight + style.PaddingTop + style.PaddingBottom
	}

	// Calculate wrapped text height
	lines := c.wrapText(pdf, text, usableWidth)
	lineHeight := fontSize * 1.2
	height := float64(lines) * lineHeight

	// Add padding
	height += style.PaddingTop + style.PaddingBottom

	// Apply min/max height constraints
	if style.MinHeight > 0 && height < style.MinHeight {
		height = style.MinHeight
	}
	if style.MaxHeight > 0 && height > style.MaxHeight {
		height = style.MaxHeight
	}

	return height
}

// wrapText wraps text to fit within the given width and returns the number of lines.
func (c *cell) wrapText(pdf *gofpdf.Fpdf, text string, width float64) int {
	if width <= 0 {
		return 1
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return 1
	}

	lineCount := 1
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		testWidth := pdf.GetStringWidth(testLine)
		if testWidth > width {
			if currentLine == "" {
				// Single word is too long, count it as one line
				return len(text) / int(width/5) + 1
			}
			lineCount++
			currentLine = word
		} else {
			currentLine = testLine
		}
	}

	return lineCount
}

// GetColSpan returns the column span.
func (c *cell) GetColSpan() int {
	return c.colSpan
}

// GetRowSpan returns the row span.
func (c *cell) GetRowSpan() int {
	return c.rowSpan
}
