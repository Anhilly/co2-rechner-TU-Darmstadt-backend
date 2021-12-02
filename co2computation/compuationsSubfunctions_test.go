package co2computation

import (
	"errors"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"io"
	"testing"
	"time"
)

func TestComputaionsSubfunctions(t *testing.T) {
	database.ConnectDatabase()
	defer database.DisconnectDatabase()

	t.Run("TestGetEnergieCO2Faktor", TestGetEnergieCO2Faktor)
	t.Run("TestZaehlerNormalfall", TestZaehlerNormalfall)
}

func TestTester2(t *testing.T) {
	database.ConnectDatabase()
	defer database.DisconnectDatabase()

	t.Run("TestZaehlerNormalfall", TestZaehlerNormalfall)
}

func TestGetEnergieCO2Faktor(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("getEnergieversorgung: ID = 1, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 1
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.NoErr(err)                   // Normalfall sollte keinen Error verursachen
		is.Equal(co2Faktor, int32(144)) // tritt ein Fehler auf wird 0 zurückgegeben
	})

	t.Run("getEnergieversorgung: ID = 2, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 2
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.NoErr(err)                   // Normalfall sollte keinen Error verursachen
		is.Equal(co2Faktor, int32(285)) // tritt ein Fehler auf wird 0 zurückgegeben
	})

	t.Run("getEnergieversorgung: ID = 3, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 3
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.NoErr(err)                  // Normalfall sollte keinen Error verursachen
		is.Equal(co2Faktor, int32(72)) // tritt ein Fehler auf wird 0 zurückgegeben
	})

	// Errortests
	t.Run("getEnergieversorgung: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 0
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.Equal(err, io.EOF)         // Datenbank wirft EOF fuer unbekannte IDs
		is.Equal(co2Faktor, int32(0)) // tritt ein Fehler auf wird 0 zurückgegeben
	})

	t.Run("getEnergieversorgung: Jahr = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 1
		var jahr int32 = 0

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.Equal(err, ErrJahrNichtVorhanden) // Funktion wirft ErrJahrNichtVorhanden fuer unbekanntes Jahr
		is.Equal(co2Faktor, int32(0))        // tritt ein Fehler auf wird 0 zurückgegeben
	})
}

