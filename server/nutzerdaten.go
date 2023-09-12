package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/keycloak"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func CheckUser(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	accessToken := strings.Split(req.Header.Get("Authorization"), " ")[1]
	userInfo, err := keycloak.KeycloakClient.GetUserInfo(ctx, accessToken, realm)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	nutzername := *userInfo.PreferredUsername // TODO: check if null pointer

	// Pruefe, ob Nutzer bereits existiert
	_, err = database.NutzerdatenFind(nutzername)
	if err == nil { // Nutzer existiert bereits
		sendResponse(res, true, nil, http.StatusOK)
		return
	}

	// Pruefe, ob Nutzer mit E-Mail bereits existiert, um Account zu migrieren
	nutzer, err := database.NutzerdatenFindByEMail(*userInfo.Email)
	if err == nil { // Nutzer fuer Migration gefunden
		nutzer.Nutzername = *userInfo.PreferredUsername // aendere Nutzername

		restorepath, err := database.CreateDump("CheckUser")
		if err != nil {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}

		err = database.NutzerdatenUpdate(nutzer)
		if err != nil {
			err2 := database.RestoreDump(restorepath)
			if err2 != nil {
				// Datenbank konnte nicht wiederhergestellt werden
				log.Println(err2)
			} else {
				err := database.RemoveDump(restorepath)
				if err != nil {
					log.Println(err)
				}
			}
			// Konnte Nutzer nicht migrieren
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}

		sendResponse(res, true, nil, http.StatusOK)
		return
	}

	// Neuen Nutzer erstellen
	restorepath, err := database.CreateDump("CheckUser")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	_, err = database.NutzerdatenInsert(*userInfo.PreferredUsername, *userInfo.Email)
	if err != nil {
		err2 := database.RestoreDump(restorepath)
		if err2 != nil {
			// Datenbank konnte nicht wiederhergestellt werden
			log.Println(err2)
		} else {
			err := database.RemoveDump(restorepath)
			if err != nil {
				log.Println(err)
			}
		}
		// Konnte keinen neuen Nutzer erstellen
		errorResponse(res, err, http.StatusConflict)
		return
	}

	err = database.RemoveDump(restorepath)
	if err != nil {
		log.Println(err)
	}

	sendResponse(res, true, nil, http.StatusCreated)
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

	sendResponse(res, true, structs.PruefeRolleRes{
		Rolle: nutzer.Rolle,
	}, http.StatusOK)
}

// DeleteNutzerdaten loescht den Nutzer, der die Loeschung angefragt hat.
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

	ctx := req.Context()

	accessToken := strings.Split(req.Header.Get("Authorization"), " ")[1]
	userInfo, err := keycloak.KeycloakClient.GetUserInfo(ctx, accessToken, realm)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	nutzername := *userInfo.PreferredUsername // TODO: check if null pointer

	// check if user is admin if they do not want to delete themselves
	if deleteNutzerdatenReq.Username != nutzername {
		nutzer, err := database.NutzerdatenFind(nutzername)
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
