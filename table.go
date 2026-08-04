package receipt

import (
	"github.com/phpdave11/gofpdf"
)

// Column represents a table column.
type Column struct {
	// Width is the column width in millimeters.
	Width float64
	// IsAuto indicates if the column width should be calculated automatically.
	IsAuto bool
}

// Columns creates a slice of Column from width values.
// Widths are in millimeters. Use 0 for auto-width columns.
func Columns(widths ...float64) []Column {
	cols := make([]Column, len(widths))
	autoCount := 0
	fixedWidth := 0.0

	// First pass: count auto columns and sum fixed widths
	for _, w := range widths {
		if w <= 0 {
			autoCount++
		} else {
			fixedWidth += w
		}
	}

	// Second pass: calculate widths
	for i, w := range widths {
		if w <= 0 {
			cols[i] = Column{Width: 0, IsAuto: true}
		} else {
			cols[i] = Column{Width: w, IsAuto: false}
		}
	}

	return cols
}

// Table represents a table with columns and rows.
type Table struct {
	// Columns holds the table columns.
	Columns []Column
	// Rows holds the table rows.
	Rows []*row
	// Styles holds table-level styles.
	Styles []Style
	// columnWidths holds the calculated widths for each column (in points).
	columnWidths []float64
	// totalWidth is the total available width for the table.
	totalWidth float64
	// Styles context for this table.
	styleContext *StyleContext
}

// TableNew creates a new Table with the given columns and optional styles.
func TableNew(columns []Column, styles ...Style) *Table {
	return &Table{
		Columns:     columns,
		Rows:        []*row{},
		Styles:      styles,
		styleContext: DefaultStyleContext(),
	}
}

// AddRow adds a row to the table.
func (t *Table) AddRow(row *row) {
	t.Rows = append(t.Rows, row)
}

// CalculateColumnWidths calculates the actual column widths based on available width.
func (t *Table) CalculateColumnWidths(availableWidth float64) {
	t.totalWidth = availableWidth

	autoCount := 0
	fixedWidth := 0.0

	// Count auto columns and sum fixed widths
	for _, col := range t.Columns {
		if col.IsAuto {
			autoCount++
		} else {
			fixedWidth += col.Width
		}
	}

	// Calculate auto column width
	autoWidth := 0.0
	if autoCount > 0 {
		autoWidth = (availableWidth - fixedWidth) / float64(autoCount)
		if autoWidth < 0 {
			autoWidth = 0
		}
	}

	// Convert mm to points for internal use
	t.columnWidths = make([]float64, len(t.Columns))
	for i, col := range t.Columns {
		if col.IsAuto {
			t.columnWidths[i] = mmToPoint(autoWidth)
		} else {
			t.columnWidths[i] = mmToPoint(col.Width)
		}
	}
}

// GetColumnWidths returns the calculated column widths in points.
func (t *Table) GetColumnWidths() []float64 {
	return t.columnWidths
}

// GetTotalWidth returns the total table width.
func (t *Table) GetTotalWidth() float64 {
	return t.totalWidth
}

