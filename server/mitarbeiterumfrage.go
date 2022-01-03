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
	r.Post("/exists", PostUmfrageExists)
	r.Post("/insertMitarbeiterUmfrage", PostMitarbeiterUmfrageInsert)

	return r
}

// PostUmfrageExists returns true if the given ID exists
func PostUmfrageExists(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfrageExistsReq := structs.UmfrageID{}
	umfrageExistsRes := structs.UmfrageID{}

	err = json.Unmarshal(s, &umfrageExistsReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	requestedUmfrageID, err := primitive.ObjectIDFromHex(umfrageExistsReq.UmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	foundUmfrage, err := database.UmfrageFind(requestedUmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// return empty string if id is nil
	if foundUmfrage.ID == primitive.NilObjectID {
		umfrageExistsRes.UmfrageID = ""
	} else {
		umfrageExistsRes.UmfrageID = foundUmfrage.ID.Hex()
	}

	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   umfrageExistsRes,
		Error:  nil,
	})
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
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
			log.Fatalln(err2)
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

	// Response
	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   umfrageExistsRes,
		Error:  nil,
	})
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
}
