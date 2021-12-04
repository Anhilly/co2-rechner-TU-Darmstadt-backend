package tests

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/co2computation"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"io"
	"testing"
)

func TestComputationsEnergie(t *testing.T) {
	database.ConnectDatabase()
	defer database.DisconnectDatabase()

	t.Run("TestBerechneEnergieverbrauch", TestBerechneEnergieverbrauch)
}

func TestBerechneEnergieverbrauch(t *testing.T) {
	is := is.NewRelaxed(t)

	// normale Berechnungen
	t.Run("BerechneEnergieverbrauch: Slice = nil", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var geaeudeFlaecheDaten []structs.GebaeudeFlaecheAPI = nil
		var jahr int32 = 2020             // muss gueltiges Jahr sein
		var idEnergieversorgung int32 = 1 // muss gueltige ID sein

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // bei normalen Berechnungen sollte kein Error geworfen werden
		is.Equal(emissionen, 0.0) // ohne Eingaben sind Emissionen = 0.0
	})

	t.Run("BerechneEnergieverbrauch: leerer Slice", func(t *testing.T) {
		is := is.NewRelaxed(t)

		geaeudeFlaecheDaten := []structs.GebaeudeFlaecheAPI{}
		var jahr int32 = 2020             // muss gueltiges Jahr sein
		var idEnergieversorgung int32 = 1 // muss gueltige ID sein

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // bei normalen Berechnungen sollte kein Error geworfen werden
		is.Equal(emissionen, 0.0) // ohne Eingaben sind Emissionen = 0.0
	})

	t.Run("BerechneEnergieverbrauch: einfache Eingabe, Einzelzaehler ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		geaeudeFlaecheDaten := []structs.GebaeudeFlaecheAPI{
			{GebaeudeNr: 1101, Flaechenanteil: 1000},
		}
		var jahr int32 = 2020             // muss gueltiges Jahr sein
		var idEnergieversorgung int32 = 1 // muss gueltige ID sein

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                     // bei normalen Berechnungen sollte kein Error geworfen werden
		is.Equal(emissionen, 6604024.85) // erwartetes Ergebnis: 6604024.85
	})

	t.Run("BerechneEnergieverbrauch: einfache Eingabe, Gebaeude mehrere Zaehler ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		geaeudeFlaecheDaten := []structs.GebaeudeFlaecheAPI{	// Zaehler 2250, 2251, 2252, 2085
			{GebaeudeNr: 1108, Flaechenanteil: 1000},
		}
		var jahr int32 = 2020             // muss gueltiges Jahr sein
		var idEnergieversorgung int32 = 1 // muss gueltige ID sein

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                     // bei normalen Berechnungen sollte kein Error geworfen werden
		is.Equal(emissionen, 23126680.04) // erwartetes Ergebnis: 23126680.04
	})

	t.Run("BerechneEnergieverbrauch: einfache Eingabe, Gruppenzaehler ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		geaeudeFlaecheDaten := []structs.GebaeudeFlaecheAPI{	// Zaehler: 3807, weiteres Gebaeude: 3016
			{GebaeudeNr: 3102, Flaechenanteil: 1000},
		}
		var jahr int32 = 2020             // muss gueltiges Jahr sein
		var idEnergieversorgung int32 = 3 // muss gueltige ID sein

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                     // bei normalen Berechnungen sollte kein Error geworfen werden
		is.Equal(emissionen, 1085282.24) // erwartetes Ergebnis: 1085282.24
	})

	t.Run("BerechneEnergieverbrauch: komplexe Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		geaeudeFlaecheDaten := []structs.GebaeudeFlaecheAPI{
			{GebaeudeNr: 1101, Flaechenanteil: 1000},
			{GebaeudeNr: 1108, Flaechenanteil: 1000},
			{GebaeudeNr: 1103, Flaechenanteil: 1000},
		}
		var jahr int32 = 2020             // muss gueltiges Jahr sein
		var idEnergieversorgung int32 = 1 // muss gueltige ID sein

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                     // bei normalen Berechnungen sollte kein Error geworfen werden
		is.Equal(emissionen, 6604024.85 + 23126680.04 + 8632077.01) // erwartetes Ergebnis: 1085282.24
	})

	t.Run("BerechneEnergieverbrauch: Gebauede ohne Zaehler", func(t *testing.T) {
		is := is.NewRelaxed(t)

		geaeudeFlaecheDaten := []structs.GebaeudeFlaecheAPI{
			{GebaeudeNr: 1101, Flaechenanteil: 100},
		}
		var jahr int32 = 2020             // muss gueltiges Jahr sein
		var idEnergieversorgung int32 = 3 // muss gueltige ID sein

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // bei normalen Berechnungen sollte kein Error geworfen werden
		is.Equal(emissionen, 0.0) // ohne Zaehlerdaten ist Emissionen = 0.0
	})

	t.Run("BerechneEnergieverbrauch: Flaechenanteil = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		geaeudeFlaecheDaten := []structs.GebaeudeFlaecheAPI{
			{GebaeudeNr: 1101, Flaechenanteil: 0},
			{GebaeudeNr: 1160, Flaechenanteil: 0},
			{GebaeudeNr: 1217, Flaechenanteil: 0},
			{GebaeudeNr: 3206, Flaechenanteil: 0},
		}
		var jahr int32 = 2020             // muss gueltiges Jahr sein
		var idEnergieversorgung int32 = 1 // muss gueltige ID sein

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // bei normalen Berechnungen sollte kein Error geworfen werden
		is.Equal(emissionen, 0.0) // ohne Flaechenanteil sind Emissionen = 0.0
	})

	// Errortests
	// auch in TestGetEnergieCO2Faktor
	t.Run("BerechneEnergieverbrauch: idEnergieversorgung = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		geaeudeFlaecheDaten := []structs.GebaeudeFlaecheAPI{}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 0

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.Equal(err, io.EOF)     // Datenbank wirft EOF
		is.Equal(emissionen, 0.0) // im Fehlerfall ist Emissionen = 0.0
	})

	// auch in TestGetEnergieCO2Faktor
	t.Run("BerechneEnergieverbrauch: Jahr = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		geaeudeFlaecheDaten := []structs.GebaeudeFlaecheAPI{}
		var jahr int32 = 0
		var idEnergieversorgung int32 = 1 // muss gueltige ID sein

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.Equal(err, co2computation.ErrJahrNichtVorhanden) // Funktion wirft ErrJahrNichtVorhanden
		is.Equal(emissionen, 0.0)                           // im Fehlerfall ist Emissionen = 0.0
	})

	t.Run("BerechneEnergieverbrauch: Gebaeude Nr = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		geaeudeFlaecheDaten := []structs.GebaeudeFlaecheAPI{
			{GebaeudeNr: 0, Flaechenanteil: 10},
		}
		var jahr int32 = 2020             // muss gueltiges Jahr sein
		var idEnergieversorgung int32 = 1 // muss gueltige ID sein

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.Equal(err, io.EOF)     // Datenbank wirft EOF
		is.Equal(emissionen, 0.0) // im Fehlerfall ist Emissionen = 0.0
	})

	t.Run("BerechneEnergieverbrauch: negativer Flaechenanteil", func(t *testing.T) {
		is := is.NewRelaxed(t)

		geaeudeFlaecheDaten := []structs.GebaeudeFlaecheAPI{
			{GebaeudeNr: 1101, Flaechenanteil: -10},
		}
		var jahr int32 = 2020             // muss gueltiges Jahr sein
		var idEnergieversorgung int32 = 1 // muss gueltige ID sein

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.Equal(err, co2computation.ErrFlaecheNegativ) // Datenbank wirft EOF
		is.Equal(emissionen, 0.0)                       // im Fehlerfall ist Emissionen = 0.0
	})
}
