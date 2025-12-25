package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/yodsakorn-so/ejspdf"
)

func main() {
	// Chrome allocator
	allocCtx, allocCancel := chromedp.NewExecAllocator(
		context.Background(),
		chromedp.DefaultExecAllocatorOptions[:]...,
	)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// Initialize ejspdf Engine
	engine, err := ejspdf.New()
	if err != nil {
		log.Fatal(err)
	}

	// Load EJS template
	tplBytes, err := os.ReadFile("invoice.ejs")
	if err != nil {
		log.Fatal(err)
	}

	// Render EJS template to PDF
	pdfBytes, err := engine.RenderEJSToPDF(ctx, string(tplBytes), map[string]any{
		"customer": "Yodsakorn",
		"total":    1234,
	})

	if err != nil {
		log.Fatal(err)
	}

	// Write PDF to output directory
	outputPath := "./output/invoice.pdf"
	if err := os.MkdirAll("./output", 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(outputPath, pdfBytes, 0644); err != nil {
		log.Fatal(err)
	}

	log.Println("invoice.pdf generated 🎉")
}
