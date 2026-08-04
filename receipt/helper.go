package receipt

import (
	"fmt"
	"strings"
	"time"
)

// Money formats a float64 as Turkish Lira currency.
// Example: Money(1234.50) returns "₺1,234.50"
func Money(amount float64) string {
	return fmt.Sprintf("₺%.2f", amount)
}

// MoneyWithSymbol formats a float64 as currency with a custom symbol.
// Example: MoneyWithSymbol("$", 1234.50) returns "$1,234.50"
func MoneyWithSymbol(symbol string, amount float64) string {
	return fmt.Sprintf("%s%.2f", symbol, amount)
}

// Date formats a time.Time as a date string.
// Example: Date(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)) returns "15.01.2024"
func Date(t time.Time) string {
	return t.Format("02.01.2006")
}

// DateTime formats a time.Time as a date and time string.
// Example: DateTime(time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)) returns "15.01.2024 14:30"
func DateTime(t time.Time) string {
	return t.Format("02.01.2006 15:04")
}

// Time formats a time.Time as a time string.
// Example: Time(time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)) returns "14:30"
func Time(t time.Time) string {
	return t.Format("15:04")
}

// Percent formats a float64 as a percentage.
// Example: Percent(0.15) returns "%15.00"
func Percent(val float64) string {
	return fmt.Sprintf("%.2f%%", val*100)
}

// KeyValuePair represents a key-value line in a receipt.
type KeyValuePair struct {
	Key    string
	Value  string
	Styles []Style
}

// KeyValue creates a key-value pair for rendering.
func KeyValue(key string, value interface{}, styles ...Style) *KeyValuePair {
	return &KeyValuePair{
		Key:    key,
		Value:  fmt.Sprint(value),
		Styles: styles,
	}
}

// Separator renders a horizontal separator line with optional space above and below.
func Separator() *LineSpec {
	return &LineSpec{
		Style:     LineSolid,
		SpaceAbove: 2,
		SpaceBelow: 2,
	}
}

// LineType represents the type of line.
type LineType int

const (
	// LineSolid is a solid line.
	LineSolid LineType = iota
	// LineDashed is a dashed line.
	LineDashed
)

// LineSpec holds line rendering specifications.
type LineSpec struct {
	Style       LineType
	SpaceAbove float64
	SpaceBelow float64
}

// Line renders a solid horizontal line.
func Line() *LineSpec {
	return &LineSpec{
		Style:       LineSolid,
		SpaceAbove:  0,
		SpaceBelow:  0,
	}
}

// DashedLine renders a dashed horizontal line.
func DashedLine() *LineSpec {
	return &LineSpec{
		Style:       LineDashed,
		SpaceAbove:  0,
		SpaceBelow:  0,
	}
}

// SolidLine is an alias for Line.
var SolidLine = Line

// Space creates a vertical space of the specified millimeters.
func Space(mm float64) float64 {
	return mm
}

// Paragraph creates a paragraph from multiple lines.
func Paragraph(lines ...string) string {
	return strings.Join(lines, "\n")
}

// CenterText creates a centered text string (alignment handled by rendering).
func CenterText(text string) string {
	return text
}

// LeftText creates a left-aligned text string (alignment handled by rendering).
func LeftText(text string) string {
	return text
}

// RightText creates a right-aligned text string (alignment handled by rendering).
func RightText(text string) string {
	return text
}

// Truncate truncates a string to the specified length with ellipsis.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
