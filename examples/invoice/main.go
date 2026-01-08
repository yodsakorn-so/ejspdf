package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/yodsakorn-so/ejspdf"
)

func main() {

	// 1. Setup Context with Timeout

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()

	// 2. Data for Invoice

	data := map[string]any{

		"customer": "John Smith",

		"total": 5623,
	}

	// 3. Render PDF directly from file

	log.Println("Generating Invoice PDF from invoice.ejs...")

	pdfBytes, err := ejspdf.RenderFromFile(ctx, "invoice.ejs", ejspdf.Options{

		Data: data,

		PageSize: "A4",

		WaitDelay: 1 * time.Second,
	})

	if err != nil {

		log.Fatalf("Render failed: %v", err)

	}

	// 4. Save to file

	outputFile := "invoice.pdf"

	if err := os.WriteFile(outputFile, pdfBytes, 0644); err != nil {

		log.Fatalf("Save file failed: %v", err)

	}

	log.Println("Invoice generated successfully: " + outputFile + " 🎉")

}
