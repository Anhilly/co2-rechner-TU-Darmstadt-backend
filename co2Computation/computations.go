package co2Computation

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
//TODO: Fahrgemeinschaften implementieren
func BerechnePendelweg(pendelwegDaten []server.PendelwegElement, tageImBuero int32) (float64, error) {
	var emissionen float64
	const arbeitstage2020 = 230 // Arbeitstage in 2020, konstant(?)

	arbeitstage := tageImBuero / 5.0 * arbeitstage2020

	for _, weg := range pendelwegDaten {
		medium, err := database.PendelwegFind(weg.IDPendelweg)
		if err != nil {
			return 0, err
		}

		if medium.Einheit == "g/Pkm" {
			emissionen += float64(arbeitstage * 2 * weg.Strecke * medium.CO2Faktor)
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
		kategorie, err := database.ITGeraeteFind(itGeraet.IDITGeraete)
		if err != nil {
			return 0, err
		}

		if kategorie.Einheit == "g/Stueck" {
			if kategorie.IDITGerate == 8 || kategorie.IDITGerate == 10 { //Druckerpatronen und Toner
				emissionen += float64(itGeraet.Anzahl * kategorie.CO2FaktorGesamt)
			} else { //alle anderen IT Geräte
				emissionen += float64(itGeraet.Anzahl * kategorie.CO2FaktorJahr)
			}
		} else {
			return 0, errors.New("BerechneITGeraete: Einheit unbekannt")
		}
	}

	return emissionen, nil
}

/**

Ergebniseinheit: g
*/
func BerechneWaerme(gebaeudeFlaecheDaten []server.GebaeudeFlaeche, jahr int32) (float64, error) {
	var emissionen float64
	var co2Faktor int32 = -1

	// Bestimmung CO2 Faktor für angegebens Jahr
	energiewerte, err := database.EnergieversorgungFind(1)
	if err != nil {
		return 0, err
	}
	for _, faktor := range energiewerte.CO2Faktor {
		if faktor.Jahr == jahr {
			co2Faktor = faktor.Wert
		}
	}
	if co2Faktor == -1 {
		return 0, errors.New("BerechneWaerme: Kein CO2 Faktor für angegebens Jahr vorhanden")
	}
	if energiewerte.Einheit != "g/kWh" { // Einheit muss immer g/kWh sein
		return 0, errors.New("BerechneWaerme: Einheit unbekannt")
	}

	// Berechnung für jedes aufgelistete Gebaeude
	for _, gebaeudeFlaeche := range gebaeudeFlaecheDaten {
		gebaeude, err := database.GebaeudeFind(gebaeudeFlaeche.GebaeudeNr)
		if err != nil {
			return 0, err
		}

		switch gebaeude.Spezialfall {
		case 1: // Normalfall
			referenzen := gebaeude.WaermeRef
			var ngf float64 = gebaeude.Flaeche.NGF
			var gesamtverbrauch float64 // Einheit: kWh

			if len(referenzen) == 0 { // Gebäude hat keinen Waermezaehler -> keine Emissionen berechnebar
				continue
			}

			for _, zaehlerID := range referenzen {
				zaehler, err := database.WaermezaehlerFind(zaehlerID)
				if err != nil {
					return 0, err
				}

				if len(zaehler.GebaeudeRef) == 0 {
					return 0, errors.New("BerechneWaerme: Zaehler " + string(zaehler.PKEnergie) + " hat keine Refernzen auf Gebaeude")
				}

				// addiere gespeicherten Verbrauch des Jahres auf Gesamtverbrauch auf
				var verbrauch float64 = -1
				for _, zaehlerstand := range zaehler.Zaehlerdaten {
					if int32(zaehlerstand.Zeitstempel.Year()) == jahr {
						verbrauch = zaehlerstand.Wert
					}
				}
				if verbrauch == -1 {
					return 0, errors.New("BerechneWaerme: Kein Verbrauch für das Jahr " + string(jahr) + ", Zaehler: " + string(zaehler.PKEnergie))
				}

				switch zaehler.Einheit {
				case "MWh":
					gesamtverbrauch += verbrauch * 1000
				case "kWh":
					gesamtverbrauch += verbrauch
				default:
					return 0, errors.New("BerechneWaerme: Einheit von Zaehler " + string(zaehler.PKEnergie) + " unbekannt")
				}

				// NGF aller referenzierten Gebaeude wird aufaddiert, um Gesamtflaeche für Verbrauch zu haben
				// fuer das oben betrachtete Gebaeude wurde die NGF schon betrachtet -> verhindert mehrfach Addition der NGF, falls Gebaeude mehrere Zaehler hat
				for _, refGebaeudeID := range zaehler.GebaeudeRef {
					if refGebaeudeID == gebaeude.Nr {
						continue
					}

					refGebaeude, err := database.GebaeudeFind(refGebaeudeID)
					if err != nil {
						return 0, err
					}
					ngf += refGebaeude.Flaeche.NGF
				}
			}

			emissionen += float64(co2Faktor) * gesamtverbrauch * float64(gebaeudeFlaeche.Flaechenanteil) / ngf

		default:
			return 0, errors.New("BerechneWaerme: Spezialfall nicht abgedeckt")
		}
	}

	return emissionen, nil
}
