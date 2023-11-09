package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
In dieser Daten sind Request und Response JSON für die API als structs aufgelistet.
*/

// Struct zum Abfragen aller Gebaeude
type AlleGebaeudeRes struct {
	Gebaeude []int32 `json:"gebaeude"`
}

// Struct zum Abfragen aller Gebaeude und eingetragenen Zählern
type AlleGebaeudeUndZaehlerRes struct {
	Gebaeude []GebaeudeNrUndZaehlerRef         `json:"gebaeude"`
	Zaehler  []ZaehlerUndZaehlerdatenVorhanden `json:"zaehler"`
}

type ZaehlerUndZaehlerdatenVorhanden struct {
	PKEnergie             int32                  `json:"pkEnergie"`
	ZaehlerdatenVorhanden []ZaehlerwertVorhanden `json:"zaehlerdatenVorhanden"`
}

type ZaehlerwertVorhanden struct {
	Jahr      int32 `json:"jahr"`
	Vorhanden bool  `json:"vorhanden"`
}

// Struct zum Abfragen ob Umfrage existiert
type UmfrageID struct {
	UmfrageID   string `json:"umfrageID"`
	Bezeichnung string `json:"bezeichnung"`
}

// Structs für Request JSONs zum Hinzufuegen und Aendern von Daten der Datenbank
type AddCO2Faktor struct {
	IDEnergieversorgung int32 `json:"idEnergieversorgung"`
	IDVertrag           int32 `json:"idVertrag"`
	Jahr                int32 `json:"jahr"`
	Wert                int32 `json:"wert"`
}

type AddZaehlerdaten struct {
	PKEnergie           int32   `json:"pkEnergie"`
	IDEnergieversorgung int32   `json:"idEnergieversorgung"`
	Jahr                int32   `json:"jahr"`
	Wert                float64 `json:"wert"`
}

type AddStandardZaehlerdaten struct {
	Jahr int32 `json:"jahr"`
}

type AddZaehlerdatenCSV struct {
	PKEnergie           []int32   `json:"pkEnergie"`
	IDEnergieversorgung []int32   `json:"idEnergieversorgung"`
	Jahr                int32     `json:"jahr"`
	Wert                []float64 `json:"wert"`
}

type InsertZaehler struct {
	PKEnergie           int32   `json:"pkEnergie"`
	IDEnergieversorgung int32   `json:"idEnergieversorgung"`
	Bezeichnung         string  `json:"bezeichnung"`
	Einheit             string  `json:"einheit"`
	GebaeudeRef         []int32 `json:"gebaeudeRef"`
}

type InsertGebaeude struct {
	Nr                   int32           `json:"nr"`
	Bezeichnung          string          `json:"bezeichnung"`
	Flaeche              GebaeudeFlaeche `json:"flaeche"`
	WaermeVersorgerJahre []int32         `json:"waerme_versorger_jahre"`
	KaelteVersorgerJahre []int32         `json:"kaelte_versorger_jahre"`
	StromVersorgerJahre  []int32         `json:"strom_versorger_jahre"`
}

type AddVersorger struct {
	Nr                  int32 `json:"nr"`
	Jahr                int32 `json:"jahr"`
	IDEnergieversorgung int32 `json:"idEnergieversorgung"`
	IDVertrag           int32 `json:"idVertrag"`
}

type AddStandardVersorger struct {
	Jahr int32 `json:"jahr"`
}

type InsertUmfrage struct {
	Bezeichnung       string             `json:"bezeichnung"`
	Mitarbeiteranzahl int32              `json:"mitarbeiteranzahl"`
	Jahr              int32              `json:"jahr"`
	Gebaeude          []UmfrageGebaeude  `json:"gebaeude"`
	ITGeraete         []UmfrageITGeraete `json:"itGeraete"`
}

type AlleUmfragen struct {
	Umfragen []Umfrage `json:"umfragen"`
}

type AlleMitarbeiterUmfragenForUmfrage struct {
	MitarbeiterUmfragen []MitarbeiterUmfrage `json:"mitarbeiterUmfragen"`
}

type UpdateUmfrage struct {
	UmfrageID         primitive.ObjectID `json:"umfrageID"`
	Bezeichnung       string             `json:"bezeichnung"`
	Mitarbeiteranzahl int32              `json:"mitarbeiteranzahl"`
	Jahr              int32              `json:"jahr"`
	Gebaeude          []UmfrageGebaeude  `json:"gebaeude"`
	ITGeraete         []UmfrageITGeraete `json:"itGeraete"`
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
}

type RequestLinkShare struct {
	UmfrageID    primitive.ObjectID `json:"umfrageID"`
	Freigabewert int32              `json:"freigabewert"`
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
	EmissionenWaerme                         float64            `json:"emissionenWaerme"`
	EmissionenStrom                          float64            `json:"emissionenStrom"`
	EmissionenKaelte                         float64            `json:"emissionenKaelte"`
	EmissionenEnergie                        float64            `json:"emissionenEnergie"`
	EmissionenITGeraeteHauptverantwortlicher float64            `json:"emissionenITGeraeteHauptverantwortlicher"`
	EmissionenITGeraeteMitarbeiter           float64            `json:"emissionenITGeraeteMitarbeiter"`
	EmissionenITGeraete                      float64            `json:"emissionenITGeraete"`
	EmissionenDienstreisen                   float64            `json:"emissionenDienstreisen"`
	EmissionenPendelwege                     float64            `json:"emissionenPendelwege"`
	EmissionenGesamt                         float64            `json:"emissionenGesamt"`
	EmissionenProMitarbeiter                 float64            `json:"emissionenProMitarbeiter"`
	EmissionenDienstreisenAufgeteilt         map[string]float64 `json:"emissionenDienstreisenAufgeteilt"`
	EmissionenPendelwegeAufgeteilt           map[int32]float64  `json:"emissionenPendelwegeAufgeteilt"`
	EmissionenITGeraeteAufgeteilt            map[int32]float64  `json:"emissionenITGeraeteAufgeteilt"`

	// Berechneter Gesamtverbrauch
	VerbrauchWearme  float64 `json:"verbrauchWaerme"`
	VerbrauchStrom   float64 `json:"verbrauchStrom"`
	VerbrauchKaelte  float64 `json:"verbrauchKaelte"`
	VerbrauchEnergie float64 `json:"verbrauchEnergie"`

	Vergleich2PersonenHaushalt float64 `json:"vergleich2PersonenHaushalt"`
	Vergleich4PersonenHaushalt float64 `json:"vergleich4PersonenHaushalt"`

	// Für Datenlücken-Visualisierung
	GebaeudeIDsUndZaehler []GebaeudeNrUndZaehlerRef         `json:"gebaeudeIDsUndZaehler"`
	Zaehler               []ZaehlerUndZaehlerdatenVorhanden `json:"zaehler"`
	UmfrageGebaeude       []UmfrageGebaeude                 `json:"umfrageGebaeude"`
}

// Responses basieren auf generischen Response Format, in dem die spezifischen Inhalte gekapselt sind
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`  // Typisch nil, wenn Error oder nichts zu reporten
	Error  interface{} `json:"error"` // Typisch nil, wenn kein Error
}

type DeleteNutzerdatenReq struct {
	Username string `json:"username"`
}

type PruefeRolleRes struct {
	Rolle int32 `json:"rolle"`
}

type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// Hilfs-Struct für die Auswertung der Umfrage
type AllePendelwege struct {
	Pendelwege  []UmfragePendelweg
	TageImBuero int32
}
