package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/config"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/keycloak"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"strings"
)

var (
	clientID     = ""
	clientSecret = ""
	realm        = ""
)

// StartServer started den Router und mounted alle Subseiten.
func StartServer(logger *lumberjack.Logger, mode string) {
	r := chi.NewRouter()

	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(logger, "", log.LstdFlags)})

	r.Use(middleware.Logger)

	if mode == "dev" { // set values for authentication middleware
		clientID = config.DevKeycloakClientID
		clientSecret = config.DevKeycloakClientSecret
		realm = config.DevKeycloakRealm
	} else if mode == "prod" {
		clientID = config.ProdKeycloakClientID
		clientSecret = config.ProdKeycloakClientSecret
		realm = config.ProdKeycloakRealm
	} else {
		log.Fatalln("Mode not specified")
	}

	r.Group(func(r chi.Router) { // authenticated routes
		r.Use(keycloakAuthMiddleware)

		r.Get("/authRoute", welcome)

		// auswertung routes
		r.Post("/auswertung/updateLinkShare", updateLinkShare)

		// nutzerdaten routes
		r.Get("/nutzer/pruefeNutzer", pruefeNutzer)
		r.Get("/nutzer/rolle", getRolle)
		r.Delete("/nutzer", deleteNutzer)

		// umfrage routes
		r.Get("/umfrage", getUmfrage)
		r.Get("/umfrage/alleUmfragenVonNutzer", getAlleUmfragenVonNutzer)
		r.Get("/umfrage/gebaeude", getAlleGebaeude)
		r.Get("/umfrage/gebaeudeUndZaehler", getAlleGebaeudeUndZaehler)
		r.Get("/umfrage/duplicate", duplicateUmfrage)
		r.Post("/umfrage/insert", postInsertUmfrage)
		r.Post("/umfrage/update", postUpdateUmfrage)
		r.Delete("/umfrage", deleteUmfrage)

	})

	r.Group(func(r chi.Router) { // admin only, authenticated routes
		r.Use(keycloakAuthMiddleware)
		r.Use(checkAdminMiddleware)

		// db routes
		r.Post("/db/addFaktor", postAddFaktor)
		r.Post("/db/addZaehlerdaten", postAddZaehlerdaten)
		r.Post("/db/addZaehlerdatenCSV", postAddZaehlerdatenCSV)
		r.Post("/db/addStandardZaehlerdaten", postAddStandardZaehlerdaten)
		r.Post("/db/addVersorger", postAddVersorger)
		r.Post("/db/addStandardVersorger", postAddStandardVersorger)
		r.Post("/db/insertZaehler", postInsertZaehler)
		r.Post("/db/insertGebaeude", postInsertGebaeude)

		// mitarbeiterumfrage routes
		r.Get("/mitarbeiterumfrage/mitarbeiterumfrageFuerUmfrage", getMitarbeiterumfrageFuerUmfrage)

		// umfrage routes
		r.Get("/umfrage/alleUmfragen", getAlleUmfragen)
	})

	// unauthenticated routes
	r.Get("/", welcome)

	// mitarbeiterUmfrage routes
	r.Get("/mitarbeiterumfrage/exists", getUmfrageExists)
	r.Post("/mitarbeiterumfrage/insert", postInsertMitarbeiterumfrage)

	// umfrage routes
	r.Get("/umfrage/jahr", getUmfrageJahr)
	r.Get("/umfrage/sharedResults", getSharedResults)

	// special routes with separate authentication in function
	r.Get("/auswertung", getAuswertung)

	log.Println("Server Started")
	log.Fatalln(http.ListenAndServe(config.Port, r))
}

func getEmailFromToken(token string, ctx context.Context) (string, error) {
	userInfo, err := keycloak.KeycloakClient.GetUserInfo(ctx, token, realm)
	if err != nil {
		return "", err
	}

	var email string
	if userInfo.Email != nil {
		email = *userInfo.Email
	} else {
		return "", errors.New("Cannot retrieve email from token")
	}

	return email, nil
}

func getUsernameFromToken(token string, ctx context.Context) (string, error) {
	userInfo, err := keycloak.KeycloakClient.GetUserInfo(ctx, token, realm)
	if err != nil {
		return "", err
	}

	var nutzername string
	if userInfo.PreferredUsername != nil {
		nutzername = *userInfo.PreferredUsername
	} else {
		return "", errors.New("Cannot retrieve username from token")
	}

	return nutzername, nil
}

func keycloakAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		authHeader := req.Header.Get("Authorization")
		if len(authHeader) < 1 {
			res.WriteHeader(401)
			return
		}

		accessToken := strings.Split(authHeader, " ")[1]

		rptResult, err := keycloak.KeycloakClient.RetrospectToken(ctx, accessToken, clientID, clientSecret, realm)
		if err != nil {
			log.Println(err)
			res.WriteHeader(403)
			return
		}

		isTokenValid := *rptResult.Active

		if !isTokenValid {
			log.Println("Invalid Token")
			res.WriteHeader(401)
			return
		}

		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func checkAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		nutzername, err := getUsernameFromToken(strings.Split(req.Header.Get("Authorization"), " ")[1], ctx)
		if err != nil {
			log.Println(err)
			res.WriteHeader(401)
			return
		}

		nutzer, err := database.NutzerdatenFind(nutzername)
		if err != nil {
			log.Println(err)
			res.WriteHeader(401)
			return
		}
		if nutzer.Rolle != 1 {
			log.Println(structs.ErrNutzerHatKeineBerechtigung)
			res.WriteHeader(401)
			return
		}

		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// sendResponse sendet Response zurueck, bei Marshal Fehler sende 500 Code Error
// @param res Writer der den Response sendet
// @param success true falls normales Response Packet, false bei Error
// @param payload ist interface welches den data bzw. error struct enthaelt
// @param code ist der HTTP Header Code
func sendResponse(res http.ResponseWriter, success bool, payload interface{}, code int32) {
	responseBuilder := structs.Response{}
	if success {
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
	_, err = res.Write(response)
	if err != nil {
		log.Println(err)
	}
}

// errorResponse sendet eine Fehlermeldung zurueck
// @param res Writer der den Response sendet
// @param err Error, welcher die Fehlernachricht enthaelt
// @param statuscode, der http Statuscode fuer den Header
func errorResponse(res http.ResponseWriter, err error, statuscode int32) {
	sendResponse(res, false, structs.Error{
		Code:    statuscode,
		Message: err.Error(),
	}, statuscode)
}
