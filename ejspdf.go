package ejspdf

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/yodsakorn-so/ejspdf/internal/pdf"
	"github.com/yodsakorn-so/ejspdf/internal/renderer"
	"github.com/yodsakorn-so/ejspdf/internal/renderer/assets"
)

// Options defines configuration for rendering EJS to PDF.
type Options struct {
	// Template is the EJS template string.
	Template string
	// Data is the data map to pass to the template.
	Data any

	// ChromePath is the custom path to Chrome/Chromium executable.
	// If empty, it will try to find Chrome automatically.
	ChromePath string

	// PageSize sets the paper size (e.g., "A4", "A3", "Letter", "Legal").
	// Default is "A4".
	PageSize string
	// Landscape sets the paper orientation. Default is false (Portrait).
	Landscape bool

	// PaperWidth overrides PageSize width if set.
	// Supports units: "mm", "cm", "in" (e.g., "80mm", "4in").
	PaperWidth string
	// PaperHeight overrides PageSize height if set.
	// Supports units: "mm", "cm", "in" (e.g., "200mm", "11in").
	PaperHeight string

	// MarginTop sets the top margin (e.g., "10mm", "0.5in"). Default is "10mm".
	MarginTop string
	// MarginBottom sets the bottom margin. Default is "10mm".
	MarginBottom string
	// MarginLeft sets the left margin. Default is "10mm".
	MarginLeft string
	// MarginRight sets the right margin. Default is "10mm".
	MarginRight string

	// DisplayHeaderFooter enables printing of headers and footers.
	DisplayHeaderFooter bool
	// HeaderTemplate is the HTML template for the print header.
	// Use classes "pageNumber", "totalPages", "date", "title", "url" to inject values.
	HeaderTemplate string
	// FooterTemplate is the HTML template for the print footer.
	FooterTemplate string

	// WaitSelector is the CSS selector to wait for before printing (e.g., "#main-content").
	// If empty, it waits for the "body" tag by default.
	WaitSelector string

	// WaitDelay is an optional additional delay after the selector is visible.
	// Useful for ensuring fonts, images, or animations are fully loaded.
	WaitDelay time.Duration
}

// Render generates a PDF from an EJS template using the provided options and context.
// If the context already contains a chromedp session, it will be reused.
func Render(ctx context.Context, opt Options) ([]byte, error) {
	if opt.Template == "" {
		return nil, fmt.Errorf("ejspdf: template is required")
	}

	// 1. Render EJS -> HTML
	rt := renderer.New()

	html, err := rt.RenderEJS(assets.EJS, opt.Template, opt.Data)
	if err != nil {
		return nil, fmt.Errorf("ejspdf: render ejs failed: %w", err)
	}

	// 2. HTML -> PDF
	chrome := pdf.New(pdf.Options{
		ChromePath:          opt.ChromePath,
		PageSize:            defaultString(opt.PageSize, "A4"),
		Landscape:           opt.Landscape,
		PaperWidth:          opt.PaperWidth,
		PaperHeight:         opt.PaperHeight,
		MarginTop:           defaultString(opt.MarginTop, "10mm"),
		MarginBottom:        defaultString(opt.MarginBottom, "10mm"),
		MarginLeft:          defaultString(opt.MarginLeft, "10mm"),
		MarginRight:         defaultString(opt.MarginRight, "10mm"),
		DisplayHeaderFooter: opt.DisplayHeaderFooter,
		HeaderTemplate:      opt.HeaderTemplate,
		FooterTemplate:      opt.FooterTemplate,
		WaitSelector:        opt.WaitSelector,
		WaitDelay:           opt.WaitDelay,
	})

	pdfBytes, err := chrome.FromHTML(ctx, html)
	if err != nil {
		return nil, fmt.Errorf("ejspdf: render pdf failed: %w", err)
	}

	return pdfBytes, nil
}

// ImageFileToBase64 reads a local image file and returns a Data URI string
// suitable for use in HTML/EJS (e.g., "data:image/png;base64,...").
func ImageFileToBase64(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("ejspdf: failed to read image file: %w", err)
	}

	mimeType := http.DetectContentType(data)
	if mimeType == "application/octet-stream" {
		ext := filepath.Ext(path)
		switch ext {
		case ".svg":
			mimeType = "image/svg+xml"
		}
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

func defaultString(v, d string) string {
	if v == "" {
		return d
	}
	return v
}
