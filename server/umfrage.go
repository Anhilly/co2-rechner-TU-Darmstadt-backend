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
	r.Post("/insertUmfrage", PostInsertUmfrage)
	r.Post("/updateUmfrage", PostUpdateUmfrage)
	r.Post("/getUmfrage", GetUmfrage)

	// Get
	r.Get("/gebaeude", GetAllGebaeude)
	r.Get("/alleUmfragen", GetAllUmfragen)
	r.Get("/GetAllUmfragenForUser", GetAllUmfragenForUser)
	r.Get("/GetUmfrageYear", GetUmfrageYear)

	// Delete
	r.Delete("/deleteUmfrage", DeleteUmfrage)

	return r
}

func isOwnerOfUmfrage(umfrageRef []primitive.ObjectID, umfrageID primitive.ObjectID) bool {
	for _, id := range umfrageRef {
		if id == umfrageID {
			return true
		}
	}
	return false
}

/**
PostUpdateUmfrage updates an umfrage with received values
*/
func PostUpdateUmfrage(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfrageReq := structs.UpdateUmfrage{}
	umfrageRes := structs.UmfrageID{}

	err = json.Unmarshal(s, &umfrageReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	if !AuthWithResponse(res, umfrageReq.Auth.Username, umfrageReq.Auth.Sessiontoken) {
		return
	}
	nutzer, _ := database.NutzerdatenFind(umfrageReq.Auth.Username)
	if nutzer.Rolle != 1 && !isOwnerOfUmfrage(nutzer.UmfrageRef, umfrageReq.UmfrageID) {
		errorResponse(res, structs.ErrNutzerHatKeineBerechtigung, http.StatusUnauthorized)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostUpdateUmfrage")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	umfrageID, err := database.UmfrageUpdate(umfrageReq)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Println(err2)
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

	sendResponse(res, true, umfrageRes, http.StatusOK)
}

/**
Postrequest sendet Umfrage struct fuer passende UmfrageID zurueck, sofern auth Eigentuemer oder Admin
*/
func GetUmfrage(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfrageReq := structs.RequestUmfrage{}

	err = json.Unmarshal(s, &umfrageReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	if !AuthWithResponse(res, umfrageReq.Auth.Username, umfrageReq.Auth.Sessiontoken) {
		return
	}
	nutzer, _ := database.NutzerdatenFind(umfrageReq.Auth.Username)
	if nutzer.Rolle != 1 && !isOwnerOfUmfrage(nutzer.UmfrageRef, umfrageReq.UmfrageID) {
		errorResponse(res, structs.ErrNutzerHatKeineBerechtigung, http.StatusUnauthorized)
		return
	}

	umfrage, err := database.UmfrageFind(umfrageReq.UmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	sendResponse(res, true, umfrage, http.StatusOK)
}

/**
GetAllGebaeude liefert alle Gebaeude in der Datenbank zurueck.
*/
func GetAllGebaeude(res http.ResponseWriter, req *http.Request) {
	gebaeudeRes := structs.AllGebaeudeRes{}

	var err error
	gebaeudeRes.Gebaeude, err = database.GebaeudeAlleNr()
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	sendResponse(res, true, gebaeudeRes, http.StatusOK)
}

/**
PostInsertUmfrage inserts the received Umfrage and returns the ID of the Umfrage-Entry
*/
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
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostInsertUmfrage")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	if !AuthWithResponse(res, umfrageReq.Hauptverantwortlicher.Username, umfrageReq.Hauptverantwortlicher.Sessiontoken) {
		return
	}

	umfrageID, err := database.UmfrageInsert(umfrageReq)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Println(err2)
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

	sendResponse(res, true, umfrageRes, http.StatusOK)
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
		errorResponse(res, err, http.StatusBadRequest) // 400
		return
	}

	// Pruefe ob Nutzer authentifiziert ist, dann ob er die zu loeschende Umfrage besitzt
	if !AuthWithResponse(res, umfrageReq.Hauptverantwortlicher.Username, umfrageReq.Hauptverantwortlicher.Sessiontoken) {
		return
	}

	err = database.UmfrageDelete(umfrageReq.Hauptverantwortlicher.Username, umfrageReq.UmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	sendResponse(res, true, nil, http.StatusOK)
}

/**
GetAllUmfragen returns all Umfragen as structs.AlleUmfragen
*/
func GetAllUmfragen(res http.ResponseWriter, req *http.Request) {
	umfragenRes := structs.AlleUmfragen{}

	var err error
	umfragenRes.Umfragen, err = database.AlleUmfragen()
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	sendResponse(res, true, umfragenRes, http.StatusOK)
}

/**
GetAllUmfragenForUser returns all Umfragen for a given user as structs.AlleUmfragen
*/
func GetAllUmfragenForUser(res http.ResponseWriter, req *http.Request) {
	user := req.URL.Query().Get("user")
	if checkValidSessionToken(user) != nil {
		return
	}
	umfragenRes := structs.AlleUmfragen{}

	// hole Umfragen aus der Datenbank
	var err error
	umfragenRes.Umfragen, err = database.AlleUmfragenForUser(user)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	sendResponse(res, true, umfragenRes, http.StatusOK)
}

/**
GetUmfrageYear returns the bilanzierungsjahr for the given umfrage
*/
func GetUmfrageYear(res http.ResponseWriter, req *http.Request) {
	var requestedUmfrageID primitive.ObjectID
	err := requestedUmfrageID.UnmarshalText([]byte(req.URL.Query().Get("id")))
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfrageJahrRes := structs.UmfrageYearRes{}

	// hole Umfrage aus der Datenbank
	umfrage, err := database.UmfrageFind(requestedUmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// set year
	umfrageJahrRes.Jahr = umfrage.Jahr

	sendResponse(res, true, umfrageJahrRes, http.StatusOK)
}
