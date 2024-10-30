package co2computation

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
)

// BerechneEnergieverbrauch berechnet für die gegeben Gebaeude, Flaechenanteile und Jahr die entsprechenden Emissionen
// hinsichtlich der uebergebenen ID fuer die entsprechende Energie (Waerme = 1, Strom = 2, Kaelte = 3).
// Ergebniseinheit: g
func BerechneEnergieverbrauch(gebaeudeFlaecheDaten []structs.UmfrageGebaeude, jahr int32, idEnergieversorgung int32) (float64, float64, error) {
	var gesamtemissionen float64
	var gesamtverbrauch float64

	co2Faktoren, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)
	if err != nil {
		return 0, 0, err
	}

	// Berechnung der Emissionen für jedes aufgelistete Gebaeude
	for _, gebaeudeFlaeche := range gebaeudeFlaecheDaten {
		gebaeude, err := database.GebaeudeFind(gebaeudeFlaeche.GebaeudeNr)
		if err != nil {
			return 0, 0, err
		}

		switch gebaeude.Spezialfall {
		case 1: // Normalfall
			emissionen, gvNutzflaeche, err := gebaeudeNormalfall(co2Faktoren, gebaeude, idEnergieversorgung, jahr, gebaeudeFlaeche.Nutzflaeche)
			if err != nil {
				return 0, 0, err
			}
			gesamtemissionen += emissionen
			gesamtverbrauch += gvNutzflaeche

		default:
			return 0, 0, structs.ErrGebaeudeSpezialfall
		}
	}

	return math.Round(gesamtemissionen*100) / 100, math.Round(gesamtverbrauch*100) / 100, nil // Ergebnisrundung auf 2 Nachkommastellen
}

// getEnergieCO2Faktor liefert den CO2 Faktor für das gegebene Jahr und Energieform zurück.
// Ergebniseinheit: g/kWh
func getEnergieCO2Faktor(id int32, jahr int32) (map[int32]int32, error) {
	var co2FaktorVertraege []structs.CO2FaktorVetrag = nil
	co2Faktoren := make(map[int32]int32) // idVertrag: co2Faktor

	// Bestimmung CO2 Faktor für angegebenes Jahr
	energiewerte, err := database.EnergieversorgungFind(id)
	if err != nil {
		return nil, err
	}
	for _, faktor := range energiewerte.CO2Faktor {
		if faktor.Jahr == jahr {
			co2FaktorVertraege = faktor.Vertraege
		}
	}
	if co2FaktorVertraege == nil {
		return nil, structs.ErrJahrNichtVorhanden
	}
	if energiewerte.Einheit != structs.EinheitgkWh { // Einheit muss immer g/kWh sein
		return nil, fmt.Errorf(structs.ErrStrEinheitUnbekannt, "getCO2FaktorEnergie", energiewerte.Einheit)
	}

	for _, vertrag := range co2FaktorVertraege { // Array zu Map
		co2Faktoren[vertrag.IDVertrag] = vertrag.Wert
	}

	return co2Faktoren, nil
}

