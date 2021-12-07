package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
)

const (
	port = ":9000"
)

func StartServer() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Diese Middleware befindet sich hier nur waehrend der Entwicklung,
	// die Verarbeitung und configuration von cors wird in Produktion von unserem Webserver uebernommen.
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, //nolint:gomnd    // Maximum value not ignored by any of major browsers
	}))

	r.Mount("/umfrage", RouteUmfrage())

	log.Println("Server Started")

	log.Fatalln(http.ListenAndServe(port, r))
}
