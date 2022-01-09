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
	r.Post("/exists", PostUmfrageExists) // TODO remove after updating frontend to GET requests
	r.Post("/insertMitarbeiterUmfrage", PostMitarbeiterUmfrageInsert)
	r.Post("/mitarbeiterUmfrageForUmfrage", PostMitarbeiterUmfrageForUmfrage) // TODO remove after updating frontend to GET requests
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
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
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

func PostMitarbeiterUmfrageForUmfrage(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	mitarbeiterUmfragenReq := structs.UmfrageID{}
	mitarbeiterUmfragenRes := structs.AlleMitarbeiterUmfragenForUmfrage{}

	err = json.Unmarshal(s, &mitarbeiterUmfragenReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	requestedUmfrageID, err := primitive.ObjectIDFromHex(mitarbeiterUmfragenReq.UmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	mitarbeiterUmfragenRes.MitarbeiterUmfragen, err = database.MitarbeiterUmfrageFindForUmfrage(requestedUmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   mitarbeiterUmfragenRes,
		Error:  nil,
	})
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
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

	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   mitarbeiterUmfragenRes,
		Error:  nil,
	})
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
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

// GetUmfrageExists returns true if the given ID exists
func GetUmfrageExists(res http.ResponseWriter, req *http.Request) {
	var requestedUmfrageID primitive.ObjectID
	err := requestedUmfrageID.UnmarshalText([]byte(req.URL.Query().Get("id")))
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfrageExistsRes := structs.UmfrageID{}

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
