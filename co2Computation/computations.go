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
Die Funktion berechnet für die gegeben Gebaeude, Flaechenanteile und Jahr die entsprechenden Emissionen hinsichtlich der
übergebenen Energie.
Ergebniseinheit: g
*/
func BerechneEnergieverbrauch(gebaeudeFlaecheDaten []server.GebaeudeFlaeche, jahr int32, idEnergieversorung int32) (float64, error) {
	var gesamtemissionen float64

	co2Faktor, err := getEnergieCO2Faktor(idEnergieversorung, jahr)
	if err != nil {
		return 0, err
	}

	// Berechnung für jedes aufgelistete Gebaeude
	for _, gebaeudeFlaeche := range gebaeudeFlaecheDaten {
		gebaeude, err := database.GebaeudeFind(gebaeudeFlaeche.GebaeudeNr)
		if err != nil {
			return 0, err
		}

		switch gebaeude.Spezialfall {
		case 1: // Normalfall
			emissionen, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorung, jahr, gebaeudeFlaeche.Flaechenanteil)
			if err != nil {
				return 0, err
			}
			gesamtemissionen += emissionen

		default:
			return 0, errors.New("BerechneEnergieverbrauch: Spezialfall nicht abgedeckt")
		}
	}

	return gesamtemissionen, nil
}

/**
Funktion liefert den CO2 Faktor für das gegeben Jahr und Energieform zurück.
Ergebniseinheit: g/kWh
*/
func getEnergieCO2Faktor(id int32, jahr int32) (int32, error) {
	var co2Faktor int32 = -1

	// Bestimmung CO2 Faktor für angegebens Jahr
	energiewerte, err := database.EnergieversorgungFind(id)
	if err != nil {
		return 0, err
	}
	for _, faktor := range energiewerte.CO2Faktor {
		if faktor.Jahr == jahr {
			co2Faktor = faktor.Wert
		}
	}
	if co2Faktor == -1 {
		return 0, errors.New("getEnergieCO2Faktor: Kein CO2 Faktor für angegebens Jahr vorhanden")
	}
	if energiewerte.Einheit != "g/kWh" { // Einheit muss immer g/kWh sein
		return 0, errors.New("getEnergieCO2Faktor: Einheit unbekannt")
	}

	return co2Faktor, nil
}

/**
Die Funktion bildet den Normalfall für die Emissionsberechnungen eines Gebaeudes und dem Flaechenanteil.
Ergebniseinheit: g
*/
func gebaeudeNormalfall(co2Faktor int32, gebaeude database.Gebaeude, idEnergieversorgung int32, jahr int32, flaechenanteil int32) (float64, error) {
	var gesamtverbrauch float64                  // Einheit: kWh
	var gesamtNGF float64 = gebaeude.Flaeche.NGF // Einheit: m^2
	var refGebaeude []int32

	switch idEnergieversorgung {
	case 1: // Waerme
		refGebaeude = gebaeude.WaermeRef
	case 2: // Strom
		refGebaeude = gebaeude.StromRef
	case 3: // Kaelte
		refGebaeude = gebaeude.KaelteRef
	}

	for _, zaehlerID := range refGebaeude {
		var zaehler database.Zaehler
		var err error

		switch idEnergieversorgung {
		case 1: // Waerme
			zaehler, err = database.WaermezaehlerFind(zaehlerID)
		case 2: // Strom
			zaehler, err = database.StromzaehlerFind(zaehlerID)
		case 3: // Kaelte
			zaehler, err = database.KaeltezaehlerFind(zaehlerID)
		}
		if err != nil {
			return 0, err
		}

		switch zaehler.Spezialfall {
		case 1:
			verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeude.Nr)
			if err != nil {
				return 0, err
			}

			gesamtverbrauch += verbrauch
			gesamtNGF += ngf

		default:
			return 0, errors.New("BerechneEnergieverbrauch: Spezialfall für Zaehler unbekannt")
		}
	}

	emissionen := float64(co2Faktor) * gesamtverbrauch * float64(flaechenanteil) / gesamtNGF

	return emissionen, nil
}

/**
Funktion stellt den Normalfall zur Bestimmung des Verbauchs und zugehöriger Gebaeudeflaeche dar.
Ergebniseinheit: kWh, m^2
*/
func zaehlerNormalfall(zaehler database.Zaehler, jahr int32, gebaudeNr int32) (float64, float64, error) {
	var ngf float64

	if len(zaehler.GebaeudeRef) == 0 {
		return 0, 0, errors.New("BerechneEnergieverbrauch: Zaehler " + string(zaehler.PKEnergie) + " hat keine Refernzen auf Gebaeude")
	}

	// addiere gespeicherten Verbrauch des Jahres auf Gesamtverbrauch auf
	var verbrauch float64 = -1
	for _, zaehlerstand := range zaehler.Zaehlerdaten {
		if int32(zaehlerstand.Zeitstempel.Year()) == jahr {
			verbrauch = zaehlerstand.Wert
		}
	}
	if verbrauch == -1 {
		return 0, 0, errors.New("BerechneEnergieverbrauch: Kein Verbrauch für das Jahr " + string(jahr) + ", Zaehler: " + string(zaehler.PKEnergie))
	}

	switch zaehler.Einheit {
	case "MWh":
		verbrauch = verbrauch * 1000
	case "kWh":
		verbrauch = verbrauch
	default:
		return 0, 0, errors.New("BerechneEnergieverbrauch: Einheit von Zaehler " + string(zaehler.PKEnergie) + " unbekannt")
	}

	// NGF aller referenzierten Gebaeude wird aufaddiert, um Gesamtflaeche für Verbrauch zu haben
	// fuer das oben betrachtete Gebaeude wurde die NGF schon betrachtet -> verhindert mehrfach Addition der NGF, falls Gebaeude mehrere Zaehler hat
	for _, refGebaeudeID := range zaehler.GebaeudeRef {
		if refGebaeudeID == gebaudeNr {
			continue
		}

		refGebaeude, err := database.GebaeudeFind(refGebaeudeID)
		if err != nil {
			return 0, 0, err
		}
		ngf += refGebaeude.Flaeche.NGF
	}

	return verbrauch, ngf, nil
}
