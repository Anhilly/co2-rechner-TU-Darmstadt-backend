package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/config"
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

	r.Mount("/umfrage", RouteUmfrage())
	r.Mount("/mitarbeiterUmfrage", RouteMitarbeiterUmfrage())
	r.Mount("/db", RouteDB())
	r.Mount("/auth", RouteAuthentication())
	r.Mount("/auswertung", RouteAuswertung())

	r.Group(func(r chi.Router) { // authenticated routes
		r.Use(keycloakAuthMiddleware)

		// temporary route for testing
		r.Get("/hello", welcome)

		// nutzerdaten routes
		r.Get("/nutzerdaten/checkUser", CheckUser)
		r.Get("/nutzerdaten/checkRolle", CheckRolle)
		r.Delete("/nutzerdaten/deleteNutzerdaten", DeleteNutzerdaten)

		// umfrage routes
		r.Get("/umfrage/gebaeudeUndZaehler", GetAllGebaeudeUndZaehler)

	})

	// unauthenticated routes
	r.Get("/", welcome)

	log.Println("Server Started")
	log.Fatalln(http.ListenAndServe(config.Port, r))
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
