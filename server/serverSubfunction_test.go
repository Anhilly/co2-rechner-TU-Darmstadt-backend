package server

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestServerSubfunctions(t *testing.T) {
	is := is.NewRelaxed(t)

	err := database.ConnectDatabase("dev")
	is.NoErr(err)
	defer func() {
		err := database.DisconnectDatabase()
		is.NoErr(err)
	}()

	// Auswertung
	t.Run("TestBinaereZahlerdatenFuerZaehler", TestBinaereZahlerdatenFuerZaehler)
}

func TestBinaereZahlerdatenFuerZaehler(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("binaereZahlerdatenFuerZaehler: ein Zaehler mit vollstaendingen Daten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID := primitive.NewObjectID()
		location, _ := time.LoadLocation("Etc/GMT")
		aktuellesJahr := int32(time.Now().Year())

		alleZahler := []structs.ZaehlerUndZaehlerdaten{
			{
				ZaehlerID:    zaehlerOID,
				Zaehlerdaten: []structs.Zaehlerwerte{},
			},
		}

		for i := structs.ErstesJahr; i <= aktuellesJahr; i++ {
			for j := 1; j <= 12; j++ {
				alleZahler[0].Zaehlerdaten = append(alleZahler[0].Zaehlerdaten, structs.Zaehlerwerte{
					Wert:        0.0,
					Zeitstempel: time.Date(int(i), time.Month(j), 01, 0, 0, 0, 0, location).UTC(),
				})
			}
		}

		compareDoc := []structs.ZaehlerUndZaehlerdatenVorhanden{
			{
				ZaehlerID:             zaehlerOID,
				ZaehlerdatenVorhanden: []structs.ZaehlerwertVorhanden{},
			},
		}
		for i := structs.ErstesJahr; i <= aktuellesJahr; i++ {
			compareDoc[0].ZaehlerdatenVorhanden = append(compareDoc[0].ZaehlerdatenVorhanden, structs.ZaehlerwertVorhanden{Jahr: i, Vorhanden: true})
		}

		result := binaereZahlerdatenFuerZaehler(alleZahler)
		is.Equal(result, compareDoc)
	})

	t.Run("binaereZahlerdatenFuerZaehler: ein Zaehler mit unvollstaendingen Daten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID := primitive.NewObjectID()

		location, _ := time.LoadLocation("Etc/GMT")
		alleZahler := []structs.ZaehlerUndZaehlerdaten{
			{
				ZaehlerID: zaehlerOID,
				Zaehlerdaten: []structs.Zaehlerwerte{
					{
						Wert:        414.61,
						Zeitstempel: time.Date(2020, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
					{
						Wert:        555.3,
						Zeitstempel: time.Date(2019, time.February, 01, 0, 0, 0, 0, location).UTC(),
					},
					{
						Wert:        169.59,
						Zeitstempel: time.Date(2018, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
					{
						Wert:        380.67,
						Zeitstempel: time.Date(2021, time.May, 01, 0, 0, 0, 0, location).UTC(),
					},
					{
						Wert:        370.39,
						Zeitstempel: time.Date(2022, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
				},
			},
		}

		aktuellesJahr := int32(time.Now().Year())
		compareDoc := []structs.ZaehlerUndZaehlerdatenVorhanden{
			{
				ZaehlerID:             zaehlerOID,
				ZaehlerdatenVorhanden: []structs.ZaehlerwertVorhanden{},
			},
		}

		for i := structs.ErstesJahr; i <= aktuellesJahr; i++ {
			compareDoc[0].ZaehlerdatenVorhanden = append(compareDoc[0].ZaehlerdatenVorhanden, structs.ZaehlerwertVorhanden{Jahr: i, Vorhanden: false})
		}
		compareDoc[0].ZaehlerdatenVorhanden[0].Vorhanden = true
		compareDoc[0].ZaehlerdatenVorhanden[1].Vorhanden = true
		compareDoc[0].ZaehlerdatenVorhanden[2].Vorhanden = true

		result := binaereZahlerdatenFuerZaehler(alleZahler)
		is.Equal(result, compareDoc)
	})

	t.Run("binaereZahlerdatenFuerZaehler: mehrere Zaehler mit unvollstaendingen Daten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID1 := primitive.NewObjectID()
		zaehlerOID2 := primitive.NewObjectID()
		zaehlerOID3 := primitive.NewObjectID()

		location, _ := time.LoadLocation("Etc/GMT")
		alleZahler := []structs.ZaehlerUndZaehlerdaten{
			{
				ZaehlerID: zaehlerOID1,
				Zaehlerdaten: []structs.Zaehlerwerte{
					{
						Wert:        380.67,
						Zeitstempel: time.Date(2021, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
					{
						Wert:        370.39,
						Zeitstempel: time.Date(2022, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
				},
			},
			{
				ZaehlerID: zaehlerOID2,
				Zaehlerdaten: []structs.Zaehlerwerte{
					{
						Wert:        169.59,
						Zeitstempel: time.Date(2018, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
					{
						Wert:        380.67,
						Zeitstempel: time.Date(2021, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
				},
			},
			{
				ZaehlerID: zaehlerOID3,
				Zaehlerdaten: []structs.Zaehlerwerte{
					{
						Wert:        370.39,
						Zeitstempel: time.Date(2022, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
					{
						Wert:        370.39,
						Zeitstempel: time.Date(3000, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
				},
			},
		}

		compareDoc := []structs.ZaehlerUndZaehlerdatenVorhanden{
			{
				ZaehlerID:             zaehlerOID1,
				ZaehlerdatenVorhanden: []structs.ZaehlerwertVorhanden{},
			},
			{
				ZaehlerID:             zaehlerOID2,
				ZaehlerdatenVorhanden: []structs.ZaehlerwertVorhanden{},
			},
			{
				ZaehlerID:             zaehlerOID3,
				ZaehlerdatenVorhanden: []structs.ZaehlerwertVorhanden{},
			},
		}

		aktuellesJahr := int32(time.Now().Year())
		for i := structs.ErstesJahr; i <= aktuellesJahr; i++ {
			compareDoc[0].ZaehlerdatenVorhanden = append(compareDoc[0].ZaehlerdatenVorhanden, structs.ZaehlerwertVorhanden{Jahr: i, Vorhanden: false})
			compareDoc[1].ZaehlerdatenVorhanden = append(compareDoc[1].ZaehlerdatenVorhanden, structs.ZaehlerwertVorhanden{Jahr: i, Vorhanden: false})
			compareDoc[2].ZaehlerdatenVorhanden = append(compareDoc[2].ZaehlerdatenVorhanden, structs.ZaehlerwertVorhanden{Jahr: i, Vorhanden: false})
		}

		compareDoc[0].ZaehlerdatenVorhanden[1].Vorhanden = true
		compareDoc[0].ZaehlerdatenVorhanden[2].Vorhanden = true

		compareDoc[1].ZaehlerdatenVorhanden[1].Vorhanden = true

		compareDoc[2].ZaehlerdatenVorhanden[2].Vorhanden = true

		result := binaereZahlerdatenFuerZaehler(alleZahler)
		is.Equal(result, compareDoc)
	})
}
