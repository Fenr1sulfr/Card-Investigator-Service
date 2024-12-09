package main

import (
	"net/http"
)

func (app *application) proccessNotificationFile(w http.ResponseWriter, r *http.Request) {
	regNum := app.readRegistryNumberParam(r)
	notification, err := app.models.Notification.GetNotificationInfo(regNum)

	if err != nil {
		app.serverErrorRespone(w, r, err)
		return
	}
	var pdfBytes []byte
	// dateFormat := "02 января 2006 года" // Русский формат для даты
	// timeFormat := "15:04"               // 24-часовой формат для времени
	// timeDate := notification.Date.Format(dateFormat)
	// timeHours := notification.Date.Format(timeFormat)
	app.wg.Add(1)
	app.background(func() {
		// data := map[string]any{
		// 	"CaseNumber":      notification.CaseNumber,
		// 	"RecipientName":   notification.RecipientName,
		// 	"Date":            timeDate,
		// 	"Time":            timeHours,
		// 	"Location":        notification.Location,
		// 	"RelationToEvent": notification.RelationToEvent,
		// 	"Investigator":    notification.Investigator,
		// 	"CodeArticle":     notification.CodeArticle,
		// }
		// pdfBytes, err = app.pdfWorker.MakePdf("notification.html", data)
		pdfBytes, err = app.pdfWorker.CreatePDFNotification(*notification)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
		app.wg.Done()
	})
	app.wg.Wait()
	err = app.writePDF(w, http.StatusCreated, pdfBytes, nil)

}
