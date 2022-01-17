package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
)

const (
	port = ":9000"
)

func StartServer() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Diese Middleware befindet sich hier nur waehrend der Entwicklung,
	// die Verarbeitung und configuration von cors wird in Produktion von unserem Webserver uebernommen.
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, //nolint:gomnd    // Maximum value not ignored by any of major browsers
	}))

	r.Mount("/umfrage", RouteUmfrage())
	r.Mount("/mitarbeiterUmfrage", RouteMitarbeiterUmfrage())
	r.Mount("/db", RouteDB())
	r.Mount("/auth", RouteAuthentication())
	r.Mount("/auswertung", RouteAuswertung())

	log.Println("Server Started")

	log.Fatalln(http.ListenAndServe(port, r))
}

/**
sendResponse sendet Response zurueck, bei Marshal Fehler sende 500 Code Error
 @param res Writer der den Response sendet
 @param data true falls normales Response Packet, false bei Error
 @param payload ist interface welches den data bzw. error struct enthaelt
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
	_, _ = res.Write(response)
}

/**
errorResponse sendet eine Fehlermeldung zurueck
 @param res Writer der den Response sendet
 @param err Error, welcher die Fehlernachricht enthaelt
 @param statuscode, der http Statuscode fuer den Header
*/
func errorResponse(res http.ResponseWriter, err error, statuscode int32) {
	sendResponse(res, false, structs.Error{
		Code:    statuscode,
		Message: err.Error(),
	}, statuscode)
}
