# Receipt PDF Generator

A fluent Go package for generating 80mm thermal POS printer receipts. Wraps `gofpdf` completely, providing a clean Compose-style API.

## Features

- UTF-8 support with Turkish characters (ç, ğ, ı, ö, ş, ü, Ç, Ğ, İ, Ö, Ş, Ü)
- Automatic text wrapping and row height calculation
- Table layout with colspan and rowspan
- Automatic page breaks for long receipts
- QR code generation
- Barcode generation (Code128, Code39)
- Key-value pair formatting
- Header and footer callbacks

## Installation

```bash
go get github.com/MustafaCebeci/receipt
```

## Quick Start

```go
package main

import (
    "github.com/MustafaCebeci/receipt"
)

func main() {
    pdf := receipt.New("receipt.pdf")

    pdf.Title("SWIFTY CAFE")
    pdf.Line()

    pdf.Table(receipt.Columns(45, 8, 17))
    pdf.Row(
        receipt.Cell("Latte", receipt.Bold()),
        receipt.Cell("2", receipt.Center()),
        receipt.Cell(receipt.Money(90.00), receipt.Right()),
    )
    pdf.Row(
        receipt.Cell("Cappuccino"),
        receipt.Cell("1", receipt.Center()),
        receipt.Cell(receipt.Money(45.00), receipt.Right()),
    )

    pdf.Line()
    pdf.KeyValue("TOPLAM", receipt.Money(135.00), receipt.Bold())

    pdf.QRCode("https://swiftycafe.com/order/12345")
    pdf.Save()
}
```

## API Reference

### Content Methods

| Method | Description |
|--------|-------------|
| `Title(text)` | Set receipt title (centered, bold) |
| `SubTitle(text)` | Set receipt subtitle |
| `Line()` | Draw solid horizontal line |
| `DashedLine()` | Draw dashed horizontal line |
| `Separator()` | Draw dashed line with space above/below |
| `Space(mm)` | Add vertical space |
| `Table(columns)` | Start a new table |
| `Row(cells...)` | Add a row to the current table |
| `KeyValue(key, value)` | Add a key-value line |
| `QRCode(content)` | Add a QR code |
| `Barcode(content, type)` | Add a barcode |
| `Footer(fn)` | Set footer callback |
| `Header(fn)` | Set header callback |
| `Save()` | Save the PDF to file |

### Style Functions

| Function | Description |
|----------|-------------|
| `Bold()` | Bold text |
| `Italic()` | Italic text |
| `Underline()` | Underlined text |
| `Size(n)` | Font size |
| `Left()` | Left align |
| `Center()` | Center align |
| `Right()` | Right align |
| `Top()` | Top align |
| `Middle()` | Middle align |
| `Bottom()` | Bottom align |
| `Padding(n)` | Cell padding |
| `Border()` | Enable all borders |
| `ColSpan(n)` | Column span |
| `RowSpan(n)` | Row span |
| `Header()` | Header style (bold, 10pt) |
| `Footer()` | Footer style (8pt) |

### Helper Functions

| Function | Description |
|----------|-------------|
| `Money(amount)` | Format as Turkish Lira (₺) |
| `Date(time)` | Format date (DD.MM.YYYY) |
| `DateTime(time)` | Format date and time |
| `Percent(val)` | Format as percentage |
| `Columns(widths...)` | Create table columns |

## Paper Size

This package is designed specifically for **80mm thermal receipt printers**. The paper width is fixed at 80mm.

## License

MIT License
