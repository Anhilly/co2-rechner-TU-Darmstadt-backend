package database

type Energieversorgung struct {
	Kategorie string `json:"kategorie" bson:"kategorie"`
	Einheit string `json:"einheit" bson:"einheit"`
	Revision int32 `json:"revision" bson:"revision"`
	CO2Faktor []CO2Energie `json:"CO2Faktor" bson:"CO2Faktor"`
}

type CO2Energie struct{
	Wert float64 `json:"wert" bson:"wert"`
	Jahr int32 `json:"jahr" bson:"jahr"`
}

// collection gebaeude
type Gebaeude struct{
	Nr int32 `json:"nr" bson:"nr"` // (index)
	Bezeichnung string `json:"bezeichnung" bson:"bezeichnung"`
	Flaeche []GebaeudeFlaeche `json:"flaeche" bson:"flaeche"`
	Einheit string `json:"einheit" bson:"einheit"`
	Revision int32 `json:"revision" bson:"revision"`

	KaelteRef []string `json:"kaelteRef" bson:"kaelteRef"` // -> Kaeltezaehler.extSystemID
	WaermeRef []string `json:"waermeRef" bson:"waermeRef"` // -> extSystemID
	StromRef []string `json:"stromRef" bson:"stromRef"` // -> pkEnergie
}

type GebaeudeFlaeche struct{
	HNF float64 `json:"hnf" bson:"hnf"`				// Hauptnutzungsfläche
	NNF float64 `json:"nnf" bson:"nnf"`				// Nebennutzungsfläche
	NGF float64 `json:"ngf" bson:"ngf"` 			// Nettogrundfläche (HNF+NNF+VF)
	FF float64 `json:"ff" bson:"ff"`				// Funktionsfläche
	VF float64 `json:"vf" bson:"vf"`				// Verkehrsfläche
	FreiF float64 `json:"freif" bson:"freif"`		// Freifläche
	GesamtF float64 `json:"gesamtf" bson:"gesamtf"`	// Gesamtfläche
}

// collection kaeltezaehler
type Kaeltezaehler struct{
	ExtSystemID string `json:"extSystemID" bson:"extSystemID"` 	// (index)
	Bezeichnung string `json:"bezeichnung" bson:"bezeichnung"`
	Zaehlerdaten []Zaehlerwerte `json:"zaehlerdaten" bson:"zaehlerdaten"`
	Einheit string `json:"einheit" bson:"einheit"`
	PKEnergie int32 `json:"pkEnergie" bson:"pkEnergie"`			// (index)
	Revision int32 `json:"revision" bson:"revision"`

	GebaeudeRef []int32 `json:"gebaeudeRef" bson:"gebaeudeRef"`	// -> Gebaeude.nr
}

type Zaehlerwerte struct {
	Wert float64 `json:"wert" bson:"wert"`
	Zeitstempel string `json:"zeitstempel" bson:"zeitstempel"`
}

// collection stromzaehler
type Stromzaehler struct{
	Bezeichnung string `json:"bezeichnung" bson:"bezeichnung"`
	Zaehlerdaten []Zaehlerwerte `json:"zaehlerdaten" bson:"zaehlerdaten"`
	Einheit string `json:"einheit" bson:"einheit"`
	PKEnergie int32 `json:"pkEnergie" bson:"pkEnergie"`		// (index)
	Revision int32 `json:"revision" bson:"revision"`

	GebaeudeRef []int32 `json:"gebaeudeRef" bson:"gebaeudeRef"`	// -> Gebaeude.nr
}

// collection kaeltezaehler
type Waermezaehler struct{
	ExtSystemID string `json:"extSystemID" bson:"extSystemID"`	// (index)
	Bezeichnung string `json:"bezeichnung" bson:"bezeichnung"`
	Zaehlerdaten []Zaehlerwerte `json:"zaehlerdaten" bson:"zaehlerdaten"`
	Einheit string `json:"einheit" bson:"einheit"`
	PKEnergie int32 `json:"pkEnergie" bson:"pkEnergie"`
	Revision int32 `json:"revision" bson:"revision"`			// (index)

	GebaeudeRef []int32 `json:"gebaeudeRef" bson:"gebaeudeRef"`	// -> Gebaeude.nr
}