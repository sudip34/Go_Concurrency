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

	//define application routes
	mux.Get("/", app.HomePage)
	return mux
}
