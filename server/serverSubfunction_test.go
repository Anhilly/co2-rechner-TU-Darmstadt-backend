package server

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
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

		location, _ := time.LoadLocation("Etc/GMT")
		alleZahler := []structs.ZaehlerUndZaehlerdaten{
			{
				PKEnergie: 101,
				Zaehlerdaten: []structs.Zaehlerwerte{
					{
						Wert:        414.61,
						Zeitstempel: time.Date(2020, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
					{
						Wert:        555.3,
						Zeitstempel: time.Date(2019, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
					{
						Wert:        169.59,
						Zeitstempel: time.Date(2018, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
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
		}

		result := binaereZahlerdatenFuerZaehler(alleZahler)
		is.Equal(result, []structs.ZaehlerUndZaehlerdatenVorhanden{
			{
				PKEnergie: 101,
				ZaehlerdatenVorhanden: []structs.ZaehlerwertVorhanden{
					{Jahr: 2018, Vorhanden: true},
					{Jahr: 2019, Vorhanden: true},
					{Jahr: 2020, Vorhanden: true},
					{Jahr: 2021, Vorhanden: true},
					{Jahr: 2022, Vorhanden: true},
					{Jahr: 2023, Vorhanden: false},
					{Jahr: 2024, Vorhanden: false},
				},
			},
		})
	})

	t.Run("binaereZahlerdatenFuerZaehler: ein Zaehler mit unvollstaendingen Daten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		location, _ := time.LoadLocation("Etc/GMT")
		alleZahler := []structs.ZaehlerUndZaehlerdaten{
			{
				PKEnergie: 101,
				Zaehlerdaten: []structs.Zaehlerwerte{
					{
						Wert:        169.59,
						Zeitstempel: time.Date(2018, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
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
		}

		result := binaereZahlerdatenFuerZaehler(alleZahler)
		is.Equal(result, []structs.ZaehlerUndZaehlerdatenVorhanden{
			{
				PKEnergie: 101,
				ZaehlerdatenVorhanden: []structs.ZaehlerwertVorhanden{
					{Jahr: 2018, Vorhanden: true},
					{Jahr: 2019, Vorhanden: false},
					{Jahr: 2020, Vorhanden: false},
					{Jahr: 2021, Vorhanden: true},
					{Jahr: 2022, Vorhanden: true},
					{Jahr: 2023, Vorhanden: false},
					{Jahr: 2024, Vorhanden: false},
				},
			},
		})
	})

	t.Run("binaereZahlerdatenFuerZaehler: mehrere Zaehler mit unvollstaendingen Daten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		location, _ := time.LoadLocation("Etc/GMT")
		alleZahler := []structs.ZaehlerUndZaehlerdaten{
			{
				PKEnergie: 101,
				Zaehlerdaten: []structs.Zaehlerwerte{
					{
						Wert:        169.59,
						Zeitstempel: time.Date(2018, time.January, 01, 0, 0, 0, 0, location).UTC(),
					},
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
				PKEnergie: 103,
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
				PKEnergie: 107,
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

		result := binaereZahlerdatenFuerZaehler(alleZahler)
		is.Equal(result, []structs.ZaehlerUndZaehlerdatenVorhanden{
			{
				PKEnergie: 101,
				ZaehlerdatenVorhanden: []structs.ZaehlerwertVorhanden{
					{Jahr: 2018, Vorhanden: true},
					{Jahr: 2019, Vorhanden: false},
					{Jahr: 2020, Vorhanden: false},
					{Jahr: 2021, Vorhanden: true},
					{Jahr: 2022, Vorhanden: true},
					{Jahr: 2023, Vorhanden: false},
					{Jahr: 2024, Vorhanden: false},
				},
			},
			{
				PKEnergie: 103,
				ZaehlerdatenVorhanden: []structs.ZaehlerwertVorhanden{
					{Jahr: 2018, Vorhanden: true},
					{Jahr: 2019, Vorhanden: false},
					{Jahr: 2020, Vorhanden: false},
					{Jahr: 2021, Vorhanden: true},
					{Jahr: 2022, Vorhanden: false},
					{Jahr: 2023, Vorhanden: false},
					{Jahr: 2024, Vorhanden: false},
				},
			},
			{
				PKEnergie: 107,
				ZaehlerdatenVorhanden: []structs.ZaehlerwertVorhanden{
					{Jahr: 2018, Vorhanden: false},
					{Jahr: 2019, Vorhanden: false},
					{Jahr: 2020, Vorhanden: false},
					{Jahr: 2021, Vorhanden: false},
					{Jahr: 2022, Vorhanden: true},
					{Jahr: 2023, Vorhanden: false},
					{Jahr: 2024, Vorhanden: false},
				},
			},
		})
	})
}
