package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RouteDB() chi.Router {
	r := chi.NewRouter()

	r.Post("/addFaktor", PostAddFaktor)
	r.Post("/addZaehlerdaten", PostAddZaehlerdaten)
	r.Post("/addZaehler", PostAddZaehler)
	r.Post("/addGebaeude", PostAddGebaeude)

	return r
}

func PostAddFaktor(res http.ResponseWriter, req *http.Request) {

}

func PostAddZaehlerdaten(res http.ResponseWriter, req *http.Request) {

}

func PostAddZaehler(res http.ResponseWriter, req *http.Request) {

}

func PostAddGebaeude(res http.ResponseWriter, req *http.Request) {

}
