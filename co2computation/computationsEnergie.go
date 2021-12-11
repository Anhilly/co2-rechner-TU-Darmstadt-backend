package co2computation

import (
	"errors"
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"math"
)

var (
	// Fehler durch Nutzereingabe
	ErrJahrNichtVorhanden = errors.New("getEnergieCO2Faktor: Kein CO2 Faktor für angegebenes Jahr vorhanden")
	// Fehler durch Nutzereingabe
	ErrFlaecheNegativ = errors.New("gebaeudeNormalfall: Flaechenanteil ist negativ")
	// Fehler durch fehlende Behandlung eines Gebaeudespezialfalls im Code
	ErrGebaeudeSpezialfall = errors.New("BerechneEnergieverbrauch: Spezialfall fuer Gebaeude nicht abgedeckt")
	// Fehler durch fehlende Behandlung eines Zaehlerspezialfalls im Code
	ErrZaehlerSpezialfall = errors.New("gebaeudeNormalfall: Spezialfall fuer Zaehler nicht abgedeckt")
	// Fehler durch falsche Daten in Datenbank
	ErrStrGebaeuderefFehlt = "%s: Zaehler %d hat keine Referenzen auf Gebaeude"
	// Fehler durch fehlende Werte in Datenbank
	ErrStrVerbrauchFehlt = "%s: Kein Verbrauch für das Jahr %d, Zaehler: %d"
	// Fehler durch nicht behandelte Einheit oder Fehler in der Datenbank
	ErrStrEinheitUnbekannt = "%s: Einheit %s unbekannt"
)

const (
	IDEnergieversorgungWaerme int32 = 1
	IDEnergieversorgungStrom  int32 = 2
	IDEnergieversorgungKaelte int32 = 3
)

/**
Die Funktion berechnet für die gegeben Gebaeude, Flaechenanteile und Jahr die entsprechenden Emissionen hinsichtlich der
übergebenen ID fuer die entsprechende Energie (Waerme = 1, Strom = 2, Kaelte = 3).
Ergebniseinheit: g
*/
func BerechneEnergieverbrauch(gebaeudeFlaecheDaten []structs.GebaeudeFlaecheAPI, jahr int32, idEnergieversorgung int32) (float64, error) {
	var gesamtemissionen float64

	co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)
	if err != nil {
		return 0, err
	}

	// Berechnung der Emissionen für jedes aufgelistete Gebaeude
	for _, gebaeudeFlaeche := range gebaeudeFlaecheDaten {
		gebaeude, err := database.GebaeudeFind(gebaeudeFlaeche.GebaeudeNr)
		if err != nil {
			return 0, err
		}

		switch gebaeude.Spezialfall {
		case 1: // Normalfall
			emissionen, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, gebaeudeFlaeche.Flaechenanteil)
			if err != nil {
				return 0, err
			}
			gesamtemissionen += emissionen

		default:
			return 0, ErrGebaeudeSpezialfall
		}
	}

	return math.Round(gesamtemissionen*100) / 100, nil // Ergebnisrundung auf 2 Nachkommastellen
}

