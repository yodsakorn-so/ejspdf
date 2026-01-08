package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yodsakorn-so/ejspdf"
)

func main() {
	// 1. Setup Context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 2. Prepare Data & Go Functions
	data := map[string]any{
		"appName": "EJS PDF Report",
		"user":    "Tester",
		"products": []map[string]any{
			{"name": "Laptop", "price": 45000.50},
			{"name": "Mouse", "price": 599.00},
			{"name": "Keyboard", "price": 1200.00},
		},
		// Go Function 1: Return current timestamp
		"now": func() string {
			return time.Now().Format("02 Jan 2006 15:04:05")
		},
		// Go Function 2: Format currency
		"formatMoney": func(amount float64) string {
			return fmt.Sprintf("฿%s", formatComma(amount))
		},
	}

	// 3. Render
	fmt.Println("Rendering PDF...")
	// Note: We render main.ejs, which automatically includes header.ejs
	pdfBytes, err := ejspdf.RenderFromFile(ctx, "main.ejs", ejspdf.Options{
		Data: data,
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}

	// 4. Save
	if err := os.WriteFile("output_v040.pdf", pdfBytes, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success! Saved to output_v040.pdf 🎉")
}

// Helper function
func formatComma(n float64) string {
	return fmt.Sprintf("%.2f", n) // Simplified
}
