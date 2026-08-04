// Package receipt, 80mm termal yazıcılar için basit, siyah-beyaz,
// tablo tabanlı PDF fiş üretimi sağlar. Tüm gofpdf detayları bu paket
// içinde gizlenir; dışarıya sadece Receipt tipi ve metodları açılır.
//
// Akış:
//
//	r := receipt.New(80) // 80mm genişlik
//	r.AddRow(
//	    receipt.NewCell("URUN", receipt.TypeTitle, receipt.AlignLeft),
//	    receipt.NewCell("AD", receipt.TypeTitle, receipt.AlignCenter),
//	    receipt.NewCell("FIYAT", receipt.TypeTitle, receipt.AlignRight),
//	)
//	r.AddLine()
//	r.AddRow(
//	    receipt.NewCell("Latte", receipt.TypeNormal, receipt.AlignLeft),
//	    receipt.NewCell("2", receipt.TypeNormal, receipt.AlignCenter),
//	    receipt.NewCell("270.00", receipt.TypeNormal, receipt.AlignRight),
//	)
//	r.Save("fis.pdf")
//
// Sayfa yüksekliği dinamiktir: içerik önce hafızada biriktirilir,
// Save() çağrıldığında toplam yükseklik hesaplanıp PDF sayfası o
// yüksekliğe göre oluşturulur. Bu yüzden çok kısa ya da çok uzun
// fişler için ayrı ayrı sayfa boyutu ayarlamaya gerek kalmaz.
package receipt

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

// ---------------------------------------------------------------------
// Tipler: hücre içeriği (başlık / normal / küçük) ve hizalama
// ---------------------------------------------------------------------

// CellType, bir hücrenin metin boyutunu/rolünü belirler.
type CellType int

const (
	// TypeTitle: büyük, kalın başlık metni (örn. "SWIFTY CAFE").
	TypeTitle CellType = iota
	// TypeNormal: standart gövde metni (ürün adı, fiyat vb.).
	TypeNormal
	// TypeSmall: dipnot / ek bilgi gibi küçük metin (adres, vergi no vb.).
	TypeSmall
)

// Align, bir hücre içindeki metnin yatay hizalamasını belirler.
type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// fontSize, her CellType için kullanılacak varsayılan punto değerini döner.
func (t CellType) fontSize() float64 {
	switch t {
	case TypeTitle:
		return 12
	case TypeSmall:
		return 7
	default: // TypeNormal
		return 9
	}
}

// gofpdfAlignStr, Align değerini gofpdf'in beklediği hizalama koduna çevirir.
func (a Align) gofpdfAlignStr() string {
	switch a {
	case AlignCenter:
		return "C"
	case AlignRight:
		return "R"
	default:
		return "L"
	}
}

// ---------------------------------------------------------------------
// Cell: bir satırdaki tek bir hücre
// ---------------------------------------------------------------------

// Cell, bir satır içindeki tek bir hücreyi temsil eder.
type Cell struct {
	Text   string
	Type   CellType
	Align  Align
	Bold   bool
	Italic bool

	// Width, hücrenin mm cinsinden genişliğidir. 0 verilirse, satırdaki
	// genişliği belirtilmemiş diğer hücrelerle birlikte kalan alanı eşit
	// paylaşır (otomatik genişlik).
	Width float64
}

// NewCell, verilen metin/tip/hizalama ile sade bir hücre oluşturur.
// Kalın/italik gibi ek stiller için WithBold / WithItalic / WithWidth
// zincirleme metodlarını kullanabilirsiniz.
func NewCell(text string, cellType CellType, align Align) *Cell {
	return &Cell{
		Text:  text,
		Type:  cellType,
		Align: align,
	}
}

// WithBold, hücreyi kalın yazı tipiyle işaretler ve zincirlemeye izin verir.
func (c *Cell) WithBold() *Cell {
	c.Bold = true
	return c
}

// WithItalic, hücreyi italik yazı tipiyle işaretler ve zincirlemeye izin verir.
func (c *Cell) WithItalic() *Cell {
	c.Italic = true
	return c
}