func TestZaehlerNormalfall(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("zaehlerNormalfall: ID = 2084, Einzelzaehler (Waermezaehler)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.WaermezaehlerFind(2084)
		var jahr int32 = 2020
		var gebaeudeNr int32 = 1101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(verbrauch, 704660.0) // erwartet 704660 (Verbrauch Jahr 2020)
		is.Equal(ngf, 0.0)            // erwartet 0.0 da keine Gruppenzaehler
	})

	t.Run("zaehlerNormalfall: ID = 2253, Gruppenzaehler (Waermezaehler)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.WaermezaehlerFind(2253)
		var jahr int32 = 2020
		var gebaeudeNr int32 = 0 // keine der Referenzen (1308, 1321) von Zaehler 2253

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.NoErr(err)            // Normalfall wirft keine Errors
		is.Equal(verbrauch, 0.0) // erwartet 0.0 (Verbrauch Jahr 2020)
		is.Equal(ngf, 3096.56)   // erwartet 3096.56 da Gruppenzaehler
	})

	t.Run("zaehlerNormalfall: Umrechnung MWh in kWh (Waermezaehler)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.WaermezaehlerFind(2255)
		var jahr int32 = 2020
		var gebaeudeNr int32 = 1314

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(verbrauch, 120200.0) // erwartet 120200 (Verbrauch Jahr 2020)
		is.Equal(ngf, 0.0)            // erwartet 0.0 da keine Gruppenzaehler
	})

	t.Run("zaehlerNormalfall: kWh beleibt kWh (Kaelterzaehler)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.KaeltezaehlerFind(6108)
		var jahr int32 = 2021
		var gebaeudeNr int32 = 1220

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.NoErr(err)                // Normalfall wirft keine Errors
		is.Equal(verbrauch, 17336.0) // erwartet 17336.0 (Verbrauch Jahr 2021)
		is.Equal(ngf, 0.0)           // erwartet 0.0 da keine Gruppenzaehler
	})

	t.Run("zaehlerNormalfall: Stromzaehler", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.StromzaehlerFind(3524)
		var jahr int32 = 2020
		var gebaeudeNr int32 = 3506

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(verbrauch, 208676.2) // erwartet 208676.2 (Verbrauch Jahr 2020)
		is.Equal(ngf, 0.0)            // erwartet 0.0 da keine Gruppenzaehler
	})

	// Errortests
	// dieser Fall sollte in der realen Welt nicht auftreten, sonst ist Fehler in den Daten
	t.Run("zaehlerNormalfall: Zaehler ohne Referenz zu Gebaeude", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler := structs.Zaehler{PKEnergie: 0}
		var jahr int32 = 2020
		var gebaeudeNr int32 = 1101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.Equal(err, errors.New("zaehlerNormalfall: Zaehler 0 hat keine Refernzen auf Gebaeude")) // Funktion wirft Error
		is.Equal(verbrauch, 0.0)                                                                   // Fehlerfall liefert 0.0
		is.Equal(ngf, 0.0)                                                                         // Fehlerfall liefert 0.0
	})

	t.Run("zaehlerNormalfall: Jahr nicht vorhanden in Zaehlerdaten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.WaermezaehlerFind(2084)
		var jahr int32 = 0
		var gebaeudeNr int32 = 1101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.Equal(err, errors.New("zaehlerNormalfall: Kein Verbrauch für das Jahr 0, Zaehler: 2084")) // Funktion wirft Error
		is.Equal(verbrauch, 0.0)                                                                     // Fehlerfall liefert 0.0
		is.Equal(ngf, 0.0)                                                                           // Fehlerfall liefert 0.0
	})

	// dieser Fall sollte in der realen Welt nicht auftreten, sonst ist Fehler in den Daten
	t.Run("zaehlerNormalfall: Zaehler ohne Referenz zu Gebaeude", func(t *testing.T) {
		is := is.NewRelaxed(t)

		location, _ := time.LoadLocation("Etc/GMT")
		zaehler := structs.Zaehler{
			PKEnergie: 0,
			Zaehlerdaten: []structs.Zaehlerwerte{
				{
					Wert:        788.66,
					Zeitstempel: time.Date(2000, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "unbekannt",
			GebaeudeRef: []int32{1},
		}
		var jahr int32 = 2000
		var gebaeudeNr int32 = 1101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.Equal(err, errors.New("zaehlerNormalfall: Einheit von Zaehler 0 unbekannt")) // Funktion wirft Error
		is.Equal(verbrauch, 0.0)                                                        // Fehlerfall liefert 0.0
		is.Equal(ngf, 0.0)                                                              // Fehlerfall liefert 0.0
	})

	// dieser Fall sollte in der realen Welt nicht auftreten, sonst ist Fehler in den Daten
	t.Run("zaehlerNormalfall: Zaehler ohne Referenz zu Gebaeude", func(t *testing.T) {
		is := is.NewRelaxed(t)

		location, _ := time.LoadLocation("Etc/GMT")
		zaehler := structs.Zaehler{
			PKEnergie: 0,
			Zaehlerdaten: []structs.Zaehlerwerte{
				{
					Wert:        788.66,
					Zeitstempel: time.Date(2000, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "MWh",
			GebaeudeRef: []int32{1},
		}
		var jahr int32 = 2000
		var gebaeudeNr int32 = 1101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.Equal(err, io.EOF)    // Datenbank wirft EOF, weil Refenrenz nicht gefunden werden kann
		is.Equal(verbrauch, 0.0) // Fehlerfall liefert 0.0
		is.Equal(ngf, 0.0)       // Fehlerfall liefert 0.0
	})
}
