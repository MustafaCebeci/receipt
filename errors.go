// Package receipt provides a fluent API for generating 80mm thermal POS receipts.
// It wraps gofpdf completely, hiding all third-party API from the developer.
package receipt

import "errors"

// Package-level errors.
var (
	// ErrInvalidWidth is returned when a column width is not positive.
	ErrInvalidWidth = errors.New("receipt: column width must be positive")

	// ErrCellNotInTable is returned when a cell is rendered outside a table.
	ErrCellNotInTable = errors.New("receipt: cell rendered outside table")

	// ErrNilContent is returned when a cell has nil content.
	ErrNilContent = errors.New("receipt: cell content is nil")

	// ErrTableNotStarted is returned when Row/Cell is called without Table.
	ErrTableNotStarted = errors.New("receipt: Row/Cell called without Table")

	// ErrRowSpanOverflow is returned when RowSpan exceeds remaining rows.
	ErrRowSpanOverflow = errors.New("receipt: RowSpan exceeds remaining rows")

	// ErrSaveFailed is returned when the PDF cannot be saved.
	ErrSaveFailed = errors.New("receipt: failed to save PDF")

	// ErrInvalidPaperSize is returned when an unsupported paper size is used.
	ErrInvalidPaperSize = errors.New("receipt: invalid paper size")

	// ErrInvalidAlignment is returned when an invalid alignment is specified.
	ErrInvalidAlignment = errors.New("receipt: invalid alignment")

	// ErrFontLoading is returned when a font cannot be loaded.
	ErrFontLoading = errors.New("receipt: failed to load font")

	// ErrQRCodeGeneration is returned when QR code generation fails.
	ErrQRCodeGeneration = errors.New("receipt: failed to generate QR code")

	// ErrBarcodeGeneration is returned when barcode generation fails.
	ErrBarcodeGeneration = errors.New("receipt: failed to generate barcode")

	// ErrBarcodeType is returned when an unsupported barcode type is used.
	ErrBarcodeType = errors.New("receipt: unsupported barcode type")
)
