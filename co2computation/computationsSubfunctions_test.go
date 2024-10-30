package co2computation

import (
	"errors"
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"testing"
	"time"
)

func TestComputationsSubfunctions(t *testing.T) {
	is := is.NewRelaxed(t)

	err := database.ConnectDatabase("dev")
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

		is.NoErr(err) // Normalfall wirft keine Errors
		is.Equal(co2Faktor, map[int32]int32{
			1: 144,
		}) // erwartetes Ergebnis: map[1:144]
	})

	t.Run("getEnergieversorgung: ID = 2, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 2
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.NoErr(err) // Normalfall wirft keine Errors
		is.Equal(co2Faktor, map[int32]int32{
			1: 285,
		}) // erwartetes Ergebnis: map[1:285]
	})

	t.Run("getEnergieversorgung: ID = 3, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 3
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.NoErr(err) // Normalfall wirft keine Errors
		is.Equal(co2Faktor, map[int32]int32{
			1: 72,
		}) // erwartetes Ergebnis: map[1:72]
	})

	// Errortests
	t.Run("getEnergieversorgung: ID = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 0
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(co2Faktor, nil)            // Fehlerfall liefert nil
	})

	t.Run("getEnergieversorgung: Jahr = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 1
		var jahr int32 = 0

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.Equal(err, structs.ErrJahrNichtVorhanden) // Funktion wirft ErrJahrNichtVorhanden
		is.Equal(co2Faktor, nil)                     // Fehlerfall liefert nil
	})
}