// gebaeudeNormalfall bildet den Normalfall für die Emissionsberechnungen eines Gebaeudes und dem Flaechenanteil.
// Ergebniseinheit: g
func gebaeudeNormalfall(co2Faktoren map[int32]int32, gebaeude structs.Gebaeude, idEnergieversorgung int32, jahr int32, nutzflaeche int32) (float64, float64, error) {
	var gesamtverbrauch float64                  // Einheit: kWh
	var gvNutzflaeche float64                    // Einheit: kWh auf Nutzflaeche runtergerechnet
	var gesamtNGF float64 = gebaeude.Flaeche.NGF // Einheit: m^2
	var zaehlerRefs []primitive.ObjectID
	var versoger []structs.Versoger
	var idVertrag int32 = -1

	if nutzflaeche == 0 {
		return 0, 0, nil
	} else if nutzflaeche < 0 {
		return 0, 0, structs.ErrFlaecheNegativ
	}

	switch idEnergieversorgung { // waehlt Zaehlerreferenzen entsprechend ID
	case structs.IDEnergieversorgungWaerme: // Waerme
		zaehlerRefs = gebaeude.WaermeRef
		versoger = gebaeude.Waermeversorger
	case structs.IDEnergieversorgungStrom: // Strom
		zaehlerRefs = gebaeude.StromRef
		versoger = gebaeude.Stromversorger
	case structs.IDEnergieversorgungKaelte: // Kaelte
		zaehlerRefs = gebaeude.KaelteRef
		versoger = gebaeude.Kaelteversorger
	}

	// finde IDVertrag für das angegebene Jahr
	for _, vertrag := range versoger {
		if vertrag.Jahr == jahr {
			idVertrag = vertrag.IDVertrag
		}
	}
	if idVertrag == -1 {
		return 0, 0, fmt.Errorf(structs.ErrStrKeinVersorger, gebaeude.Nr, jahr)
	}

	// Betrachte alle im Gebaeude referenzierten Zaehler
	for _, zaehlerID := range zaehlerRefs {
		zaehler, err := database.ZaehlerFindOID(zaehlerID, idEnergieversorgung) // holt Zaehler aus der Datenbank
		if err != nil {
			return 0, 0, err
		}

		switch zaehler.Spezialfall { // Behandlung des Zaehlers nach Spezialfallwert
		case 1: // Normalfall
			verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeude.GebaeudeID)
			if err != nil {
				return 0, 0, err
			}

			gesamtverbrauch += verbrauch
			gesamtNGF += ngf

		case 2: // Spezialfall für Kaeltezaehler L202XXXXXaKA000XXXXXXZ50CO00001 (und L402XXXXXXKA000XXXXXXZ50CO00001)
			verbrauch, err := zaehlerSpezialfall(zaehler, jahr, "L402XXXXXXKA000XXXXXXZ50CO00001")
			if err != nil {
				return 0, 0, err
			}

			gesamtverbrauch += verbrauch

		case 3: // Spezialfall für Kaeltezaehler L204XXXXXXKA000XXXXXXZ50CO00001 (und L206XXXXXXKA000XXXXXXZ50CO00001)
			verbrauch, err := zaehlerSpezialfall(zaehler, jahr, "L206XXXXXXKA000XXXXXXZ50CO00001")
			if err != nil {
				return 0, 0, err
			}

			gesamtverbrauch += verbrauch

		default:
			return 0, 0, structs.ErrZaehlerSpezialfall
		}
	}

	var emissionen float64
	if gesamtNGF <= 0 {
		emissionen = 0
	} else {
		co2Faktor, ok := co2Faktoren[idVertrag]
		if !ok {
			return 0, 0, fmt.Errorf(structs.ErrStrKeinFaktorFuerVertrag, jahr, idEnergieversorgung, idVertrag)
		}
		emissionen = float64(co2Faktor) * gesamtverbrauch * float64(nutzflaeche) / gesamtNGF
	}

	gvNutzflaeche = gesamtverbrauch * float64(nutzflaeche) / gesamtNGF

	return emissionen, gvNutzflaeche, nil
}