// WithWidth, hücreye sabit bir mm genişliği atar (0 = otomatik).
func (c *Cell) WithWidth(mm float64) *Cell {
	c.Width = mm
	return c
}

// fontStyle, gofpdf'in SetFont metoduna verilecek stil string'ini üretir.
func (c *Cell) fontStyle() string {
	style := ""
	if c.Bold {
		style += "B"
	}
	if c.Italic {
		style += "I"
	}
	return style
}

// ---------------------------------------------------------------------
// İç element modeli: Save() çağrılana kadar biriktirilen satırlar/çizgiler
// ---------------------------------------------------------------------

// elementKind, biriktirilen bir elementin türünü belirtir.
type elementKind int

const (
	kindRow elementKind = iota
	kindLine
	kindSpace
)

// element, PDF'e çizilecek tek bir satır/çizgi/boşluk birimidir.
type element struct {
	kind  elementKind
	cells []*Cell // kindRow için doludur
	spaceMM float64 // kindSpace için doludur
}

// ---------------------------------------------------------------------
// Receipt: dışarıya açılan ana sınıf
// ---------------------------------------------------------------------

// Receipt, biriktirilmiş fiş içeriğini tutar ve Save() çağrıldığında
// gerçek PDF'i üretir.
type Receipt struct {
	widthMM     float64 // kağıt genişliği (örn. 80)
	marginMM    float64 // sol/sağ/üst/alt kenar boşluğu
	rowHeightMM float64 // her satırın (tek satırlık, kaydırmasız) yüksekliği

	elements []element
}

// New, verilen kağıt genişliğinde (mm) yeni ve boş bir Receipt oluşturur.
// Termal yazıcı senaryosu için genelde 80 kullanılır.
func New(widthMM float64) *Receipt {
	return &Receipt{
		widthMM:     widthMM,
		marginMM:    4,
		rowHeightMM: 5,
		elements:    make([]element, 0, 32),
	}
}

// contentWidthMM, kenar boşlukları düşüldükten sonra kalan yazılabilir genişliktir.
func (r *Receipt) contentWidthMM() float64 {
	return r.widthMM - 2*r.marginMM
}

// AddRow, verilen hücrelerle yeni bir satır ekler. Satırdaki hücre sayısı
// tamamen bu çağrıda verilen cells parametresinin uzunluğuna göre belirlenir;
// yani "kaç hücre olacağı" çağıran taraf tarafından dinamik olarak seçilir.
func (r *Receipt) AddRow(cells ...*Cell) *Receipt {
	if len(cells) == 0 {
		// Hücresiz satır anlamsızdır; sessizce yok say.
		return r
	}
	r.elements = append(r.elements, element{kind: kindRow, cells: cells})
	return r
}

// AddLine, o ana kadarki en son satırın altına yatay bir ayraç çizgisi ekler.
func (r *Receipt) AddLine() *Receipt {
	r.elements = append(r.elements, element{kind: kindLine})
	return r
}

// AddSpace, mm cinsinden dikey boşluk ekler (satırlar arasını açmak için).
func (r *Receipt) AddSpace(mm float64) *Receipt {
	r.elements = append(r.elements, element{kind: kindSpace, spaceMM: mm})
	return r
}

// ---------------------------------------------------------------------
// Yükseklik hesaplama (dinamik sayfa boyu için ön geçiş)
// ---------------------------------------------------------------------

// calculateTotalHeightMM, mevcut tüm elementlerin kaplayacağı toplam
// dikey alanı (mm) hesaplar. PDF sayfası bu değere göre oluşturulur,
// böylece sayfa uzunluğu içerik miktarına göre otomatik ayarlanmış olur.
func (r *Receipt) calculateTotalHeightMM() float64 {
	total := r.marginMM * 2 // üst + alt boşluk
	for _, el := range r.elements {
		switch el.kind {
		case kindRow:
			total += r.rowHeightMM
		case kindLine:
			total += r.rowHeightMM * 0.6
		case kindSpace:
			total += el.spaceMM
		}
	}
	return total
}

// ---------------------------------------------------------------------
// Render / Save
// ---------------------------------------------------------------------

