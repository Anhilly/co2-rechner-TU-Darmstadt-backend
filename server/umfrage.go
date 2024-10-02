package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/keycloak"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// isOwnerOfUmfrage prueft, ob die Umfrage in der Liste der Umfragen des Nutzers auftaucht.
// @param umfrageID Umfrage deren Nutzer bestimmt werden soll
// @param umfrageRef Liste aller dem Nutzer gehoerende Umfragen
func isOwnerOfUmfrage(umfrageRef []primitive.ObjectID, umfrageID primitive.ObjectID) bool {
	for _, id := range umfrageRef {
		if id == umfrageID {
			return true
		}
	}
	return false
}

// postUpdateUmfrage updated die Werte der Umfrage mit UmfrageID,
// wenn der authentifizierte Nutzer die Umfrage besitzt oder Admin ist.
// Liefert im Erfolgsfall die gleichgebliebene UmfrageID zurueck.
func postUpdateUmfrage(res http.ResponseWriter, req *http.Request) {
	umfrageReq := structs.UpdateUmfrage{}
	umfrageRes := structs.UmfrageID{}

	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	nutzername, err := keycloak.GetUsernameFromToken(strings.Split(req.Header.Get("Authorization"), " ")[1], req.Context())
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(s, &umfrageReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	nutzer, _ := database.NutzerdatenFind(nutzername)
	if nutzer.Rolle != 1 && !isOwnerOfUmfrage(nutzer.UmfrageRef, umfrageReq.UmfrageID) {
		errorResponse(res, structs.ErrNutzerHatKeineBerechtigung, http.StatusUnauthorized)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("p")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	umfrageID, err := database.UmfrageUpdate(umfrageReq)
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

	// return empty umfrage string if umfrageID is invalid
	if umfrageID == primitive.NilObjectID {
		umfrageRes.UmfrageID = ""
	} else {
		umfrageRes.UmfrageID = umfrageID.Hex()
	}

	sendResponse(res, true, umfrageRes, http.StatusOK)
}

// getUmfrage empfaengt GET Request und sendet Umfrage struct fuer passende UmfrageID zurueck,
// sofern auth Eigentuemer oder Admin
func getUmfrage(res http.ResponseWriter, req *http.Request) {
	nutzername, err := keycloak.GetUsernameFromToken(strings.Split(req.Header.Get("Authorization"), " ")[1], req.Context())
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	var requestedUmfrageID primitive.ObjectID
	err = requestedUmfrageID.UnmarshalText([]byte(req.URL.Query().Get("id")))
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// pruefe Zugriffsrechte
	nutzer, _ := database.NutzerdatenFind(nutzername)
	if nutzer.Rolle != 1 && !isOwnerOfUmfrage(nutzer.UmfrageRef, requestedUmfrageID) {
		errorResponse(res, structs.ErrNutzerHatKeineBerechtigung, http.StatusUnauthorized)
		return
	}

	umfrage, err := database.UmfrageFind(requestedUmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	sendResponse(res, true, umfrage, http.StatusOK)
}

// getAlleGebaeude sendet Response mit allen Gebaeuden in der Datenbank zurueck.
func getAlleGebaeude(res http.ResponseWriter, req *http.Request) {
	var err error
	gebaeudeRes := structs.AlleGebaeudeRes{}

	gebaeudeRes.Gebaeude, err = database.GebaeudeAlleNr()
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	sendResponse(res, true, gebaeudeRes, http.StatusOK)
}

// getAlleGebaeudeUndZaehler sendet Response mit allen Gebaeuden Nummern und den eingetragenen Zaehlern in der Datenbank zurueck.
// Zusaetzlich werden alle Zaehler mit Angabe, ob ein Wert für jedes von 2018 bis zum aktuellen Jahr vorhanden ist.
func getAlleGebaeudeUndZaehler(res http.ResponseWriter, req *http.Request) {

	// TODO: Überprüfung im Frontend

	var err error
	gebaeudeRes := structs.AlleGebaeudeUndZaehlerRes{}

	gebaeudeRes.Gebaeude, err = database.GebaeudeAlleNrUndZaehlerRef()
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	alleZaehler, err := database.ZaehlerAlleZaehlerUndDaten()
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	gebaeudeRes.Zaehler = binaereZahlerdatenFuerZaehler(alleZaehler)

	sendResponse(res, true, gebaeudeRes, http.StatusOK)
}

// postInsertUmfrage fuegt die empfangene Umfrage in die Datenbank ein
// sendet ein structs.UmfrageID mit DB ID gesetzt zurueck
func postInsertUmfrage(res http.ResponseWriter, req *http.Request) {
	umfrageReq := structs.InsertUmfrage{}
	umfrageRes := structs.UmfrageID{}

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

	err = json.Unmarshal(s, &umfrageReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("postInsertUmfrage")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	umfrageID, err := database.UmfrageInsert(umfrageReq, nutzername)
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

	// return empty umfrage string if umfrageID is invalid
	if umfrageID == primitive.NilObjectID {
		umfrageRes.UmfrageID = ""
	} else {
		umfrageRes.UmfrageID = umfrageID.Hex()
	}

	sendResponse(res, true, umfrageRes, http.StatusOK)
}

// duplicateUmfrage dupliziert die Umfrage mit übergebener ObjectID
// und sendet structs.UmfrageID mit neuer ObjectID zurück
func duplicateUmfrage(res http.ResponseWriter, req *http.Request) {
	duplicateReq := structs.DuplicateUmfrageReq{}
	umfrageRes := structs.UmfrageID{}

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

	err = json.Unmarshal(s, &duplicateReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	nutzer, _ := database.NutzerdatenFind(nutzername)
	if nutzer.Rolle != 1 && !isOwnerOfUmfrage(nutzer.UmfrageRef, duplicateReq.UmfrageID) {
		errorResponse(res, structs.ErrNutzerHatKeineBerechtigung, http.StatusUnauthorized)
		return
	}

	umfrage, err := database.UmfrageFind(duplicateReq.UmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	duplizierteUmfrage := structs.InsertUmfrage{
		Bezeichnung:       umfrage.Bezeichnung + " " + duplicateReq.Suffix,
		Mitarbeiteranzahl: umfrage.Mitarbeiteranzahl,
		Jahr:              umfrage.Jahr,
		Gebaeude:          umfrage.Gebaeude,
		ITGeraete:         umfrage.ITGeraete,
	}

	ordner, err := database.CreateDump("PostDuplicateUmfrage")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	umfrageID, err := database.UmfrageInsert(duplizierteUmfrage, nutzername)
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

	// return empty umfrage string if umfrageID is invalid
	if umfrageID == primitive.NilObjectID {
		umfrageRes.UmfrageID = ""
	} else {
		umfrageRes.UmfrageID = umfrageID.Hex()
	}

	sendResponse(res, true, umfrageRes, http.StatusOK)
}

// deleteUmfrage loescht eine Umfrage mit der empfangenen UmfrageID aus der Datenbank.
// Sendet im Erfolgsfall null zurueck
func deleteUmfrage(res http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	umfrageReq := structs.UmfrageIDRequest{}

	nutzername, err := keycloak.GetUsernameFromToken(strings.Split(req.Header.Get("Authorization"), " ")[1], req.Context())
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(s, &umfrageReq)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest) // 400
		return
	}

	// pruefe Zugriffsrechte
	nutzer, _ := database.NutzerdatenFind(nutzername)
	if nutzer.Rolle != 1 && !isOwnerOfUmfrage(nutzer.UmfrageRef, umfrageReq.UmfrageID) {
		errorResponse(res, structs.ErrNutzerHatKeineBerechtigung, http.StatusUnauthorized)
		return
	}

	err = database.UmfrageDelete(nutzername, umfrageReq.UmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	sendResponse(res, true, nil, http.StatusOK)
}

// getAlleUmfragen sendet alle Umfragen aus der DB in structs.AlleUmfragen zurueck
func getAlleUmfragen(res http.ResponseWriter, req *http.Request) {
	var err error
	umfragenRes := structs.AlleUmfragen{}

	umfragenRes.Umfragen, err = database.AlleUmfragen()
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	sendResponse(res, true, umfragenRes, http.StatusOK)
}

// getAlleUmfragenVonNutzer sendet alle Umfragen, die dem authentifizierten Nutzer gehoeren
// als structs.AlleUmfragen zurueck
func getAlleUmfragenVonNutzer(res http.ResponseWriter, req *http.Request) {
	nutzername, err := keycloak.GetUsernameFromToken(strings.Split(req.Header.Get("Authorization"), " ")[1], req.Context())
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfragenRes := structs.AlleUmfragen{}

	// hole Umfragen aus der Datenbank
	umfragenRes.Umfragen, err = database.AlleUmfragenForUser(nutzername)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	sendResponse(res, true, umfragenRes, http.StatusOK)
}

// getUmfrageJahr sendet das Bilanzierungsjahr fuer den GET Parameter UmfrageID zurueck.
// Diese Funktion muss ohne Authentifizierung funktionieren, da sie fuer die Mitarbeiterumfrage benoetigt wird
func getUmfrageJahr(res http.ResponseWriter, req *http.Request) {
	var requestedUmfrageID primitive.ObjectID
	err := requestedUmfrageID.UnmarshalText([]byte(req.URL.Query().Get("id")))
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfrageJahrRes := structs.UmfrageYearRes{}

	// hole Umfrage aus der Datenbank
	umfrage, err := database.UmfrageFind(requestedUmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// set year
	umfrageJahrRes.Jahr = umfrage.Jahr

	sendResponse(res, true, umfrageJahrRes, http.StatusOK)
}

// getSharedResults sendet zurueck ob die Auswertung der Umfrage mit der ID im GET Parameter UmfrageID zum teilen
// freigegeben ist.
// Diese Funktion muss ohne Authentifizierung funktionieren, da sie von beliebigen unregistrierten Nutzern aufgerufen wird.
func getSharedResults(res http.ResponseWriter, req *http.Request) {
	var requestedUmfrageID primitive.ObjectID
	err := requestedUmfrageID.UnmarshalText([]byte(req.URL.Query().Get("id")))
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	umfrageSharedRes := structs.UmfrageSharedResultsRes{}

	// hole Umfrage aus der Datenbank
	umfrage, err := database.UmfrageFind(requestedUmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	umfrageSharedRes.Freigegeben = umfrage.AuswertungFreigegeben

	sendResponse(res, true, umfrageSharedRes, http.StatusOK)
}

// postShareUmfrage empfaengt POST Request mit UmfrageID und fügt die Umfrage dem Nutzer des Requests hinzu
func postShareUmfrage(res http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	umfrageReq := structs.UmfrageIDRequest{}

	nutzername, err := keycloak.GetUsernameFromToken(strings.Split(req.Header.Get("Authorization"), " ")[1], req.Context())
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(s, &umfrageReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// pruefe, ob Umfrage existiert
	umfrage, err := database.UmfrageFind(umfrageReq.UmfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}
	umfrageNameResponse := structs.UmfrageShareRes{
		Bezeichnung:  umfrage.Bezeichnung,
		Hinzugefuegt: false,
	}

	// pruefe, ob Nutzer schon Zugriff auf Umfrage hat
	nutzer, err := database.NutzerdatenFind(nutzername)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	if isOwnerOfUmfrage(nutzer.UmfrageRef, umfrageReq.UmfrageID) {
		sendResponse(res, true, umfrageNameResponse, http.StatusOK)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("postShareUmfrage")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.NutzerdatenAddUmfrageref(nutzername, umfrageReq.UmfrageID)
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

	umfrageNameResponse.Hinzugefuegt = true

	sendResponse(res, true, umfrageNameResponse, http.StatusOK)
}
