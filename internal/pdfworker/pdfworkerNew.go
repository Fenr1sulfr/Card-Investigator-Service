package pdfworker

import (
	"api/internal/data"
	"fmt"
	"log"
	"path"
	"reflect"
	"time"

	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
)

type NewPdfWoker struct{}

var FontPath = path.Join("./temporary_files", "test.ttf")

func (n NewPdfWoker) CreatePDFNotification(d data.Notification) ([]byte, error) {
	// Initialize PDF with Portrait and A4 size
	cfg := config.NewBuilder().
		WithOrientation(orientation.Vertical).
		WithPageSize(pagesize.A4).
		WithLeftMargin(15).
		WithTopMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).
		Build()
	m := maroto.New(cfg)

	// Build sections of the PDF
	// 1. Header

	// addHeader(m)
	// 2. Invoice Number
	// addInvoiceDetails(m)
	// GetHeader(m)
	// 3. Item List
	// err := addItemList(m, data)
	// if err != nil {
	// 	return nil, err
	// }
	addNotificationContent(m, d)

	// 4. Footer - Signature and QR code
	// addFooter(m)

	// Save the PDF file
	document, err := m.Generate()
	if err != nil {
		log.Fatal(err.Error())
	}

	// err = document.Save("./invoice_sample.pdf")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// log.Println("PDF saved successfully.")
	return document.GetBytes(), nil
}

func addNotificationContent(m core.Maroto, notification data.Notification) error {

	var test string = "Судебная повестка"
	m.AddRow(20,
		text.NewCol(12, test,
			props.Text{
				Top:    5,
				Align:  align.Center,
				Size:   16,
				Family: FontPath,
			}))
	m.AddRow(20,
		text.NewCol(12, "по делу N - "+notification.CaseNumber, props.Text{
			Top:    5,
			Style:  fontstyle.Bold,
			Align:  align.Center,
			Size:   12,
			Family: FontPath,
		}),
	)

	m.AddRow(10,
		text.NewCol(6, "Date: "+time.Now().Format("02 Jan 2006"), props.Text{
			Align:  align.Left,
			Size:   10,
			Family: FontPath,
		}),
		text.NewCol(6, "Investigator - "+notification.Investigator, props.Text{
			Align:  align.Right,
			Size:   10,
			Family: FontPath,
		}),
	)
	m.AddRow(10, line.NewCol(12))
	// Title
	m.AddRow(10,
		text.NewCol(12, "Tit "+notification.CaseNumber, props.Text{
			Size:   12,
			Align:  align.Left,
			Family: FontPath,
		}),
	)

	// Message
	m.AddRow(10,
		text.NewCol(12, "Message: Гражданин(ка) "+notification.RecipientName+"Вам необходимо явиться "+notification.Date.String()+" В "+notification.Date.GoString()+" В качестве "+notification.RelationToEvent, props.Text{
			Size:   10,
			Align:  align.Left,
			Family: FontPath,
		}),
	)

	// Date
	m.AddRow(10,
		text.NewCol(12, "Date: "+notification.Date.String(), props.Text{
			Size:   10,
			Align:  align.Left,
			Family: FontPath,
		}),
	)

	// URL
	m.AddRow(10,
		text.NewCol(12, "Явка строго обязательна. В случае неявки Вы можете быть привлечены к уголовной ответственности по статье "+notification.CodeArticle+" уголовного кодекса РК", props.Text{
			Size:   10,
			Align:  align.Left,
			Family: FontPath,
		}),
	)

	return nil
}

// Adds a header to the PDF
// func addHeader(m core.Maroto) {

// 	m.AddRow(20,
// 		text.NewCol(12, "Судебная повестка",
// 			props.Text{
// 				Top:   5,
// 				Style: fontstyle.Bold,
// 				Align: align.Center,
// 				Size:  16,
// 			}))
// 	m.AddRow(20,
// 		text.NewCol(12, "по делу N - ", props.Text{
// 			Top:   5,
// 			Style: fontstyle.Bold,
// 			Align: align.Center,
// 			Size:  12,
// 		}),
// 	)
// }

// Adds invoice details
// func addInvoiceDetails(m core.Maroto) {
// 	m.AddRow(10,
// 		text.NewCol(6, "Date: "+time.Now().Format("02 Jan 2006"), props.Text{
// 			Align: align.Left,
// 			Size:  10,
// 		}),
// 		text.NewCol(6, "Investigator - "+, props.Text{
// 			Align: align.Right,
// 			Size:  10,
// 		}),
// 	)
// 	m.AddRow(10, line.NewCol(12))
// }

