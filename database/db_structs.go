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