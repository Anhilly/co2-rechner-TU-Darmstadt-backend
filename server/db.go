package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RouteDB() chi.Router {
	r := chi.NewRouter()

	r.Post("/addFaktor", PostAddFaktor)
	r.Post("/addZaehlerdaten", PostAddZaehlerdaten)
	r.Post("/insertZaehler", PostInsertZaehler)
	r.Post("/insertGebaeude", PostInsertGebaeude)

	return r
}

func PostAddFaktor(res http.ResponseWriter, req *http.Request) {

}

func PostAddZaehlerdaten(res http.ResponseWriter, req *http.Request) {

}

func PostInsertZaehler(res http.ResponseWriter, req *http.Request) {

}

func PostInsertGebaeude(res http.ResponseWriter, req *http.Request) {

}
