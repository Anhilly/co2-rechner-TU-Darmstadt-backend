package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
)

func RouteDB() chi.Router {
	r := chi.NewRouter()

	r.Post("/addFaktor", PostAddFaktor)
	r.Post("/addZaehlerdaten", PostAddZaehlerdaten)
	r.Post("/insertZaehler", PostInsertZaehler)
	r.Post("/insertGebaeude", PostInsertGebaeude)

	return r
}

func PostAddFaktor(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.AddCO2Faktor{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	if !AuthWithResponse(res, data.Auth.Username, data.Auth.Sessiontoken) {
		return
	}
	nutzer, _ := database.NutzerdatenFind(data.Auth.Username)
	if nutzer.Rolle != 1 {
		errorResponse(res, err, http.StatusUnauthorized)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostAddFaktor")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.EnergieversorgungAddFaktor(data)
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

func PostAddZaehlerdaten(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.AddZaehlerdaten{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	if !AuthWithResponse(res, data.Auth.Username, data.Auth.Sessiontoken) {
		return
	}
	nutzer, _ := database.NutzerdatenFind(data.Auth.Username)
	if nutzer.Rolle != 1 {
		errorResponse(res, err, http.StatusUnauthorized)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostAddZaehlerdaten")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.ZaehlerAddZaehlerdaten(data)
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

func PostInsertZaehler(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.InsertZaehler{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	if !AuthWithResponse(res, data.Auth.Username, data.Auth.Sessiontoken) {
		return
	}
	nutzer, _ := database.NutzerdatenFind(data.Auth.Username)
	if nutzer.Rolle != 1 {
		errorResponse(res, err, http.StatusUnauthorized)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostInsertZaehler")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.ZaehlerInsert(data)
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

func PostInsertGebaeude(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.InsertGebaeude{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	if !AuthWithResponse(res, data.Auth.Username, data.Auth.Sessiontoken) {
		return
	}
	nutzer, _ := database.NutzerdatenFind(data.Auth.Username)
	if nutzer.Rolle != 1 {
		errorResponse(res, err, http.StatusUnauthorized)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostInsertGebaeude")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.GebaeudeInsert(data)
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
