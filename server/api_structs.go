package server

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
/* For testing:
//ein temporärer JSON für eine monotlitische Umfrage
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
}

type UmfrageHauptverantwortlicherReq struct {
	Gebaeude          []GebaeudeFlaeche `json:"gebaeude"`
	AnzahlMitarbeiter int32             `json:"anzahlMitarbeiter"`
}

type GebaeudeFlaeche struct {
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
