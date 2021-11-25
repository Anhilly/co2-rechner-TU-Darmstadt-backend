package database

//Collection energieversorgung
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

//Collection itGeraete
type ITGeraete struct {
	Kategorie string `json:"kategorie" bson:"kategorie"`
	CO2FaktorGesamt int32 `json:"CO2FaktorGesamt" bson:"CO2FaktorGesamt"`
	CO2FaktorJahr int32 `json:"CO2FaktorJahr" bson:"CO2FaktorJahr"`
	Einheit string `json:"einheit" bson:"einheit"`
	Revision int32 `json:"revision" bson:"revision"`
}

//Collection dienstreisen
type Dienstreisen struct {
	Medium string `json:"medium" bson:"medium"`
	Einheit string `json:"einheit" bson:"einheit"`
	Revision int32 `json:"revision" bson:"revision"`
	CO2Faktor []CO2Dienstreisen `json:"CO2Faktor" bson:"CO2Faktor"`
}

type CO2Dienstreisen struct {
	Tankart string `json:"Tankart" bson:"Tankart"`
	Strecke int32 `json:"Strecke" bson:"Strecke"`
	Wert int32 `json:"Wert" bson:"Wert"`
}

//Collection pendelweg
type Pendelweg struct {
	Medium string `json:"medium" bson:"medium"`
	CO2Faktor int32 `json:"CO2Faktor" bson:"CO2Faktor"`
	Einheit int32 `json:"einheit" bson:"einheit"`
	Revision int32 `json:"revision" bson:"revision"`
}