type InvoiceItem struct {
	Item            string
	Description     string
	Quantity        string
	Price           string
	DiscountedPrice string
	Total           string
}

func GetHeader(data any) core.Row {
	// Reflect the type of the first element in the slice
	t := reflect.TypeOf(data)

	// If the input is a slice, get the element type
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	// Ensure we are working with a struct
	if t.Kind() != reflect.Struct {
		panic("GetHeader: data must be a struct or a slice of structs")
	}

	// Create a new row for the header
	headerRow := row.New(10)

	// Iterate through the fields of the struct
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		field := t.Field(i)

		// Add each field name as a column header
		headerRow.Add(text.NewCol(12/numFields, field.Name, props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}))
	}

	return headerRow
}

// func (o InvoiceItem) GetContent(i int) core.Row {
// 	r := row.New(5).Add(
// 		text.NewCol(2, o.Item),
// 		text.NewCol(3, o.Description),
// 		text.NewCol(1, o.Quantity),
// 		text.NewCol(2, o.Price),
// 		text.NewCol(2, o.DiscountedPrice),
// 		text.NewCol(2, o.Total),
// 	)

// 	if i%2 == 0 {
// 		r.WithStyle(&props.Cell{
// 			BackgroundColor: &props.Color{Red: 240, Green: 240, Blue: 240},
// 		})
// 	}

// 	return r
// }

// func getObjects() []InvoiceItem {
// 	var items []InvoiceItem
// 	contents := [][]string{
// 		{"Laptop", "14-inch, 16GB RAM", "1", "$1200", "$1000", "$1000"},
// 		{"Mouse", "Wireless optical mouse", "2", "$25", "$20", "$40"},
// 		{"Keyboard", "Mechanical, RGB", "1", "$75", "$60", "$60"},
// 	}
// 	for i := 0; i < len(contents); i++ {
// 		items = append(items, InvoiceItem{
// 			Item:            contents[i][0],
// 			Description:     contents[i][1],
// 			Quantity:        contents[i][2],
// 			Price:           contents[i][3],
// 			DiscountedPrice: contents[i][4],
// 			Total:           contents[i][5],
// 		})
// 	}
// 	return items
// }

func addItemList(m core.Maroto, data any) error {
	v := reflect.ValueOf(data)

	// Check if the data is a slice
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("data must be a slice")
	}

	// Get the type of the first element
	if v.Len() == 0 {
		return fmt.Errorf("data slice is empty")
	}

	firstElem := v.Index(0).Interface()
	firstElemType := reflect.TypeOf(firstElem)

	// Create Headers
	headerRow := row.New(10)
	for i := 0; i < firstElemType.NumField(); i++ {
		field := firstElemType.Field(i)
		headerRow.Add(text.NewCol(12/firstElemType.NumField(), field.Name, props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}))
	}
	m.AddRows(headerRow)

	// Add Content Rows
	for i := 0; i < v.Len(); i++ {
		contentRow := row.New(10)
		elem := v.Index(i).Interface()
		elemValue := reflect.ValueOf(elem)

		// Iterate over fields
		for j := 0; j < elemValue.NumField(); j++ {
			fieldValue := elemValue.Field(j).Interface()
			contentRow.Add(text.NewCol(12/elemValue.NumField(), fmt.Sprintf("%v", fieldValue), props.Text{
				Align: align.Center,
			}))
		}

		// Alternate background color for rows
		if i%2 == 0 {
			contentRow.WithStyle(&props.Cell{
				BackgroundColor: &props.Color{Red: 240, Green: 240, Blue: 240},
			})
		}
		m.AddRows(contentRow)
	}

	return nil
}

// Adds a footer with total and signature
// func addFooter(m core.Maroto) {
// 	m.AddRow(15,
// 		text.NewCol(8, "Total Amount", props.Text{
// 			Top:   5,
// 			Style: fontstyle.Bold,
// 			Size:  10,
// 			Align: align.Right,
// 		}),
// 		text.NewCol(4, "$1100", props.Text{
// 			Top:   5,
// 			Style: fontstyle.Bold,
// 			Size:  10,
// 			Align: align.Center,
// 		}),
// 	)

// 	m.AddRow(40,
// 		signature.NewCol(6, "Authorized Signatory", props.Signature{FontFamily: fontfamily.Courier}),
// 		code.NewQrCol(6, "https://codeheim.io", props.Rect{
// 			Percent: 75,
// 			Center:  true,
// 		}),
// 	)
// }
