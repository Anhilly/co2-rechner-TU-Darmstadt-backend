package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
In dieser Daten sind Request und Response JSON für die API als structs aufgelistet.
*/

// For testing:
type UmfrageMitarbeiterReq struct {
	Pendelweg   []UmfragePendelweg   `json:"pendelweg"`
	TageImBuero int32                `json:"tageImBuero"`
	Dienstreise []UmfrageDienstreise `json:"dienstreise"`
	ITGeraete   []UmfrageITGeraete   `json:"itGeraete"`
}

type UmfrageMitarbeiterRes struct {
	PendelwegeEmissionen   float64 `json:"pendelwegeEmissionen"`
	DienstreisenEmissionen float64 `json:"dienstreisenEmissionen"`
	ITGeraeteEmissionen    float64 `json:"itGeraeteEmissionen"`
}

// Struct zum Abfragen aller Gebaeudedaten
type AllGebaeudeRes struct {
	Gebaeude []int32 `json:"gebaeude"`
}

// Struct zum Abfragen ob Umfrage existiert
type UmfrageID struct {
	UmfrageID string `json:"umfrageID"`
}

// Structs für Request JSONs zum Hinzufuegen und Aendern von Daten der Datenbank
type AddCO2Faktor struct {
	IDEnergieversorgung int32 `json:"idEnergieversorgung"`
	Jahr                int32 `json:"jahr"`
	Wert                int32 `json:"wert"`
}

type AddZaehlerdaten struct {
	PKEnergie           int32   `json:"pkEnergie"`
	IDEnergieversorgung int32   `json:"idEnergieversorgung"`
	Jahr                int32   `json:"jahr"`
	Wert                float64 `json:"wert"`
}

type InsertZaehler struct {
	PKEnergie           int32   `json:"pkEnergie"`
	IDEnergieversorgung int32   `json:"idEnergieversorgung"`
	Bezeichnung         string  `json:"bezeichnung"`
	Einheit             string  `json:"einheit"`
	GebaeudeRef         []int32 `json:"gebaeudeRef"`
}

type InsertGebaeude struct {
	Nr          int32           `json:"nr"`
	Bezeichnung string          `json:"bezeichnung"`
	Flaeche     GebaeudeFlaeche `json:"flaeche"`
}

type InsertUmfrage struct {
	Mitarbeiteranzahl int32              `json:"mitarbeiteranzahl"`
	Jahr              int32              `json:"jahr"`
	Gebaeude          []UmfrageGebaeude  `json:"gebaeude"`
	ITGeraete         []UmfrageITGeraete `json:"itGeraete"`
	NutzerEmail       string             `json:"nutzerEmail"`
}

type InsertUmfrageRes struct {
	UmfrageID           string  `json:"umfrageID"`
	KaelteEmissionen    float64 `json:"kaelteEmissionen"`
	WaermeEmissionen    float64 `json:"waermeEmissionen"`
	StromEmissionen     float64 `json:"stromEmissionen"`
	ITGeraeteEmissionen float64 `json:"itGeraeteEmissionen"`
}

type InsertMitarbeiterUmfrage struct {
	Pendelweg   []UmfragePendelweg   `json:"pendelweg"`
	TageImBuero int32                `json:"tageImBuero"`
	Dienstreise []UmfrageDienstreise `json:"dienstreise"`
	ITGeraete   []UmfrageITGeraete   `json:"itGeraete"`
	IDUmfrage   primitive.ObjectID   `json:"idUmfrage"`
}

// Request fuer Umfragenauswertung
type AuswertungReq struct {
	UmfrageIDTemp primitive.ObjectID `json:"umfrageID"`
}

// Requests zur Authentifizierung und Abmeldung
type AuthReq struct { // wird fuer Anmeldung und Registrierung verwendet
	Username string `json:"username"`
	Passwort string `json:"password"`
}

type AbmeldungReq struct {
	Username string `json:"username"`
}

// Responses basieren auf generischen Response Format, in dem die spezifischen Inhalte gekapselt sind
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`  // Typisch nil, wenn Error oder nichts zu reporten
	Error  interface{} `json:"error"` // Typisch nil, wenn kein Error
}

type AbmeldeRes struct {
	Message string `json:"message"`
}

type AuthRes struct {
	Message      string `json:"message"`
	Sessiontoken string `json:"sessiontoken"`
}

type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
