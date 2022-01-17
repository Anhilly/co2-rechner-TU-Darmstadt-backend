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
GeneriereSessionToken generiert einen Cookie Token, welcher den Nutzer authentifiziert und speichert ihn in Map.
Dabei findet keine Authentifizierung statt!
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
checkValidSessionToken ueberprueft ob die email einen gueltigen Sessiontoken registriert hat, falls ja nil, sonst error
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
loescheSessionToken loescht den Cookie Token welcher den Nutzer mit email authentifiziert
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
Authenticate authentifiziert einen Nutzer mit email und returned nil bei Erfolg, sonst error
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

/**
AuthWithResponse prueft ob fuer die uebergeben Anmeldedaten ein valider Benutzer registriert ist.
Im Fehlerfall sendet es Unauthorized mit res Writer an Frontend und returned false, sonst nichts und gibt true zurueck
*/
func AuthWithResponse(res http.ResponseWriter, email string, token string) bool {
	//Authentication
	errAuth := Authenticate(email, token)
	//Falls kein valider Session Token vorhanden.
	if errAuth != nil {
		errorResponse(res, errAuth, http.StatusUnauthorized)
		return false
	}
	return true
}

/**
PostPruefeSession prueft ob ein gueltiger Sessiontoken registriert ist und prueft diesen mit dem Request ab
*/
func PostPruefeSession(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}
	sessionReq := structs.PruefeSessionReq{}
	err = json.Unmarshal(s, &sessionReq)

	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}
	//Authentifiziere Nutzer
	if !AuthWithResponse(res, sessionReq.Username, sessionReq.Sessiontoken) {
		return
	}
	//Falls ein valider Session Token vorhanden ist
	sendResponse(res, true, nil, http.StatusOK)
}

/**
Die Funktion liefert einen Response welcher bei valider Benutzereingabe den Nutzer authentisiert, sonst Fehler
*/
func PostAnmeldung(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}
	anmeldeReq := structs.AuthReq{}
	err = json.Unmarshal(s, &anmeldeReq)

	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	nutzerdaten, err := database.NutzerdatenFind(anmeldeReq.Username)

	if errors.Is(err, mongo.ErrNoDocuments) {
		// Es existiert kein Account mit dieser Email
		// Sende genauere Fehlermeldung zurueck, statt ErrNoDocuments
		errorResponse(res, structs.ErrFalschesPasswortError, http.StatusUnauthorized)
		return
	} else if err != nil {
		errorResponse(res, err, http.StatusUnauthorized)
		return
	}

	// Vergleiche Passwort mit gespeichertem Hash
	evaluationError := bcrypt.CompareHashAndPassword([]byte(nutzerdaten.Passwort), []byte(anmeldeReq.Passwort))

	if evaluationError != nil {
		// Falsches Passwort
		errorResponse(res, structs.ErrFalschesPasswortError, http.StatusUnauthorized)
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
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	registrierungReq := structs.AuthReq{}

	err = json.Unmarshal(s, &registrierungReq)

	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}
	var path = "registrierung"
	restorepath, err := database.CreateDump(path)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.NutzerdatenInsert(registrierungReq)
	if err != nil {
		err2 := database.RestoreDump(restorepath)
		if err2 != nil {
			// Datenbank konnte nicht wiederhergestellt werden
			log.Fatal(err2)
		}
		// Konnte keinen neuen Nutzer erstellen
		errorResponse(res, err, http.StatusConflict)
		return
	}

	// Generiere Cookie Token
	token := GeneriereSessionToken(registrierungReq.Username)
	sendResponse(res, true, structs.AuthRes{
		Message:      "Der neue Nutzeraccount wurde erstellt",
		Sessiontoken: token,
	}, http.StatusCreated)
}

/**
PostPruefeNutzerRolle ueberprueft die Nutzerrolle (Admin, User) eines authentifizierten Nutzers und liefert die Kennung zurueck
*/
func PostPruefeNutzerRolle(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}
	sessionReq := structs.PruefeSessionReq{}
	err = json.Unmarshal(s, &sessionReq)

	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}
	//Falls kein valider Session Token vorhanden.
	if !AuthWithResponse(res, sessionReq.Username, sessionReq.Sessiontoken) {
		return
	}

	nutzer, err := database.NutzerdatenFind(sessionReq.Username)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
	}
	// Falls ein valider Session Token vorhanden ist
	sendResponse(res, true, nutzer.Rolle, http.StatusOK)
}

/**
Die Funktion liefert einen HTTP Response zurueck, welcher den Nutzer abmeldet, sonst Fehler
*/
func DeleteAbmeldung(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	abmeldungReq := structs.AbmeldungReq{}
	err = json.Unmarshal(s, &abmeldungReq)

	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error -> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	err = loescheSessionToken(abmeldungReq.Username)

	if err != nil {
		// Konnte nicht loeschen
		errorResponse(res, err, http.StatusConflict)
		return
	}
	// session Token geloescht
	sendResponse(res, true, structs.AbmeldeRes{
		Message: "Der Session Token wurde gelöscht"}, http.StatusOK)
}
