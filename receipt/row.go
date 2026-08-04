package receipt

import (
	"github.com/phpdave11/gofpdf"
)

// row represents a table row with cells.
type row struct {
	// Cells holds the cells in this row.
	Cells []*cell
	// Styles holds row-level styles.
	Styles []Style
	// height is the calculated height of this row.
	height float64
	// isHeader indicates if this is a header row.
	isHeader bool
}

// NewRow creates a new Row with the given cells and optional styles.
func NewRow(cells []*cell, styles ...Style) *row {
	return &row{
		Cells:  cells,
		Styles: styles,
	}
}

// calculateHeight calculates the maximum height required by any cell in the row.
func (r *row) calculateHeight(pdf *gofpdf.Fpdf, columnWidths []float64, style *StyleContext) float64 {
	if len(r.Cells) == 0 {
		return 0
	}

	maxHeight := 0.0

	for i, cell := range r.Cells {
		if cell.rowSpan > 1 {
			// Skip cells that will be filled by row span from previous rows
			continue
		}

		var colWidth float64
		if i < len(columnWidths) {
			colWidth = columnWidths[i]
		} else {
			colWidth = 50 // Default fallback
		}

		cellStyle := &StyleContext{}
		cellStyle.Merge(style)
		Apply(cellStyle, cell.Styles...)

		// Handle colSpan
		if cell.colSpan > 1 && i+cell.colSpan <= len(columnWidths) {
			for j := i + 1; j < i+cell.colSpan; j++ {
				colWidth += columnWidths[j]
			}
		}

		height := cell.calculateHeight(pdf, colWidth, cellStyle)
		if height > maxHeight {
			maxHeight = height
		}
	}

	// Apply row-level height constraints
	if r.height > 0 && r.height > maxHeight {
		maxHeight = r.height
	}

	return maxHeight
}

// GetHeight returns the calculated height.
func (r *row) GetHeight() float64 {
	return r.height
}

// SetHeight sets the row height.
func (r *row) SetHeight(h float64) {
	r.height = h
}

// IsHeader returns true if this is a header row.
func (r *row) IsHeader() bool {
	return r.isHeader
}

// SetHeader sets whether this is a header row.
func (r *row) SetHeader(isHeader bool) {
	r.isHeader = isHeader
}
