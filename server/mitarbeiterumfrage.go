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

// PostUpdateMitarbeiterUmfrage updated eine Mitarbeiterumfrage mit den empfangenen Daten
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
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	if !AuthWithResponse(res, umfrageReq.Auth.Username, umfrageReq.Auth.Sessiontoken) {
		return
	}
	nutzer, _ := database.NutzerdatenFind(umfrageReq.Auth.Username)
	if nutzer.Rolle != 1 {
		errorResponse(res, err, http.StatusUnauthorized)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostUpdateMitarbeiterUmfrage")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	umfrageID, err := database.MitarbeiterUmfrageUpdate(umfrageReq)
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

	// Response
	sendResponse(res, true, umfrageRes, http.StatusOK)
}

// GetMitarbeiterUmfrageForUmfrage liefert alle Mitarbeiterumfragen,
// welche mit der Umfrage mit der ID UmfrageID assoziiert sind, zurueck.
func GetMitarbeiterUmfrageForUmfrage(res http.ResponseWriter, req *http.Request) {
	var requestedUmfrageID primitive.ObjectID
	err := requestedUmfrageID.UnmarshalText([]byte(req.URL.Query().Get("id")))
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	mitarbeiterUmfragenRes := structs.AlleMitarbeiterUmfragenForUmfrage{}

	mitarbeiterUmfragenRes.MitarbeiterUmfragen, err = database.MitarbeiterUmfrageFindForUmfrage(requestedUmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	sendResponse(res, true, mitarbeiterUmfragenRes, http.StatusOK)
}

// GetUmfrageExists liefert einen structs.UmfrageExistsRes zurueck, falls die Umfrage existiert,
// dabei wird auch zurueck geliefert, ob die Umfrage durch alle Mitarbeiter ausgefuellt wurde.
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

// PostMitarbeiterUmfrageInsert fuegt die empfangene Mitarbeiterumfrage in die DB ein und
// sendet null zurueck, wenn das Einfuegen erfolgreich war.
func PostMitarbeiterUmfrageInsert(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfrageExistsReq := structs.InsertMitarbeiterUmfrage{}

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
	_, err = database.MitarbeiterUmfrageInsert(umfrageExistsReq)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Println(err2)
		}
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	sendResponse(res, true, nil, http.StatusOK)
}
