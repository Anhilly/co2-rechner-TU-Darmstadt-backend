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
Die Funktion liefert einen Response welcher bei valider Benutzereingabe den Nutzer authentisiert, sonst Fehler
*/
func PostAnmeldung(res http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	anmeldeReq := structs.AuthReq{}
	anmeldeRes := structs.Response{}
	json.Unmarshal(s, &anmeldeReq)

	nutzerdaten, err := database.NutzerdatenFind(anmeldeReq.Username)

	if err != nil {
		anmeldeRes.Status = structs.ResponseError
		anmeldeRes.Data = nil
		tmpError := structs.Error{}

		if err == io.EOF {
			// Es existiert kein Account mit dieser Email
			// Sende genauere Fehlermeldung zurück, statt EOF
			tmpError.Message = structs.ErrNichtExistenteEmail.Error()
		} else {
			tmpError.Message = err.Error()
		}
		//Schreibe Errorcode aus HTML Header in JSON
		tmpError.Code = 401
		anmeldeRes.Error = tmpError

		response, _ := json.Marshal(anmeldeRes)
		res.WriteHeader(http.StatusUnauthorized) // 401 Unauthorisierter Nutzer
		res.Write(response)
		return
	}

	// Vergleiche Passwort mit gespeichertem Hash
	evaluation := bcrypt.CompareHashAndPassword([]byte(nutzerdaten.Passwort), []byte(anmeldeReq.Passwort))

	if evaluation == nil {
		//Korrektes Passwort authentifiziere den Nutzer
		anmeldeRes.Status = structs.ResponseSuccess
		anmeldeRes.Error = nil

		// Generiere Cookie Token
		token := generiereSessionToken(anmeldeReq.Username)

		anmeldeRes.Data = structs.AuthRes{
			Sessiontoken: token,
			Message:      "Nutzer authentifiziert",
		}

		response, _ := json.Marshal(anmeldeRes)
		res.WriteHeader(http.StatusOK) // 200
		res.Write(response)
	} else {
		// Falsches Passwort
		anmeldeRes.Status = structs.ResponseError
		anmeldeRes.Data = nil
		anmeldeRes.Error = structs.Error{
			Code:    401,
			Message: structs.ErrFalschesPasswortError.Error(),
		}

		response, _ := json.Marshal(anmeldeRes)
		res.WriteHeader(http.StatusUnauthorized) // 401
		res.Write(response)
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
	println(abmeldungReq.Username, " test ", AuthMap[abmeldungReq.Username].Sessiontoken)
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
