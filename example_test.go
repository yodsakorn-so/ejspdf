package ejspdf_test

import (
	"context"
	"fmt"
	"log"

	"github.com/yodsakorn-so/ejspdf"
)

func ExampleRender() {
	ctx := context.Background()

	pdfBytes, err := ejspdf.Render(ctx, ejspdf.Options{
		Template: "<h1>Hello <%= name %></h1>",
		Data: map[string]any{
			"name": "World",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(pdfBytes) > 0)
	// Output: true
}
