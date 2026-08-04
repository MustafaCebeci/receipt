package receipt

// Style represents a composable style that can be applied to cells.
type Style interface {
	apply(s *StyleContext)
}

// StyleFunc is a function-style style that modifies a StyleContext.
type StyleFunc func(s *StyleContext)

func (f StyleFunc) apply(s *StyleContext) {
	f(s)
}

// StyleContext holds the current style state.
type StyleContext struct {
	// Font
	FontFamily string
	FontSize   float64
	Bold       bool
	Italic     bool
	Underline  bool
	Strike     bool

	// Alignment
	HorizontalAlign HorizontalAlign
	VerticalAlign   VerticalAlign

	// Spacing
	PaddingLeft   float64
	PaddingRight  float64
	PaddingTop    float64
	PaddingBottom float64

	MarginLeft   float64
	MarginRight  float64
	MarginTop    float64
	MarginBottom float64

	// Border
	BorderLeft   bool
	BorderRight  bool
	BorderTop    bool
	BorderBottom bool

	// Dimensions
	Width    float64
	Height   float64
	MinHeight float64
	MaxHeight float64

	// Span
	ColSpan int
	RowSpan int

	// Text
	Wrap    bool
	NoWrap  bool

	// Size presets
	Small   bool
	Normal  bool
	Header  bool
	Footer  bool
	Title   bool
	Subtitle bool
}

// DefaultStyleContext returns a StyleContext with default values.
func DefaultStyleContext() *StyleContext {
	return &StyleContext{
		FontSize:        10,
		HorizontalAlign: AlignLeft,
		VerticalAlign:   AlignTop,
		Wrap:            true,
		PaddingLeft:     2,
		PaddingRight:    2,
		PaddingTop:      2,
		PaddingBottom:   2,
		MarginLeft:      0,
		MarginRight:     0,
		MarginTop:       0,
		MarginBottom:    0,
		ColSpan:         1,
		RowSpan:         1,
	}
}

// Merge merges other style contexts into this one (later ones override).
func (s *StyleContext) Merge(others ...*StyleContext) {
	for _, other := range others {
		if other == nil {
			continue
		}
		if other.FontFamily != "" {
			s.FontFamily = other.FontFamily
		}
		if other.FontSize > 0 {
			s.FontSize = other.FontSize
		}
		if other.Bold {
			s.Bold = true
		}
		if other.Italic {
			s.Italic = true
		}
		if other.Underline {
			s.Underline = true
		}
		if other.Strike {
			s.Strike = true
		}
		if other.HorizontalAlign != 0 {
			s.HorizontalAlign = other.HorizontalAlign
		}
		if other.VerticalAlign != 0 {
			s.VerticalAlign = other.VerticalAlign
		}
		if other.PaddingLeft > 0 {
			s.PaddingLeft = other.PaddingLeft
		}
		if other.PaddingRight > 0 {
			s.PaddingRight = other.PaddingRight
		}
		if other.PaddingTop > 0 {
			s.PaddingTop = other.PaddingTop
		}
		if other.PaddingBottom > 0 {
			s.PaddingBottom = other.PaddingBottom
		}
		if other.MarginLeft > 0 {
			s.MarginLeft = other.MarginLeft
		}
		if other.MarginRight > 0 {
			s.MarginRight = other.MarginRight
		}
		if other.MarginTop > 0 {
			s.MarginTop = other.MarginTop
		}
		if other.MarginBottom > 0 {
			s.MarginBottom = other.MarginBottom
		}
		if other.BorderLeft {
			s.BorderLeft = true
		}
		if other.BorderRight {
			s.BorderRight = true
		}
		if other.BorderTop {
			s.BorderTop = true
		}
		if other.BorderBottom {
			s.BorderBottom = true
		}
		if other.Width > 0 {
			s.Width = other.Width
		}
		if other.Height > 0 {
			s.Height = other.Height
		}
		if other.MinHeight > 0 {
			s.MinHeight = other.MinHeight
		}
		if other.MaxHeight > 0 {
			s.MaxHeight = other.MaxHeight
		}
		if other.ColSpan > 1 {
			s.ColSpan = other.ColSpan
		}
		if other.RowSpan > 1 {
			s.RowSpan = other.RowSpan
		}
		if other.Wrap {
			s.Wrap = true
		}
		if other.NoWrap {
			s.NoWrap = true
		}
		if other.Small {
			s.Small = true
		}
		if other.Normal {
			s.Normal = true
		}
		if other.Header {
			s.Header = true
		}
		if other.Footer {
			s.Footer = true
		}
		if other.Title {
			s.Title = true
		}
		if other.Subtitle {
			s.Subtitle = true
		}
	}
}