// zaehlerNormalfall stellt den Normalfall zur Bestimmung des Verbrauchs und zugehöriger Gebaeudeflaeche dar.
// Ergebniseinheit: kWh, m^2
func zaehlerNormalfall(zaehler structs.Zaehler, jahr int32, gebaeudeID primitive.ObjectID) (float64, float64, error) {
	var ngf float64

	if len(zaehler.GebaeudeRef) == 0 {
		return 0, 0, fmt.Errorf(structs.ErrStrGebaeuderefFehlt, "zaehlerNormalfall", zaehler.DPName)
	}

	// addiere gespeicherten Verbrauch des Jahres auf Gesamtverbrauch auf
	var verbrauch float64 = 0
	var datenVorhanden bool = false
	for _, zaehlerstand := range zaehler.Zaehlerdaten { // vorerst aggregiert alle Werte des Jahres, wird potentiell noch geändert
		if int32(zaehlerstand.Zeitstempel.Year()) == jahr {
			datenVorhanden = true
			verbrauch += zaehlerstand.Wert
		}
	}
	if !datenVorhanden {
		return 0, 0, fmt.Errorf(structs.ErrStrVerbrauchFehlt, "zaehlerNormalfall", jahr, zaehler.DPName)
	}

	switch zaehler.Einheit {
	case structs.EinheitMWh:
		verbrauch *= 1000
	case structs.EinheitkWh:
		// da Verbrauch schon in kWh muss nichts gemacht werden
	default:
		return 0, 0, fmt.Errorf(structs.ErrStrEinheitUnbekannt, "zaehlerNormalfall", zaehler.Einheit)
	}

	// NGF aller referenzierten Gebaeude wird aufaddiert, um die Gesamtflaeche fuer den Verbrauch zu bekommen
	// Die Flaeche des Gebaeudes, der diesen Zaehler referenziert hat, wurde schon behandelt.
	// Dies verhindert, dass die Flaeche bei einem Gebaeude mit mehreren Zaehlern mehrfach addiert wird
	for _, refGebaeudeID := range zaehler.GebaeudeRef {
		if refGebaeudeID == gebaeudeID {
			continue
		}

		refGebaeude, err := database.GebaeudeFindOID(refGebaeudeID)
		if err != nil {
			return 0, 0, err
		}
		ngf += refGebaeude.Flaeche.NGF
	}

	return verbrauch, ngf, nil
}

// zaehlerSpezialfall stellt den Spezialfall 2 und 3 für die Kaeltezaehler 6691 und 3622 dar.
// Es ist eine abgewandelte Version des Normalfalls und genau auf diese Zaehler zugeschnitten.
// Ergebniseinheit: kWh
func zaehlerSpezialfall(zaehler structs.Zaehler, jahr int32, andererZaehlerDPName string) (float64, error) {
	var verbrauch float64 = 0 // Verbrauch des Gruppenzaehlers
	var datenVorhanden bool = false
	for _, zaehlerstand := range zaehler.Zaehlerdaten {
		if int32(zaehlerstand.Zeitstempel.Year()) == jahr {
			datenVorhanden = true
			verbrauch += zaehlerstand.Wert
		}
	}
	if !datenVorhanden {
		return 0, fmt.Errorf(structs.ErrStrVerbrauchFehlt, "zaehlerSpezialfall", jahr, zaehler.DPName)
	}

	subtraktionszaehler, err := database.ZaehlerFindDPName(andererZaehlerDPName, structs.IDEnergieversorgungKaelte)
	if err != nil {
		return 0, err
	}
	var subtraktionsverbrauch float64 = 0 // Verbrauch des Zaehlers, der subtrahiert werden muss
	datenVorhanden = false
	for _, zaehlerstand := range subtraktionszaehler.Zaehlerdaten {
		if int32(zaehlerstand.Zeitstempel.Year()) == jahr {
			datenVorhanden = true
			subtraktionsverbrauch += zaehlerstand.Wert
		}
	}
	if !datenVorhanden {
		return 0, fmt.Errorf(structs.ErrStrVerbrauchFehlt, "zaehlerSpezialfall", jahr, zaehler.DPName)
	}

	// Umrechnung der Einheiten in kWh
	switch zaehler.Einheit {
	case structs.EinheitMWh:
		verbrauch *= 1000
	case structs.EinheitkWh:
		// da Verbrauch schon in kWh muss nichts gemacht werden
	default:
		return 0, fmt.Errorf(structs.ErrStrEinheitUnbekannt, "zaehlerNormalfall", zaehler.Einheit)
	}

	switch subtraktionszaehler.Einheit {
	case structs.EinheitMWh:
		subtraktionsverbrauch *= 1000
	case structs.EinheitkWh:
		// da Verbrauch schon in kWh muss nichts gemacht werden
	default:
		return 0, fmt.Errorf(structs.ErrStrEinheitUnbekannt, "zaehlerNormalfall", zaehler.Einheit)
	}

	differenz := verbrauch - subtraktionsverbrauch
	if differenz < 0 { // Wert wird auf 0 gesetzt, falls er negativ ist, um Berechnungen nicht zu verfaelschen
		differenz = 0
	}

	return differenz, nil
}
