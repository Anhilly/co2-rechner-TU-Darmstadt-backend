package co2Computation

import (
	"errors"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/server"
)

/**
Die Funktion berechnet die Gesamtemissionen für den übergebenen Slice an Dienstreisen.
*/
func berechneDienstreisen(dienstreisenDaten []server.DienstreiseElement) (float64, error) {
	var emission float64

	for _, dienstreise := range dienstreisenDaten {
		var co2Faktor int32 = -1 // zur Überprüfung, ob der CO2Faktor umgesetzt wurde

		medium, err := database.DienstreisenFind(dienstreise.IDDienstreise)
		if err != nil {
			return 0, err
		}

		switch medium.IDDienstreisen { // muss explizit behandelt werden, da je nach Medium der CO2 Faktor anders bestimmt wird
		case 1: // Bahn
			co2Faktor = medium.CO2Faktor[0].Wert
		case 2: // Auto
			for _, faktor := range medium.CO2Faktor {
				if faktor.Tankart == dienstreise.Tankart {
					co2Faktor = faktor.Wert
				}
			}
			if co2Faktor == -1 {
				return 0, errors.New("Tankart nicht vorhanden")
			}
		case 3: // Flugzeug
			for _, faktor := range medium.CO2Faktor {
				if faktor.Streckentyp == dienstreise.Streckentyp {
					co2Faktor = faktor.Wert
				}
			}
			if co2Faktor == -1 {
				return 0, errors.New("Streckentyp nicht vorhanden")
			}
		default:
			return 0, errors.New("ID nicht vorhanden")
		}

		emission += float64(co2Faktor * dienstreise.Strecke * 2)
	}

	return emission, nil
}

/**
Die Funktion berechnet die Gesamtemissionen auf Basis der gegeben Pendelwege und der Tage im Büro.
*/
//TODO: Fahrgemeinschaften implementieren
func berechnePendelweg(pendelwegDaten []server.PendelwegElement, tageImBuero int32) (float64, error) {
	var emissionen float64
	const arbeitstage2020 = 230 // Arbeitstage in 2020, konstant(?)

	arbeitstage := tageImBuero / 5.0 * arbeitstage2020

	for _, weg := range pendelwegDaten {
		medium, err := database.PendelwegFind(weg.IDPendelweg)
		if err != nil {
			return 0, err
		}

		emissionen += float64(arbeitstage * 2 * weg.Strecke * medium.CO2Faktor)
	}

	return emissionen, nil
}
