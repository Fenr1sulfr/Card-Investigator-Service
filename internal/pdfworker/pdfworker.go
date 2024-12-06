package pdfworker

import (
	"bytes"
	"embed"
	"html/template"
	"log"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

//go:embed templates/*
var pdfTemplateFS embed.FS

type PdfWorker struct {
}

func (p PdfWorker) MakePdf(templateFile string, data any) ([]byte, error) {

	tmpl, err := template.ParseFS(pdfTemplateFS, "templates/"+templateFile)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		return nil, err
	}
	var htmlBuffer bytes.Buffer
	err = tmpl.Execute(&htmlBuffer, data)
	if err != nil {
		return nil, err
	}

	// Generate the PDF
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		panic(err)
	}

	// Add the HTML as a page
	page := wkhtmltopdf.NewPageReader(&htmlBuffer)
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)
	pdfg.AddPage(page)

	// Generate the PDF
	err = pdfg.Create()
	if err != nil {
		panic(err)
	}

	return pdfg.Bytes(), nil
	// pdfg.Bytes()

	// // Save the PDF file
	// err = pdfg.WriteFile("temporary_files/styled_output.pdf")
	// if err != nil {
	// 	panic(err)
	// }
	// return nil,nil
}
