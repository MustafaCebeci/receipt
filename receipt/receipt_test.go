package receipt

import (
	"os"
	"testing"
	"time"
)

// TestMoney tests the Money formatting function.
func TestMoney(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "₺0.00"},
		{1, "₺1.00"},
		{10.5, "₺10.50"},
		{1234.56, "₺1234.56"},
		{-50, "₺-50.00"},
	}

	for _, tt := range tests {
		result := Money(tt.input)
		if result != tt.expected {
			t.Errorf("Money(%v) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

// TestDate tests the Date formatting function.
func TestDate(t *testing.T) {
	timestamp := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	expected := "15.03.2024"
	result := Date(timestamp)
	if result != expected {
		t.Errorf("Date() = %v, want %v", result, expected)
	}
}

// TestDateTime tests the DateTime formatting function.
func TestDateTime(t *testing.T) {
	timestamp := time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC)
	expected := "15.03.2024 14:30"
	result := DateTime(timestamp)
	if result != expected {
		t.Errorf("DateTime() = %v, want %v", result, expected)
	}
}

// TestPercent tests the Percent formatting function.
func TestPercent(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "0.00%"},
		{0.15, "15.00%"},
		{1, "100.00%"},
		{0.1234, "12.34%"},
	}

	for _, tt := range tests {
		result := Percent(tt.input)
		if result != tt.expected {
			t.Errorf("Percent(%v) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

// TestCell tests cell creation.
func TestCell(t *testing.T) {
	c := Cell("Test", Bold(), Size(12))
	if c == nil {
		t.Fatal("Cell() returned nil")
	}
	if c.String() != "Test" {
		t.Errorf("Cell().String() = %v, want 'Test'", c.String())
	}
	if c.colSpan != 1 {
		t.Errorf("Cell().colSpan = %v, want 1", c.colSpan)
	}
}

// TestCellWithSpan tests cell creation with colspan.
func TestCellWithSpan(t *testing.T) {
	c := Cell("Test", ColSpan(3))
	if c.colSpan != 3 {
		t.Errorf("Cell().colSpan = %v, want 3", c.colSpan)
	}
}

// TestColumns tests the Columns helper.
func TestColumns(t *testing.T) {
	cols := Columns(50, 20, 30)
	if len(cols) != 3 {
		t.Fatalf("len(Columns(50, 20, 30)) = %d, want 3", len(cols))
	}
	if cols[0].Width != 50 || cols[1].Width != 20 || cols[2].Width != 30 {
		t.Errorf("Column widths incorrect")
	}
	if cols[0].IsAuto || cols[1].IsAuto || cols[2].IsAuto {
		t.Errorf("Fixed columns should not be marked as auto")
	}
}

// TestColumnsWithAuto tests the Columns helper with auto width.
func TestColumnsWithAuto(t *testing.T) {
	cols := Columns(50, 0, 30)
	if len(cols) != 3 {
		t.Fatalf("len(Columns(50, 0, 30)) = %d, want 3", len(cols))
	}
	if !cols[1].IsAuto {
		t.Errorf("Middle column should be marked as auto")
	}
}

// TestStyleFunctions tests that style functions return valid styles.
func TestStyleFunctions(t *testing.T) {
	styles := []Style{
		Bold(),
		Italic(),
		Underline(),
		Strike(),
		Size(12),
		Font("Arial"),
		Left(),
		Center(),
		Right(),
		Top(),
		Middle(),
		Bottom(),
		Padding(5),
		PaddingHorizontal(3),
		PaddingVertical(4),
		Margin(2),
		MarginHorizontal(1),
		MarginVertical(1),
		Border(),
		BorderTop(),
		BorderBottom(),
		BorderLeft(),
		BorderRight(),
		ColSpan(2),
		RowSpan(3),
		Width(50),
		Height(10),
		MinHeight(5),
		MaxHeight(20),
		Wrap(),
		NoWrap(),
		Small(),
		Normal(),
		Header(),
		Footer(),
	}

	for i, s := range styles {
		if s == nil {
			t.Errorf("Style function %d returned nil", i)
		}
	}
}

// TestStyleContextMerge tests merging style contexts.
func TestStyleContextMerge(t *testing.T) {
	ctx1 := &StyleContext{FontSize: 10, Bold: false}
	ctx2 := &StyleContext{FontSize: 12, Bold: true}

	ctx1.Merge(ctx2)

	if ctx1.FontSize != 12 {
		t.Errorf("Merge: FontSize = %v, want 12", ctx1.FontSize)
	}
	if !ctx1.Bold {
		t.Errorf("Merge: Bold = %v, want true", ctx1.Bold)
	}
}

// TestAlignmentConstants tests alignment constants.
func TestAlignmentConstants(t *testing.T) {
	if AlignLeft >= AlignCenter || AlignCenter >= AlignRight {
		t.Error("Alignment constants not in ascending order")
	}
	if AlignTop >= AlignMiddle || AlignMiddle >= AlignBottom {
		t.Error("Vertical alignment constants not in ascending order")
	}
}

// TestExampleReceipt creates a simple receipt and saves it.
func TestExampleReceipt(t *testing.T) {
	filename := "test_output_example.pdf"
	pdf := New(filename)

	pdf.Title("TEST MARKET")
	pdf.SubTitle("Alışveriş Fişi")
	pdf.Line()

	pdf.Table(Columns(50, 10, 20))
	pdf.Row(
		Cell("Ürün", Header()),
		Cell("Adet", Header()),
		Cell("Fiyat", Header()),
	)
	pdf.Row(
		Cell("Süt 1L"),
		Cell("2"),
		Cell(Money(45.00)),
	)
	pdf.Row(
		Cell("Ekmek"),
		Cell("1"),
		Cell(Money(15.00)),
	)

	pdf.Line()
	pdf.KeyValue("Ara Toplam", Money(105.00))
	pdf.KeyValue("KDV (%8)", Money(8.40))
	pdf.Line()
	pdf.KeyValue("GENEL TOPLAM", Money(113.40), Bold())

	pdf.Save()

	// Verify file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("PDF file was not created")
	} else {
		os.Remove(filename) // Clean up
	}
}

// TestRestaurantReceipt creates a restaurant order receipt.
func TestRestaurantReceipt(t *testing.T) {
	filename := "test_output_restaurant.pdf"
	pdf := New(filename)

	pdf.Title("SWIFTY CAFE")
	pdf.Space(2)
	pdf.Line()
	pdf.Space(2)

	pdf.Table(Columns(45, 8, 17))
	pdf.Row(
		Cell("Latte", Bold()),
		Cell("2", Center()),
		Cell(Money(90.00), Right()),
	)
	pdf.Row(
		Cell("Cappuccino"),
		Cell("1", Center()),
		Cell(Money(45.00), Right()),
	)
	pdf.Row(
		Cell("Türk Kahvesi"),
		Cell("1", Center()),
		Cell(Money(35.00), Right()),
	)

	pdf.Space(3)
	pdf.Line()
	pdf.Space(2)
	pdf.KeyValue("Ara Toplam", Money(170.00))
	pdf.KeyValue("Servis Ücreti", Money(17.00))
	pdf.KeyValue("KDV (%10)", Money(18.70))
	pdf.Line()
	pdf.KeyValue("TOPLAM", Money(205.70), Bold())

	pdf.Space(3)
	pdf.Separator()
	pdf.Space(2)

	pdf.QRCode("https://swiftycafe.com/order/12345")

	pdf.Save()

	// Verify file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("PDF file was not created")
	} else {
		os.Remove(filename) // Clean up
	}
}
