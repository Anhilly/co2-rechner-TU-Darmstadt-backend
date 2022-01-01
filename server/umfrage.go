package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
)

func RouteUmfrage() chi.Router {
	r := chi.NewRouter()

	// Posts
	//r.Post("/mitarbeiter", PostMitarbeiter)
	r.Post("/insertUmfrage", PostInsertUmfrage)

	// Get
	r.Get("/gebaeude", GetAllGebaeude)

	// Delete
	r.Delete("/deleteUmfrage", DeleteUmfrage)

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
//func PostMitarbeiter(res http.ResponseWriter, req *http.Request) {
//	s, _ := ioutil.ReadAll(req.Body)
//	umfrageReq := structs.UmfrageMitarbeiterReq{}
//	umfrageRes := structs.UmfrageMitarbeiterRes{}
//	json.Unmarshal(s, &umfrageReq)
//	umfrageRes.DienstreisenEmissionen, _ = co2computation.BerechneDienstreisen(umfrageReq.Dienstreise)
//	umfrageRes.PendelwegeEmissionen, _ = co2computation.BerechnePendelweg(umfrageReq.Pendelweg, umfrageReq.TageImBuero)
//	umfrageRes.ITGeraeteEmissionen, _ = co2computation.BerechneITGeraete(umfrageReq.ITGeraete)
//
//	response, _ := json.Marshal(umfrageRes)
//
//	res.WriteHeader(http.StatusOK)
//	res.Write(response)
//}

func PostInsertUmfrage(res http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	umfrageReq := structs.InsertUmfrage{}
	umfrageRes := structs.UmfrageID{}
	json.Unmarshal(s, &umfrageReq)

	var umfrageID primitive.ObjectID

	// TODO check if umfrage is valid before inserting
	umfrageID, _ = database.UmfrageInsert(umfrageReq)

	// return empty umfrage string if umfrageID is invalid
	if umfrageID == primitive.NilObjectID {
		umfrageRes.UmfrageID = ""
	} else {
		umfrageRes.UmfrageID = umfrageID.Hex()
	}

	response, _ := json.Marshal(umfrageRes)

	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

/**
Loescht eine uebermittelte Umfrage, gegeben durch die UmfrageID
*/

func DeleteUmfrage(res http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	umfrageReq := structs.DeleteUmfrage{}

	err := json.Unmarshal(s, &umfrageReq)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// Pruefe ob Nutzer authentifiziert ist, dann ob er die zu loeschende Umfrage besitzt
	err = Authenticate(umfrageReq.Hauptverantwortlicher.Username, umfrageReq.Hauptverantwortlicher.Sessiontoken)
	if err != nil {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusUnauthorized,
			Message: "Ungueltige Anmeldedaten",
		}, http.StatusUnauthorized)
	}

	err = database.UmfrageDelete(umfrageReq.UmfrageID)
	if err != nil {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusNotFound,
			Message: "Datenbankeintrag nicht gefunden",
		}, http.StatusNotFound)
	}

	sendResponse(res, true, nil, http.StatusOK)
}
