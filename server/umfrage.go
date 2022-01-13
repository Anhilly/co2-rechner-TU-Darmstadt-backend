package server

import (
	"encoding/json"
	"fmt"
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
	r.Post("/updateUmfrage", PostUpdateUmfrage)

	// Get
	r.Get("/gebaeude", GetAllGebaeude)
	r.Get("/alleUmfragen", GetAllUmfragen)
	r.Get("/GetAllUmfragenForUser", GetAllUmfragenForUser)

	// Delete
	r.Delete("/deleteUmfrage", DeleteUmfrage)

	return r
}

func isOwnerOfUmfrage(umfrageRef []primitive.ObjectID, umfrageID primitive.ObjectID) bool {
	//TODO Error checking statt verwerfen
	for _, id := range umfrageRef {
		fmt.Println(id, " == ", umfrageID, "InLoop")
		if id == umfrageID {
			fmt.Println(id, " == ", umfrageID)
			return true
		}
	}
	return false
}

// PostUpdateUmfrage updates an umfrage with received values
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
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	if !AuthWithResponse(res, req, umfrageReq.Auth.Username, umfrageReq.Auth.Sessiontoken) {
		return
	}
	nutzer, _ := database.NutzerdatenFind(umfrageReq.Auth.Username)

	if nutzer.Rolle != 1 || !isOwnerOfUmfrage(nutzer.UmfrageRef, umfrageReq.UmfrageID) {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostUpdateUmfrage")
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

	umfrageID, err := database.UmfrageUpdate(umfrageReq)
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
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostAddFaktor")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = Authenticate(umfrageReq.Hauptverantwortlicher.Username, umfrageReq.Hauptverantwortlicher.Sessiontoken)
	if err != nil {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusUnauthorized,
			Message: "Ungueltige Anmeldedaten",
		}, http.StatusUnauthorized)
	}

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
		}, http.StatusBadRequest) // 400
		return
	}

	// Pruefe ob Nutzer authentifiziert ist, dann ob er die zu loeschende Umfrage besitzt
	err = Authenticate(umfrageReq.Hauptverantwortlicher.Username, umfrageReq.Hauptverantwortlicher.Sessiontoken)
	if err != nil {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusUnauthorized,
			Message: "Ungueltige Anmeldedaten",
		}, http.StatusUnauthorized) // 401
	}

	err = database.UmfrageDelete(umfrageReq.Hauptverantwortlicher.Username, umfrageReq.UmfrageID)
	if err != nil {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusInternalServerError,
			Message: "Datenbankeintrag nicht gefunden",
		}, http.StatusInternalServerError) // Nicht vorhanden -> 500
	}

	sendResponse(res, true, nil, http.StatusOK)
}

// GetAllUmfragen returns all Umfragen as structs.AlleUmfragen
func GetAllUmfragen(res http.ResponseWriter, req *http.Request) {
	umfragenRes := structs.AlleUmfragen{}

	var err error
	umfragenRes.Umfragen, err = database.AlleUmfragen()
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   umfragenRes,
		Error:  nil,
	})

	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
}

// GetAllUmfragenForUser returns all Umfragen for a given user as structs.AlleUmfragen
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

	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   umfragenRes,
		Error:  nil,
	})

	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
}
