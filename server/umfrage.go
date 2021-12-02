package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RouteUmfrage() chi.Router{
	r := chi.NewRouter()

	r.Post("/mitarbeiter", PostMitarbeiter)

	return r
}

func PostMitarbeiter(res http.ResponseWriter , req *http.Request){
	
}