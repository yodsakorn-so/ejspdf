package pdf

import (
	"context"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func HTMLToPDF(ctx context.Context, html string) ([]byte, error) {
	var pdfBuf []byte

	// แปลง HTML เป็น data URL
	dataURL := "data:text/html;charset=utf-8," + url.PathEscape(html)

	err := chromedp.Run(ctx,
		// เปิด HTML
		chromedp.Navigate(dataURL),

		// รอให้ DOM และ layout render เสร็จจริง
		chromedp.WaitReady("body", chromedp.ByQuery),

		// เผื่อเวลาโหลด font / css (สำคัญมาก)
		chromedp.Sleep(5*time.Second),

		// สั่งพิมพ์เป็น PDF
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).
				WithPaperHeight(11.69).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				Do(ctx)
			if err != nil {
				return err
			}

			pdfBuf = buf
			return nil
		}),
	)

	return pdfBuf, err
}
