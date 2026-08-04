package receipt

// Paper80WidthMM is the width of 80mm thermal paper in millimeters.
const Paper80WidthMM = 80.0

// Paper80WidthPoints is the width of 80mm thermal paper in PDF points.
const Paper80WidthPoints = Paper80WidthMM * 72 / 25.4 // ~226.77

// Margins holds the document margin settings.
type Margins struct {
	Left   float64
	Right  float64
	Top    float64
	Bottom float64
}

// DefaultMargins returns the default margins for 80mm receipt.
func DefaultMargins() Margins {
	return Margins{
		Left:   5,
		Right:  5,
		Top:    10,
		Bottom: 10,
	}
}

// ContentWidth returns the usable content width in millimeters.
func (m Margins) ContentWidth() float64 {
	return Paper80WidthMM - m.Left - m.Right
}

// ContentWidthPoints returns the usable content width in points.
func (m Margins) ContentWidthPoints() float64 {
	return mmToPoint(m.ContentWidth())
}

// Document holds document-level settings.
type Document struct {
	// Paper size (fixed to 80mm for this package)
	PaperWidth float64
	// Page height (0 means unlimited)
	PageHeight float64
	// Orientation
	Orientation Orientation
	// Margins
	Margins Margins
	// AutoPageBreak enables automatic page breaks
	AutoPageBreak bool
	// PageBreakMargin is the margin at page break
	PageBreakMargin float64
}

// Orientation represents page orientation.
type Orientation string

const (
	// Portrait orientation
	Portrait Orientation = "P"
	// Landscape orientation
	Landscape Orientation = "L"
)

// NewDocument creates a new document with default settings for 80mm thermal paper.
func NewDocument() *Document {
	return &Document{
		PaperWidth:      Paper80WidthMM,
		PageHeight:      0, // Unlimited
		Orientation:     Portrait,
		Margins:          DefaultMargins(),
		AutoPageBreak:    true,
		PageBreakMargin:  10,
	}
}

// SetOrientation sets the page orientation.
func (d *Document) SetOrientation(o Orientation) {
	d.Orientation = o
}

// SetMargins sets the document margins.
func (d *Document) SetMargins(left, right, top, bottom float64) {
	d.Margins = Margins{
		Left:   left,
		Right:  right,
		Top:    top,
		Bottom: bottom,
	}
}
