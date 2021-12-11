package server

import (
	"encoding/json"
	"errors"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
)

var (
	falschesPasswortError = errors.New("Die Kombination aus Passwort und Email stimmt nicht überein")
)

func RouteAuthentication() chi.Router {
	r := chi.NewRouter()

	r.Post("/anmeldung", PostAnmeldung)
	r.Post("/registrierung", PostRegistrierung)

	return r
}

/**
Generiert einen Cookie Token welcher den Nutzer authentifiziert
*/
func generiereCookieToken(anmeldeReq structs.AnmeldungReq) string {
	sessionToken := uuid.NewString()
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	return sessionToken
}

/**
Die Funktion liefert einen Response welcher bei valider Benutzereingabe den Nutzer authentisiert, sonst Fehler
*/
func PostAnmeldung(res http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	anmeldeReq := structs.AnmeldungReq{}
	anmeldeRes := structs.AnmeldungRes{}
	json.Unmarshal(s, &anmeldeReq)

	nutzerdaten, err := database.NutzerdatenFind(anmeldeReq.Email)

	if err != nil {
		// Es existiert kein Account mit dieser Email
		anmeldeRes.Message = err.Error()
		anmeldeRes.Success = false
		response, _ := json.Marshal(anmeldeRes)
		res.WriteHeader(http.StatusUnauthorized) // 401
		res.Write(response)
	}
	// Vergleiche Passwort mit gespeichertem Hash
	evaluation := bcrypt.CompareHashAndPassword([]byte(nutzerdaten.Passwort), []byte(anmeldeReq.Passwort))

	if evaluation == nil {
		//Korrektes Passwort authentifiziere den Nutzer
		anmeldeRes.Message = "Benutzer authentifiziert"
		anmeldeRes.Success = true

		// Generiere Cookie Token
		anmeldeRes.Cookietoken = generiereCookieToken(anmeldeReq)

		response, _ := json.Marshal(anmeldeRes)
		res.WriteHeader(http.StatusOK) // 200
		res.Write(response)
	} else {
		// Falsches Passwort
		anmeldeRes.Message = falschesPasswortError.Error()
		anmeldeRes.Success = false
		response, _ := json.Marshal(anmeldeRes)
		res.WriteHeader(http.StatusUnauthorized) // 401
		res.Write(response)
	}
}

/**
Die Funktion liefert einen HTTP Response zurück, welcher den neuen Nutzer authentifiziert, oder eine Fehlermeldung liefert
*/
func PostRegistrierung(res http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	registrierungReq := structs.AnmeldungReq{}
	registrierungRes := structs.AnmeldungRes{}
	json.Unmarshal(s, &registrierungReq)

	err := database.NutzerdatenInsert(registrierungReq)
	if err == nil {
		// Created new user
		registrierungRes.Success = true
		registrierungRes.Message = "Der neue Nutzeraccount wurde erstellt."

		// Generiere Cookie Token
		registrierungRes.Cookietoken = generiereCookieToken(registrierungReq)

		response, _ := json.Marshal(registrierungRes)
		res.WriteHeader(http.StatusCreated) // 201 account created
		res.Write(response)
	} else {
		// Failed to create new user
		registrierungRes.Success = false
		registrierungRes.Message = err.Error()
		response, _ := json.Marshal(registrierungRes)
		res.WriteHeader(http.StatusConflict) // 409 Conflict
		res.Write(response)
	}
}
