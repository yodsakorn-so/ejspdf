package pdf

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Options defines PDF rendering options.
type Options struct {
	ChromePath string

	PageSize  string
	Landscape bool

	// Custom size (overrides PageSize)
	PaperWidth  string
	PaperHeight string

	MarginTop    string
	MarginBottom string
	MarginLeft   string
	MarginRight  string

	// Header/Footer
	DisplayHeaderFooter bool
	HeaderTemplate      string
	FooterTemplate      string

	// Wait options
	WaitSelector string
	WaitDelay    time.Duration

	// Print options
	Scale            float64
	PageRanges       string
	IgnoreBackground bool
}

// Chrome represents a Chrome-based PDF renderer.
type Chrome struct {
	opt Options
}

// New creates a new Chrome PDF renderer.
func New(opt Options) *Chrome {
	return &Chrome{opt: opt}
}

// FromHTML converts an HTML string into a PDF document.
func (c *Chrome) FromHTML(ctx context.Context, html string) ([]byte, error) {
	// 1. Validation & Unit Conversion
	mt, mb, ml, mr, err := c.parseAllMargins()
	if err != nil {
		return nil, err
	}

	width, height, err := c.calculateDimensions()
	if err != nil {
		return nil, err
	}

	// 2. Chrome Setup (Reuse or Create)
	// Check if the provided context already has a chromedp session
	var chromeCtx context.Context
	var cancel context.CancelFunc

	if chromedp.FromContext(ctx) != nil {
		// Reuse existing session, but create a new tab (context)
		chromeCtx, cancel = chromedp.NewContext(ctx)
	} else {
		// Create new allocator and session
		// Default options usually include Headless, DisableGPU, etc.
		// We append NoSandbox to support running in CI/Docker environments.
		allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.NoSandbox,
		)
		
		if c.opt.ChromePath != "" {
			allocOpts = append(allocOpts, chromedp.ExecPath(c.opt.ChromePath))
		}
		allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, allocOpts...)
		defer allocCancel()

		chromeCtx, cancel = chromedp.NewContext(allocCtx)
	}
	defer cancel()

	var pdfBytes []byte

	// 3. Prepare Content
	encodedHTML := base64.StdEncoding.EncodeToString([]byte(html))
	dataURL := "data:text/html;charset=utf-8;base64," + encodedHTML

	// 4. Build Actions
	actions := []chromedp.Action{
		chromedp.Navigate(dataURL),
	}

	if c.opt.WaitSelector != "" {
		actions = append(actions, chromedp.WaitVisible(c.opt.WaitSelector))
	} else {
		actions = append(actions, chromedp.WaitReady("body"))
	}

	if c.opt.WaitDelay > 0 {
		actions = append(actions, chromedp.Sleep(c.opt.WaitDelay))
	}

	// Handle Header/Footer defaults
	headerTpl := c.opt.HeaderTemplate
	footerTpl := c.opt.FooterTemplate
	if c.opt.DisplayHeaderFooter {
		if headerTpl == "" {
			headerTpl = "<span> </span>"
		}
		if footerTpl == "" {
			footerTpl = "<span> </span>"
		}
	}

	// Handle Scale default
	scale := c.opt.Scale
	if scale <= 0 {
		scale = 1.0
	}

	// Print Action
	actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		pdfBytes, _, err = page.PrintToPDF().
			WithPrintBackground(!c.opt.IgnoreBackground).
			WithLandscape(c.opt.Landscape).
			WithPaperWidth(width).
			WithPaperHeight(height).
			WithMarginTop(mt).
			WithMarginBottom(mb).
			WithMarginLeft(ml).
			WithMarginRight(mr).
			WithDisplayHeaderFooter(c.opt.DisplayHeaderFooter).
			WithHeaderTemplate(headerTpl).
			WithFooterTemplate(footerTpl).
			WithScale(scale).
			WithPageRanges(c.opt.PageRanges).
			Do(ctx)
		return err
	}))

	// 5. Execute
	if err := chromedp.Run(chromeCtx, actions...); err != nil {
		return nil, fmt.Errorf("chromedp run failed: %w", err)
	}

	return pdfBytes, nil
}

func (c *Chrome) parseAllMargins() (mt, mb, ml, mr float64, err error) {
	if mt, err = parseMargin(c.opt.MarginTop); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid margin top: %w", err)
	}
	if mb, err = parseMargin(c.opt.MarginBottom); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid margin bottom: %w", err)
	}
	if ml, err = parseMargin(c.opt.MarginLeft); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid margin left: %w", err)
	}
	if mr, err = parseMargin(c.opt.MarginRight); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid margin right: %w", err)
	}
	return
}

func (c *Chrome) calculateDimensions() (float64, float64, error) {
	if c.opt.PaperWidth != "" && c.opt.PaperHeight != "" {
		w, err := parseMargin(c.opt.PaperWidth)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid paper width: %w", err)
		}
		h, err := parseMargin(c.opt.PaperHeight)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid paper height: %w", err)
		}
		return w, h, nil
	}
	w, h := getPageDimensions(c.opt.PageSize)
	return w, h, nil
}

func getPageDimensions(format string) (width, height float64) {
	switch format {
	case "A3":
		return 11.69, 16.54
	case "A5":
		return 5.83, 8.27
	case "Letter":
		return 8.5, 11.0
	case "Legal":
		return 8.5, 14.0
	case "Tabloid":
		return 11.0, 17.0
	case "A4":
		fallthrough
	default:
		return 8.27, 11.69
	}
}