/**
Funktion liefert den CO2 Faktor für das gegebene Jahr und Energieform zurück.
Ergebniseinheit: g/kWh
*/
func getEnergieCO2Faktor(id int32, jahr int32) (int32, error) {
	var co2Faktor int32 = -1

	// Bestimmung CO2 Faktor für angegebenes Jahr
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
		return 0, fmt.Errorf(ErrStrEinheitUnbekannt, "getCO2FaktorEnergie", energiewerte.Einheit)
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

	switch idEnergieversorgung { // waehlt Zaehlerreferenzen entsprechend ID
	case IDEnergieversorgungWaerme: // Waerme
		refGebaeude = gebaeude.WaermeRef
	case IDEnergieversorgungStrom: // Strom
		refGebaeude = gebaeude.StromRef
	case IDEnergieversorgungKaelte: // Kaelte
		refGebaeude = gebaeude.KaelteRef
	}

	// Betrachte alle im Gebaeude referenzierten Zaehler
	for _, zaehlerID := range refGebaeude {
		zaehler, err := database.ZaehlerFind(zaehlerID, idEnergieversorgung) // holt Zaehler aus der Datenbank
		if err != nil {
			return 0, err
		}

		switch zaehler.Spezialfall { // Behandlung des Zaehlers nach Spezialfallwert
		case 1: // Normalfall
			verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeude.Nr)
			if err != nil {
				return 0, err
			}

			gesamtverbrauch += verbrauch
			gesamtNGF += ngf
		case 2: // Spezialfall für Kaeltezaehler 3621 (und 3619)
			verbrauch, err := zaehlerSpezialfall(zaehler, jahr, 3619)
			if err != nil {
				return 0, err
			}

			gesamtverbrauch += verbrauch

		case 3: // Spezialfall für Kaeltezaehler 3622 (und 3620)
			verbrauch, err := zaehlerSpezialfall(zaehler, jahr, 3620)
			if err != nil {
				return 0, err
			}

			gesamtverbrauch += verbrauch

		default:
			return 0, ErrZaehlerSpezialfall
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
Funktion stellt den Normalfall zur Bestimmung des Verbrauchs und zugehöriger Gebaeudeflaeche dar.
Ergebniseinheit: kWh, m^2
*/
func zaehlerNormalfall(zaehler structs.Zaehler, jahr int32, gebaeudeNr int32) (float64, float64, error) {
	var ngf float64

	if len(zaehler.GebaeudeRef) == 0 {
		return 0, 0, fmt.Errorf(ErrStrGebaeuderefFehlt, "zaehlerNormalfall", zaehler.PKEnergie)
	}

	// addiere gespeicherten Verbrauch des Jahres auf Gesamtverbrauch auf
	var verbrauch float64 = -1
	for _, zaehlerstand := range zaehler.Zaehlerdaten {
		if int32(zaehlerstand.Zeitstempel.Year()) == jahr {
			verbrauch = zaehlerstand.Wert
		}
	}
	if verbrauch == -1 {
		return 0, 0, fmt.Errorf(ErrStrVerbrauchFehlt, "zaehlerNormalfall", jahr, zaehler.PKEnergie)
	}

	switch zaehler.Einheit {
	case "MWh":
		verbrauch *= 1000
	case "kWh":
		// da Verbrauch schon in kWh muss nichts gemacht werden
	default:
		return 0, 0, fmt.Errorf(ErrStrEinheitUnbekannt, "zaehlerNormalfall", zaehler.Einheit)
	}

	// NGF aller referenzierten Gebaeude wird aufaddiert, um die Gesamtflaeche fuer den Verbrauch zu bekommen
	// Die Flaeche des Gebaeudes, der diesen Zaehler referenziert hat, wurde schon behandelt.
	// Dies verhindert, dass die Flaeche bei einem Gebaeude mit mehreren Zaehlern mehrfach addiert wird
	for _, refGebaeudeID := range zaehler.GebaeudeRef {
		if refGebaeudeID == gebaeudeNr {
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
des Normalfalls und genau auf diese Zaehler zugeschnitten.
Ergebniseinheit: kWh
*/
func zaehlerSpezialfall(zaehler structs.Zaehler, jahr int32, andereZaehlerID int32) (float64, error) {
	var verbrauch float64 = -1 // Verbrauch des Gruppenzaehlers
	for _, zaehlerstand := range zaehler.Zaehlerdaten {
		if int32(zaehlerstand.Zeitstempel.Year()) == jahr {
			verbrauch = zaehlerstand.Wert
		}
	}
	if verbrauch == -1 {
		return 0, fmt.Errorf(ErrStrVerbrauchFehlt, "zaehlerSpezialfall", jahr, zaehler.PKEnergie)
	}

	subtraktionszaehler, err := database.ZaehlerFind(andereZaehlerID, IDEnergieversorgungKaelte)
	if err != nil {
		return 0, err
	}
	var subtraktionsverbrauch float64 = -1 // Verbrauch des Zaehlers, der subtrahiert werden muss
	for _, zaehlerstand := range subtraktionszaehler.Zaehlerdaten {
		if int32(zaehlerstand.Zeitstempel.Year()) == jahr {
			subtraktionsverbrauch = zaehlerstand.Wert
		}
	}
	if subtraktionsverbrauch == -1 {
		return 0, fmt.Errorf(ErrStrVerbrauchFehlt, "zaehlerSpezialfall", jahr, zaehler.PKEnergie)
	}

	differenz := verbrauch - subtraktionsverbrauch
	if differenz > 0 { // Wert wird auf 0 gesetzt, falls er negativ ist, um Berechnungen nicht zu verfaelschen
		differenz *= 1000
	} else {
		differenz = 0
	}

	return differenz, nil
}
