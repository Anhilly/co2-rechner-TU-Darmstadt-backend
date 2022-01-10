package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
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

// GetAllGebaeude returns all gebaeude as []int32
func GetAllGebaeude(res http.ResponseWriter, req *http.Request) {
	gebaeudeRes := structs.AllGebaeudeRes{}

	var err error
	gebaeudeRes.Gebaeude, err = database.GebaeudeAlleNr()
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   gebaeudeRes,
		Error:  nil,
	})
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
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

// PostInsertUmfrage inserts the received Umfrage and returns the ID of the Umfrage-Entry
func PostInsertUmfrage(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfrageReq := structs.InsertUmfrage{}
	umfrageRes := structs.UmfrageID{}

	err = json.Unmarshal(s, &umfrageReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostAddFaktor")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// TODO check if umfrage is valid before inserting
	umfrageID, err := database.UmfrageInsert(umfrageReq)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Fatalln(err2)
		}
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// return empty umfrage string if umfrageID is invalid
	if umfrageID == primitive.NilObjectID {
		umfrageRes.UmfrageID = ""
	} else {
		umfrageRes.UmfrageID = umfrageID.Hex()
	}

	// Response
	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   umfrageRes,
		Error:  nil,
	})
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
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
