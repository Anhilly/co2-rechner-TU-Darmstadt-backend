package co2computation

import (
	"errors"
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"io"
	"math"
	"testing"
	"time"
)

func TestComputationsSubfunctions(t *testing.T) {
	is := is.NewRelaxed(t)

	err := database.ConnectDatabase()
	is.NoErr(err)
	defer func() {
		err := database.DisconnectDatabase()
		is.NoErr(err)
	}()

	t.Run("TestGetEnergieCO2Faktor", TestGetEnergieCO2Faktor)
	t.Run("TestZaehlerNormalfall", TestZaehlerNormalfall)
	t.Run("TestZaehlerSpezialfall", TestZaehlerSpezialfall)
	t.Run("TestGebaeudeNormalfall", TestGebaeudeNormalfall)
}

func TestGetEnergieCO2Faktor(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("getEnergieversorgung: ID = 1, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 1
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.NoErr(err)                   // Normalfall wirft keine Errors
		is.Equal(co2Faktor, int32(144)) // erwartetes Ergebnis: 144
	})

	t.Run("getEnergieversorgung: ID = 2, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 2
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.NoErr(err)                   // Normalfall wirft keine Errors
		is.Equal(co2Faktor, int32(285)) // erwartetes Ergebnis: 285
	})

	t.Run("getEnergieversorgung: ID = 3, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 3
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.NoErr(err)                  // Normalfall wirft keine Errors
		is.Equal(co2Faktor, int32(72)) // erwartetes Ergebnis: 72
	})

	// Errortests
	t.Run("getEnergieversorgung: ID = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 0
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.Equal(err, io.EOF)         // Datenbank wirft EOF
		is.Equal(co2Faktor, int32(0)) // Fehlerfall liefert 0.0
	})

	t.Run("getEnergieversorgung: Jahr = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 1
		var jahr int32 = 0

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.Equal(err, ErrJahrNichtVorhanden) // Funktion wirft ErrJahrNichtVorhanden
		is.Equal(co2Faktor, int32(0))        // Fehlerfall liefert 0.0
	})
}

