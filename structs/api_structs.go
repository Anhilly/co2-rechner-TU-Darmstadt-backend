package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
In dieser Daten sind Request und Response JSON für die API als structs aufgelistet.
*/

// For testing:
//type UmfrageMitarbeiterReq struct {
//	Pendelweg   []UmfragePendelweg   `json:"pendelweg"`
//	TageImBuero int32                `json:"tageImBuero"`
//	Dienstreise []UmfrageDienstreise `json:"dienstreise"`
//	ITGeraete   []UmfrageITGeraete   `json:"itGeraete"`
//}
//
//type UmfrageMitarbeiterRes struct {
//	PendelwegeEmissionen   float64 `json:"pendelwegeEmissionen"`
//	DienstreisenEmissionen float64 `json:"dienstreisenEmissionen"`
//	ITGeraeteEmissionen    float64 `json:"itGeraeteEmissionen"`
//}

// Struct zum Abfragen aller Gebaeudedaten
type AllGebaeudeRes struct {
	Gebaeude []int32 `json:"gebaeude"`
}

// Struct zum Abfragen ob Umfrage existiert
type UmfrageID struct {
	UmfrageID   string `json:"umfrageID"`
	Bezeichnung string `json:"bezeichnung"`
}

// Structs für Request JSONs zum Hinzufuegen und Aendern von Daten der Datenbank
type AddCO2Faktor struct {
	Auth                AuthToken `json:"authToken"`
	IDEnergieversorgung int32     `json:"idEnergieversorgung"`
	Jahr                int32     `json:"jahr"`
	Wert                int32     `json:"wert"`
}

type AddZaehlerdaten struct {
	Auth                AuthToken `json:"authToken"`
	PKEnergie           int32     `json:"pkEnergie"`
	IDEnergieversorgung int32     `json:"idEnergieversorgung"`
	Jahr                int32     `json:"jahr"`
	Wert                float64   `json:"wert"`
}

type InsertZaehler struct {
	Auth                AuthToken `json:"authToken"`
	PKEnergie           int32     `json:"pkEnergie"`
	IDEnergieversorgung int32     `json:"idEnergieversorgung"`
	Bezeichnung         string    `json:"bezeichnung"`
	Einheit             string    `json:"einheit"`
	GebaeudeRef         []int32   `json:"gebaeudeRef"`
}

type InsertGebaeude struct {
	Auth        AuthToken       `json:"authToken"`
	Nr          int32           `json:"nr"`
	Bezeichnung string          `json:"bezeichnung"`
	Flaeche     GebaeudeFlaeche `json:"flaeche"`
}

type InsertUmfrage struct {
	Bezeichnung           string             `json:"bezeichnung"`
	Mitarbeiteranzahl     int32              `json:"mitarbeiteranzahl"`
	Jahr                  int32              `json:"jahr"`
	Gebaeude              []UmfrageGebaeude  `json:"gebaeude"`
	ITGeraete             []UmfrageITGeraete `json:"itGeraete"`
	Hauptverantwortlicher AuthToken          `json:"hauptverantwortlicher"`
}

type AlleUmfragen struct {
	Umfragen []Umfrage `json:"umfragen"`
}

type AlleMitarbeiterUmfragenForUmfrage struct {
	MitarbeiterUmfragen []MitarbeiterUmfrage `json:"mitarbeiterUmfragen"`
}

// Nutzer Authentifikation Token
type AuthToken struct {
	Username     string `json:"username"`
	Sessiontoken string `json:"sessiontoken"`
}

// TODO wird das noch benötigt?
//type InsertUmfrageRes struct {
//	UmfrageID           string  `json:"umfrageID"`
//	KaelteEmissionen    float64 `json:"kaelteEmissionen"`
//	WaermeEmissionen    float64 `json:"waermeEmissionen"`
//	StromEmissionen     float64 `json:"stromEmissionen"`
//	ITGeraeteEmissionen float64 `json:"itGeraeteEmissionen"`
//}

