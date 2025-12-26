package ejspdf_test

import (
	"context"
	"testing"
	"time"

	"github.com/yodsakorn-so/ejspdf"
)

func TestRender_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	t.Run("Basic Render", func(t *testing.T) {
		pdfBytes, err := ejspdf.Render(ctx, ejspdf.Options{
			Template: "<h1>Integration Test</h1><p><%= value %></p>",
			Data: map[string]any{
				"value": "Hello World",
			},
		})

		if err != nil {
			t.Fatalf("Failed to render: %v", err)
		}

		if len(pdfBytes) < 1000 {
			t.Errorf("PDF output too small, got %d bytes", len(pdfBytes))
		}
	})

	t.Run("Landscape and Custom Size", func(t *testing.T) {
		pdfBytes, err := ejspdf.Render(ctx, ejspdf.Options{
			Template:    "<h1>Landscape Test</h1>",
			Landscape:   true,
			PaperWidth:  "100mm",
			PaperHeight: "100mm",
		})

		if err != nil {
			t.Fatalf("Failed to render: %v", err)
		}

		if len(pdfBytes) == 0 {
			t.Error("PDF output is empty")
		}
	})
}
