package server

import (
	"encoding/json"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
)

func RouteMitarbeiterUmfrage() chi.Router {
	r := chi.NewRouter()

	// POST
	r.Post("/exists", PostUmfrageExists)
	r.Post("/insertMitarbeiterUmfrage", PostMitarbeiterUmfrageInsert)

	return r
}

// returns true if the given ID exists
func PostUmfrageExists(res http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	umfrageExistsReq := structs.UmfrageID{}
	umfrageExistsRes := structs.UmfrageID{}
	json.Unmarshal(s, &umfrageExistsReq)

	var requestedUmfrageID primitive.ObjectID
	requestedUmfrageID, _ = primitive.ObjectIDFromHex(umfrageExistsReq.UmfrageID)

	var foundUmfrage structs.Umfrage
	foundUmfrage, _ = database.UmfrageFind(requestedUmfrageID)

	// return empty string if id is nil
	if foundUmfrage.ID == primitive.NilObjectID {
		umfrageExistsRes.UmfrageID = ""
	} else {
		umfrageExistsRes.UmfrageID = foundUmfrage.ID.Hex()
	}

	response, _ := json.Marshal(umfrageExistsRes)

	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func PostMitarbeiterUmfrageInsert(res http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	umfrageExistsReq := structs.InsertMitarbeiterUmfrage{}
	umfrageExistsRes := structs.UmfrageID{}
	json.Unmarshal(s, &umfrageExistsReq)

	var umfrageID primitive.ObjectID
	umfrageID, _ = database.MitarbeiterUmfrageInsert(umfrageExistsReq)

	// return empty string if id is nil
	if umfrageID == primitive.NilObjectID {
		umfrageExistsRes.UmfrageID = ""
	} else {
		umfrageExistsRes.UmfrageID = umfrageID.Hex()
	}

	response, _ := json.Marshal(umfrageExistsRes)

	res.WriteHeader(http.StatusOK)
	res.Write(response)
}