func TestZaehlerNormalfall(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("zaehlerNormalfall: OID = 6710a1ca3a47b613426ff656, Einzelzaehler (Waermezaehler)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff656") // S101XXXXXXHE000XXXXXXZ40CO00001
		zaehler, _ := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungWaerme)

		var jahr int32 = 2020
		gebaeudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a286") // 1101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeOID)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(verbrauch, 704659.0) // erwartetes Ergebnis: 704659 (Verbrauch Jahr 2020)
		is.Equal(ngf, 0.0)            // erwartetes Ergebnis: 0.0 (keine Gruppenzaehler = keine weitere Flaeche)
	})

	t.Run("zaehlerNormalfall: OID = 6710a1c93a47b613426ff62b, Gruppenzaehler (Waermezaehler)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1c93a47b613426ff62b") // B101XXXXXXHE000XXXXXXZ40CO00001
		zaehler, _ := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungWaerme)

		var jahr int32 = 2020
		gebaeudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a2cf") // 2101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeOID)

		is.NoErr(err)                             // Normalfall wirft keine Errors
		is.Equal(math.Round(verbrauch), 788658.0) // erwartetes Ergebnis: 788658.99 (Verbrauch Jahr 2020)
		is.Equal(ngf, 1772.13)                    // erwartetes Ergebnis: 1772.13 (Gruppenzaehler)
	})

	t.Run("zaehlerNormalfall: Stromzaehler", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff695") // L101_Verbrauch_Gesamt
		zaehler, _ := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungStrom)

		var jahr int32 = 2023
		gebaeudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a2e7") // 3101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeOID)

		is.NoErr(err)                    // Normalfall wirft keine Errors
		is.Equal(verbrauch, 2279367.734) // erwartetes Ergebnis: 2279367.734 (Verbrauch Jahr 2020)
		is.Equal(ngf, 0.0)               // erwartetes Ergebnis: 0.0 (kein Gruppenzaehler)
	})

	// Errortests
	// Fehler tritt nur durch Datenfehler in der Datenbank auf
	t.Run("zaehlerNormalfall: Zaehler ohne Referenz zu Gebaeude", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehler := structs.Zaehler{
			ZaehlerID: primitive.NilObjectID,
			DPName:    "xxxx",
		}
		var jahr int32 = 2020

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, primitive.NilObjectID)

		is.Equal(err, fmt.Errorf(structs.ErrStrGebaeuderefFehlt, "zaehlerNormalfall", "xxxx")) // Funktion wirft ErrStrGebaeuderefFehlt
		is.Equal(verbrauch, 0.0)                                                               // Fehlerfall liefert 0.0
		is.Equal(ngf, 0.0)                                                                     // Fehlerfall liefert 0.0
	})

	t.Run("zaehlerNormalfall: Jahr nicht vorhanden in Zaehlerdaten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff656") // S101XXXXXXHE000XXXXXXZ40CO00001
		zaehler, _ := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungWaerme)

		var jahr int32 = 0
		gebaeudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a286") // 1101

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, gebaeudeOID)

		is.Equal(err, fmt.Errorf(structs.ErrStrVerbrauchFehlt, "zaehlerNormalfall", jahr, "S101XXXXXXHE000XXXXXXZ40CO00001")) // Funktion wirft ErrStrVerbrauchFehlt
		is.Equal(verbrauch, 0.0)                                                                                              // Fehlerfall liefert 0.0
		is.Equal(ngf, 0.0)                                                                                                    // Fehlerfall liefert 0.0
	})

	// Fehler tritt nur durch Datenfehler in der Datenbank auf
	t.Run("zaehlerNormalfall: Einheit in Zaehler unbekannt", func(t *testing.T) {
		is := is.NewRelaxed(t)

		location, _ := time.LoadLocation("Etc/GMT")
		zaehler := structs.Zaehler{
			ZaehlerID: primitive.NewObjectID(),
			DPName:    "xxxx",
			Zaehlerdaten: []structs.Zaehlerwerte{
				{
					Wert:        788.66,
					Zeitstempel: time.Date(2000, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "tV",
			GebaeudeRef: []primitive.ObjectID{primitive.NewObjectID()},
		}
		var jahr int32 = 2000

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, primitive.NilObjectID)

		is.Equal(err, errors.New("zaehlerNormalfall: Einheit tV unbekannt")) // Funktion wirft ErrStrEinheitUnbekannt
		is.Equal(verbrauch, 0.0)                                             // Fehlerfall liefert 0.0
		is.Equal(ngf, 0.0)                                                   // Fehlerfall liefert 0.0
	})

	// Fehler tritt nur durch Datenfehler in der Datenbank auf
	t.Run("zaehlerNormalfall: referenziertes Gebaeude nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeOID := primitive.NewObjectID()

		location, _ := time.LoadLocation("Etc/GMT")
		zaehler := structs.Zaehler{
			ZaehlerID: primitive.NewObjectID(),
			DPName:    "xxxx",
			Zaehlerdaten: []structs.Zaehlerwerte{
				{
					Wert:        788.66,
					Zeitstempel: time.Date(2000, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "MWh",
			GebaeudeRef: []primitive.ObjectID{gebaeudeOID},
		}
		var jahr int32 = 2000

		verbrauch, ngf, err := zaehlerNormalfall(zaehler, jahr, primitive.NilObjectID)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(verbrauch, 0.0)            // Fehlerfall liefert 0.0
		is.Equal(ngf, 0.0)                  // Fehlerfall liefert 0.0
	})
}

func TestZaehlerSpezialfall(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("zaehlerSpezialfall: Spezialfall = 2, ID = 6691, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff682") // L202XXXXXaKA000XXXXXXZ50CO00001
		zaehler, _ := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungKaelte)

		is.Equal(zaehler.Spezialfall, int32(2))

		var jahr int32 = 2020
		verbrauch, err := zaehlerSpezialfall(zaehler, jahr, "L402XXXXXXKA000XXXXXXZ50CO00001")

		is.NoErr(err)            // Normalfall wirft keine Errors
		is.Equal(verbrauch, 0.0) // erwartetes Ergebnis: 0.0 (Verbrauch Jahr 2020)
	})

	t.Run("zaehlerSpezialfall: Spezialfall = 3, ID = 3622, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff685") // L204XXXXXXKA000XXXXXXZ50CO00001
		zaehler, _ := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungKaelte)

		is.Equal(zaehler.Spezialfall, int32(3))

		var jahr int32 = 2020
		verbrauch, err := zaehlerSpezialfall(zaehler, jahr, "L206XXXXXXKA000XXXXXXZ50CO00001")

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(verbrauch, 958260.0) // erwartetes Ergebnis: 958260.0 (Verbrauch Jahr 2020)
	})

	// Errortests
	t.Run("zaehlerSpezialfall: Jahr = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff682") // L202XXXXXaKA000XXXXXXZ50CO00001
		zaehler, _ := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungKaelte)

		is.Equal(zaehler.Spezialfall, int32(2))

		var jahr int32 = 0
		verbrauch, err := zaehlerSpezialfall(zaehler, jahr, "L402XXXXXXKA000XXXXXXZ50CO00001")

		is.Equal(err, errors.New("zaehlerSpezialfall: Kein Verbrauch für das Jahr 0, Zaehler: L202XXXXXaKA000XXXXXXZ50CO00001")) // Funktion wirft ErrStrVerbrauchFehlt
		is.Equal(verbrauch, 0.0)                                                                                                 // Fehlerfall liefert 0.0
	})
}

