package structs

/*
In dieser Daten sind Request und Response JSON für die API als structs aufgelistet.
*/

type UmfrageRes struct {
	KaelteEmissionen       float64 `json:"kaelteEmissionen"`
	WaermeEmissionen       float64 `json:"waermeEmissionen"`
	StromEmissionen        float64 `json:"stromEmissionen"`
	ITGeraeteEmissionen    float64 `json:"itGeraeteEmissionen"`
	DienstreisenEmissionen float64 `json:"dienstreisenEmissionen"`
	PendelwegeEmissionen   float64 `json:"pendelwegeEmissionen"`
}

// For testing:
type UmfrageMitarbeiterRes struct {
	PendelwegeEmissionen   float64 `json:"pendelwegeEmissionen"`
	DienstreisenEmissionen float64 `json:"dienstreisenEmissionen"`
	ITGeraeteEmissionen    float64 `json:"itGeraeteEmissionen"`
}

/* For testing:
//ein temporärer JSON für eine monolithische Umfrage
type UmfrageReq struct {
	// Hauptverantwortlicher
	Gebaeude          []GebaeudeFlaeche `json:"gebaeude"`
	AnzahlMitarbeiter int32             `json:"anzahlMitarbeiter"`

	// geteilt Über Mitarbeiter und Hauptverantwortlicher
	ITGeraete []ITGeraeteAnzahl `json:"itGeraete"`

	// Mitarbeiter
	Pendelweg   []PendelwegElement   `json:"pendelweg"`
	TageImBuero int32                `json:"tageImBuero"`
	Dienstreise []DienstreiseElement `json:"dienstreise"`
}
*/
type UmfrageMitarbeiterReq struct {
	Pendelweg   []PendelwegElement   `json:"pendelweg"`
	TageImBuero int32                `json:"tageImBuero"`
	Dienstreise []DienstreiseElement `json:"dienstreise"`
	ITGeraete   []ITGeraeteAnzahl    `json:"itGeraete"`
}

type UmfrageHauptverantwortlicherReq struct {
	Gebaeude          []GebaeudeFlaecheAPI `json:"gebaeude"`
	AnzahlMitarbeiter int32                `json:"anzahlMitarbeiter"`
	ITGeraete         []ITGeraeteAnzahl    `json:"itGeraete"`
}

type UmfrageHauptverantwortlicherRes struct {
	KaelteEmissionen    float64 `json:"kaelteEmissionen"`
	WaermeEmissionen    float64 `json:"waermeEmissionen"`
	StromEmissionen     float64 `json:"stromEmissionen"`
	ITGeraeteEmissionen float64 `json:"itGeraeteEmissionen"`
}

type GebaeudeFlaecheAPI struct {
	GebaeudeNr     int32 `json:"gebaeudeNr"`
	Flaechenanteil int32 `json:"flaechenanteil"`
}

type ITGeraeteAnzahl struct {
	IDITGeraete int32 `json:"idITGeraete"`
	Anzahl      int32 `json:"anzahl"`
}

type PendelwegElement struct {
	Strecke        int32 `json:"strecke"`
	IDPendelweg    int32 `json:"idPendelweg"`
	Personenanzahl int32 `json:"personenanzahl"`
}

type DienstreiseElement struct {
	IDDienstreise int32  `json:"idDienstreise"`
	Streckentyp   string `json:"streckentyp"`
	Strecke       int32  `json:"strecke"`
	Tankart       string `json:"tankart"`
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

type AnmeldungReq struct {
	Email    string `json:"username"`
	Passwort string `json:"password"`
}

type AnmeldungRes struct {
	Message     string `json:"message"`
	Success     bool   `json:"success"`
	Cookietoken string `json:"cookietoken"`
}
