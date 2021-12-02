package co2computation

import (
	"errors"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"strconv"
)

var (
	ErrJahrNichtVorhanden = errors.New("getEnergieCO2Faktor: Kein CO2 Faktor für angegebens Jahr vorhanden")
	ErrFlaecheNegativ     = errors.New("gebaeudeNormalfall: Flaechenanteil ist negativ")
)

/**
Die Funktion berechnet für die gegeben Gebaeude, Flaechenanteile und Jahr die entsprechenden Emissionen hinsichtlich der
übergebenen Energie.
Ergebniseinheit: g
*/
func BerechneEnergieverbrauch(gebaeudeFlaecheDaten []structs.GebaeudeFlaecheAPI, jahr int32, idEnergieversorung int32) (float64, error) {
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
Funktion liefert den CO2 Faktor für das gegebene Jahr und Energieform zurück.
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
		return 0, ErrJahrNichtVorhanden
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
func gebaeudeNormalfall(co2Faktor int32, gebaeude structs.Gebaeude, idEnergieversorgung int32, jahr int32, flaechenanteil int32) (float64, error) {
	var gesamtverbrauch float64                  // Einheit: kWh
	var gesamtNGF float64 = gebaeude.Flaeche.NGF // Einheit: m^2
	var refGebaeude []int32

	if flaechenanteil == 0 {
		return 0, nil
	} else if flaechenanteil < 0 {
		return 0, ErrFlaecheNegativ
	}

	switch idEnergieversorgung {
	case 1: // Waerme
		refGebaeude = gebaeude.WaermeRef
	case 2: // Strom
		refGebaeude = gebaeude.StromRef
	case 3: // Kaelte
		refGebaeude = gebaeude.KaelteRef
	}

	for _, zaehlerID := range refGebaeude {
		var zaehler structs.Zaehler
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
		case 1: // Normalfall
			verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeude.Nr)
			if err != nil {
				return 0, err
			}

			gesamtverbrauch += verbrauch
			gesamtNGF += ngf
		case 2: // Spezialfall für Kaeltezaehler 3621 (und 3619)
			verbrauch, err := zaehlerSpezialfallZweiDrei(zaehler, jahr, 3619)
			if err != nil {
				return 0, err
			}

			gesamtverbrauch += verbrauch

		case 3: // Spezialfall für Kaeltezaehler 3622 (und 3620)
			verbrauch, err := zaehlerSpezialfallZweiDrei(zaehler, jahr, 3620)
			if err != nil {
				return 0, err
			}

			gesamtverbrauch += verbrauch

		default:
			return 0, errors.New("BerechneEnergieverbrauch: Spezialfall für Zaehler unbekannt")
		}
	}

	var emissionen float64
	if gesamtNGF <= 0 {
		emissionen = 0
	} else {
		emissionen = float64(co2Faktor) * gesamtverbrauch * float64(flaechenanteil) / gesamtNGF
	}

	return emissionen, nil
}

/**
Funktion stellt den Normalfall zur Bestimmung des Verbauchs und zugehöriger Gebaeudeflaeche dar.
Ergebniseinheit: kWh, m^2
*/
func zaehlerNormalfall(zaehler structs.Zaehler, jahr int32, gebaudeNr int32) (float64, float64, error) {
	var ngf float64

	if len(zaehler.GebaeudeRef) == 0 {
		return 0, 0, errors.New("zaehlerNormalfall: Zaehler " + strconv.FormatInt(int64(zaehler.PKEnergie), 10) + " hat keine Refernzen auf Gebaeude")
	}

	// addiere gespeicherten Verbrauch des Jahres auf Gesamtverbrauch auf
	var verbrauch float64 = -1
	for _, zaehlerstand := range zaehler.Zaehlerdaten {
		if int32(zaehlerstand.Zeitstempel.Year()) == jahr {
			verbrauch = zaehlerstand.Wert
		}
	}
	if verbrauch == -1 {
		return 0, 0, errors.New("zaehlerNormalfall: Kein Verbrauch für das Jahr " + strconv.FormatInt(int64(jahr), 10) + ", Zaehler: " + strconv.FormatInt(int64(zaehler.PKEnergie), 10))
	}

	switch zaehler.Einheit {
	case "MWh":
		verbrauch *= 1000
	case "kWh":
		verbrauch = verbrauch
	default:
		return 0, 0, errors.New("zaehlerNormalfall: Einheit von Zaehler " + strconv.FormatInt(int64(zaehler.PKEnergie), 10) + " unbekannt")
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

/**
Die Funktion stellt den Spezialfall 2 und 3 für die Kaeltezaehler 3621 und 3622 dar. Es ist eine abgewandelte Version
des Normalfalls und genau auf diese Zaehler zugeschnitte.
Ergebniseinheit: kWh
*/
func zaehlerSpezialfallZweiDrei(zaehler structs.Zaehler, jahr int32, andereZaehlerID int32) (float64, error) {
	var verbrauch float64 = -1 //Verbauch des Gruppenzaehlers
	for _, zaehlerstand := range zaehler.Zaehlerdaten {
		if int32(zaehlerstand.Zeitstempel.Year()) == jahr {
			verbrauch = zaehlerstand.Wert
		}
	}
	if verbrauch == -1 {
		return 0, errors.New("zaehlerSpezialfallZweiDrei: Kein Verbrauch für das Jahr " + strconv.FormatInt(int64(jahr), 10) + ", Zaehler: " + strconv.FormatInt(int64(zaehler.PKEnergie), 10))
	}

	subtraktionszaehler, err := database.KaeltezaehlerFind(andereZaehlerID)
	if err != nil {
		return 0, err
	}
	var subtraktionsverbrauch float64 = -1 // Verbauch des Zaehlers, der substrahiert werden muss
	for _, zaehlerstand := range subtraktionszaehler.Zaehlerdaten {
		if int32(zaehlerstand.Zeitstempel.Year()) == jahr {
			subtraktionsverbrauch = zaehlerstand.Wert
		}
	}
	if subtraktionsverbrauch == -1 {
		return 0, errors.New("zaehlerSpezialfallZweiDrei: Kein Verbrauch für das Jahr " + strconv.FormatInt(int64(jahr), 10) + ", Zaehler: " + strconv.FormatInt(int64(zaehler.PKEnergie), 10))
	}

	differenz := verbrauch - subtraktionsverbrauch
	if differenz > 0 { // Wert wird auf 0 gesetzt, falls er negativ ist, um Berechnungen nicht zu verfaelschen
		differenz *= 1000
	} else {
		differenz = 0
	}

	return differenz, nil
}