// Save, biriktirilmiş tüm içeriği gerçek bir PDF dosyasına çizer ve diske yazar.
func (r *Receipt) Save(filename string) error {
	heightMM := r.calculateTotalHeightMM()
	// Çok kısa fişlerde bile makul bir minimum sayfa uzunluğu bırakalım.
	if heightMM < 40 {
		heightMM = 40
	}

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size: gofpdf.SizeType{
			Wd: r.widthMM,
			Ht: heightMM,
		},
	})
	pdf.SetMargins(r.marginMM, r.marginMM, r.marginMM)
	pdf.SetAutoPageBreak(false, 0) // fişin tamamı tek, dinamik boyutlu sayfaya sığar
	pdf.AddUTF8Font("DejaVuSans", "", "assets/fonts/DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVuSans", "B", "assets/fonts/DejaVuSans-Bold.ttf")
	pdf.AddUTF8Font("DejaVuSans", "I", "assets/fonts/DejaVuSans-Oblique.ttf")
	pdf.AddUTF8Font("DejaVuSans", "BI", "assets/fonts/DejaVuSans-BoldOblique.ttf")
	pdf.AddPage()

	for _, el := range r.elements {
		switch el.kind {
		case kindRow:
			r.drawRow(pdf, el.cells)
		case kindLine:
			r.drawLine(pdf)
		case kindSpace:
			pdf.SetY(pdf.GetY() + el.spaceMM)
		}
	}

	return pdf.OutputFileAndClose(filename)
}

// drawRow, tek bir satırı (bir veya daha fazla hücreyi) çizer ve
// imleci (cursor) bir sonraki satıra ilerletir.
func (r *Receipt) drawRow(pdf *gofpdf.Fpdf, cells []*Cell) {
	contentWidth := r.contentWidthMM()

	// Sabit genişliği belirtilmiş hücrelerin toplamını çıkar,
	// kalanı otomatik genişlikli hücreler arasında eşit paylaştır.
	fixedWidth := 0.0
	autoCount := 0
	for _, c := range cells {
		if c.Width > 0 {
			fixedWidth += c.Width
		} else {
			autoCount++
		}
	}
	remaining := contentWidth - fixedWidth
	if remaining < 0 {
		remaining = 0
	}
	autoWidth := 0.0
	if autoCount > 0 {
		autoWidth = remaining / float64(autoCount)
	}

	y := pdf.GetY()
	x := r.marginMM
	pdf.SetXY(x, y)

	for _, c := range cells {
		w := c.Width
		if w <= 0 {
			w = autoWidth
		}

		style := c.fontStyle()
		if c.Type == TypeTitle && !c.Bold {
			style += "B"
		}
		pdf.SetFont("DejaVuSans", style, c.Type.fontSize())

		pdf.CellFormat(w, r.rowHeightMM, c.Text, "", 0, c.Align.gofpdfAlignStr(), false, 0, "")
	}

	// İmleci bir sonraki satıra indir (satır başına dön + aşağı in).
	pdf.SetXY(r.marginMM, y+r.rowHeightMM)
}

// drawLine, geçerli Y konumuna kesikli-olmayan düz bir ayraç çizgisi çizer
// ve imleci çizginin altına ilerletir.
func (r *Receipt) drawLine(pdf *gofpdf.Fpdf) {
	y := pdf.GetY()
	x1 := r.marginMM
	x2 := r.widthMM - r.marginMM

	pdf.SetLineWidth(0.2)
	pdf.Line(x1, y, x2, y)

	pdf.SetXY(r.marginMM, y+r.rowHeightMM*0.6)
}

// ---------------------------------------------------------------------
// Küçük yardımcı: para/etiket gibi tek satırlık kısayollar isteyen
// kullanıcılar için opsiyonel bir yardımcı fonksiyon.
// ---------------------------------------------------------------------

// Money, bir sayıyı "125.00" formatında (para birimi eki olmadan) döner.
// Kullanıcı isterse kendi para birimi ekini (örn. " TL") string'e ekleyebilir.
func Money(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
