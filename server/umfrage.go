package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/co2computation"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
)

func RouteUmfrage() chi.Router {
	r := chi.NewRouter()

	r.Post("/mitarbeiter", PostMitarbeiter)

	return r
}

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
