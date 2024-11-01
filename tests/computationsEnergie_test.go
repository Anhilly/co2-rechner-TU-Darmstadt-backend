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

	err := database.ConnectDatabase("dev")
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

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
		is.Equal(verbrauch, 0.0)  // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
	})

	t.Run("BerechneEnergieverbrauch: leerer Slice", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
		is.Equal(verbrauch, 0.0)  // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
	})

	t.Run("BerechneEnergieverbrauch: einfache Eingabe, Einzelzaehler ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 1101, Nutzflaeche: 1000},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                    // Normalfall wirft keine Errors
		is.Equal(emissionen, 6604015.48) // erwartetes Ergebnis: 6604015.48
		is.Equal(verbrauch, 45861.22)    // erwartetes Ergebnis: 45861.22
	})

	t.Run("BerechneEnergieverbrauch: einfache Eingabe, Gebaeude mehrere Zaehler ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 3260, Nutzflaeche: 1000},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 3

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                       // Normalfall wirft keine Errors
		is.Equal(emissionen, 3486568639.15) // erwartetes Ergebnis: 3486568639.15
		is.Equal(verbrauch, 48424564.43)
	})

	t.Run("BerechneEnergieverbrauch: einfache Eingabe, Gruppenzaehler ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{ // Zaehler: 3807, weiteres Gebaeude: 3016
			{GebaeudeNr: 3102, Nutzflaeche: 1000},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 3

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                    // Normalfall wirft keine Errors
		is.Equal(emissionen, 1085282.24) // erwartetes Ergebnis: 1085282.24
		is.Equal(verbrauch, 15073.36)
	})

	// Test ueberprueft, ob Referenzen von Gebaeude 1321 korrekt angepasst
	t.Run("BerechneEnergieverbrauch: einfache Eingabe, Gebaeude 1321 hat nur noch ein Zaehler ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 1321, Nutzflaeche: 1000},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 2

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                    // Normalfall wirft keine Errors
		is.Equal(emissionen, 8804939.48) // erwartetes Ergebnis: 8804939.48
		is.Equal(verbrauch, 30894.52)
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

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                     // Normalfall wirft keine Errors
		is.Equal(emissionen, 22347183.87) // erwartetes Ergebnis: 22347183.87
		is.Equal(verbrauch, 155188.78)
	})

	t.Run("BerechneEnergieverbrauch: Gebaeude ohne Zaehler", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 1101, Nutzflaeche: 100},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 3

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (kein Zaehler = kein berechenbarer Verbrauch)
		is.Equal(verbrauch, 0.0)
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

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (kein Nutzflaeche = keine Emissionen)
		is.Equal(verbrauch, 0.0)
	})

	// Spezialfaelle
	t.Run("BerechneEnergieverbrauch: Kaeltezaehler, Spezialfall 2", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 3202, Nutzflaeche: 100},
		}
		var jahr int32 = 2020
		var idEnergieversorgung = structs.IDEnergieversorgungKaelte

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // Spezialfall 2 wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0
		is.Equal(verbrauch, 0.0)
	})

	t.Run("BerechneEnergieverbrauch: Kaeltezaehler, Spezialfall 3", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 3204, Nutzflaeche: 100},
		}
		var jahr int32 = 2020
		var idEnergieversorgung = structs.IDEnergieversorgungKaelte

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)                   // Spezialfall 3 wirft keine Errors
		is.Equal(emissionen, 659215.24) // erwartetes Ergebnis: 659215.24
		is.Equal(verbrauch, 9155.77)
	})

	// Errortests
	t.Run("BerechneEnergieverbrauch: idEnergieversorgung = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 0

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(emissionen, 0.0)           // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
		is.Equal(verbrauch, 0.0)
	})

	t.Run("BerechneEnergieverbrauch: Jahr = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{}
		var jahr int32 = 0
		var idEnergieversorgung int32 = 1

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.Equal(err, structs.ErrJahrNichtVorhanden) // Funktion wirft ErrJahrNichtVorhanden
		is.Equal(emissionen, 0.0)                    // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
		is.Equal(verbrauch, 0.0)
	})

	t.Run("BerechneEnergieverbrauch: Gebaeude Nr = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 0, Nutzflaeche: 10},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(emissionen, 0.0)           // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
		is.Equal(verbrauch, 0.0)
	})

	t.Run("BerechneEnergieverbrauch: negativer Nutzflaeche eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeFlaecheDaten := []structs.UmfrageGebaeude{
			{GebaeudeNr: 1101, Nutzflaeche: -10},
		}
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, verbrauch, err := co2computation.BerechneEnergieverbrauch(gebaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.Equal(err, structs.ErrFlaecheNegativ) // Funktion wirft ErrFlaecheNegativ
		is.Equal(emissionen, 0.0)                // im Fehlerfall ist Emissionen = 0.0
		is.Equal(verbrauch, 0.0)
	})

	// Fehler GebaeudeSpezialfall --> Was ist 'Spezialfall' für ein Attribut bzw. warum ist es immer ==1?
	// Fehler ErrStrEinheitUnbekannt momentan nicht abpruefbar. Benoetigt falschen Datensatz in Datenbank
}
