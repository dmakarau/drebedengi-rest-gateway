package main

import (
	"log"
	"net/http"

	"github.com/drebedengi-rest/internal/config"
	"github.com/drebedengi-rest/internal/handlers"
	"github.com/drebedengi-rest/internal/soap"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	h := &handlers.Handler{
		SOAP: soap.NewClient(cfg.APIId, cfg.Login, cfg.Password, cfg.SoapURL),
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(contentTypeJSON)

	r.Mount("/api/v1", handlers.NewRouter(h))

	log.Printf("Listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}

func contentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
