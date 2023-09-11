package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/keycloak"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// RouteNutzerdaten mounted alle aufrufbaren API Endpunkte unter */nutzerdaten
func RouteNutzerdaten() chi.Router {
	r := chi.NewRouter()

	r.Delete("/deleteNutzerdaten", DeleteNutzerdaten)

	return r
}

func CheckUser(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	accessToken := strings.Split(req.Header.Get("Authorization"), " ")[1]
	userInfo, err := keycloak.KeycloakClient.GetUserInfo(ctx, accessToken, realm)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	nutzername := *userInfo.PreferredUsername // TODO: check if null pointer

	// check if user already exists
	_, err = database.NutzerdatenFind(nutzername)
	if err == nil {
		sendResponse(res, true, nil, http.StatusOK)
		return
	}

	// check if there is an account with same email for migration
	// TODO: Accounnt migration

	// create new user
	// TODO: create new user

	//restorepath, err := database.CreateDump("PostRegistrierung")
	//if err != nil {
	//	errorResponse(res, err, http.StatusInternalServerError)
	//	return
	//}
	//id, err := database.NutzerdatenInsert()
	//if err != nil {
	//	err2 := database.RestoreDump(restorepath)
	//	if err2 != nil {
	//		// Datenbank konnte nicht wiederhergestellt werden
	//		log.Println(err2)
	//	} else {
	//		err := database.RemoveDump(restorepath)
	//		if err != nil {
	//			log.Println(err)
	//		}
	//	}
	//	// Konnte keinen neuen Nutzer erstellen
	//	errorResponse(res, err, http.StatusConflict)
	//	return
	//}
	//
	//err = database.RemoveDump(restorepath)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//sendResponse(res, true, nil, http.StatusCreated)
}

type RolleRes struct {
	Rolle int32 `json:"rolle"`
}

func CheckRolle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	accessToken := strings.Split(req.Header.Get("Authorization"), " ")[1]
	userInfo, err := keycloak.KeycloakClient.GetUserInfo(ctx, accessToken, realm)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	nutzername := *userInfo.PreferredUsername

	nutzer, err := database.NutzerdatenFind(nutzername)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	sendResponse(res, true, RolleRes{
		Rolle: nutzer.Rolle,
	}, http.StatusOK)
}

// PostNutzerdatenDelete loescht den Nutzer, der die Loeschung angefragt hat.
// Gibt auftretende Errors zurück, bspw. interne Berechnungsfehler oder unauthorized access.
func DeleteNutzerdaten(res http.ResponseWriter, req *http.Request) {
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
	// check if user is admin if they do not want to delete themselves
	if deleteNutzerdatenReq.Username != deleteNutzerdatenReq.Auth.Username {
		nutzer, err := database.NutzerdatenFind(deleteNutzerdatenReq.Username)
		if err != nil {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}

		// if user is not an admin, return unauthorized error
		if nutzer.Rolle != structs.IDRolleAdmin {
			errorResponse(res, err, http.StatusUnauthorized)
			return
		}
	}

	ordner, err := database.CreateDump("PostDeleteNutzerdaten")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// delete user
	err = database.NutzerdatenDelete(deleteNutzerdatenReq.Username)
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

	sendResponse(res, true, deleteNutzerdatenReq.Username, http.StatusOK)
}