func TestGebaeudeNormalfall(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("gebaeudeNormalfall: Flaechenanteil = 0 eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, _ := database.GebaeudeFind(1101)
		co2Faktor := map[int32]int32{
			0: 0,
		}
		var idEnergieversorgung int32 = 0
		var jahr int32 = 0
		var flaechenanteil int32 = 0

		emissionen, gvNutzflaeche, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.NoErr(err)                // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0)    // erwartetes Ergebnis: 0.0 (kein Flaechenanteil = keine Emissionen)
		is.Equal(gvNutzflaeche, 0.0) // erwartetes Ergebnis: 0.0
	})

	t.Run("gebaeudeNormalfall: keine Zaehler von bestimmten Typ", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, _ := database.GebaeudeFind(1101)
		co2Faktor := map[int32]int32{
			1: 100,
		}
		var idEnergieversorgung int32 = 2
		var jahr int32 = 2020
		var flaechenanteil int32 = 1000

		emissionen, gvNutzflaeche, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.NoErr(err)                // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0)    // erwartetes Ergebnis: 0.0 (kein Zaehler = keine berechenbaren Emissionen)
		is.Equal(gvNutzflaeche, 0.0) // erwartetes Ergebnis: 0.0
	})

	t.Run("gebaeudeNormalfall: einfach Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, _ := database.GebaeudeFind(1101)
		co2Faktor := map[int32]int32{
			1: 144,
		}
		var idEnergieversorgung int32 = 1
		var jahr int32 = 2020
		var flaechenanteil int32 = 1000

		emissionen, gvNutzflaeche, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.NoErr(err)                                            // Normalfall wirft keine Errors
		is.Equal(math.Round(emissionen*1000)/1000, 6604015.482)  // erwartetes Ergebnis: 6604015.482
		is.Equal(math.Round(gvNutzflaeche*1000)/1000, 45861.219) // erwartetes Ergebnis: 45861.219
	})

	t.Run("gebaeudeNormalfall: Gebaeude mit mehreren Zaehlern", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, _ := database.GebaeudeFind(3260)
		co2Faktor := map[int32]int32{
			1: 144,
		}
		var idEnergieversorgung int32 = structs.IDEnergieversorgungKaelte
		var jahr int32 = 2020
		var flaechenanteil int32 = 1000

		emissionen, gvNutzflaeche, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.NoErr(err)                                               // Normalfall wirft keine Errors
		is.Equal(math.Round(emissionen*1000)/1000, 6973137278.292)  // erwartetes Ergebnis: 6973137278.292
		is.Equal(math.Round(gvNutzflaeche*1000)/1000, 48424564.433) // erwartetes Ergebnis: 48424564.433
	})

	t.Run("gebaeudeNormalfall: Gebaeude mit Gruppenzaehler", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, _ := database.GebaeudeFind(2101) // Gruppe mit 2102, 2108
		co2Faktor := map[int32]int32{
			1: 144,
		}
		var idEnergieversorgung int32 = 1
		var jahr int32 = 2020
		var flaechenanteil int32 = 1000

		emissionen, gvNutzflaeche, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.NoErr(err)                                          // Normalfall wirft keine Errors
		is.Equal(math.Round(emissionen*100)/100, 22709126.5)   // erwartetes Ergebnis: 22709126.5
		is.Equal(math.Round(gvNutzflaeche*100)/100, 157702.27) // erwartetes Ergebnis: 157702.27
	})

	// Errortests
	t.Run("gebaeudeNormalfall: negativer Flaechenanteil eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, _ := database.GebaeudeFind(1101)
		co2Faktor := map[int32]int32{
			0: 0,
		}
		var idEnergieversorgung int32 = 0
		var jahr int32 = 0
		var flaechenanteil int32 = -10

		emissionen, gvNutzflaeche, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.Equal(err, structs.ErrFlaecheNegativ) // Funktion wirft ErrFlaecheNegativ
		is.Equal(emissionen, 0.0)                // Fehlerfall liefert 0.0
		is.Equal(gvNutzflaeche, 0.0)             // Fehlerfall liefert 0.0
	})

	// Fehler tritt nur durch Datenfehler in der Datenbank auf
	t.Run("gebaeudeNormalfall: Gebaeude mit ungültiger Referenz", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude := structs.Gebaeude{
			Nr:        0,
			WaermeRef: []primitive.ObjectID{primitive.NilObjectID},
			Waermeversorger: []structs.Versoger{
				{Jahr: 0, IDVertrag: 1},
			},
		}
		co2Faktor := map[int32]int32{
			0: 0,
		}
		var idEnergieversorgung int32 = 1
		var jahr int32 = 0
		var flaechenanteil int32 = 100

		emissionen, gvNutzflaeche, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(emissionen, 0.0)           // Fehlerfall liefert 0.0
		is.Equal(gvNutzflaeche, 0.0)        // Fehlerfall liefert 0.0
	})

	// Fehler tritt nur durch Datenfehler in der Datenbank auf
	t.Run("gebaeudeNormalfall: Gebaeude ohne Versorger", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude := structs.Gebaeude{
			Nr:        0,
			WaermeRef: []primitive.ObjectID{primitive.NilObjectID},
		}
		co2Faktor := map[int32]int32{
			0: 0,
		}
		var idEnergieversorgung int32 = 1
		var jahr int32 = 0
		var flaechenanteil int32 = 100

		emissionen, gvNutzflaeche, err := gebaeudeNormalfall(co2Faktor, gebaeude, idEnergieversorgung, jahr, flaechenanteil)

		is.Equal(err, fmt.Errorf(structs.ErrStrKeinVersorger, gebaeude.Nr, jahr)) // Datenbank wirft ErrNoDocuments
		is.Equal(emissionen, 0.0)                                                 // Fehlerfall liefert 0.0
		is.Equal(gvNutzflaeche, 0.0)                                              // Fehlerfall liefert 0.0
	})
}
