package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/co2computation"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

// RouteAuswertung mounted alle aufrufbaren API Endpunkte unter */auswertung
func RouteAuswertung() chi.Router {
	r := chi.NewRouter()

	r.Post("/", PostAuswertung)
	r.Post("/updateSetLinkShare", UpdateSetLinkShare)
	return r
}

// PostAuswertung fuehrt die CO2-Emissionen Berechnung fuer die uebertragene Umfrage durch und sendet einen
// structs.AuswertungRes zurueck.
func PostAuswertung(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	auswertungReq := structs.RequestUmfrage{}
	err = json.Unmarshal(s, &auswertungReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	var umfrageID = auswertungReq.UmfrageID

	// hole Umfragen aus der Datenbank
	umfrage, err := database.UmfrageFind(umfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// Wenn Auswertung nicht fuers teilen freigegeben muss Nutzer authoritaet geprueft werden
	if umfrage.AuswertungFreigegeben == 0 {
		// Authentifizierung
		if !AuthWithResponse(res, auswertungReq.Auth.Username, auswertungReq.Auth.Sessiontoken) {
			return
		}
		nutzer, _ := database.NutzerdatenFind(auswertungReq.Auth.Username)
		if nutzer.Rolle != 1 && !isOwnerOfUmfrage(nutzer.UmfrageRef, auswertungReq.UmfrageID) {
			errorResponse(res, structs.ErrNutzerHatKeineBerechtigung, http.StatusUnauthorized)
			return
		}
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
	auswertung.Bezeichnung = umfrage.Bezeichnung
	auswertung.Jahr = umfrage.Jahr
	auswertung.Mitarbeiteranzahl = umfrage.Mitarbeiteranzahl
	auswertung.Umfragenanzahl = int32(len(mitarbeiterumfragen))
	if auswertung.Mitarbeiteranzahl > 0 {
		auswertung.Umfragenanteil = float64(auswertung.Umfragenanzahl) / float64(auswertung.Mitarbeiteranzahl)
	}
	auswertung.AuswertungFreigeben = umfrage.AuswertungFreigegeben

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

	// Hochrechnung der Mitarbeiteremissionen
	if auswertung.Umfragenanzahl != 0 { // Hochrechnung nur falls Mitarbeiterumfragen vorhanden
		factor := (float64(auswertung.Mitarbeiteranzahl) / float64(auswertung.Umfragenanzahl))

		auswertung.EmissionenITGeraeteMitarbeiter *= factor
		auswertung.EmissionenPendelwege *= factor
		auswertung.EmissionenDienstreisen *= factor
	}

	auswertung.EmissionenITGeraete = auswertung.EmissionenITGeraeteMitarbeiter + auswertung.EmissionenITGeraeteHauptverantwortlicher
	auswertung.EmissionenEnergie = auswertung.EmissionenWaerme + auswertung.EmissionenStrom + auswertung.EmissionenKaelte
	auswertung.EmissionenGesamt = auswertung.EmissionenPendelwege + auswertung.EmissionenITGeraete + auswertung.EmissionenDienstreisen + auswertung.EmissionenEnergie

	if auswertung.Mitarbeiteranzahl > 0 {
		auswertung.EmissionenProMitarbeiter = auswertung.EmissionenGesamt / float64(auswertung.Mitarbeiteranzahl)
	}

	auswertung.Vergleich2PersonenHaushalt = auswertung.EmissionenGesamt / structs.Verbrauch2PersonenHaushalt
	auswertung.Vergleich4PersonenHaushalt = auswertung.EmissionenGesamt / structs.Verbrauch4PersonenHaushalt

	sendResponse(res, true, auswertung, http.StatusOK)
}

// UpdateSetLinkShare empfaengt ein POST Update Request und setzt den LinkSharing Status auf den empfangenen Wert.
// 0 = Link Share deaktiviert, 1 = aktiviert.
// Kann nicht durch einen Administrator geaendert werden, nur durch besitzenden Nutzer.
func UpdateSetLinkShare(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	linkshareReq := structs.RequestLinkShare{}
	err = json.Unmarshal(s, &linkshareReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Pruefe ob uebermittelter LinkShare Value korrekt ist, d.h. 0 oder 1
	if linkshareReq.LinkShare != 0 && linkshareReq.LinkShare != 1 {
		var err = errors.New("Anfrage ungueltig")
		errorResponse(res, err, http.StatusBadRequest)
	}

	// Authentifizierung
	if !AuthWithResponse(res, linkshareReq.Auth.Username, linkshareReq.Auth.Sessiontoken) {
		return
	}
	nutzer, _ := database.NutzerdatenFind(linkshareReq.Auth.Username)
	if nutzer.Rolle != 1 && !isOwnerOfUmfrage(nutzer.UmfrageRef, linkshareReq.UmfrageID) {
		errorResponse(res, structs.ErrNutzerHatKeineBerechtigung, http.StatusUnauthorized)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostUpdateActivateLinkShare")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	_, err = database.UmfrageUpdateLinkShare(linkshareReq.LinkShare, linkshareReq.UmfrageID)
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