func TestZaehlerNormalfall(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("zaehlerNormalfall: ID = 2084, Einzelzaehler (Waermezaehler)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.WaermezaehlerFind(2084)
		var jahr int32 = 2020
		var gebaeudeNr int32 = 1101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(verbrauch, 704660.0) // erwartetes Ergebnis: 704660 (Verbrauch Jahr 2020)
		is.Equal(ngf, 0.0)            // erwartetes Ergebnis: 0.0 (keine Gruppenzaehler = keine weitere Flaeche)
	})

	t.Run("zaehlerNormalfall: ID = 2253, Gruppenzaehler (Waermezaehler)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.WaermezaehlerFind(2253)
		var jahr int32 = 2020
		var gebaeudeNr int32 = 0 // keine der Referenzen (1308, 1321) von Zaehler 2253

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.NoErr(err)            // Normalfall wirft keine Errors
		is.Equal(verbrauch, 0.0) // erwartetes Ergebnis: 0.0 (Verbrauch Jahr 2020)
		is.Equal(ngf, 3096.56)   // erwartetes Ergebnis: 3096.56 (Gruppenzaehler)
	})

	t.Run("zaehlerNormalfall: Umrechnung MWh in kWh (Waermezaehler)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.WaermezaehlerFind(2255)
		var jahr int32 = 2020
		var gebaeudeNr int32 = 1314

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(verbrauch, 120200.0) // erwartetes Ergebnis: 120200 (Verbrauch Jahr 2020)
		is.Equal(ngf, 0.0)            // erwartetes Ergebnis: 0.0 (kein Gruppenzaehler)
	})

	t.Run("zaehlerNormalfall: kWh beleibt kWh (Kaeltezaehler)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.KaeltezaehlerFind(6108)
		var jahr int32 = 2021
		var gebaeudeNr int32 = 1220

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.NoErr(err)                // Normalfall wirft keine Errors
		is.Equal(verbrauch, 17336.0) // erwartetes Ergebnis: 17336.0 (Verbrauch Jahr 2021)
		is.Equal(ngf, 0.0)           // erwartetes Ergebnis: 0.0 (kein Gruppenzaehler)
	})

	t.Run("zaehlerNormalfall: Stromzaehler", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.StromzaehlerFind(3524)
		var jahr int32 = 2020
		var gebaeudeNr int32 = 3506

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(verbrauch, 208676.2) // erwartetes Ergebnis: 208676.2 (Verbrauch Jahr 2020)
		is.Equal(ngf, 0.0)            // erwartetes Ergebnis: 0.0 (kein Gruppenzaehler)
	})

	// Errortests
	// Fehler tritt nur durch Datenfehler in der Datenbank auf
	t.Run("zaehlerNormalfall: Zaehler ohne Referenz zu Gebaeude", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 0
		zaehler := structs.Zaehler{PKEnergie: pkEnergie}
		var jahr int32 = 2020
		var gebaeudeNr int32 = 1101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.Equal(err, fmt.Errorf(ErrStrGebaeuderefFehlt, "zaehlerNormalfall", pkEnergie)) // Funktion wirft ErrStrGebaeuderefFehlt
		is.Equal(verbrauch, 0.0)                                                          // Fehlerfall liefert 0.0
		is.Equal(ngf, 0.0)                                                                // Fehlerfall liefert 0.0
	})

	t.Run("zaehlerNormalfall: Jahr nicht vorhanden in Zaehlerdaten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.WaermezaehlerFind(2084)
		var jahr int32 = 0
		var gebaeudeNr int32 = 1101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.Equal(err, fmt.Errorf(ErrStrVerbrauchFehlt, "zaehlerNormalfall", jahr, 2084)) // Funktion wirft ErrStrVerbrauchFehlt
		is.Equal(verbrauch, 0.0)                                                         // Fehlerfall liefert 0.0
		is.Equal(ngf, 0.0)                                                               // Fehlerfall liefert 0.0
	})

	// Fehler tritt nur durch Datenfehler in der Datenbank auf
	t.Run("zaehlerNormalfall: Einheit in Zaehler unbekannt", func(t *testing.T) {
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
			Einheit:     "tV",
			GebaeudeRef: []int32{1},
		}
		var jahr int32 = 2000
		var gebaeudeNr int32 = 1101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeNr)

		is.Equal(err, errors.New("zaehlerNormalfall: Einheit tV unbekannt")) // Funktion wirft ErrStrEinheitUnbekannt
		is.Equal(verbrauch, 0.0)                                             // Fehlerfall liefert 0.0
		is.Equal(ngf, 0.0)                                                   // Fehlerfall liefert 0.0
	})

	// Fehler tritt nur durch Datenfehler in der Datenbank auf
	t.Run("zaehlerNormalfall: referenziertes Gebaeude nicht vorhanden", func(t *testing.T) {
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

		is.Equal(err, io.EOF)    // Datenbank wirft EOF
		is.Equal(verbrauch, 0.0) // Fehlerfall liefert 0.0
		is.Equal(ngf, 0.0)       // Fehlerfall liefert 0.0
	})
}

func TestZaehlerSpezialfall(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("zaehlerSpezialfall: Spezialfall = 2, ID = 3621, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.KaeltezaehlerFind(3621)
		var jahr int32 = 2020
		var andereZaehlerID int32 = 3619

		verbrauch, err := zaehlerSpezialfall(zaehler, jahr, andereZaehlerID)

		is.NoErr(err)            // Normalfall wirft keine Errors
		is.Equal(verbrauch, 0.0) // erwartetes Ergebnis: 0.0 (Verbrauch Jahr 2020)
	})

	t.Run("zaehlerSpezialfall: Spezialfall = 3, ID = 3622, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.KaeltezaehlerFind(3622)
		var jahr int32 = 2020
		var andereZaehlerID int32 = 3620

		verbrauch, err := zaehlerSpezialfall(zaehler, jahr, andereZaehlerID)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(verbrauch, 958260.0) // erwartetes Ergebnis: 958260.0 (Verbrauch Jahr 2020)
	})

	t.Run("zaehlerSpezialfall: Spezialfall = 2, ID = 3621, Jahr = 2018", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.KaeltezaehlerFind(3621)
		var jahr int32 = 2018
		var andereZaehlerID int32 = 3619

		verbrauch, err := zaehlerSpezialfall(zaehler, jahr, andereZaehlerID)

		is.NoErr(err)                                    // Normalfall wirft keine Errors
		is.Equal(math.Round(verbrauch*100)/100, 33100.0) // erwartetes Ergebnis: 33100.0 (Verbrauch Jahr 2020)
	})

	// Errortests
	t.Run("zaehlerSpezialfall: Jahr = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler, _ := database.KaeltezaehlerFind(3622)
		var jahr int32 = 0
		var andereZaehlerID int32 = 3620

		verbrauch, err := zaehlerSpezialfall(zaehler, jahr, andereZaehlerID)

		is.Equal(err, errors.New("zaehlerSpezialfall: Kein Verbrauch für das Jahr 0, Zaehler: 3622")) // Funktion wirft ErrStrVerbrauchFehlt
		is.Equal(verbrauch, 0.0)                                                                      // Fehlerfall liefert 0.0
	})
}

