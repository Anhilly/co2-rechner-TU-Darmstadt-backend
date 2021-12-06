package co2computation

import (
	"errors"
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
)

var (
	// Fehler durch Nutzereingabe
	ErrPersonenzahlZuKlein = errors.New("BerechnePendelweg: Personenzahl ist kleiner als 1")
	// Fehler durch Nutzereingabe (oder Wert fehlt in Datenbank)
	ErrTankartUnbekannt = errors.New("BerechneDienstreisen: Tankart nicht vorhanden")
	// Fehler durch Nutzereingabe (oder Wert fehlt in Datenbank)
	ErrStreckentypUnbekannt = errors.New("BerechneDienstreisen: Streckentyp nicht vorhanden")
	// Fehler durch Nutzereingabe
	ErrStreckeNegativ = errors.New("Berechne_: Strecke ist negativ")
	// Fehler durch Nutzereingabe
	ErrAnzahlNegativ = errors.New("BerechneITGeraete: Anzahl an IT-Geraeten ist negativ")
	// Fehler durch fehlende Implemetierung einer Berechnung
	ErrBerechnungUnbekannt = errors.New("BerechneDienstreisen: Keine Berechnung fuer angegeben ID vorhanden")
)

/**
Die Funktion berechnet die Gesamtemissionen für den übergebenen Slice an Dienstreisen.
Ergebniseinheit: g
*/
func BerechneDienstreisen(dienstreisenDaten []structs.DienstreiseElement) (float64, error) {
	var emission float64

	for _, dienstreise := range dienstreisenDaten {
		var co2Faktor int32 = -1 // zur Überprüfung, ob der CO2Faktor umgesetzt wurde

		if dienstreise.Strecke == 0 {
			continue
		} else if dienstreise.Strecke < 0 {
			return 0, ErrStreckeNegativ
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
				return 0, ErrTankartUnbekannt
			}
		case 3: // Flugzeug
			for _, faktor := range medium.CO2Faktor {
				if faktor.Streckentyp == dienstreise.Streckentyp {
					co2Faktor = faktor.Wert
				}
			}
			if co2Faktor == -1 {
				return 0, ErrStreckentypUnbekannt
			}
		default:
			return 0, ErrBerechnungUnbekannt
		}

		if medium.Einheit == "g/Pkm" {
			emission += float64(co2Faktor * dienstreise.Strecke * 2)
		} else {
			return 0, fmt.Errorf(ErrStrEinheitUnbekannt, "BerechneDienstreisen", medium.Einheit)
		}
	}

	return emission, nil
}

/**
Die Funktion berechnet die Gesamtemissionen auf Basis der gegeben Pendelwege und der Tage im Büro.
Ergebniseinheit: g
*/
func BerechnePendelweg(pendelwegDaten []structs.PendelwegElement, tageImBuero int32) (float64, error) {
	var emissionen float64
	const arbeitstage2020 = 230 // Arbeitstage in 2020, konstant(?)

	if tageImBuero == 0 {
		return 0, nil
	}

	arbeitstage := int32(float64(tageImBuero) / 5.0 * arbeitstage2020)

	for _, weg := range pendelwegDaten {
		if weg.Strecke == 0 {
			continue
		} else if weg.Strecke < 0 {
			return 0, ErrStreckeNegativ
		}

		if weg.Personenanzahl < 1 {
			return 0, ErrPersonenzahlZuKlein
		}

		medium, err := database.PendelwegFind(weg.IDPendelweg)
		if err != nil {
			return 0, err
		}

		if medium.Einheit == "g/Pkm" {
			emissionen += float64(arbeitstage*2*weg.Strecke*medium.CO2Faktor) / float64(weg.Personenanzahl)
		} else {
			return 0, fmt.Errorf(ErrStrEinheitUnbekannt, "BerechnePendelweg", medium.Einheit)
		}
	}

	return emissionen, nil
}

/**
Die Funktion berechnet Emissionen pro Jahr für den Slice an IT-Geräten.
Ergebniseinheit: g
*/
func BerechneITGeraete(itGeraeteDaten []structs.ITGeraeteAnzahl) (float64, error) {
	var emissionen float64

	for _, itGeraet := range itGeraeteDaten {
		if itGeraet.Anzahl == 0 {
			continue
		} else if itGeraet.Anzahl < 0 {
			return 0, ErrAnzahlNegativ
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
			return 0, fmt.Errorf(ErrStrEinheitUnbekannt, "BerechneITGeraete", kategorie.Einheit)
		}
	}

	return emissionen, nil
}
