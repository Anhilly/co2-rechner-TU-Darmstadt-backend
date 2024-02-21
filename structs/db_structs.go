package structs

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

/*
Im dieser Datei sind die Dokumententype zur Emissionsberechnung als structs aufgeführt.
*/

// collection gebaeude
type Gebaeude struct {
	Nr          int32           `json:"nr" bson:"nr"` // (index)
	Bezeichnung string          `json:"bezeichnung" bson:"bezeichnung"`
	Flaeche     GebaeudeFlaeche `json:"flaeche" bson:"flaeche"`
	Einheit     string          `json:"einheit" bson:"einheit"`
	Spezialfall int32           `json:"spezialfall" bson:"spezialfall"`
	Revision    int32           `json:"revision" bson:"revision"`

	Stromversorger  []Versoger `json:"stromversorger" bson:"stromversorger"`
	Waermeversorger []Versoger `json:"waermeversorger" bson:"waermeversorger"`
	Kaelteversorger []Versoger `json:"kaelteversorger" bson:"kaelteversorger"`

	KaelteRef []int32 `json:"kaelteRef" bson:"kaelteRef"` // -> Kaeltezaehler.pkEnergie
	WaermeRef []int32 `json:"waermeRef" bson:"waermeRef"` // -> Waermezaehler.pkEnergie
	StromRef  []int32 `json:"stromRef" bson:"stromRef"`   // -> Stromzaehler.pkEnergie
}

type GebaeudeFlaeche struct {
	HNF     float64 `json:"hnf" bson:"hnf"`         // Hauptnutzungsfläche
	NNF     float64 `json:"nnf" bson:"nnf"`         // Nebennutzungsfläche
	NGF     float64 `json:"ngf" bson:"ngf"`         // Nettogrundfläche (HNF+NNF+VF)
	FF      float64 `json:"ff" bson:"ff"`           // Funktionsfläche
	VF      float64 `json:"vf" bson:"vf"`           // Verkehrsfläche
	FreiF   float64 `json:"freif" bson:"freif"`     // Freifläche
	GesamtF float64 `json:"gesamtf" bson:"gesamtf"` // Gesamtfläche
}

type Versoger struct {
	Jahr      int32 `json:"jahr" bson:"jahr"`
	IDVertrag int32 `json:"idVertrag" bson:"idVertrag"`
}

type GebaeudeNrUndZaehlerRef struct {
	Nr        int32   `json:"nr" bson:"nr"`
	KaelteRef []int32 `json:"kaelteRef" bson:"kaelteRef"`
	WaermeRef []int32 `json:"waermeRef" bson:"waermeRef"`
	StromRef  []int32 `json:"stromRef" bson:"stromRef"`
}

// Uebertyp fuer Kaeltezaehler, Waermezaehler und Stromzaehler
type Zaehler struct {
	Zaehlertyp   string         `json:"-" bson:"-"`                 // Feld wird nur in Go benutzt
	PKEnergie    int32          `json:"pkEnergie" bson:"pkEnergie"` // (index)
	Bezeichnung  string         `json:"bezeichnung" bson:"bezeichnung"`
	Zaehlerdaten []Zaehlerwerte `json:"zaehlerdaten" bson:"zaehlerdaten"`
	Einheit      string         `json:"einheit" bson:"einheit"`
	Spezialfall  int32          `json:"spezialfall" bson:"spezialfall"`
	Revision     int32          `json:"revision" bson:"revision"`

	GebaeudeRef []int32 `json:"gebaeudeRef" bson:"gebaeudeRef"` // -> Gebaeude.nr
}

// Uebertyp fuer Kaeltezaehler, Waermezaehler und Stromzaehler
type ZaehlerUndZaehlerdaten struct {
	PKEnergie    int32          `json:"pkEnergie" bson:"pkEnergie"`
	Zaehlerdaten []Zaehlerwerte `json:"zaehlerdaten" bson:"zaehlerdaten"`
}

type Zaehlerwerte struct {
	Wert        float64   `json:"wert" bson:"wert"`
	Zeitstempel time.Time `json:"zeitstempel" bson:"zeitstempel"`
}

// Collection energieversorgung
type Energieversorgung struct {
	IDEnergieversorgung int32        `json:"idEnergieversorgung" bson:"idEnergieversorgung"` // (index)
	Kategorie           string       `json:"kategorie" bson:"kategorie"`
	Einheit             string       `json:"einheit" bson:"einheit"`
	Revision            int32        `json:"revision" bson:"revision"`
	CO2Faktor           []CO2Energie `json:"CO2Faktor" bson:"CO2Faktor"`
}

type CO2Energie struct {
	Jahr      int32             `json:"jahr" bson:"jahr"`
	Vertraege []CO2FaktorVetrag `json:"vertraege" bson:"vertraege"`
}

type CO2FaktorVetrag struct {
	Wert      int32 `json:"wert" bson:"wert"`
	IDVertrag int32 `json:"idVertrag" bson:"idVertrag"`
}