type UpdateUmfrage struct {
	Auth              AuthToken          `json:"authToken"`
	UmfrageID         primitive.ObjectID `json:"umfrageID"`
	Bezeichnung       string             `json:"bezeichnung"`
	Mitarbeiteranzahl int32              `json:"mitarbeiteranzahl"`
	Jahr              int32              `json:"jahr"`
	Gebaeude          []UmfrageGebaeude  `json:"gebaeude"`
	ITGeraete         []UmfrageITGeraete `json:"itGeraete"`
}

type UpdateMitarbeiterUmfrage struct {
	UmfrageID   primitive.ObjectID   `json:"umfrageID"`
	Pendelweg   []UmfragePendelweg   `json:"pendelweg"`
	TageImBuero int32                `json:"tageImBuero"`
	Dienstreise []UmfrageDienstreise `json:"dienstreise"`
	ITGeraete   []UmfrageITGeraete   `json:"itGeraete"`
}

type InsertMitarbeiterUmfrage struct {
	Pendelweg   []UmfragePendelweg   `json:"pendelweg"`
	TageImBuero int32                `json:"tageImBuero"`
	Dienstreise []UmfrageDienstreise `json:"dienstreise"`
	ITGeraete   []UmfrageITGeraete   `json:"itGeraete"`
	IDUmfrage   primitive.ObjectID   `json:"idUmfrage"`
}

type UmfrageExistsRes struct {
	UmfrageID string `json:"umfrageID"`
	Complete  bool   `json:"complete"`
}

type UmfrageYearRes struct {
	Jahr int32 `json:"jahr"`
}

type DeleteUmfrage struct {
	UmfrageID             primitive.ObjectID `json:"umfrageID"`
	Hauptverantwortlicher AuthToken          `json:"hauptverantwortlicher"`
}

type RequestUmfrage struct {
	UmfrageID primitive.ObjectID `json:"umfrageID"`
	Auth      AuthToken          `json:"authToken"`
}

// Request fuer Umfragenauswertung
//type AuswertungReq struct {
//	UmfrageIDTemp primitive.ObjectID `json:"umfrageID"`
//}

type AuswertungRes struct {
	// Information von Umfrage
	ID                primitive.ObjectID `json:"id"`
	Bezeichnung       string             `json:"bezeichnung"`
	Mitarbeiteranzahl int32              `json:"mitarbeiteranzahl"`
	Jahr              int32              `json:"jahr"`
	Umfragenanzahl    int32              `json:"umfragenanzahl"`
	Umfragenanteil    float64            `json:"umfragenanteil"`

	// Berechnete Werte fuer Auswertung
	EmissionenWaerme                         float64 `json:"emissionenWaerme"`
	EmissionenStrom                          float64 `json:"emissionenStrom"`
	EmissionenKaelte                         float64 `json:"emissionenKaelte"`
	EmissionenEnergie                        float64 `json:"emissionenEnergie"`
	EmissionenITGeraeteHauptverantwortlicher float64 `json:"emissionenITGeraeteHauptverantwortlicher"`
	EmissionenITGeraeteMitarbeiter           float64 `json:"emissionenITGeraeteMitarbeiter"`
	EmissionenITGeraete                      float64 `json:"emissionenITGeraete"`
	EmissionenDienstreisen                   float64 `json:"emissionenDienstreisen"`
	EmissionenPendelwege                     float64 `json:"emissionenPendelwege"`
	EmissionenGesamt                         float64 `json:"emissionenGesamt"`
	EmissionenProMitarbeiter                 float64 `json:"emissionenProMitarbeiter"`

	Vergleich2PersonenHaushalt float64 `json:"vergleich2PersonenHaushalt"`
	Vergleich4PersonenHaushalt float64 `json:"vergleich4PersonenHaushalt"`
}

// Requests zur Authentifizierung und Abmeldung
type AuthReq struct { // wird fuer Anmeldung und Registrierung verwendet
	Username string `json:"username"`
	Passwort string `json:"password"`
}

type AbmeldungReq struct {
	Username string `json:"username"`
}

type PruefeSessionReq struct {
	Username     string `json:"username"`
	Sessiontoken string `json:"sessiontoken"`
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
