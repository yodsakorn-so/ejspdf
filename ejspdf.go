// Package ejspdf provides EJS to PDF rendering using Chrome.
package ejspdf

import (
	"context"

	"github.com/yodsakorn-so/ejspdf/pdf"
	"github.com/yodsakorn-so/ejspdf/renderer"
)

type Engine struct {
	rt *renderer.Runtime
}

// New creates a new EJS PDF rendering engine.
func New() (*Engine, error) {
	rt, err := renderer.New()
	if err != nil {
		return nil, err
	}
	return &Engine{rt: rt}, nil
}

// RenderEJSToPDF renders an EJS template into a PDF document.
func (e *Engine) RenderEJSToPDF(
	ctx context.Context,
	tpl string,
	data any,
) ([]byte, error) {
	html, err := e.rt.Render(tpl, data)
	if err != nil {
		return nil, err
	}

	return pdf.HTMLToPDF(ctx, html)
}
