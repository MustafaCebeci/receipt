# Receipt PDF Generator

80mm termal yazıcılar için, tek bir Go dosyasında實現された fluent API ile PDF fiş üretimi sağlayan paket.

## Özellikler

- **UTF-8 desteği** — Türkçe karakterler (ç, ğ, ı, ö, ş, ü) sorunsuz
- **Dinamik sayfa yüksekliği** — içerik otomatik hesaplanır, kısa/uzun fişler için ayrı boyut gerekmez
- **Basit API** — sadece 4 temel metod: `AddRow`, `AddLine`, `AddSpace`, `Save`
- **Düşük bağımlılık** — sadece `gofpdf`

## Kurulum

```bash
go get github.com/MustafaCebeci/receipt
```

## Hızlı Başlangıç

```go
r := receipt.New(receipt.DefaultWidth) // 72.1mm (576 dots @ 203 DPI printable area)

r.AddRow(receipt.NewCell("SWIFTY CAFE", receipt.TypeTitle, receipt.AlignCenter))
r.AddRow(receipt.NewCell("Adana / Turkiye", receipt.TypeSmall, receipt.AlignCenter))
r.AddLine()

r.AddRow(
    receipt.NewCell("Fis No", receipt.TypeNormal, receipt.AlignLeft),
    receipt.NewCell("000154", receipt.TypeNormal, receipt.AlignRight),
)

r.AddLine()

r.AddRow(
    receipt.NewCell("URUN", receipt.TypeNormal, receipt.AlignLeft).WithBold().WithWidth(40),
    receipt.NewCell("ADET", receipt.TypeNormal, receipt.AlignCenter).WithBold().WithWidth(15),
    receipt.NewCell("FIYAT", receipt.TypeNormal, receipt.AlignRight).WithBold(),
)
r.AddLine()

r.AddRow(
    receipt.NewCell("Latte", receipt.TypeNormal, receipt.AlignLeft).WithWidth(40),
    receipt.NewCell("2", receipt.TypeNormal, receipt.AlignCenter).WithWidth(15),
    receipt.NewCell(receipt.Money(270)+" TL", receipt.TypeNormal, receipt.AlignRight),
)

r.AddLine()

r.AddRow(
    receipt.NewCell("TOPLAM", receipt.TypeNormal, receipt.AlignLeft).WithBold(),
    receipt.NewCell(receipt.Money(430)+" TL", receipt.TypeNormal, receipt.AlignRight).WithBold(),
)

r.Save("fis.pdf")
```

## API

### Tipler

| Tip | Açıklama |
|-----|----------|
| `CellType` | Hücre metin boyutu: `TypeTitle` (12pt), `TypeNormal` (9pt), `TypeSmall` (7pt) |
| `Align` | Yatay hizalama: `AlignLeft`, `AlignCenter`, `AlignRight` |

### Receipt Metodları

| Metod | Açıklama |
|-------|----------|
| `New(widthMM)` | Verilen mm genişliğinde yeni fiş oluşturur. Çoğu 80mm yazıcı için `DefaultWidth` (72.1mm) önerilir |
| `SetMargin(mm)` | Kenar boşluğunu ayarlar (varsayılan: 2mm) |
| `AddRow(cells...)` | Hücrelerle yeni bir satır ekler |
| `AddLine()` | Yatay ayraç çizgisi ekler |
| `AddSpace(mm)` | mm cinsinden dikey boşluk ekler |
| `Save(filename)` | PDF'i diske yazar |

### Cell Metodları

| Metod | Açıklama |
|-------|----------|
| `WithBold()` | Kalın yazı tipi |
| `WithItalic()` | İtalik yazı tipi |
| `WithWidth(mm)` | Hücre genişliği (0 = otomatik) |

### Yardımcı Fonksiyonlar

| Fonksiyon | Açıklama |
|-----------|----------|
| `Money(amount)` | `125.00` formatında döner |
| `NewCell(text, type, align)` | Yeni hücre oluşturur |

## Kağıt Boyutu

80mm termal yazıcılar için tasarlanmıştır. Yükseklik dinamik olarak içerik miktarına göre hesaplanır.

### Sabitler

| Sabit | Değer | Açıklama |
|-------|-------|----------|
| `DefaultWidth` | 72.1mm | Sunlux RP8020 ve benzeri yazıcılar için printable area (576 dots @ 203 DPI) |
| `DefaultMargin` | 2.0mm | Her kenardan varsayılan boşluk (toplam 4mm) |

**Önemli:** 80mm yazıcıların kağıt genişliği 80mm olsa da, gerçek baskı alanı 72.1mm'dir. `DefaultWidth` (72.1mm) kullanmanız önerilir. PDF her zaman portrait orientation'da oluşturulur.

## Font Bağımlılığı

Paket, DejaVu Sans TrueType fontlarını kullanır ve font dosyaları `assets/fonts/` klasöründe bulunur. Paketi kullanan projelerde bu klasörün korunması gerekir.

```bash
# Projeyi klonladığınızda fontlar otomatik gelir
git clone https://github.com/MustafaCebeci/receipt.git
cd receipt
```

Fontlar `go get` ile çekildiğinde otomatik olarak gelir. Ancak projeyi başka bir konumdan çalıştırıyorsanız, `assets/fonts/` klasörünün çalışma dizininde bulunduğundan emin olun.

## Lisans

MIT License
