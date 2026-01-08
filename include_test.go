package ejspdf_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yodsakorn-so/ejspdf"
)

// TestIncludeFeature tests the EJS include functionality.
func TestIncludeFeature(t *testing.T) {
	// 1. Create partial and main template files
	tmpDir := t.TempDir()
	
	headerPath := filepath.Join(tmpDir, "header.ejs")
	mainPath := filepath.Join(tmpDir, "main.ejs")

	err := os.WriteFile(headerPath, []byte("<h1>THIS IS HEADER</h1>"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Use include with relative path
	err = os.WriteFile(mainPath, []byte(`<%- include('header.ejs') %> <p>Body Content</p>`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. Render
	pdfBytes, err := ejspdf.RenderFromFile(ctx, mainPath, ejspdf.Options{
		Data: nil,
	})

	// 3. Expect success
	if err != nil {
		t.Fatalf("Include feature failed: %v", err)
	}
	
	if len(pdfBytes) == 0 {
		t.Error("PDF output is empty")
	}
}

// TestFunctionInjection tests passing Go functions into EJS templates.
func TestFunctionInjection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	upperFunc := func(s string) string {
		return "GO-" + s
	}

	pdfBytes, err := ejspdf.Render(ctx, ejspdf.Options{
		Template: `<p>Result: <%= toUpper('test') %></p>`,
		Data: map[string]any{
			"toUpper": upperFunc,
		},
	})

	if err != nil {
		t.Fatalf("Function injection failed: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Error("PDF empty")
	}
}