// Apply applies styles to a context.
func Apply(s *StyleContext, styles ...Style) {
	for _, style := range styles {
		if style != nil {
			style.apply(s)
		}
	}
}

// HorizontalAlign represents horizontal text alignment.
type HorizontalAlign int

const (
	// AlignLeft aligns text to the left.
	AlignLeft HorizontalAlign = iota
	// AlignCenter aligns text to the center.
	AlignCenter
	// AlignRight aligns text to the right.
	AlignRight
)

// VerticalAlign represents vertical text alignment.
type VerticalAlign int

const (
	// AlignTop aligns content to the top.
	AlignTop VerticalAlign = iota
	// AlignMiddle aligns content to the middle.
	AlignMiddle
	// AlignBottom aligns content to the bottom.
	AlignBottom
)

// Bold returns a style that makes text bold.
func Bold() Style {
	return StyleFunc(func(s *StyleContext) {
		s.Bold = true
	})
}

// Italic returns a style that makes text italic.
func Italic() Style {
	return StyleFunc(func(s *StyleContext) {
		s.Italic = true
	})
}

// Underline returns a style that underlines text.
func Underline() Style {
	return StyleFunc(func(s *StyleContext) {
		s.Underline = true
	})
}

// Strike returns a style that strikes through text.
func Strike() Style {
	return StyleFunc(func(s *StyleContext) {
		s.Strike = true
	})
}

// Size returns a style that sets the font size.
func Size(size float64) Style {
	return StyleFunc(func(s *StyleContext) {
		s.FontSize = size
	})
}

// Font returns a style that sets the font family.
func Font(family string) Style {
	return StyleFunc(func(s *StyleContext) {
		s.FontFamily = family
	})
}

// TextColor returns a style that sets the text color (placeholder for compatibility).
func TextColor(c Color) Style {
	return StyleFunc(func(s *StyleContext) {
		// Colors are not used in text-only mode, but kept for API compatibility
	})
}

// Fill returns a style that sets the fill color (placeholder for compatibility).
func Fill(c Color) Style {
	return StyleFunc(func(s *StyleContext) {
		// Fill colors not supported in text-only mode
	})
}

// Left returns a style that left-aligns text.
func Left() Style {
	return StyleFunc(func(s *StyleContext) {
		s.HorizontalAlign = AlignLeft
	})
}

// Center returns a style that centers text.
func Center() Style {
	return StyleFunc(func(s *StyleContext) {
		s.HorizontalAlign = AlignCenter
	})
}

// Right returns a style that right-aligns text.
func Right() Style {
	return StyleFunc(func(s *StyleContext) {
		s.HorizontalAlign = AlignRight
	})
}

// Top returns a style that top-aligns content.
func Top() Style {
	return StyleFunc(func(s *StyleContext) {
		s.VerticalAlign = AlignTop
	})
}

// Middle returns a style that middle-aligns content.
func Middle() Style {
	return StyleFunc(func(s *StyleContext) {
		s.VerticalAlign = AlignMiddle
	})
}

// Bottom returns a style that bottom-aligns content.
func Bottom() Style {
	return StyleFunc(func(s *StyleContext) {
		s.VerticalAlign = AlignBottom
	})
}

// Padding returns a style that sets all padding values.
func Padding(padding float64) Style {
	return StyleFunc(func(s *StyleContext) {
		s.PaddingLeft = padding
		s.PaddingRight = padding
		s.PaddingTop = padding
		s.PaddingBottom = padding
	})
}

// PaddingHorizontal returns a style that sets left and right padding.
func PaddingHorizontal(padding float64) Style {
	return StyleFunc(func(s *StyleContext) {
		s.PaddingLeft = padding
		s.PaddingRight = padding
	})
}

// PaddingVertical returns a style that sets top and bottom padding.
func PaddingVertical(padding float64) Style {
	return StyleFunc(func(s *StyleContext) {
		s.PaddingTop = padding
		s.PaddingBottom = padding
	})
}

// Margin returns a style that sets all margin values.
func Margin(margin float64) Style {
	return StyleFunc(func(s *StyleContext) {
		s.MarginLeft = margin
		s.MarginRight = margin
		s.MarginTop = margin
		s.MarginBottom = margin
	})
}

