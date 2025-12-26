# ejspdf

[![Go Reference](https://pkg.go.dev/badge/github.com/yodsakorn-so/ejspdf.svg)](https://pkg.go.dev/github.com/yodsakorn-so/ejspdf)
[![Go Report Card](https://goreportcard.com/badge/github.com/yodsakorn-so/ejspdf)](https://goreportcard.com/report/github.com/yodsakorn-so/ejspdf)
[![Build Status](https://github.com/yodsakorn-so/ejspdf/actions/workflows/go.yml/badge.svg)](https://github.com/yodsakorn-so/ejspdf/actions)
[![License](https://img.shields.io/github/license/yodsakorn-so/ejspdf)](LICENSE)

Render **EJS templates to PDF** using **Go + Chrome (chromedp)**.

`ejspdf` is a Go library that allows you to render [EJS](https://ejs.co/) templates into PDF documents using a real headless Chrome browser. This ensures perfect rendering of modern HTML, CSS, and Flexbox layouts.

Ideal for generating:
- 📄 Invoices & Receipts
- 🏥 Medical Certificates
- 🇹🇭 Thai Government Documents (Full Thai Font Support)
- 📊 Reports with Charts/Graphs

---

## ✨ Features

- ✅ **EJS Template Support**: Use standard EJS syntax (`<%= var %>`, `<% if %>`, etc.).
- ✅ **Chrome-based Rendering**: What you see in the browser is exactly what you get in PDF.
- ✅ **Base64 Data URI**: Robust handling of special characters (e.g., `#` in CSS colors) and Thai language.
- ✅ **Custom Page Sizes**: Support for standard (A4, A3) and custom dimensions (e.g., "80mm").
- ✅ **Header & Footer**: Inject dynamic headers/footers with page numbers.
- ✅ **Context Aware**: Full control over timeouts and cancellation via `context.Context`.

---

## 📦 Installation

```bash
go get github.com/yodsakorn-so/ejspdf
```

**Prerequisites:**
- **Go** 1.22+
- **Google Chrome** or **Chromium** installed on the host machine.

---

## 🚀 Quick Start

```go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/yodsakorn-so/ejspdf"
)

func main() {
	// 1. Prepare EJS Template
	template := `
    <html>
      <body>
        <h1>Invoice for <%= name %></h1>
        <p>Total: <b><%= total %> THB</b></p>
      </body>
    </html>
    `

	// 2. Setup Context with Timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 3. Render PDF
	pdfBytes, err := ejspdf.Render(ctx, ejspdf.Options{
		Template: template,
		Data: map[string]any{
			"name":  "John Doe",
			"total": "5,000",
		},
		PageSize:  "A4",
		Landscape: false,
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}

	// 4. Save to file
	if err := os.WriteFile("output.pdf", pdfBytes, 0644); err != nil {
		log.Fatal(err)
	}
	
	log.Println("PDF generated successfully! 🎉")
}
```

---

## ⚙️ Options

| Option | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `Template` | `string` | **Required** | The EJS template string. |
| `Data` | `any` | `nil` | Data object/map passed to the template. |
| `ChromePath` | `string` | Auto | Custom path to Chrome executable. |
| `PageSize` | `string` | `"A4"` | `A4`, `A3`, `A5`, `Letter`, `Legal`, `Tabloid`. |
| `Landscape` | `bool` | `false` | Set to `true` for landscape orientation. |
| `PaperWidth` | `string` | `""` | Custom width (e.g., "80mm", "4in"). Overrides `PageSize`. |
| `PaperHeight` | `string` | `""` | Custom height (e.g., "200mm"). Overrides `PageSize`. |
| `MarginTop` | `string` | `"10mm"` | Top margin (e.g., "0.4in", "1cm"). |
| `Margin...` | `string` | `"10mm"` | Bottom, Left, Right margins. |
| `WaitSelector` | `string` | `""` | CSS selector to wait for before printing (e.g., `"#app"`). |
| `WaitDelay` | `time.Duration` | `0` | Additional delay time (e.g., `500 * time.Millisecond`). |
| `DisplayHeaderFooter` | `bool` | `false` | Enable header and footer printing. |
| `HeaderTemplate` | `string` | `""` | HTML template for header. |
| `FooterTemplate` | `string` | `""` | HTML template for footer. |

---

## 💡 Tips

### Landscape Mode
When using `Landscape: true`, ensure your CSS does not enforce fixed portrait dimensions (e.g., `width: 210mm; height: 297mm`). Let the content flow naturally or use percentages.

### Header & Footer Templates
Chromium supports specific classes to inject values:
- `date`: Formatted print date
- `title`: Document title
- `url`: Document location
- `pageNumber`: Current page number
- `totalPages`: Total pages

**Example:**
```go
HeaderTemplate: `<div style="font-size:10px; width:100%; text-align:center;">
    Page <span class="pageNumber"></span> of <span class="totalPages"></span>
</div>`,
```

Note: Margins must be large enough to accommodate the header/footer, or they might be clipped.

---

## 🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

[MIT](LICENSE)
