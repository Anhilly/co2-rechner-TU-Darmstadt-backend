package co2Computation

import (
	"errors"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/server"
)

func computeDienstreisen(dienstreisenData []server.DienstreiseElement) (float64, error){
	var emission float64 = 0.0

	for _, item := range dienstreisenData{
		var co2Faktor int32

		medium, err := database.DienstreisenFind(item.IDDienstreise)
		if err != nil {
			return 0, nil
		}

		switch medium.IDDienstreisen{
		case 1:	// Bahn
			co2Faktor = medium.CO2Faktor[0].Wert
		case 2: // Auto
			if item.Tankart == "Benzin" {
				co2Faktor = medium.CO2Faktor[0].Wert
			} else if item.Tankart == "Diesel" {
				co2Faktor = medium.CO2Faktor[1].Wert
			} else{
				return 0, errors.New("Tankart nicht vorhanden")
			}
		case 3: // Flugzeug
			if item.Streckentyp == "Langstrecke"{
				co2Faktor = medium.CO2Faktor[0].Wert
			} else if item.Streckentyp == "Kurzstrecke"{
				co2Faktor = medium.CO2Faktor[1].Wert
			} else{
				return 0, errors.New("Streckentyp nicht vorhanden")
			}
		default:
			return 0, errors.New("ID nicht vorhanden")
		}

		emission = emission + float64(co2Faktor * item.Strecke)
	}
	
	return emission, nil
}