// MarginHorizontal returns a style that sets left and right margin.
func MarginHorizontal(margin float64) Style {
	return StyleFunc(func(s *StyleContext) {
		s.MarginLeft = margin
		s.MarginRight = margin
	})
}

// MarginVertical returns a style that sets top and bottom margin.
func MarginVertical(margin float64) Style {
	return StyleFunc(func(s *StyleContext) {
		s.MarginTop = margin
		s.MarginBottom = margin
	})
}

// Border returns a style that enables all borders.
func Border() Style {
	return StyleFunc(func(s *StyleContext) {
		s.BorderLeft = true
		s.BorderRight = true
		s.BorderTop = true
		s.BorderBottom = true
	})
}

// BorderTop returns a style that enables top border.
func BorderTop() Style {
	return StyleFunc(func(s *StyleContext) {
		s.BorderTop = true
	})
}

// BorderBottom returns a style that enables bottom border.
func BorderBottom() Style {
	return StyleFunc(func(s *StyleContext) {
		s.BorderBottom = true
	})
}

// BorderLeft returns a style that enables left border.
func BorderLeft() Style {
	return StyleFunc(func(s *StyleContext) {
		s.BorderLeft = true
	})
}

// BorderRight returns a style that enables right border.
func BorderRight() Style {
	return StyleFunc(func(s *StyleContext) {
		s.BorderRight = true
	})
}

// BorderColor returns a style that sets border color (placeholder for compatibility).
func BorderColor(c Color) Style {
	return StyleFunc(func(s *StyleContext) {
		// Border colors not supported in text-only mode
	})
}

// ColSpan returns a style that sets the column span.
func ColSpan(n int) Style {
	return StyleFunc(func(s *StyleContext) {
		s.ColSpan = n
	})
}

// RowSpan returns a style that sets the row span.
func RowSpan(n int) Style {
	return StyleFunc(func(s *StyleContext) {
		s.RowSpan = n
	})
}

// Width returns a style that sets the width.
func Width(w float64) Style {
	return StyleFunc(func(s *StyleContext) {
		s.Width = w
	})
}

// Height returns a style that sets the height.
func Height(h float64) Style {
	return StyleFunc(func(s *StyleContext) {
		s.Height = h
	})
}

// MinHeight returns a style that sets the minimum height.
func MinHeight(h float64) Style {
	return StyleFunc(func(s *StyleContext) {
		s.MinHeight = h
	})
}

// MaxHeight returns a style that sets the maximum height.
func MaxHeight(h float64) Style {
	return StyleFunc(func(s *StyleContext) {
		s.MaxHeight = h
	})
}

// Wrap returns a style that enables text wrapping.
func Wrap() Style {
	return StyleFunc(func(s *StyleContext) {
		s.Wrap = true
		s.NoWrap = false
	})
}

// NoWrap returns a style that disables text wrapping.
func NoWrap() Style {
	return StyleFunc(func(s *StyleContext) {
		s.NoWrap = true
		s.Wrap = false
	})
}

// Small returns a style that sets small font size (8pt).
func Small() Style {
	return StyleFunc(func(s *StyleContext) {
		s.Small = true
		s.FontSize = 8
	})
}

// Normal returns a style that sets normal font size (10pt).
func Normal() Style {
	return StyleFunc(func(s *StyleContext) {
		s.Normal = true
		s.FontSize = 10
	})
}

// Header returns a style suitable for table headers.
func Header() Style {
	return StyleFunc(func(s *StyleContext) {
		s.Header = true
		s.Bold = true
		s.FontSize = 10
	})
}

// Footer returns a style suitable for page footers.
func Footer() Style {
	return StyleFunc(func(s *StyleContext) {
		s.Footer = true
		s.FontSize = 8
	})
}

// TitleStyle returns a style suitable for titles.
func TitleStyle() Style {
	return StyleFunc(func(s *StyleContext) {
		s.Title = true
		s.Bold = true
		s.FontSize = 14
		s.HorizontalAlign = AlignCenter
	})
}

// SubtitleStyle returns a style suitable for subtitles.
func SubtitleStyle() Style {
	return StyleFunc(func(s *StyleContext) {
		s.Subtitle = true
		s.Bold = true
		s.FontSize = 12
		s.HorizontalAlign = AlignCenter
	})
}
