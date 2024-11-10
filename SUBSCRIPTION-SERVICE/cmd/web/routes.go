package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	//setup middleware
	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)

	//define application routes
	mux.Get("/", app.HomePage)

	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.Logout)
	mux.Get("/register", app.RegisterPage)
	mux.Post("/register", app.PostRegisterPage)
	mux.Get("/activate-account", app.ActivateAccount)

	// mux.Get("/test-email", func(w http.ResponseWriter, r *http.Request) {
	// 	mail := Mail{
	// 		Domain:      "localhost",
	// 		Host:        "localhost",
	// 		Port:        8025,
	// 		Encryption:  "none",
	// 		FromAddress: "info@saha.com",
	// 		FromName:    "Info",
	// 		ErrorChan:   make(chan error),
	// 	}

	// 	msg := Message{
	// 		To:      "me@here.com",
	// 		Subject: "TEST email",
	// 		Data:    "Hello ! Frist Test email sent",
	// 	}

	// 	mail.sendMail(msg, make(chan error))
	// })

	return mux
}
