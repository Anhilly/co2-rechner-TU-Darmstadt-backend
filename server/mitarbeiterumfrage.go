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

func RouteMitarbeiterUmfrage() chi.Router {
	r := chi.NewRouter()

	// POST
	r.Post("/insertMitarbeiterUmfrage", PostMitarbeiterUmfrageInsert)
	r.Post("/updateMitarbeiterUmfrage", PostUpdateMitarbeiterUmfrage)

	// GET
	r.Get("/exists", GetUmfrageExists)
	r.Get("/mitarbeiterUmfrageForUmfrage", GetMitarbeiterUmfrageForUmfrage)

	return r
}

// PostUpdateMitarbeiterUmfrage updates an mitarbeiterUmfrage with received values
func PostUpdateMitarbeiterUmfrage(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfrageReq := structs.UpdateMitarbeiterUmfrage{}
	umfrageRes := structs.UmfrageID{}

	err = json.Unmarshal(s, &umfrageReq)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostUpdateMitarbeiterUmfrage")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// TODO authentication does not work here? says email would not have a valid session token?
	//err = Authenticate(umfrageReq.Hauptverantwortlicher.Username, umfrageReq.Hauptverantwortlicher.Sessiontoken)
	//if err != nil {
	//	sendResponse(res, false, structs.Error{
	//		Code:    http.StatusUnauthorized,
	//		Message: "Ungueltige Anmeldedaten",
	//	}, http.StatusUnauthorized)
	//}

	umfrageID, err := database.MitarbeiterUmfrageUpdate(umfrageReq)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Fatal(err2)
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
	sendResponse(res, true, umfrageRes, http.StatusOK)
}

// GetMitarbeiterUmfrageForUmfrage returns all MitarbeiterUmfragen belonging to a certain given umfrageID
func GetMitarbeiterUmfrageForUmfrage(res http.ResponseWriter, req *http.Request) {
	var requestedUmfrageID primitive.ObjectID
	err := requestedUmfrageID.UnmarshalText([]byte(req.URL.Query().Get("id")))
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	mitarbeiterUmfragenRes := structs.AlleMitarbeiterUmfragenForUmfrage{}

	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	mitarbeiterUmfragenRes.MitarbeiterUmfragen, err = database.MitarbeiterUmfrageFindForUmfrage(requestedUmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	sendResponse(res, true, mitarbeiterUmfragenRes, http.StatusOK)
}

// GetUmfrageExists returns the umfrageID if the umfrage exists and whether it is already complete or not.
func GetUmfrageExists(res http.ResponseWriter, req *http.Request) {
	var requestedUmfrageID primitive.ObjectID
	err := requestedUmfrageID.UnmarshalText([]byte(req.URL.Query().Get("id")))
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfrageExistsRes := structs.UmfrageExistsRes{}

	umfrage, err := database.UmfrageFind(requestedUmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// return empty string if id is nil
	if umfrage.ID == primitive.NilObjectID {
		umfrageExistsRes.UmfrageID = ""

		sendResponse(res, true, umfrageExistsRes, http.StatusOK)
		return

	} else {
		umfrageExistsRes.UmfrageID = umfrage.ID.Hex()
		umfrageExistsRes.Bezeichnung = umfrage.Bezeichnung
	}

	mitarbeiterumfragen, err := database.MitarbeiterUmfrageFindMany(umfrage.MitarbeiterUmfrageRef)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	mitarbeiterMax := umfrage.Mitarbeiteranzahl
	umfragenFilled := int32(len(mitarbeiterumfragen))

	// check if umfrage is complete
	if umfragenFilled < mitarbeiterMax {
		umfrageExistsRes.Complete = false
	} else {
		umfrageExistsRes.Complete = true
	}
	sendResponse(res, true, umfrageExistsRes, http.StatusOK)
}

// PostMitarbeiterUmfrageInsert inserts the received Umfrage and returns the ID of the inserted Umfrage-Entry.
func PostMitarbeiterUmfrageInsert(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfrageExistsReq := structs.InsertMitarbeiterUmfrage{}
	umfrageExistsRes := structs.UmfrageID{}

	err = json.Unmarshal(s, &umfrageExistsReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostMitarbeiterUmfrageInsert")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	umfrageID, err := database.MitarbeiterUmfrageInsert(umfrageExistsReq)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Fatal(err2)
		}
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// return empty string if id is nil
	if umfrageID == primitive.NilObjectID {
		umfrageExistsRes.UmfrageID = ""
	} else {
		umfrageExistsRes.UmfrageID = umfrageID.Hex()
	}

	sendResponse(res, true, umfrageExistsRes, http.StatusOK)
}
