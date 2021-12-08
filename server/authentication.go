package server

import (
	"github.com/go-chi/chi/v5"
)

func RouteAuthentication() chi.Router {
	r := chi.NewRouter()

	//r.Post("/anmeldung", PostAnmeldung)
	//r.Post("/registrierung", PostRegistrierung)

	return r
}

/*
//Temporaere Funktion zum testen des Frontends
func PostAnmeldung(res http.ResponseWriter, req *http.Request) {
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
*/
