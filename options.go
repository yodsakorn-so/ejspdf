package ejspdf

type Options struct {
	Template Template
	Data     any
	PDF      PDFOptions
}

type Template struct {
	Path string
	Text string
}

func TemplateFile(path string) Template {
	return Template{Path: path}
}

func TemplateText(text string) Template {
	return Template{Text: text}
}

type PDFOptions struct {
	Format string // A4
}
