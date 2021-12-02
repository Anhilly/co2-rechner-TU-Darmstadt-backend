package co2computation

import (
	"errors"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/server"
)

/**
Die Funktion berechnet die Gesamtemissionen für den übergebenen Slice an Dienstreisen.
Ergebniseinheit: g
*/
func BerechneDienstreisen(dienstreisenDaten []server.DienstreiseElement) (float64, error) {
	var emission float64

	for _, dienstreise := range dienstreisenDaten {
		var co2Faktor int32 = -1 // zur Überprüfung, ob der CO2Faktor umgesetzt wurde

		if dienstreise.Strecke == 0 {
			continue
		}

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
				return 0, errors.New("BerechneDienstreisen: Tankart nicht vorhanden")
			}
		case 3: // Flugzeug
			for _, faktor := range medium.CO2Faktor {
				if faktor.Streckentyp == dienstreise.Streckentyp {
					co2Faktor = faktor.Wert
				}
			}
			if co2Faktor == -1 {
				return 0, errors.New("BerechneDienstreisen: Streckentyp nicht vorhanden")
			}
		default:
			return 0, errors.New("BerechneDienstreisen: ID nicht vorhanden")
		}

		if medium.Einheit == "g/Pkm" {
			emission += float64(co2Faktor * dienstreise.Strecke * 2)
		} else {
			return 0, errors.New("BerechneDienstreisen: Einheit unbekannt")
		}
	}

	return emission, nil
}

/**
Die Funktion berechnet die Gesamtemissionen auf Basis der gegeben Pendelwege und der Tage im Büro.
Ergebniseinheit: g
*/
func BerechnePendelweg(pendelwegDaten []server.PendelwegElement, tageImBuero int32) (float64, error) {
	var emissionen float64
	const arbeitstage2020 = 230 // Arbeitstage in 2020, konstant(?)

	if tageImBuero == 0 {
		return 0, nil
	}

	arbeitstage := tageImBuero / 5.0 * arbeitstage2020

	for _, weg := range pendelwegDaten {
		if weg.Strecke == 0 {
			continue
		}

		if weg.Personenanzahl < 1 {
			return 0, errors.New("BerechnePendelweg: Personenzahl ist kleiner als 1")
		}

		medium, err := database.PendelwegFind(weg.IDPendelweg)
		if err != nil {
			return 0, err
		}

		if medium.Einheit == "g/Pkm" {
			emissionen += float64(arbeitstage*2*weg.Strecke*medium.CO2Faktor) / float64(weg.Personenanzahl)
		} else {
			return 0, errors.New("BerechnePendelweg: Einheit unbekannt")
		}
	}

	return emissionen, nil
}

/**
Die Funktion berechnet Emissionen pro Jahr für den Slice an IT-Geräten.
Ergebniseinheit: g
*/
func BerechneITGeraete(itGeraeteDaten []server.ITGeraeteAnzahl) (float64, error) {
	var emissionen float64

	for _, itGeraet := range itGeraeteDaten {
		if itGeraet.Anzahl == 0 {
			continue
		}

		kategorie, err := database.ITGeraeteFind(itGeraet.IDITGeraete)
		if err != nil {
			return 0, err
		}

		if kategorie.Einheit == "g/Stueck" {
			if kategorie.IDITGerate == 8 || kategorie.IDITGerate == 10 { // Druckerpatronen und Toner
				emissionen += float64(itGeraet.Anzahl * kategorie.CO2FaktorGesamt)
			} else { // alle anderen IT Geräte
				emissionen += float64(itGeraet.Anzahl * kategorie.CO2FaktorJahr)
			}
		} else {
			return 0, errors.New("BerechneITGeraete: Einheit unbekannt")
		}
	}

	return emissionen, nil
}