// Collection itGeraete
type ITGeraete struct {
	IDITGerate      int32  `json:"idITGeraete" bson:"idITGeraete"` // (index)
	Kategorie       string `json:"kategorie" bson:"kategorie"`
	CO2FaktorGesamt int32  `json:"CO2FaktorGesamt" bson:"CO2FaktorGesamt"`
	CO2FaktorJahr   int32  `json:"CO2FaktorJahr" bson:"CO2FaktorJahr"`
	Einheit         string `json:"einheit" bson:"einheit"`
	Revision        int32  `json:"revision" bson:"revision"`
}

// Collection dienstreisen
type Dienstreisen struct {
	IDDienstreisen int32             `json:"idDienstreisen" bson:"idDienstreisen"` // (index)
	Medium         string            `json:"medium" bson:"medium"`
	Einheit        string            `json:"einheit" bson:"einheit"`
	Revision       int32             `json:"revision" bson:"revision"`
	CO2Faktor      []CO2Dienstreisen `json:"CO2Faktor" bson:"CO2Faktor"`
}

type CO2Dienstreisen struct {
	Tankart     string `json:"tankart" bson:"tankart"`
	Streckentyp string `json:"streckentyp" bson:"streckentyp"`
	Wert        int32  `json:"wert" bson:"wert"`
}

// Collection pendelweg
type Pendelweg struct {
	IDPendelweg int32  `json:"idPendelweg" bson:"idPendelweg"` // (index)
	Medium      string `json:"medium" bson:"medium"`
	CO2Faktor   int32  `json:"CO2Faktor" bson:"CO2Faktor"`
	Einheit     string `json:"einheit" bson:"einheit"`
	Revision    int32  `json:"revision" bson:"revision"`
}

// Collection nutzerdaten
type Nutzerdaten struct {
	NutzerID   primitive.ObjectID   `json:"_id" bson:"_id"`
	Nutzername string               `json:"nutzername" bson:"nutzername"`
	EMail      string               `json:"email" bson:"email"`
	Rolle      int32                `json:"rolle" bson:"rolle"`
	Revision   int32                `json:"revision" bson:"revision"`
	UmfrageRef []primitive.ObjectID `json:"umfrageRef" bson:"umfrageRef"`
}

// Collection umfrage
type Umfrage struct {
	ID                    primitive.ObjectID   `json:"_id" bson:"_id"`
	Bezeichnung           string               `json:"bezeichnung" bson:"bezeichnung"`
	Mitarbeiteranzahl     int32                `json:"mitarbeiteranzahl" bson:"mitarbeiteranzahl"`
	Jahr                  int32                `json:"jahr" bson:"jahr"`
	Gebaeude              []UmfrageGebaeude    `json:"gebaeude" bson:"gebaeude"`
	ITGeraete             []UmfrageITGeraete   `json:"itGeraete" bson:"itGeraete"`
	AuswertungFreigegeben int32                `json:"auswertungFreigegeben" bson:"auswertungFreigegeben"`
	Revision              int32                `json:"revision" bson:"revision"`
	MitarbeiterumfrageRef []primitive.ObjectID `json:"mitarbeiterUmfrageRef" bson:"mitarbeiterUmfrageRef"`
}

type UmfrageGebaeude struct {
	GebaeudeNr  int32 `json:"gebaeudeNr" bson:"gebaeudeNr"` // -> Nr in Gebaeude
	Nutzflaeche int32 `json:"nutzflaeche" bson:"nutzflaeche"`
}

type UmfrageITGeraete struct {
	IDITGeraete int32 `json:"idITGeraete" bson:"idITGeraete"` // -> IDITGeraete in ITGereate
	Anzahl      int32 `json:"anzahl" bson:"anzahl"`
}

// Collection mitarbeiterUmfrage
type MitarbeiterUmfrage struct {
	ID          primitive.ObjectID   `json:"_id" bson:"_id"`
	Pendelweg   []UmfragePendelweg   `json:"pendelweg" bson:"pendelweg"`
	TageImBuero int32                `json:"tageImBuero" bson:"tageImBuero"`
	Dienstreise []UmfrageDienstreise `json:"dienstreise" bson:"dienstreise"`
	ITGeraete   []UmfrageITGeraete   `json:"itGereate" bson:"itGeraete"`
	Revision    int32                `json:"revision" bson:"revision"`
}

type UmfragePendelweg struct {
	IDPendelweg    int32 `json:"idPendelweg" bson:"idPendelweg"` // -> IDPendelweg in Pendelweg
	Strecke        int32 `json:"strecke" bson:"strecke"`
	Personenanzahl int32 `json:"personenanzahl" bson:"personenanzahl"`
}

type UmfrageDienstreise struct {
	IDDienstreise int32  `json:"idDienstreise" bson:"idDienstreise"` // -> IDDienstreisen in Dienstreisen
	Streckentyp   string `json:"streckentyp" bson:"streckentyp"`
	Strecke       int32  `json:"strecke" bson:"strecke"`
	Tankart       string `json:"tankart" bson:"tankart"`
	Klasse        string `json:"klasse" bson:"klasse"`
}
