package structs

import "time"

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

type Zaehlerwerte struct {
	Wert        float64   `json:"wert" bson:"wert"`
	Zeitstempel time.Time `json:"zeitstempel" bson:"zeitstempel"`
}

/*
// collection kaeltezaehler
type Kaeltezaehler struct {
	PKEnergie    int32          `json:"pkEnergie" bson:"pkEnergie"` // (index)
	ExtSystemID  string         `json:"extSystemID" bson:"extSystemID"`
	Bezeichnung  string         `json:"bezeichnung" bson:"bezeichnung"`
	Zaehlerdaten []Zaehlerwerte `json:"zaehlerdaten" bson:"zaehlerdaten"`
	Einheit      string         `json:"einheit" bson:"einheit"`
	Spezialfall  int32          `json:"spezialfall" bson:"spezialfall"`
	Revision     int32          `json:"revision" bson:"revision"`

	GebaeudeRef []int32 `json:"gebaeudeRef" bson:"gebaeudeRef"` // -> Gebaeude.nr
}

// collection stromzaehler
type Stromzaehler struct {
	PKEnergie    int32          `json:"pkEnergie" bson:"pkEnergie"` // (index)
	Bezeichnung  string         `json:"bezeichnung" bson:"bezeichnung"`
	Zaehlerdaten []Zaehlerwerte `json:"zaehlerdaten" bson:"zaehlerdaten"`
	Einheit      string         `json:"einheit" bson:"einheit"`
	Spezialfall  int32          `json:"spezialfall" bson:"spezialfall"`
	Revision     int32          `json:"revision" bson:"revision"`

	GebaeudeRef []int32 `json:"gebaeudeRef" bson:"gebaeudeRef"` // -> Gebaeude.nr
}

// collection kaeltezaehler
type Waermezaehler struct {
	PKEnergie    int32          `json:"pkEnergie" bson:"pkEnergie"` // (index)
	ExtSystemID  string         `json:"extSystemID" bson:"extSystemID"`
	Bezeichnung  string         `json:"bezeichnung" bson:"bezeichnung"`
	Zaehlerdaten []Zaehlerwerte `json:"zaehlerdaten" bson:"zaehlerdaten"`
	Einheit      string         `json:"einheit" bson:"einheit"`
	Spezialfall  int32          `json:"spezialfall" bson:"spezialfall"`
	Revision     int32          `json:"revision" bson:"revision"`

	GebaeudeRef []int32 `json:"gebaeudeRef" bson:"gebaeudeRef"` // -> Gebaeude.nr
} */

//Collection energieversorgung
type Energieversorgung struct {
	IDEnergieversorgung int32        `json:"idEnergieversorgung" bson:"idEnergieversorgung"` // (index)
	Kategorie           string       `json:"kategorie" bson:"kategorie"`
	Einheit             string       `json:"einheit" bson:"einheit"`
	Revision            int32        `json:"revision" bson:"revision"`
	CO2Faktor           []CO2Energie `json:"CO2Faktor" bson:"CO2Faktor"`
}

type CO2Energie struct {
	Wert int32 `json:"wert" bson:"wert"`
	Jahr int32 `json:"jahr" bson:"jahr"`
}

//Collection itGeraete
type ITGeraete struct {
	IDITGerate      int32  `json:"idITGeraete" bson:"idITGeraete"` // (index)
	Kategorie       string `json:"kategorie" bson:"kategorie"`
	CO2FaktorGesamt int32  `json:"CO2FaktorGesamt" bson:"CO2FaktorGesamt"`
	CO2FaktorJahr   int32  `json:"CO2FaktorJahr" bson:"CO2FaktorJahr"`
	Einheit         string `json:"einheit" bson:"einheit"`
	Revision        int32  `json:"revision" bson:"revision"`
}

//Collection dienstreisen
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

//Collection pendelweg
type Pendelweg struct {
	IDPendelweg int32  `json:"idPendelweg" bson:"idPendelweg"` // (index)
	Medium      string `json:"medium" bson:"medium"`
	CO2Faktor   int32  `json:"CO2Faktor" bson:"CO2Faktor"`
	Einheit     string `json:"einheit" bson:"einheit"`
	Revision    int32  `json:"revision" bson:"revision"`
}