// Render renders the table to the PDF.
func (t *Table) Render(pdf *gofpdf.Fpdf, x, y float64) (float64, error) {
	if len(t.Rows) == 0 || len(t.Columns) == 0 {
		return y, nil
	}

	// Apply table styles
	Apply(t.styleContext, t.Styles...)

	// Track current position
	currentY := y
	currentX := x

	// Create row span tracking: tracks how many rows a cell spans
	rowSpanRemaining := make([]int, len(t.Columns))

	for rowIdx, row := range t.Rows {
		// Calculate row height
		rowHeight := row.calculateHeight(pdf, t.columnWidths, t.styleContext)

		// Check if we need a page break
		_, pageHeight := pdf.GetPageSize()
		remainingHeight := pageHeight - pdf.GetY() - mmToPoint(20) // 20mm bottom margin

		if rowHeight > remainingHeight {
			// Need new page
			pdf.AddPage()
			currentY = pdf.GetY()
			// Reset row span tracking for new page
			rowSpanRemaining = make([]int, len(t.Columns))
		}

		// Process cells in this row
		cellX := currentX
		spanCols := make([]bool, len(t.Columns)) // Track which columns are spanned

		for colIdx := 0; colIdx < len(t.Columns); colIdx++ {
			// Skip if this column is spanned from a previous row
			if rowSpanRemaining[colIdx] > 0 {
				rowSpanRemaining[colIdx]--
				// Add to cellX to skip this column
				if colIdx < len(t.columnWidths) {
					cellX += t.columnWidths[colIdx]
				}
				spanCols[colIdx] = true
				continue
			}
			spanCols[colIdx] = false

			if rowIdx >= len(t.Rows) || colIdx >= len(t.Rows[rowIdx].Cells) {
				// No cell for this position
				cellX += t.columnWidths[colIdx]
				continue
			}

			cell := t.Rows[rowIdx].Cells[colIdx]
			if cell == nil {
				cellX += t.columnWidths[colIdx]
				continue
			}

			// Get cell style
			cellStyle := &StyleContext{}
			cellStyle.Merge(t.styleContext)
			Apply(cellStyle, row.Styles...)
			Apply(cellStyle, cell.Styles...)

			// Calculate cell width (handle colSpan)
			cellWidth := t.columnWidths[colIdx]
			if cell.colSpan > 1 {
				for j := colIdx + 1; j < colIdx+cell.colSpan && j < len(t.columnWidths); j++ {
					cellWidth += t.columnWidths[j]
				}
			}

			// Set up font for height calculation
			fontStyle := ""
			if cellStyle.Bold {
				fontStyle += "B"
			}
			if cellStyle.Italic {
				fontStyle += "I"
			}

			fontSize := cellStyle.FontSize
			if cellStyle.Small {
				fontSize = 8
			} else if cellStyle.Header {
				fontSize = 10
			}

			pdf.SetFont(cellStyle.FontFamily, fontStyle, fontSize)

			// Calculate content height
			contentHeight := cell.calculateHeight(pdf, pointToMM(cellWidth), cellStyle)

			// Apply vertical alignment
			cellY := currentY
			if cellStyle.VerticalAlign == AlignMiddle {
				cellY = currentY + (rowHeight-contentHeight)/2
			} else if cellStyle.VerticalAlign == AlignBottom {
				cellY = currentY + rowHeight - contentHeight
			}

			// Draw cell background (white for now)
			pdf.SetFillColor(255, 255, 255)
			pdf.Rect(cellX, currentY, cellWidth, rowHeight, "F")

			// Draw borders
			if cellStyle.BorderLeft || cellStyle.BorderTop || cellStyle.BorderRight || cellStyle.BorderBottom {
				pdf.SetDrawColor(0, 0, 0)
				if cellStyle.BorderLeft {
					pdf.Line(cellX, currentY, cellX, currentY+rowHeight)
				}
				if cellStyle.BorderTop {
					pdf.Line(cellX, currentY, cellX+cellWidth, currentY)
				}
				if cellStyle.BorderRight {
					pdf.Line(cellX+cellWidth, currentY, cellX+cellWidth, currentY+rowHeight)
				}
				if cellStyle.BorderBottom {
					pdf.Line(cellX, currentY+rowHeight, cellX+cellWidth, currentY+rowHeight)
				}
			}

			// Calculate text position with padding
			textX := cellX + mmToPoint(cellStyle.PaddingLeft)
			textWidth := cellWidth - mmToPoint(cellStyle.PaddingLeft) - mmToPoint(cellStyle.PaddingRight)
			textY := cellY + mmToPoint(cellStyle.PaddingTop)

			// Calculate horizontal text position
			text := cell.String()

			if cellStyle.HorizontalAlign == AlignCenter {
				textX = cellX + cellWidth/2 - pdf.GetStringWidth(text)/2
			} else if cellStyle.HorizontalAlign == AlignRight {
				textX = cellX + cellWidth - pdf.GetStringWidth(text) - mmToPoint(cellStyle.PaddingRight)
			}

			// Draw text
			if cellStyle.Underline {
				pdf.SetFont(cellStyle.FontFamily, fontStyle+"U", fontSize)
			}
			if cellStyle.Strike {
				pdf.SetFont(cellStyle.FontFamily, fontStyle+"D", fontSize)
			}

			pdf.SetTextColor(0, 0, 0)
			pdf.SetXY(textX, textY)
			pdf.Cell(textWidth, mmToPoint(fontSize*1.2/2), text)

			// Update row span tracking
			if cell.rowSpan > 1 {
				rowSpanRemaining[colIdx] = cell.rowSpan - 1
			}

			// Move to next column position
			cellX += cellWidth
		}

		// Update Y position for next row
		currentY += rowHeight
	}

	// Update PDF Y position
	pdf.SetY(currentY)

	return currentY, nil
}
