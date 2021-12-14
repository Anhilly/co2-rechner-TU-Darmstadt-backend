package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Session struct {
	Sessiontoken string
	GenTime      time.Time
}

var (
	// AuthMap Speichert Authentication Daten für Benutzer
	// key:= email -> {Sessiontoken, TTL}
	AuthMap = make(map[string]Session)
)

func RouteAuthentication() chi.Router {
	r := chi.NewRouter()

	r.Post("/anmeldung", PostAnmeldung)
	r.Post("/registrierung", PostRegistrierung)
	r.Delete("/abmeldung", DeleteAbmeldung)

	return r
}

/**
Generiert einen Cookie Token, welcher den Nutzer authentifiziert und speichert ihn in Map
*/
func generiereSessionToken(email string) string {
	sessionToken := uuid.NewString()
	AuthMap[email] = Session{
		Sessiontoken: sessionToken,
		GenTime:      time.Now(),
	}
	return sessionToken
}

/**
Überprüft ob die email einen gültigen Sessiontoken registriert hat, falls ja nil
*/
func checkValidSessionToken(email string) error {
	//TTL ist 2 Stunden
	ttl := 2
	if (AuthMap[email] == Session{}) {
		return structs.ErrNutzerHatKeinenSessiontoken
	}
	entry := AuthMap[email]
	genTimePlusTTL := entry.GenTime.Add(time.Hour * time.Duration(ttl))
	if genTimePlusTTL.Before(time.Now()) {
		return structs.ErrAbgelaufenerSessiontoken
	}
	return nil
}

/**
Löscht den Cookie Token welcher den Nutzer mit email authentifiziert
*/
func loescheSessionToken(email string) error {
	err := checkValidSessionToken(email)
	if err == nil || err == structs.ErrAbgelaufenerSessiontoken {
		AuthMap[email] = Session{}
		return nil
	}
	return err
}

/**
sendResponse sendet Response zurück, bei Marshal Fehler sende 500 Code Error
 @param res Writer der den Response sendet
 @param data true falls normales Response Packet, false bei Error
 @param payload ist interface welches den data bzw. error struct enthält
 @param code ist der HTTP Header Code
*/
func sendResponse(res http.ResponseWriter, data bool, payload interface{}, code int32) {
	responseBuilder := structs.Response{}
	if data {
		responseBuilder.Status = structs.ResponseSuccess
		responseBuilder.Error = nil
		responseBuilder.Data = payload
	} else {
		responseBuilder.Status = structs.ResponseError
		responseBuilder.Data = nil
		responseBuilder.Error = payload
	}
	response, err := json.Marshal(responseBuilder)
	if err == nil {
		res.WriteHeader(int(code))
	} else {
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.Write(response)
}

/**
Die Funktion liefert einen Response welcher bei valider Benutzereingabe den Nutzer authentisiert, sonst Fehler
*/
func PostAnmeldung(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	anmeldeReq := structs.AuthReq{}
	err = json.Unmarshal(s, &anmeldeReq)

	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	nutzerdaten, err := database.NutzerdatenFind(anmeldeReq.Username)

	if err == io.EOF {
		// Es existiert kein Account mit dieser Email
		// Sende genauere Fehlermeldung zurück, statt EOF
		sendResponse(res, false, structs.Error{
			Code:    http.StatusUnauthorized,
			Message: structs.ErrNichtExistenteEmail.Error(),
		}, http.StatusUnauthorized)
		return
	}
	if err != nil {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	// Vergleiche Passwort mit gespeichertem Hash
	evaluation := bcrypt.CompareHashAndPassword([]byte(nutzerdaten.Passwort), []byte(anmeldeReq.Passwort))

	if evaluation == nil {
		// Korrektes Passwort authentifiziere den Nutzer
		// Generiere Cookie Token
		token := generiereSessionToken(anmeldeReq.Username)
		sendResponse(res, true, structs.AuthRes{
			Message:      "Nutzer authentifiziert",
			Sessiontoken: token,
		}, http.StatusOK)
	} else {
		// Falsches Passwort
		sendResponse(res, false, structs.Error{
			Code:    http.StatusUnauthorized,
			Message: structs.ErrFalschesPasswortError.Error(),
		}, http.StatusUnauthorized)
	}
}

/**
Die Funktion liefert einen HTTP Response zurück, welcher den neuen Nutzer authentifiziert, oder eine Fehlermeldung liefert
*/
func PostRegistrierung(res http.ResponseWriter, req *http.Request) {
	//TODO DB save and restore im Fehlerfall
	//TODO Error handling für ReadAll und Umarshal
	s, err := ioutil.ReadAll(req.Body)
	registrierungReq := structs.AuthReq{}
	registrierungRes := structs.Response{}
	err = json.Unmarshal(s, &registrierungReq)

	err = database.NutzerdatenInsert(registrierungReq)
	if err == nil {
		// Neuen Nutzer erstellt

		registrierungRes.Status = structs.ResponseSuccess
		registrierungRes.Error = nil

		// Generiere Cookie Token
		token := generiereSessionToken(registrierungReq.Username)

		registrierungRes.Data = structs.AuthRes{
			Message:      "Der neue Nutzeraccount wurde erstellt.",
			Sessiontoken: token,
		}

		response, _ := json.Marshal(registrierungRes)
		res.WriteHeader(http.StatusCreated) // 201 account created
		res.Write(response)
	} else {
		// Konnte keinen neuen Nutzer erstellen
		registrierungRes.Status = structs.ResponseError
		registrierungRes.Data = nil
		registrierungRes.Error = structs.Error{
			Code:    409,
			Message: err.Error(),
		}
		response, _ := json.Marshal(registrierungRes)
		res.WriteHeader(http.StatusConflict) // 409 Conflict
		res.Write(response)
	}
}

/**
Die Funktion liefert einen HTTP Response zurück, welcher den Nutzer abmeldet, sonst Fehler
*/
func DeleteAbmeldung(res http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	abmeldungReq := structs.AbmeldungReq{}
	abmeldungRes := structs.Response{}
	json.Unmarshal(s, &abmeldungReq)

	err := loescheSessionToken(abmeldungReq.Username)

	if err == nil {
		// Session Token gelöscht
		abmeldungRes.Status = structs.ResponseSuccess
		abmeldungRes.Data = structs.AbmeldeRes{Message: "Der Session Token wurde gelöscht"}
		abmeldungRes.Error = nil

		response, _ := json.Marshal(abmeldungRes)
		res.WriteHeader(http.StatusOK) // 200 Enacted and return message
		res.Write(response)
	} else {
		// Konnte nicht löschen
		abmeldungRes.Status = structs.ResponseError
		abmeldungRes.Error = structs.Error{
			Code:    409,
			Message: err.Error(),
		}
		response, _ := json.Marshal(abmeldungRes)
		res.WriteHeader(http.StatusConflict) // 409 Conflict
		res.Write(response)
	}
}
