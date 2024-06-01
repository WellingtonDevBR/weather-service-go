package main

import (
	"net/http"
	"weather-service/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/weather/{cep}", handlers.GetWeather)

	http.ListenAndServe(":8080", r)
}
