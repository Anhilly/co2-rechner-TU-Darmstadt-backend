package co2computation

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
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
			return 0, structs.ErrStreckeNegativ
		}

		medium, err := database.DienstreisenFind(dienstreise.IDDienstreise)
		if err != nil {
			return 0, err
		}

		// muss explizit behandelt werden, da je nach Medium der CO2 Faktor anders bestimmt wird
		switch medium.IDDienstreisen {
		case structs.IDDienstreiseBahn: // Bahn
			co2Faktor = medium.CO2Faktor[0].Wert
		case structs.IDDienstreiseAuto: // Auto
			for _, faktor := range medium.CO2Faktor {
				if faktor.Tankart == dienstreise.Tankart {
					co2Faktor = faktor.Wert
				}
			}
			if co2Faktor == -1 {
				return 0, structs.ErrTankartUnbekannt
			}
		case structs.IDDienstreiseFlugzeug: // Flugzeug
			for _, faktor := range medium.CO2Faktor {
				if faktor.Streckentyp == dienstreise.Streckentyp {
					co2Faktor = faktor.Wert
				}
			}
			if co2Faktor == -1 {
				return 0, structs.ErrStreckentypUnbekannt
			}
		default:
			return 0, structs.ErrBerechnungUnbekannt
		}

		if medium.Einheit == structs.EinheitgPkm {
			emission += float64(co2Faktor * dienstreise.Strecke * 2) //nolint:gomnd
		} else {
			return 0, fmt.Errorf(structs.ErrStrEinheitUnbekannt, "BerechneDienstreisen", medium.Einheit)
		}
	}

	return emission, nil
}

/**
Die Funktion berechnet die Gesamtemissionen auf Basis der gegebenen Pendelwege und der Tage im Büro.
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
			return 0, structs.ErrStreckeNegativ
		}

		if weg.Personenanzahl < 1 {
			return 0, structs.ErrPersonenzahlZuKlein
		}

		medium, err := database.PendelwegFind(weg.IDPendelweg)
		if err != nil {
			return 0, err
		}

		if medium.Einheit == structs.EinheitgPkm {
			emissionen += float64(arbeitstage*2*weg.Strecke*medium.CO2Faktor) / float64(weg.Personenanzahl)
		} else {
			return 0, fmt.Errorf(structs.ErrStrEinheitUnbekannt, "BerechnePendelweg", medium.Einheit)
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
			return 0, structs.ErrAnzahlNegativ
		}

		kategorie, err := database.ITGeraeteFind(itGeraet.IDITGeraete)
		if err != nil {
			return 0, err
		}

		if kategorie.Einheit == structs.EinheitgStueck {
			if kategorie.IDITGerate == 8 || kategorie.IDITGerate == 10 { // Druckerpatronen und Toner
				emissionen += float64(itGeraet.Anzahl * kategorie.CO2FaktorGesamt)
			} else { // alle anderen IT Geräte
				emissionen += float64(itGeraet.Anzahl * kategorie.CO2FaktorJahr)
			}
		} else {
			return 0, fmt.Errorf(structs.ErrStrEinheitUnbekannt, "BerechneITGeraete", kategorie.Einheit)
		}
	}

	return emissionen, nil
}
