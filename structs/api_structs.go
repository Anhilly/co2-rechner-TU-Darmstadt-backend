package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
In dieser Daten sind Request und Response JSON für die API als structs aufgelistet.
*/

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
	IDVertrag           int32     `json:"idVertrag"`
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
	Auth                 AuthToken       `json:"authToken"`
	Nr                   int32           `json:"nr"`
	Bezeichnung          string          `json:"bezeichnung"`
	Flaeche              GebaeudeFlaeche `json:"flaeche"`
	WaermeVersorgerJahre []int32         `json:"waerme_versorger_jahre"`
	KaelteVersorgerJahre []int32         `json:"kaelte_versorger_jahre"`
	StromVersorgerJahre  []int32         `json:"strom_versorger_jahre"`
}

type AddVersorger struct {
	Auth                AuthToken `json:"authToken"`
	Nr                  int32     `json:"nr"`
	Jahr                int32     `json:"jahr"`
	IDEnergieversorgung int32     `json:"idEnergieversorgung"`
	IDVertrag           int32     `json:"idVertrag"`
}

type AddStandardVersorger struct {
	Auth AuthToken `json:"authToken"`
	Jahr int32     `json:"jahr"`
}

type InsertUmfrage struct {
	Bezeichnung       string             `json:"bezeichnung"`
	Mitarbeiteranzahl int32              `json:"mitarbeiteranzahl"`
	Jahr              int32              `json:"jahr"`
	Gebaeude          []UmfrageGebaeude  `json:"gebaeude"`
	ITGeraete         []UmfrageITGeraete `json:"itGeraete"`
	Auth              AuthToken          `json:"authToken"`
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
	Auth        AuthToken            `json:"authToken"`
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
	UmfrageID   string `json:"umfrageID"`
	Bezeichnung string `json:"bezeichnung"`
	Complete    bool   `json:"complete"`
}

type UmfrageYearRes struct {
	Jahr int32 `json:"jahr"`
}

type UmfrageSharedResultsRes struct {
	Freigegeben int32 `json:"freigegeben"`
}

type DeleteUmfrage struct {
	UmfrageID primitive.ObjectID `json:"umfrageID"`
	Auth      AuthToken          `json:"authToken"`
}

type RequestUmfrage struct {
	UmfrageID primitive.ObjectID `json:"umfrageID"`
	Auth      AuthToken          `json:"authToken"`
}

type DuplicateUmfrage struct {
	UmfrageID primitive.ObjectID `json:"umfrageID"`
	Auth      AuthToken          `json:"authToken"`
}

type RequestLinkShare struct {
	UmfrageID    primitive.ObjectID `json:"umfrageID"`
	Freigabewert int32              `json:"freigabewert"`
	Auth         AuthToken          `json:"authToken"`
}

type RequestAuth struct {
	Auth AuthToken `json:"authToken"`
}

type AuswertungRes struct {
	// Information von Umfrage
	ID                    primitive.ObjectID `json:"id"`
	Bezeichnung           string             `json:"bezeichnung"`
	Mitarbeiteranzahl     int32              `json:"mitarbeiteranzahl"`
	Jahr                  int32              `json:"jahr"`
	Umfragenanzahl        int32              `json:"umfragenanzahl"`
	Umfragenanteil        float64            `json:"umfragenanteil"`
	AuswertungFreigegeben int32              `json:"auswertungFreigegeben"`

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

type PasswortAendernReq struct { // wird fuer Anmeldung und Registrierung verwendet
	Auth          AuthToken `json:"authToken"`
	Passwort      string    `json:"passwort"`
	NeuesPasswort string    `json:"neuesPasswort"`
}

type PasswortVergessenReq struct {
	Username string `json:"username"`
}

type EmailBestaetigung struct {
	UserID primitive.ObjectID `json:"nutzerID"`
}

type AbmeldungReq struct {
	Username string `json:"username"`
}

type PruefeSessionReq struct {
	Username     string `json:"username"`
	Sessiontoken string `json:"sessiontoken"`
}

type PruefeSessionRes struct {
	Rolle           int32 `json:"rolle"`
	EmailBestaetigt int32 `json:"emailBestaetigt"`
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

type RegistrierungRes struct {
	Message string `json:"message"`
}

type AuthRes struct {
	Message      string `json:"message"`
	Sessiontoken string `json:"sessiontoken"`
	Rolle        int32  `json:"rolle"`
}

type DeleteNutzerdatenReq struct {
	Auth     AuthToken `json:"authToken"`
	Username string    `json:"username"`
}

type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
