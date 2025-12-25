package ejspdf

import (
	"context"

	"github.com/yodsakorn-so/ejspdf/pdf"
	"github.com/yodsakorn-so/ejspdf/renderer"
)

type Engine struct {
	rt *renderer.Runtime
}

func New() (*Engine, error) {
	rt, err := renderer.New()
	if err != nil {
		return nil, err
	}
	return &Engine{rt: rt}, nil
}

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
