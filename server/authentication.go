package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type session struct {
	Sessiontoken string
	GenTime      time.Time
}

var (
	// authMap Speichert Authentication Daten fuer Benutzer
	// key:= email -> {Sessiontoken, TTL}
	authMap = make(map[string]session)
)

func RouteAuthentication() chi.Router {
	r := chi.NewRouter()

	r.Post("/anmeldung", PostAnmeldung)
	r.Post("/registrierung", PostRegistrierung)
	r.Post("/pruefeSession", PostPruefeSession)
	r.Post("/pruefeNutzerRolle", PostPruefeNutzerRolle)
	r.Delete("/abmeldung", DeleteAbmeldung)

	return r
}

/**
Generiert einen Cookie Token, welcher den Nutzer authentifiziert und speichert ihn in Map
*/
func GeneriereSessionToken(email string) string {
	sessionToken := uuid.NewString()
	authMap[email] = session{
		Sessiontoken: sessionToken,
		GenTime:      time.Now(),
	}
	return sessionToken
}

/**
Ueberprueft ob die email einen gueltigen Sessiontoken registriert hat, falls ja nil
*/
func checkValidSessionToken(email string) error {
	// TTL ist 2 Stunden
	ttl := 2
	if (authMap[email] == session{}) {
		return structs.ErrNutzerHatKeinenSessiontoken
	}
	entry := authMap[email]
	genTimePlusTTL := entry.GenTime.Add(time.Hour * time.Duration(ttl))
	if genTimePlusTTL.Before(time.Now()) {
		return structs.ErrAbgelaufenerSessiontoken
	}
	return nil
}

/**
Loescht den Cookie Token welcher den Nutzer mit email authentifiziert
*/
func loescheSessionToken(email string) error {
	err := checkValidSessionToken(email)
	if err == nil || errors.Is(err, structs.ErrAbgelaufenerSessiontoken) {
		authMap[email] = session{}
		return nil
	}
	return err
}

/**
Authentifiziert einen Nutzer mit email und returned nil bei Erfolg, sonst error
*/
func Authenticate(email string, token string) error {
	err := checkValidSessionToken(email)
	if err != nil {
		// Kein valider Token registriert
		return err
	}
	if authMap[email].Sessiontoken != token {
		// Falscher Token fuer Nutzer
		return structs.ErrFalscherSessiontoken
	}
	return nil
}

// Returnt true zurück falls kein Fehler besteht, falls ein fehler besteht,
// wird ein StatusUnauthorized gesendet und falls zurueckgegeben
func AuthWithResponse(res http.ResponseWriter, req *http.Request, email string, token string) bool {
	_, err := ioutil.ReadAll(req.Body)
	if err != nil {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, http.StatusInternalServerError)
	}

	//Authentication
	errAuth := Authenticate(email, token)
	//Falls kein valider Session Token vorhanden.
	if errAuth != nil {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusUnauthorized,
			Message: errAuth.Error(),
		}, http.StatusUnauthorized)
		return false
	}
	return true
}

func PostPruefeSession(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	sessionReq := structs.PruefeSessionReq{}
	err = json.Unmarshal(s, &sessionReq)

	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	//Falls kein valider Session Token vorhanden.
	err = Authenticate(sessionReq.Username, sessionReq.Sessiontoken)
	if err != nil {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		}, http.StatusUnauthorized)
		return
	} else {
		//Falls ein valider Session Token vorhanden ist
		sendResponse(res, true, nil, 200)
		return
	}
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

	if errors.Is(err, mongo.ErrNoDocuments) {
		// Es existiert kein Account mit dieser Email
		// Sende genauere Fehlermeldung zurueck, statt ErrNoDocuments
		sendResponse(res, false, structs.Error{
			Code:    http.StatusUnauthorized,
			Message: structs.ErrFalschesPasswortError.Error(),
		}, http.StatusUnauthorized)
		return
	} else if err != nil {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	// Vergleiche Passwort mit gespeichertem Hash
	evaluationError := bcrypt.CompareHashAndPassword([]byte(nutzerdaten.Passwort), []byte(anmeldeReq.Passwort))

	if evaluationError != nil {
		// Falsches Passwort
		sendResponse(res, false, structs.Error{
			Code:    http.StatusUnauthorized,
			Message: structs.ErrFalschesPasswortError.Error(),
		}, http.StatusUnauthorized)
		return
	}

	// Korrektes Passwort authentifiziere den Nutzer
	// Generiere Cookie Token
	token := GeneriereSessionToken(anmeldeReq.Username)
	sendResponse(res, true, structs.AuthRes{
		Message:      "Nutzer authentifiziert",
		Sessiontoken: token,
	}, http.StatusOK)
}

/**
Die Funktion liefert einen HTTP Response zurueck, welcher den neuen Nutzer authentifiziert, oder eine Fehlermeldung liefert
*/
func PostRegistrierung(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	registrierungReq := structs.AuthReq{}

	err = json.Unmarshal(s, &registrierungReq)

	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	restorepath, err := database.CreateDump("PostRegistrierung")
	if err != nil {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, http.StatusInternalServerError)
	}

	err = database.NutzerdatenInsert(registrierungReq)
	if err != nil {
		// Konnte keinen neuen Nutzer erstellen
		sendResponse(res, false, structs.Error{
			Code:    http.StatusConflict,
			Message: err.Error(),
		}, http.StatusConflict)
		err = database.RestoreDump(restorepath)
		if err != nil {
			// Datenbank konnte nicht wiederhergestellt werden
			log.Fatal(err)
		}
		return
	}

	// Generiere Cookie Token
	token := GeneriereSessionToken(registrierungReq.Username)
	sendResponse(res, true, structs.AuthRes{
		Message:      "Der neue Nutzeraccount wurde erstellt",
		Sessiontoken: token,
	}, http.StatusCreated)
}

func PostPruefeNutzerRolle(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	sessionReq := structs.PruefeSessionReq{}
	err = json.Unmarshal(s, &sessionReq)

	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	//Falls kein valider Session Token vorhanden.
	if Authenticate(sessionReq.Username, sessionReq.Sessiontoken) != nil {
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	} else {
		nutzer, _ := database.NutzerdatenFind(sessionReq.Username)
		// Falls ein valider Session Token vorhanden ist
		sendResponse(res, true, nutzer.Rolle, 200)
		return
	}
}

/**
Die Funktion liefert einen HTTP Response zurueck, welcher den Nutzer abmeldet, sonst Fehler
*/
func DeleteAbmeldung(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	abmeldungReq := structs.AbmeldungReq{}
	err = json.Unmarshal(s, &abmeldungReq)

	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		sendResponse(res, false, structs.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	err = loescheSessionToken(abmeldungReq.Username)

	if err != nil {
		// Konnte nicht loeschen
		sendResponse(res, false, structs.Error{
			Code:    http.StatusConflict,
			Message: err.Error(),
		}, http.StatusConflict)
		return
	}
	// session Token geloescht
	sendResponse(res, true, structs.AbmeldeRes{
		Message: "Der session Token wurde gelöscht"}, http.StatusOK)
}
