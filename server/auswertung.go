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
	"math"
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
		auswertung.Umfragenanteil = math.Round(float64(auswertung.Umfragenanzahl)/float64(auswertung.Mitarbeiteranzahl)*1000) / 1000
	}
	auswertung.AuswertungFreigegeben = umfrage.AuswertungFreigegeben

	auswertung.EmissionenWaerme, err = co2computation.BerechneEnergieverbrauch(umfrage.Gebaeude, umfrage.Jahr, structs.IDEnergieversorgungWaerme)
	if err != nil {
		if errors.Is(err, structs.ErrJahrNichtVorhanden) {
			auswertung.EmissionenWaerme = -1 // Fuer Frontend zum Hinweis, dass keine Auswertung fuer Jahr moeglich
		} else {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
	}

	auswertung.EmissionenStrom, err = co2computation.BerechneEnergieverbrauch(umfrage.Gebaeude, umfrage.Jahr, structs.IDEnergieversorgungStrom)
	if err != nil {
		if errors.Is(err, structs.ErrJahrNichtVorhanden) {
			auswertung.EmissionenStrom = -1 // Fuer Frontend zum Hinweis, dass keine Auswertung fuer Jahr moeglich
		} else {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
	}

	auswertung.EmissionenKaelte, err = co2computation.BerechneEnergieverbrauch(umfrage.Gebaeude, umfrage.Jahr, structs.IDEnergieversorgungKaelte)
	if err != nil {
		if errors.Is(err, structs.ErrJahrNichtVorhanden) {
			auswertung.EmissionenKaelte = -1 // Fuer Frontend zum Hinweis, dass keine Auswertung fuer Jahr moeglich
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

		auswertung.EmissionenITGeraeteMitarbeiter = math.Round(auswertung.EmissionenITGeraeteMitarbeiter*factor*100) / 100 // Rundung auf 2 Nachkommastellen
		auswertung.EmissionenPendelwege = math.Round(auswertung.EmissionenPendelwege*factor*100) / 100
		auswertung.EmissionenDienstreisen = math.Round(auswertung.EmissionenDienstreisen*factor*100) / 100
	}

	auswertung.EmissionenITGeraete = auswertung.EmissionenITGeraeteMitarbeiter + auswertung.EmissionenITGeraeteHauptverantwortlicher
	auswertung.EmissionenEnergie = auswertung.EmissionenWaerme + auswertung.EmissionenStrom + auswertung.EmissionenKaelte
	auswertung.EmissionenGesamt = auswertung.EmissionenPendelwege + auswertung.EmissionenITGeraete + auswertung.EmissionenDienstreisen + auswertung.EmissionenEnergie

	if auswertung.Mitarbeiteranzahl > 0 {
		auswertung.EmissionenProMitarbeiter = math.Round(auswertung.EmissionenGesamt/float64(auswertung.Mitarbeiteranzahl)*100) / 100
	}

	auswertung.Vergleich2PersonenHaushalt = math.Round(auswertung.EmissionenGesamt/structs.Verbrauch2PersonenHaushalt*100) / 100
	auswertung.Vergleich4PersonenHaushalt = math.Round(auswertung.EmissionenGesamt/structs.Verbrauch4PersonenHaushalt*100) / 100

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

	// Pruefe ob uebermittelter Freigabewert korrekt ist, d.h. 0 oder 1
	if linkshareReq.Freigabewert != 0 && linkshareReq.Freigabewert != 1 {
		var err = errors.New("Anfrage ungueltig")
		errorResponse(res, err, http.StatusBadRequest)
	}

	// Authentifizierung
	if !AuthWithResponse(res, linkshareReq.Auth.Username, linkshareReq.Auth.Sessiontoken) {
		return
	}
	nutzer, _ := database.NutzerdatenFind(linkshareReq.Auth.Username)
	if nutzer.Rolle != structs.IDRolleAdmin && !isOwnerOfUmfrage(nutzer.UmfrageRef, linkshareReq.UmfrageID) {
		errorResponse(res, structs.ErrNutzerHatKeineBerechtigung, http.StatusUnauthorized)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostUpdateActivateLinkShare")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	_, err = database.UmfrageUpdateLinkShare(linkshareReq.Freigabewert, linkshareReq.UmfrageID)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Println(err2)
		} else {
			err := database.RemoveDump(ordner)
			if err != nil {
				log.Println(err)
			}
		}
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.RemoveDump(ordner)
	if err != nil {
		log.Println(err)
	}

	sendResponse(res, true, nil, http.StatusOK)
}
