package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

const (
	port = ":9000"
)

// StartServer started den Router und mounted alle Subseiten.
func StartServer() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Mount("/umfrage", RouteUmfrage())
	r.Mount("/mitarbeiterUmfrage", RouteMitarbeiterUmfrage())
	r.Mount("/db", RouteDB())
	r.Mount("/auth", RouteAuthentication())
	r.Mount("/auswertung", RouteAuswertung())
	r.Mount("/nutzerdaten", RouteNutzerdaten())

	log.Println("Server Started")

	log.Fatalln(http.ListenAndServe(port, r))
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
