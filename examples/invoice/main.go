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
	templateBytes, err := os.ReadFile("invoice.ejs")
	if err != nil {
		log.Fatalf("read template failed: %v", err)
	}

	// 2. Setup Context with longer Timeout (60s)
	// Chrome startup or external resources (fonts/CDN) might take some time.
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 3. Render PDF
	log.Println("Generating PDF (this may take a while if downloading fonts)...")
	pdfBytes, err := ejspdf.Render(ctx, ejspdf.Options{
		Template: string(templateBytes),
		Data: map[string]any{
			"customer": "John Smith",
			"total":    5623,
		},
		PageSize:  "A4",
		Landscape: false,
		// Leaving WaitSelector empty uses the default (WaitReady("body")),
		// which is more stable for general web pages.
		WaitSelector: "",
		// Add a small delay to ensure fonts and styles are fully rendered.
		WaitDelay: 1 * time.Second,
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}

	// 4. Save to file
	if err := os.WriteFile("invoice.pdf", pdfBytes, 0644); err != nil {
		log.Fatalf("Save file failed: %v", err)
	}

	log.Println("PDF generated successfully: invoice.pdf 🎉")
}
