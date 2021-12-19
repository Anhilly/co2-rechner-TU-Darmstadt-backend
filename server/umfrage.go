package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/co2computation"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
)

func RouteUmfrage() chi.Router {
	r := chi.NewRouter()

	// Posts
	r.Post("/mitarbeiter", PostMitarbeiter)
	r.Post("/hauptverantwortlicher", PostHauptverantwortlicher)

	// Get
	r.Get("/gebaeude", GetAllGebaeude)

	return r
}

// returns all gebaeude as []int32
func GetAllGebaeude(res http.ResponseWriter, req *http.Request) {
	gebaeudeRes := structs.AllGebaeudeRes{}

	gebaeudeRes.Gebaeude, _ = database.GebaeudeAlleNr()
	response, _ := json.Marshal(gebaeudeRes)

	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

//Temporaere Funktion zum testen des Frontends
func PostMitarbeiter(res http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	umfrageReq := structs.UmfrageMitarbeiterReq{}
	umfrageRes := structs.UmfrageMitarbeiterRes{}
	json.Unmarshal(s, &umfrageReq)
	umfrageRes.DienstreisenEmissionen, _ = co2computation.BerechneDienstreisen(umfrageReq.Dienstreise)
	umfrageRes.PendelwegeEmissionen, _ = co2computation.BerechnePendelweg(umfrageReq.Pendelweg, umfrageReq.TageImBuero)
	umfrageRes.ITGeraeteEmissionen, _ = co2computation.BerechneITGeraete(umfrageReq.ITGeraete)

	response, _ := json.Marshal(umfrageRes)

	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

//Temporaere Funktion zum testen des Frontends
func PostHauptverantwortlicher(res http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	umfrageReq := structs.UmfrageHauptverantwortlicherReq{}
	umfrageRes := structs.UmfrageHauptverantwortlicherRes{}
	json.Unmarshal(s, &umfrageReq)
	// TODO Jahr soll nicht hardcoded sein, sondern als parameter mitübergeben werden.
	umfrageRes.WaermeEmissionen, _ = co2computation.BerechneEnergieverbrauch(umfrageReq.Gebaeude, 2020, 1)
	umfrageRes.StromEmissionen, _ = co2computation.BerechneEnergieverbrauch(umfrageReq.Gebaeude, 2020, 2)
	umfrageRes.KaelteEmissionen, _ = co2computation.BerechneEnergieverbrauch(umfrageReq.Gebaeude, 2020, 3)
	umfrageRes.ITGeraeteEmissionen, _ = co2computation.BerechneITGeraete(umfrageReq.ITGeraete)

	response, _ := json.Marshal(umfrageRes)

	res.WriteHeader(http.StatusOK)
	res.Write(response)
}