func TestGebaeudeNormalfall(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("gebaeudeNormalfall: Flaechenanteil = 0 eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, _ := database.GebaeudeFind(1101)
		var co2Faktor int32 = 0
		var idEnergieversorgung int32 = 0
		var jahr int32 = 0
		var flaechenanteil int32 = 0

		emissionen, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (kein Flaechenanteil = keine Emissionen)
	})

	t.Run("gebaeudeNormalfall: keine Zaehler von bestimmten Typ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, _ := database.GebaeudeFind(1101)
		var co2Faktor int32 = 0
		var idEnergieversorgung int32 = 2
		var jahr int32 = 0
		var flaechenanteil int32 = 1000

		emissionen, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (kein Zaehler = keine berechenbaren Emissionen)
	})

	t.Run("gebaeudeNormalfall: einfach Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, _ := database.GebaeudeFind(1101)
		var co2Faktor int32 = 144
		var idEnergieversorgung int32 = 1
		var jahr int32 = 2020
		var flaechenanteil int32 = 1000

		emissionen, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.NoErr(err)                                           // Normalfall wirft keine Errors
		is.Equal(math.Round(emissionen*1000)/1000, 6604024.854) // erwartetes Ergebnis: 6604024.854
	})

	t.Run("gebaeudeNormalfall: Gebaeude mit mehreren Zaehlern", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, _ := database.GebaeudeFind(1103) // referenziert Zaehler 2349, 2350, 2351, 2352, 2353, 2354
		var co2Faktor int32 = 144
		var idEnergieversorgung int32 = 1
		var jahr int32 = 2020
		var flaechenanteil int32 = 1000

		emissionen, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.NoErr(err)                                           // Normalfall wirft keine Errors
		is.Equal(math.Round(emissionen*1000)/1000, 8632077.005) // erwartetes Ergebnis: 8632077.005
	})

	t.Run("gebaeudeNormalfall: Gebaeude mit Gruppenzaehler", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, _ := database.GebaeudeFind(2101) // Gruppe mit 2102, 2108
		var co2Faktor int32 = 144
		var idEnergieversorgung int32 = 1
		var jahr int32 = 2020
		var flaechenanteil int32 = 1000

		emissionen, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.NoErr(err)                                         // Normalfall wirft keine Errors
		is.Equal(math.Round(emissionen*100)/100, 22709184.09) // erwartetes Ergebnis: 22709184.09
	})

	// Errortests
	t.Run("gebaeudeNormalfall: negativer Flaechenanteil eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, _ := database.GebaeudeFind(1101)
		var co2Faktor int32 = 0
		var idEnergieversorgung int32 = 0
		var jahr int32 = 0
		var flaechenanteil int32 = -10

		emissionen, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.Equal(err, ErrFlaecheNegativ) // Funktion wirft ErrFlaecheNegativ
		is.Equal(emissionen, 0.0)        // Fehlerfall liefert 0.0
	})

	// Fehler tritt nur durch Datenfehler in der Datenbank auf
	t.Run("gebaeudeNormalfall: Gebaeude mit ungültiger Referenz", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude := structs.Gebaeude{
			Nr:        0,
			WaermeRef: []int32{0},
		}
		var co2Faktor int32 = 0
		var idEnergieversorgung int32 = 1
		var jahr int32 = 0
		var flaechenanteil int32 = 100

		emissionen, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.Equal(err, io.EOF)     // Datenbank wirft EOF
		is.Equal(emissionen, 0.0) // Fehlerfall liefert 0.0
	})
}
