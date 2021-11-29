package server

/*
In dieser Daten sind Request und Response JSON für die API als structs aufgelistet.
*/

type UmfrageResponse struct {
	KaelteEmissionen       float64 `json:"kaelteEmissionen"`
	WaermeEmissionen       float64 `json:"waermeEmissionen"`
	StromEmissionen        float64 `json:"stromEmissionen"`
	ITGeraeteEmissionen    float64 `json:"itGeraeteEmissionen"`
	DienstreisenEmissionen float64 `json:"dienstreisenEmissionen"`
	PendelwegeEmissionen   float64 `json:"pendelwegeEmissionen"`
}

//ein temporärer JSON für eine monotlitische Umfrage
type UmfrageRequest struct {
	//Hauptverantwortlicher
	Gebaeude          []GebaeudeFlaeche `json:"gebaeude"`
	AnzahlMitarbeiter int32             `json:"anzahlMitarbeiter"`

	//geteilt Über Mitarbeiter und Hauptverantwortlicher
	ITGeraete []ITGeraeteAnzahl `json:"itGeraete"`

	//Mitarbeiter
	Pendelweg       []PendelwegElement   `json:"pendelweg"`
	TageImBuero     int32                `json:"tageImBuero"`
	Dienstreise     []DienstreiseElement `json:"dienstreise"`
	Papierverbrauch int32                `json:"papierverbrauch"`
}

type GebaeudeFlaeche struct {
	GebaeudeNr     string `json:"gebaeudeNr"`
	Flaechenanteil int32  `json:"flaechenanteil"`
}

type ITGeraeteAnzahl struct {
	IDITGeraete string `json:"idITGeraete"`
	Anzahl      int32  `json:"anzahl"`
}

type PendelwegElement struct {
	Strecke     int32 `json:"strecke"`
	IDPendelweg int32 `json:"idPendelweg"`
}

type DienstreiseElement struct {
	IDDienstreise int32  `json:"idDienstreise"`
	Streckentyp   string `json:"streckentyp"`
	Strecke       int32  `json:"strecke"`
	Tankart       string `json:"tankart"`
}