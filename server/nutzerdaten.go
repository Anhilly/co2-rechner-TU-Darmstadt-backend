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

// RouteNutzerdaten mounted alle aufrufbaren API Endpunkte unter */nutzerdaten
func RouteNutzerdaten() chi.Router {
	r := chi.NewRouter()

	r.Post("/deleteNutzerdaten", PostNutzerdatenDelete)

	return r
}

// PostNutzerdatenDelete loescht den Nutzer, der die Loeschung angefragt hat.
// Gibt auftretende Errors zurück, bspw. interne Berechnungsfehler oder unauthorized access.
func PostNutzerdatenDelete(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	deleteNutzerdatenReq := structs.DeleteNutzerdatenReq{}
	err = json.Unmarshal(s, &deleteNutzerdatenReq)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	if !AuthWithResponse(res, deleteNutzerdatenReq.Auth.Username, deleteNutzerdatenReq.Auth.Sessiontoken) {
		return
	}

	ordner, err := database.CreateDump("PostDeleteNutzerdaten")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// check if user is admin if they do not want to delete themselves
	if deleteNutzerdatenReq.Username != deleteNutzerdatenReq.Auth.Username {
		nutzer, err := database.NutzerdatenFind(deleteNutzerdatenReq.Auth.Username)
		if err != nil {
			errorResponse(res, err, http.StatusInternalServerError)
		}

		// if user is not an admin, return unauthorized error
		if nutzer.Rolle != structs.IDRolleAdmin {
			errorResponse(res, err, http.StatusUnauthorized)
		}
	}

	// delete user
	err = database.NutzerdatenDelete(deleteNutzerdatenReq.Username)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Println(err2)
		}
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	sendResponse(res, true, deleteNutzerdatenReq.Auth.Username, http.StatusOK)
}
