package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
)

// getMitarbeiterUmfrageFuerUmfrage liefert alle Mitarbeiterumfragen,
// welche mit der Umfrage mit der ID UmfrageID assoziiert sind, zurueck.
func getMitarbeiterUmfrageFuerUmfrage(res http.ResponseWriter, req *http.Request) {
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

// getUmfrageExists liefert einen structs.UmfrageExistsRes zurueck, falls die Umfrage existiert,
// dabei wird auch zurueck geliefert, ob die Umfrage durch alle Mitarbeiter ausgefuellt wurde.
// Diese Funktion hat keine Authentifizierung, da sie fuer die Mitarbeiterumfrage benoetigt wird.
func getUmfrageExists(res http.ResponseWriter, req *http.Request) {
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

// postInsertMitarbeiterumfrage fuegt die empfangene Mitarbeiterumfrage in die DB ein und
// sendet null zurueck, wenn das Einfuegen erfolgreich war.
func postInsertMitarbeiterumfrage(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	insertMitarbeiterumfrageReq := structs.InsertMitarbeiterUmfrage{}

	err = json.Unmarshal(s, &insertMitarbeiterumfrageReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("postInsertMitarbeiterumfrage")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	_, err = database.MitarbeiterUmfrageInsert(insertMitarbeiterumfrageReq)
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
