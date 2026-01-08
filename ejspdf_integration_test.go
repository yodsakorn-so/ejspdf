package ejspdf_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yodsakorn-so/ejspdf"
)

func TestRender_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	t.Run("Basic Render", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

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
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

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

	t.Run("RenderFromFile with Options", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		tmpFile := "test_temp.ejs"
		err := os.WriteFile(tmpFile, []byte("<h1>File Test <%= name %></h1>"), 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile)

		pdfBytes, err := ejspdf.RenderFromFile(ctx, tmpFile, ejspdf.Options{
			Data:             map[string]any{"name": "EJS"},
			Scale:            0.5,
			IgnoreBackground: true,
		})
		if err != nil {
			t.Fatalf("Failed to render from file: %v", err)
		}
		if len(pdfBytes) == 0 {
			t.Error("PDF output is empty")
		}
	})
}
