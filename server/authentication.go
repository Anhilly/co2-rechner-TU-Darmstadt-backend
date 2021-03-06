package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-password/password"
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
	// key:= username -> {Sessiontoken, TTL}
	authMap = make(map[string]session)
)

// RouteAuthentication mounted alle aufrufbaren API Endpunkte unter */auth
func RouteAuthentication() chi.Router {
	r := chi.NewRouter()

	r.Post("/anmeldung", PostAnmeldung)
	r.Post("/registrierung", PostRegistrierung)
	r.Post("/pruefeSession", PostPruefeSession)
	r.Post("/emailBestaetigung", PostEmailBestaetigung)
	r.Post("/passwortVergessen", PostPasswortVergessen)
	r.Post("/passwortAendern", PostPasswortAendern)

	r.Delete("/abmeldung", DeleteAbmeldung)

	return r
}

// GeneriereSessionToken generiert einen Cookie Token, welcher den Nutzer authentifiziert und speichert ihn in Map.
// Dabei findet keine Authentifizierung statt!
func GeneriereSessionToken(username string) string {
	sessionToken := uuid.NewString()
	authMap[username] = session{
		Sessiontoken: sessionToken,
		GenTime:      time.Now(),
	}
	return sessionToken
}

// checkValidSessionToken ueberprueft, ob der username einen gueltigen Sessiontoken registriert hat,
// falls ja nil, sonst error
func checkValidSessionToken(username string) error {
	// TTL ist 2 Stunden
	ttl := 2
	if (authMap[username] == session{}) {
		return structs.ErrNutzerHatKeinenSessiontoken
	}
	entry := authMap[username]
	genTimePlusTTL := entry.GenTime.Add(time.Hour * time.Duration(ttl))
	if genTimePlusTTL.Before(time.Now()) {
		return structs.ErrAbgelaufenerSessiontoken
	}
	return nil
}

// loescheSessionToken loescht den Cookie Token welcher den Nutzer mit username authentifiziert
func loescheSessionToken(username string) error {
	err := checkValidSessionToken(username)
	if err == nil || errors.Is(err, structs.ErrAbgelaufenerSessiontoken) {
		authMap[username] = session{}
		return nil
	}
	return err
}

// Authenticate authentifiziert einen Nutzer mit username und returned nil bei Erfolg, sonst error
func Authenticate(username string, token string) error {
	err := checkValidSessionToken(username)
	if err != nil {
		// Kein valider Token registriert
		return err
	}
	if authMap[username].Sessiontoken != token {
		// Falscher Token fuer Nutzer
		return structs.ErrFalscherSessiontoken
	}
	return nil
}

// AuthWithResponse prueft ob fuer die uebergeben Anmeldedaten ein valider Benutzer registriert ist.
// Im Fehlerfall sendet es Unauthorized mit res Writer und returned false, sonst nichts und gibt true zurueck
func AuthWithResponse(res http.ResponseWriter, username string, token string) bool {
	// Authentication
	errAuth := Authenticate(username, token)
	// Falls kein valider Session Token vorhanden.
	if errAuth != nil {
		errorResponse(res, errAuth, http.StatusUnauthorized)
		return false
	}
	return true
}

