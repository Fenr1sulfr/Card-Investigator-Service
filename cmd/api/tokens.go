package main

import (
	"api/internal/data"
	"api/internal/data/validator"
	"errors"
	"net/http"
	"time"
)

func (app *application) createPasswordResetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()

	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			v.AddError("email", "no matchin email address found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorRespone(w, r, err)
		}
		return
	}
	if !user.Activated {
		v.AddError("email", "user account must be activated")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	token, err := app.models.Tokens.New(user.ID, 45*time.Minute, data.ScopePasswordReset)
	if err != nil {
		app.serverErrorRespone(w, r, err)
		return
	}
	app.background(func() {
		data := map[string]any{
			"passwordResetToken": token.Plaintext,
		}
		err = app.mailer.Send(user.Email, "token_password_reset.html", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	env := envelope{"message": "an email will be sent to you containing password reset instructions"}

	err = app.writeJSON(w, http.StatusAccepted, env, nil)
	if err != nil {
		app.serverErrorRespone(w, r, err)
	}

}

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePassowrdPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	app.logger.PrintInfo(input.Email, nil)
	user, err := app.models.Users.GetByEmail(input.Email)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.invalidCredentialResponse(w, r)
		default:
			app.serverErrorRespone(w, r, err)
		}
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorRespone(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialResponse(w, r)
		return
	}
	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorRespone(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorRespone(w, r, err)
		return
	}

}
