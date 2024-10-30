package server

import (
	"encoding/json"
	"errors"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/co2computation"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/keycloak"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

// getAuswertung fuehrt die CO2-Emissionsberechnung fuer die uebertragene Umfrage durch und sendet einen
// structs.AuswertungRes zurueck.
func getAuswertung(res http.ResponseWriter, req *http.Request) {
	var requestedUmfrageID primitive.ObjectID
	err := requestedUmfrageID.UnmarshalText([]byte(req.URL.Query().Get("id")))
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// hole Umfragen aus der Datenbank
	umfrage, err := database.UmfrageFind(requestedUmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// Wenn Auswertung nicht fuers teilen freigegeben muss Nutzer authentifiziert sein
	var nutzername string
	if umfrage.AuswertungFreigegeben == 0 {
		// Keycloak Authentifizierung
		ctx := req.Context()

		authHeader := req.Header.Get("Authorization")
		if len(authHeader) < 1 {
			res.WriteHeader(401)
			return
		}

		accessToken := strings.Split(authHeader, " ")[1]

		rptResult, err := keycloak.KeycloakClient.RetrospectToken(ctx, accessToken, clientID, clientSecret, realm)
		if err != nil {
			log.Println(err)
			res.WriteHeader(403)
			return
		}

		isTokenValid := *rptResult.Active

		if !isTokenValid {
			log.Println("Invalid Token")
			res.WriteHeader(401)
			return
		}

		// Pruefe, ob Nutzer Admin ist oder der Besitzer der Umfrage
		nutzername, err = keycloak.GetUsernameFromToken(accessToken, req.Context())
		if err != nil {
			errorResponse(res, err, http.StatusBadRequest)
			return
		}

		nutzer, _ := database.NutzerdatenFind(nutzername)
		if nutzer.Rolle != 1 && !isOwnerOfUmfrage(nutzer.UmfrageRef, requestedUmfrageID) {
			errorResponse(res, structs.ErrNutzerHatKeineBerechtigung, http.StatusUnauthorized)
			return
		}
	}

	mitarbeiterumfragen, err := database.MitarbeiterumfrageFindMany(umfrage.MitarbeiterumfrageRef)
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

	auswertung.EmissionenWaerme, auswertung.VerbrauchWearme, err = co2computation.BerechneEnergieverbrauch(umfrage.Gebaeude, umfrage.Jahr, structs.IDEnergieversorgungWaerme)
	if err != nil {
		if errors.Is(err, structs.ErrJahrNichtVorhanden) {
			auswertung.EmissionenWaerme = -1 // Fuer Frontend zum Hinweis, dass keine Auswertung fuer Jahr moeglich
		} else {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
	}

	auswertung.EmissionenStrom, auswertung.VerbrauchStrom, err = co2computation.BerechneEnergieverbrauch(umfrage.Gebaeude, umfrage.Jahr, structs.IDEnergieversorgungStrom)
	if err != nil {
		if errors.Is(err, structs.ErrJahrNichtVorhanden) {
			auswertung.EmissionenStrom = -1 // Fuer Frontend zum Hinweis, dass keine Auswertung fuer Jahr moeglich
		} else {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
	}

	auswertung.EmissionenKaelte, auswertung.VerbrauchKaelte, err = co2computation.BerechneEnergieverbrauch(umfrage.Gebaeude, umfrage.Jahr, structs.IDEnergieversorgungKaelte)
	if err != nil {
		if errors.Is(err, structs.ErrJahrNichtVorhanden) {
			auswertung.EmissionenKaelte = -1 // Fuer Frontend zum Hinweis, dass keine Auswertung fuer Jahr moeglich
		} else {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
	}

	// andere Emissionen
	var emissionenAufgeteiltITGeraeteHauptverantwortlicher map[int32]float64

	auswertung.EmissionenITGeraeteHauptverantwortlicher, emissionenAufgeteiltITGeraeteHauptverantwortlicher, err = co2computation.BerechneITGeraete(umfrage.ITGeraete)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	auswertung.EmissionenITGeraeteAufgeteilt = emissionenAufgeteiltITGeraeteHauptverantwortlicher

	// aggregiere alle Daten aus den Mitarbeiterumfragen
	var alleDiensreisen []structs.UmfrageDienstreise
	var alleITGerate []structs.UmfrageITGeraete
	var allePendelwege []structs.AllePendelwege

	for _, mitarbeiterumfrage := range mitarbeiterumfragen {
		alleITGerate = append(alleITGerate, mitarbeiterumfrage.ITGeraete...)
		alleDiensreisen = append(alleDiensreisen, mitarbeiterumfrage.Dienstreise...)
		allePendelwege = append(allePendelwege, structs.AllePendelwege{mitarbeiterumfrage.Pendelweg, mitarbeiterumfrage.TageImBuero})
	}

	// IT-Geraete
	emission, emissionenAufgeteiltITGeraeteMitarbeiter, err := co2computation.BerechneITGeraete(alleITGerate)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	auswertung.EmissionenITGeraeteMitarbeiter = emission

	// Dienstreisen
	emission, emissionenAufgeteiltDienstreisen, err := co2computation.BerechneDienstreisen(alleDiensreisen)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	auswertung.EmissionenDienstreisen = emission
	auswertung.EmissionenDienstreisenAufgeteilt = emissionenAufgeteiltDienstreisen

	// Pendelwege
	emission, emissionenAufgeteiltPendelwege, err := co2computation.BerechnePendelweg(allePendelwege)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	auswertung.EmissionenPendelwege = emission
	auswertung.EmissionenPendelwegeAufgeteilt = emissionenAufgeteiltPendelwege

	// Hochrechnung der Mitarbeiteremissionen
	if auswertung.Umfragenanzahl != 0 { // Hochrechnung nur falls Mitarbeiterumfragen vorhanden
		factor := float64(auswertung.Mitarbeiteranzahl) / float64(auswertung.Umfragenanzahl)

		auswertung.EmissionenITGeraeteMitarbeiter = math.Round(auswertung.EmissionenITGeraeteMitarbeiter*factor*100) / 100 // Rundung auf 2 Nachkommastellen
		auswertung.EmissionenPendelwege = math.Round(auswertung.EmissionenPendelwege*factor*100) / 100
		auswertung.EmissionenDienstreisen = math.Round(auswertung.EmissionenDienstreisen*factor*100) / 100

		for key, value := range auswertung.EmissionenDienstreisenAufgeteilt {
			auswertung.EmissionenDienstreisenAufgeteilt[key] = math.Round(value*factor*100) / 100
		}
		for key, value := range auswertung.EmissionenPendelwegeAufgeteilt {
			auswertung.EmissionenPendelwegeAufgeteilt[key] = math.Round(value*factor*100) / 100
		}
		for key, value := range emissionenAufgeteiltITGeraeteMitarbeiter {
			emissionenAufgeteiltITGeraeteMitarbeiter[key] = math.Round(value*factor*100) / 100
		}
	}

	// Zusammenfassung der Emissionen
	for key, value := range emissionenAufgeteiltITGeraeteMitarbeiter {
		e, ok := auswertung.EmissionenITGeraeteAufgeteilt[key]
		if ok {
			auswertung.EmissionenITGeraeteAufgeteilt[key] = e + value
		} else {
			auswertung.EmissionenITGeraeteAufgeteilt[key] = value
		}
	}

	if auswertung.EmissionenWaerme > 0 {
		auswertung.EmissionenEnergie += auswertung.EmissionenWaerme
	}
	if auswertung.EmissionenStrom > 0 {
		auswertung.EmissionenEnergie += auswertung.EmissionenStrom
	}
	if auswertung.EmissionenKaelte > 0 {
		auswertung.EmissionenEnergie += auswertung.EmissionenKaelte
	}

	auswertung.EmissionenITGeraete = auswertung.EmissionenITGeraeteMitarbeiter + auswertung.EmissionenITGeraeteHauptverantwortlicher
	auswertung.EmissionenGesamt = auswertung.EmissionenPendelwege + auswertung.EmissionenITGeraete + auswertung.EmissionenDienstreisen + auswertung.EmissionenEnergie

	auswertung.VerbrauchEnergie = auswertung.VerbrauchKaelte + auswertung.VerbrauchStrom + auswertung.VerbrauchWearme

	if auswertung.Mitarbeiteranzahl > 0 {
		auswertung.EmissionenProMitarbeiter = math.Round(auswertung.EmissionenGesamt/float64(auswertung.Mitarbeiteranzahl)*100) / 100
	}

	auswertung.Vergleich2PersonenHaushalt = math.Round(auswertung.EmissionenGesamt/structs.Verbrauch2PersonenHaushalt*100) / 100
	auswertung.Vergleich4PersonenHaushalt = math.Round(auswertung.EmissionenGesamt/structs.Verbrauch4PersonenHaushalt*100) / 100

	// Datenlücken-Visualisierung
	auswertung.GebaeudeIDsUndZaehler, err = database.GebaeudeAlleNrUndZaehlerRef()
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	alleZaehler, err := database.ZaehlerAlleZaehlerUndDaten()
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	auswertung.Zaehler = binaereZahlerdatenFuerZaehler(alleZaehler)

	auswertung.UmfrageGebaeude = umfrage.Gebaeude

	sendResponse(res, true, auswertung, http.StatusOK)
}

// binaereZahlerdatenFuerZaehler gibt fuer eine Liste an Zaehlern mit Zaehlerdaten eine binaere Liste fuer jeden
// Zaehler, ob Daten fuer dieses Jahr vorhanden sind
func binaereZahlerdatenFuerZaehler(alleZaehler []structs.ZaehlerUndZaehlerdaten) []structs.ZaehlerUndZaehlerdatenVorhanden {
	var zaehlerUndZaehlerdaten []structs.ZaehlerUndZaehlerdatenVorhanden
	aktuellesJahr := int32(time.Now().Year())

	for _, zaehler := range alleZaehler {
		var new_zaehler structs.ZaehlerUndZaehlerdatenVorhanden
		new_zaehler.ZaehlerID = zaehler.ZaehlerID

		for i := structs.ErstesJahr; i <= aktuellesJahr; i++ {
			found := false

			for _, zaehlerwert := range zaehler.Zaehlerdaten {
				if i == int32(zaehlerwert.Zeitstempel.Year()) {
					found = true
				}
			}

			new_zaehler.ZaehlerdatenVorhanden = append(
				new_zaehler.ZaehlerdatenVorhanden,
				structs.ZaehlerwertVorhanden{
					Jahr:      int32(i),
					Vorhanden: found,
				})
		}

		zaehlerUndZaehlerdaten = append(zaehlerUndZaehlerdaten, new_zaehler)
	}

	return zaehlerUndZaehlerdaten
}

// updateLinkShare empfaengt ein POST Update Request und setzt den LinkSharing Status auf den empfangenen Wert.
// 0 = Link Share deaktiviert, 1 = aktiviert.
// Kann nicht durch einen Administrator geaendert werden, nur durch besitzenden Nutzer.
func updateLinkShare(res http.ResponseWriter, req *http.Request) {
	nutzername, err := keycloak.GetUsernameFromToken(strings.Split(req.Header.Get("Authorization"), " ")[1], req.Context())
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

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
	nutzer, _ := database.NutzerdatenFind(nutzername)
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
