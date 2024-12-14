package main

import (
	"api/internal/data"
	"errors"
	"net/http"
)

func (app *application) GetAllUsersByCard(w http.ResponseWriter, r *http.Request) {
	regnum := app.readRegistryNumberParam(r)
	users, err := app.models.Confirmations.GetAllUsersByCard(regnum)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorRespone(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		app.serverErrorRespone(w, r, err)
	}
}
func (app *application) GetAllCardsByEmail(w http.ResponseWriter, r *http.Request) {
	regnum := app.readEmailParam(r)
	users, err := app.models.Confirmations.GetAllCardsByUser(regnum)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorRespone(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		app.serverErrorRespone(w, r, err)
	}
}

func (app *application) ConfirmCard(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RegNum string `json:"regnum"`
		UserID int64  `json:"user_id"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorRespone(w, r, err)
		return
	}
	err = app.models.Confirmations.Confirm(input.UserID, input.RegNum)
	if err != nil {
		app.serverErrorRespone(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusAccepted, envelope{"Confirmed": "Yes"}, nil)
}