// PostEmailBestaetigung setzt die E-Mail auf bestaetigt(1)
func PostEmailBestaetigung(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.EmailBestaetigung{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	ordner, err := database.CreateDump("PostEmailBestaetigung")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.NutzerdatenUpdateMailBestaetigung(data.UserID, structs.IDEmailBestaetigt)

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

// PostPruefeSession prueft ob ein gueltiger Sessiontoken registriert ist und prueft diesen mit dem Request ab
func PostPruefeSession(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	sessionReq := structs.PruefeSessionReq{}
	err = json.Unmarshal(s, &sessionReq)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Authentifiziere Nutzer
	if !AuthWithResponse(res, sessionReq.Username, sessionReq.Sessiontoken) {
		return
	}

	nutzer, _ := database.NutzerdatenFind(sessionReq.Username)

	data := structs.PruefeSessionRes{Rolle: nutzer.Rolle, EmailBestaetigt: nutzer.EmailBestaetigt}

	// Falls ein valider Session Token vorhanden ist
	sendResponse(res, true, data, http.StatusOK)
}

func PostPasswortAendern(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	passwortReq := structs.PasswortAendernReq{}
	err = json.Unmarshal(s, &passwortReq)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Authentifiziere Nutzer
	if !AuthWithResponse(res, passwortReq.Auth.Username, passwortReq.Auth.Sessiontoken) {
		return
	}

	nutzerdaten, err := database.NutzerdatenFind(passwortReq.Auth.Username)
	if errors.Is(err, mongo.ErrNoDocuments) {
		// Es existiert kein Account mit diesem username
		// Sende genauere Fehlermeldung zurueck, statt ErrNoDocuments
		errorResponse(res, structs.ErrFalschesPasswortError, http.StatusUnauthorized)
		return
	} else if err != nil {
		errorResponse(res, err, http.StatusUnauthorized)
		return
	}

	// Vergleiche Passwort mit gespeichertem Hash
	evaluationError := bcrypt.CompareHashAndPassword([]byte(nutzerdaten.Passwort), []byte(passwortReq.Passwort))
	if evaluationError != nil {
		// Falsches Passwort
		errorResponse(res, structs.ErrFalschesPasswortError, http.StatusUnauthorized)
		return
	}

	passworthash, err := bcrypt.GenerateFromPassword([]byte(passwortReq.NeuesPasswort), bcrypt.DefaultCost)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	nutzerdaten.Passwort = string(passworthash)

	ordner, err := database.CreateDump("PostPasswortAendern")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.NutzerdatenUpdate(nutzerdaten)

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

// PostAnmeldung liefert eine Response welcher bei valider Benutzereingabe den Nutzer authentisiert, sonst Fehler
func PostAnmeldung(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	anmeldeReq := structs.AuthReq{}
	err = json.Unmarshal(s, &anmeldeReq)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	nutzerdaten, err := database.NutzerdatenFind(anmeldeReq.Username)
	if errors.Is(err, mongo.ErrNoDocuments) {
		// Es existiert kein Account mit diesem username
		// Sende genauere Fehlermeldung zurueck, statt ErrNoDocuments
		errorResponse(res, structs.ErrFalschesPasswortError, http.StatusUnauthorized)
		return
	} else if err != nil {
		errorResponse(res, err, http.StatusUnauthorized)
		return
	}

	if nutzerdaten.EmailBestaetigt == structs.IDEmailNichtBestaetigt {
		errorResponse(res, structs.ErrNutzerUnbestaetigteMail, http.StatusUnauthorized)
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
		Rolle:        nutzerdaten.Rolle,
	}, http.StatusOK)
}

func PostPasswortVergessen(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	pwVergessen := structs.PasswortVergessenReq{}
	err = json.Unmarshal(s, &pwVergessen)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	nutzer, err := database.NutzerdatenFind(pwVergessen.Username)
	if err != nil {
		sendResponse(res, true, nil, 200)
		return
	}

	passwort, err := password.Generate(10, 3, 0, false, false)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	passwordhash, err := bcrypt.GenerateFromPassword([]byte(passwort), bcrypt.DefaultCost)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostPasswortVergessen")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	nutzer.Passwort = string(passwordhash)

	// SendMail
	err = database.NutzerdatenUpdate(nutzer)
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

	err = SendePasswortVergessenMail(nutzer.Nutzername, passwort)

	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// Generiere Cookie Token
	sendResponse(res, true, nil, http.StatusCreated)
}

// PostRegistrierung liefert einen HTTP Response zurueck, welcher den neuen Nutzer authentifiziert,
// oder eine Fehlermeldung liefert
func PostRegistrierung(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	registrierungReq := structs.AuthReq{}
	err = json.Unmarshal(s, &registrierungReq)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	restorepath, err := database.CreateDump("PostRegistrierung")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	id, err := database.NutzerdatenInsert(registrierungReq)
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

	err = SendeBestaetigungsMail(id, registrierungReq.Username)

	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// Generiere Cookie Token
	sendResponse(res, true, nil, http.StatusCreated)
}

// DeleteAbmeldung liefert einen HTTP Response zurueck, welcher den Nutzer abmeldet, sonst Fehler
func DeleteAbmeldung(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	abmeldungReq := structs.AbmeldungReq{}
	err = json.Unmarshal(s, &abmeldungReq)
	if err != nil {
		// Konnte Body der Request nicht lesen, daher Client error --> 400
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
	sendResponse(res, true, structs.AbmeldeRes{Message: "Der Session Token wurde gel??scht"}, http.StatusOK)
}
