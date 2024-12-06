package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	//registration and authentication

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password-reset", app.updateUserPasswordHandler)
	router.HandlerFunc(http.MethodGet, "/v1/profile/:email", app.requireAuthenticatedUser(app.getUserHomeInfo))
	router.HandlerFunc(http.MethodPost, "/v1/cards/create", app.requirePermission("employee", app.createCard))
	router.HandlerFunc(http.MethodGet, "/v1/cards/:regnum", app.getCard)
	router.HandlerFunc(http.MethodGet, "/v1/list", app.listCardsByRegion)

	//new
	router.HandlerFunc(http.MethodGet, "/v1/notification/create/:regnum", app.proccessNotificationFile)

	//
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
