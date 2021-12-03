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

	t.Run("BerechneEnergieverbrauch: einfache Eingabe (ein Gebaeude mit einem Zaehler) ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		geaeudeFlaecheDaten := []structs.GebaeudeFlaecheAPI{
			{GebaeudeNr: 1101, Flaechenanteil: 1000},
		}
		var jahr int32 = 2020             // muss gueltiges Jahr sein
		var idEnergieversorgung int32 = 3 // muss gueltige ID sein

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                     // bei normalen Berechnungen sollte kein Error geworfen werden
		is.Equal(emissionen, 3302012.427) // erwartetes Ergebnis: 3302012.427
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
