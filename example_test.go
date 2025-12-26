package ejspdf_test

import (
	"context"
	"os"
	"time"

	"github.com/yodsakorn-so/ejspdf"
)

func Example() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	engine, _ := ejspdf.New()

	pdfBytes, _ := engine.RenderEJSToPDF(
		ctx,
		"<h1>Hello <%= name %></h1>",
		map[string]any{
			"name": "World",
		},
	)

	os.WriteFile("out.pdf", pdfBytes, 0644)
}
