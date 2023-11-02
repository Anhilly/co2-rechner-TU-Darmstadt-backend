package co2computation

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
)

// BerechneDienstreisen berechnet die Gesamtemissionen für den übergebenen Slice an Dienstreisen.
// Ergebniseinheit: g
func BerechneDienstreisen(dienstreisenDaten []structs.UmfrageDienstreise) (float64, error) {
	var emissionen float64 = 0
	var emissionenGesamt float64 = 0
	var emissionenAufgeteilt = make(map[string]float64)

	// alle Daten zu Dienstreisen aus der Datenbank holen
	dienstreisenMedien, err := database.DienstreisenFindAll()
	if err != nil {
		return 0, err
	}
	var medien = make(map[int32]structs.Dienstreisen)
	for _, dienstreiseMedium := range dienstreisenMedien {
		medien[dienstreiseMedium.IDDienstreisen] = dienstreiseMedium
	}

	for _, dienstreise := range dienstreisenDaten {
		var co2Faktor int32 = -1 // zur Überprüfung, ob der CO2Faktor umgesetzt wurde

		if dienstreise.Strecke == 0 {
			continue
		} else if dienstreise.Strecke < 0 {
			return 0, structs.ErrStreckeNegativ
		}

		medium, ok := medien[dienstreise.IDDienstreise]
		if !ok {
			return 0, mongo.ErrNoDocuments
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
			emissionen = float64(co2Faktor * dienstreise.Strecke * 2) //nolint:gomnd
			emissionenGesamt += emissionen

			// Aufteilung der Emissionen nach Medium und Streckentyp
			var identifier string
			switch medium.IDDienstreisen {
			case structs.IDDienstreiseBahn:
				identifier = medium.Medium
			case structs.IDDienstreiseAuto:
				identifier = fmt.Sprintf("%s-%s", medium.Medium, dienstreise.Tankart)
			case structs.IDDienstreiseFlugzeug:
				identifier = fmt.Sprintf("%s-%s", medium.Medium, dienstreise.Streckentyp)
			}

			e, ok := emissionenAufgeteilt[identifier]
			if ok {
				emissionenAufgeteilt[identifier] = e + emissionen
			} else {
				emissionenAufgeteilt[identifier] = emissionen
			}
		} else {
			return 0, fmt.Errorf(structs.ErrStrEinheitUnbekannt, "BerechneDienstreisen", medium.Einheit)
		}
	}

	return math.Round(emissionenGesamt*100) / 100, nil
}

// BerechnePendelweg berechnet die Gesamtemissionen auf Basis der gegebenen Pendelwege und der Tage im Büro.
// Ergebniseinheit: g
func BerechnePendelweg(allePendelwege []structs.AllePendelwege) (float64, error) {
	var emissionen float64
	const arbeitstage2020 = 230 // Arbeitstage in 2020, konstant(?)

	// alle Daten zu Pendelwegen aus der Datenbank holen
	allePendelwegMedien, err := database.PendelwegFindAll()
	if err != nil {
		return 0, err
	}
	var medien = make(map[int32]structs.Pendelweg)
	for _, pendelwegMedium := range allePendelwegMedien {
		medien[pendelwegMedium.IDPendelweg] = pendelwegMedium
	}

	for _, pendelwegDaten := range allePendelwege {
		if pendelwegDaten.TageImBuero == 0 {
			return 0, nil
		}

		arbeitstage := int32(float64(pendelwegDaten.TageImBuero) / 5.0 * arbeitstage2020)

		for _, weg := range pendelwegDaten.Pendelwege {
			if weg.Strecke == 0 {
				continue
			} else if weg.Strecke < 0 {
				return 0, structs.ErrStreckeNegativ
			}

			if weg.Personenanzahl < 1 {
				return 0, structs.ErrPersonenzahlZuKlein
			}

			medium, ok := medien[weg.IDPendelweg]
			if !ok {
				return 0, mongo.ErrNoDocuments
			}

			if medium.Einheit == structs.EinheitgPkm {
				emissionen += float64(arbeitstage*2*weg.Strecke*medium.CO2Faktor) / float64(weg.Personenanzahl)
			} else {
				return 0, fmt.Errorf(structs.ErrStrEinheitUnbekannt, "BerechnePendelweg", medium.Einheit)
			}
		}
	}

	return math.Round(emissionen*100) / 100, nil
}

// BerechneITGeraete berechnet Emissionen pro Jahr für den Slice an IT-Geräten.
// Ergebniseinheit: g
func BerechneITGeraete(itGeraeteDaten []structs.UmfrageITGeraete) (float64, error) {
	var emissionen float64

	// alle Daten zu IT-Geraeten aus der Datenbank holen
	alleITGeraete, err := database.ITGeraeteFindAll()
	if err != nil {
		return 0, err
	}
	var kategorien = make(map[int32]structs.ITGeraete)
	for _, itGeraete := range alleITGeraete {
		kategorien[itGeraete.IDITGerate] = itGeraete
	}

	for _, itGeraet := range itGeraeteDaten {
		if itGeraet.Anzahl == 0 {
			continue
		} else if itGeraet.Anzahl < 0 {
			return 0, structs.ErrAnzahlNegativ
		}

		kategorie, ok := kategorien[itGeraet.IDITGeraete]
		if !ok {
			return 0, mongo.ErrNoDocuments
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

	return math.Round(emissionen*100) / 100, nil
}
