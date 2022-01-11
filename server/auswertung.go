package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/co2computation"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func RouteAuswertung() chi.Router {
	r := chi.NewRouter()

	r.Get("/", GetAuswertung)

	return r
}

func GetAuswertung(res http.ResponseWriter, req *http.Request) {
	var umfrageID primitive.ObjectID
	err := umfrageID.UnmarshalText([]byte(req.URL.Query().Get("id")))
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// hole Umfragen aus der Datenbank
	umfrage, err := database.UmfrageFind(umfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	mitarbeiterumfragen, err := database.MitarbeiterUmfrageFindMany(umfrage.MitarbeiterUmfrageRef)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// Auswertung der Daten
	var auswertung structs.AuswertungRes

	// allgemeine Information der Umfrage
	auswertung.ID = umfrage.ID
	auswertung.Jahr = umfrage.Jahr
	auswertung.Mitarbeiteranzahl = umfrage.Mitarbeiteranzahl
	auswertung.Umfragenanzahl = int32(len(mitarbeiterumfragen))

	auswertung.EmissionenWaerme, err = co2computation.BerechneEnergieverbrauch(umfrage.Gebaeude, umfrage.Jahr, structs.IDEnergieversorgungWaerme)
	if err != nil {
		if errors.Is(err, structs.ErrJahrNichtVorhanden) {
			auswertung.EmissionenWaerme = -1
		} else {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
	}

	auswertung.EmissionenStrom, err = co2computation.BerechneEnergieverbrauch(umfrage.Gebaeude, umfrage.Jahr, structs.IDEnergieversorgungStrom)
	if err != nil {
		if errors.Is(err, structs.ErrJahrNichtVorhanden) {
			auswertung.EmissionenWaerme = -1
		} else {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
	}

	auswertung.EmissionenKaelte, err = co2computation.BerechneEnergieverbrauch(umfrage.Gebaeude, umfrage.Jahr, structs.IDEnergieversorgungKaelte)
	if err != nil {
		if errors.Is(err, structs.ErrJahrNichtVorhanden) {
			auswertung.EmissionenWaerme = -1
		} else {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
	}

	// andere Emissionen
	auswertung.EmissionenITGeraeteHauptverantwortlicher, err = co2computation.BerechneITGeraete(umfrage.ITGeraete)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	for _, mitarbeiterumfrage := range mitarbeiterumfragen {
		// IT-Geraete
		emission, err := co2computation.BerechneITGeraete(mitarbeiterumfrage.ITGeraete)
		if err != nil {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
		auswertung.EmissionenITGeraeteMitarbeiter += emission

		// Dienstreisen
		emission, err = co2computation.BerechneDienstreisen(mitarbeiterumfrage.Dienstreise)
		if err != nil {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
		auswertung.EmissionenDienstreisen += emission

		// Pendelwege
		emission, err = co2computation.BerechnePendelweg(mitarbeiterumfrage.Pendelweg, mitarbeiterumfrage.TageImBuero)
		if err != nil {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
		auswertung.EmissionenPendelwege += emission
	}

	// Response
	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   auswertung,
		Error:  nil,
	})
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
}
