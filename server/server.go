package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/co2computation"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	r.Get("/auswertung", GetAuswertung)

	r.Mount("/umfrage", RouteUmfrage())
	r.Mount("/mitarbeiterUmfrage", RouteMitarbeiterUmfrage())
	r.Mount("/db", RouteDB())
	r.Mount("/auth", RouteAuthentication())

	log.Println("Server Started")

	log.Fatalln(http.ListenAndServe(port, r))
}

func GetAuswertung(res http.ResponseWriter, req *http.Request) {
	var umfrageID primitive.ObjectID
	err := umfrageID.UnmarshalText([]byte(req.URL.Query().Get("id")))
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// hole Umfragen aus der Datenbank
	umfrage, err := database.UmfrageFind(umfrageID)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	mitarbeiterumfragen, err := database.MitarbeiterUmfrageFindMany(umfrage.MitarbeiterUmfrageRef)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// Auswertung der Daten
	var auswertung structs.AuswertungRes

	// allgemeine Information der Umfrage
	auswertung.ID = umfrage.ID
	auswertung.Jahr = umfrage.Jahr
	auswertung.Mitarbeiteranzahl = umfrage.Mitarbeiteranzahl

	// Energie	TODO: Wie damit umgehen, falls kein Energiefaktor für gegebenes Jahr
	auswertung.EmissionenWaerme, err = co2computation.BerechneEnergieverbrauch(umfrage.Gebaeude, umfrage.Jahr, structs.IDEnergieversorgungWaerme)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	auswertung.EmissionenStrom, err = co2computation.BerechneEnergieverbrauch(umfrage.Gebaeude, umfrage.Jahr, structs.IDEnergieversorgungStrom)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	auswertung.EmissionenKaelte, err = co2computation.BerechneEnergieverbrauch(umfrage.Gebaeude, umfrage.Jahr, structs.IDEnergieversorgungKaelte)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	// andere Emissionen
	emission, err := co2computation.BerechneITGeraete(umfrage.ITGeraete)
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}
	auswertung.EmissionenITGeraete += emission

	for _, mitarbeiterumfrage := range mitarbeiterumfragen {
		// IT-Geraete
		emission, err := co2computation.BerechneITGeraete(mitarbeiterumfrage.ITGeraete)
		if err != nil {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
		auswertung.EmissionenITGeraete += emission

		// Dienstreisen
		emission, err = co2computation.BerechneDienstreisen(mitarbeiterumfrage.Dienstreise)
		if err != nil {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
		auswertung.EmissionenDienstreisen += emission

		// Pendelwege
		emission, err = co2computation.BerechnePendelweg(mitarbeiterumfrage.Pendelweg, mitarbeiterumfrage.TageImBuero)
		if err != nil {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}
		auswertung.EmissionenPendelwege += emission
	}

	// Response
	response, err := json.Marshal(structs.Response{
		Status: structs.ResponseSuccess,
		Data:   auswertung,
		Error:  nil,
	})
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
}
