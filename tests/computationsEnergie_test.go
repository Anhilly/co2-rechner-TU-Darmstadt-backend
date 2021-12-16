package tests

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/co2computation"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestComputationsEnergie(t *testing.T) {
	is := is.NewRelaxed(t)

	err := database.ConnectDatabase()
	is.NoErr(err)
	defer func() {
		err := database.DisconnectDatabase()
		is.NoErr(err)
	}()

	t.Run("TestBerechneEnergieverbrauch", TestBerechneEnergieverbrauch)
}

func TestBerechneEnergieverbrauch(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("BerechneEnergieverbrauch: Slice = nil", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var gebaeudeFlaecheDaten []structs.UmfrageGebaeude = nil
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
	})

	t.Run("BerechneEnergieverbrauch: leerer Slice", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
	})

	t.Run("BerechneEnergieverbrauch: einfache Eingabe, Einzelzaehler ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 1101, Nutzflaeche: 1000},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                    // Normalfall wirft keine Errors
		is.Equal(emissionen, 6604024.85) // erwartetes Ergebnis: 6604024.85
	})

	t.Run("BerechneEnergieverbrauch: einfache Eingabe, Gebaeude mehrere Zaehler ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 1108, Nutzflaeche: 1000}, // Zaehler 2250, 2251, 2252, 2085
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                     // Normalfall wirft keine Errors
		is.Equal(emissionen, 23126680.04) // erwartetes Ergebnis: 23126680.04
	})

	t.Run("BerechneEnergieverbrauch: einfache Eingabe, Gruppenzaehler ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{ // Zaehler: 3807, weiteres Gebaeude: 3016
			{GebaeudeNr: 3102, Nutzflaeche: 1000},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 3

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                    // Normalfall wirft keine Errors
		is.Equal(emissionen, 1085282.24) // erwartetes Ergebnis: 1085282.24
	})

	// Test ueberprueft, ob Referenzen von Gebaeude 1321 korrekt angepasst
	t.Run("BerechneEnergieverbrauch: einfache Eingabe, Gebaeude 1321 hat nur noch ein Zaehler ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 1321, Nutzflaeche: 1000},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 2

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                    // Normalfall wirft keine Errors
		is.Equal(emissionen, 8804937.41) // erwartetes Ergebnis: 8804937.41
	})

	t.Run("BerechneEnergieverbrauch: komplexe Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 1101, Nutzflaeche: 1000},
			{GebaeudeNr: 1108, Nutzflaeche: 1000},
			{GebaeudeNr: 1103, Nutzflaeche: 1000},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                    // Normalfall wirft keine Errors
		is.Equal(emissionen, 38362781.9) // erwartetes Ergebnis: 38362781.9
	})

	t.Run("BerechneEnergieverbrauch: Gebaeude ohne Zaehler", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 1101, Nutzflaeche: 100},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 3

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (kein Zaehler = kein berechenbarer Verbrauch)
	})

	t.Run("BerechneEnergieverbrauch: Nutzflaeche = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 1101, Nutzflaeche: 0},
			{GebaeudeNr: 1160, Nutzflaeche: 0},
			{GebaeudeNr: 1217, Nutzflaeche: 0},
			{GebaeudeNr: 3206, Nutzflaeche: 0},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (kein Nutzflaeche = keine Emissionen)
	})

	// Errortests
	t.Run("BerechneEnergieverbrauch: idEnergieversorgung = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 0

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(emissionen, 0.0)           // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
	})

	t.Run("BerechneEnergieverbrauch: Jahr = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{}
		var jahr int32 = 0
		var idEnergieversorgung int32 = 1

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.Equal(err, structs.ErrJahrNichtVorhanden) // Funktion wirft ErrJahrNichtVorhanden
		is.Equal(emissionen, 0.0)                    // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
	})

	t.Run("BerechneEnergieverbrauch: Gebaeude Nr = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 0, Nutzflaeche: 10},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(emissionen, 0.0)           // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
	})

	t.Run("BerechneEnergieverbrauch: negativer Nutzflaeche eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 1101, Nutzflaeche: -10},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.Equal(err, structs.ErrFlaecheNegativ) // Funktion wirft ErrFlaecheNegativ
		is.Equal(emissionen, 0.0)                // im Fehlerfall ist Emissionen = 0.0
	})
}
