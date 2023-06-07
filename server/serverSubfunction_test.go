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

	err := database.ConnectDatabase()
	is.NoErr(err)
	defer func() {
		err := database.DisconnectDatabase()
		is.NoErr(err)
	}()

	// Auswertung
	t.Run("TestBinaereZahlerdatenFuerZaehler", TestBinaereZahlerdatenFuerZaehler)

	// Authentication
	t.Run("TestCheckValidSessionToken", TestCheckValidSessionToken)
	t.Run("TestLoescheSessionToken", TestLoescheSessionToken)
	t.Run("TestGeneriereSessionToken", TestGeneriereSessionToken)
	t.Run("TestAuthenticate", TestAuthenticate)
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
				},
			},
		})
	})
}

func TestCheckValidSessionToken(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("checkValidSessionToken: username='test1'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "test1"
		GeneriereSessionToken(username)
		err := checkValidSessionToken(username)
		is.NoErr(err) // Normalfall wirft keine Errors
	})

	t.Run("checkValidSessionToken: username='name@tu-darmstadt.de'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "name@tu-darmstadt.de"
		GeneriereSessionToken(username)
		err := checkValidSessionToken(username)
		is.NoErr(err) // Normalfall wirft keine Errors
	})

	// Errortests
	t.Run("checkValidSessionToken: username='abcdefg' (nicht vorhanden)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "abcdefg"
		err := checkValidSessionToken(username)
		is.Equal(err, structs.ErrNutzerHatKeinenSessiontoken) // Nicht vorhandener Eintrag
	})

	t.Run("checkValidSessionToken: username='test@stud.tu-darmstadt.de' (abgelaufen)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "test@stud.tu-darmstadt.de"
		authMap[username] = session{
			Sessiontoken: "test",
			GenTime:      time.Date(2020, 12, 20, 11, 12, 12, 54, time.UTC),
		}
		err := checkValidSessionToken(username)
		is.Equal(err, structs.ErrAbgelaufenerSessiontoken)
	})
}

func TestLoescheSessionToken(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("loescheSessionToken: username='name@stud.tu-darmstadt.de'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "name@stud.tu-darmstadt.de"
		GeneriereSessionToken(username)
		err := loescheSessionToken(username)
		is.NoErr(err) // Nomalfall wirft keinen Fehler
	})

	t.Run("loescheSessionToken: username='test@stud.tu-darmstadt.de' (abgelaufener)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "test@stud.tu-darmstadt.de"
		authMap[username] = session{
			Sessiontoken: "test",
			GenTime:      time.Date(2020, 12, 20, 11, 12, 12, 54, time.UTC),
		}
		err := loescheSessionToken(username)
		is.NoErr(err) // Normalfall wirft keinen Fehler
	})

	// Fehlerfall
	t.Run("loescheSessionToken: username='nichtvorhanden' (nicht vorhanden)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "nichtvorhanden"
		err := loescheSessionToken(username)
		is.Equal(err, structs.ErrNutzerHatKeinenSessiontoken)
	})
}

func TestGeneriereSessionToken(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("GeneriereSessionToken: username='felix@stud.tu-darmstadt.de'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "felix@stud.tu-darmstadt.de"
		token := GeneriereSessionToken(username)
		is.Equal(token, authMap[username].Sessiontoken)
	})

	t.Run("GeneriereSessionToken: username='repeatedEntry'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "repeatedEntry"
		GeneriereSessionToken(username)
		token := GeneriereSessionToken(username) // Ueberschreibe alten Token
		is.Equal(token, authMap[username].Sessiontoken)
	})
}

func TestAuthenticate(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("Authenticate: username='anton@tobi.com'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "anton@tobi.com"
		token := GeneriereSessionToken(username)
		err := Authenticate(username, token)

		is.NoErr(err) // Im Normalfall wird kein Fehler geworfen
	})

	// Errorfall
	t.Run("Authenticate: username='test123' token='keinToken' (kein Token)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "test123"
		var token = "keinToken"
		err := Authenticate(username, token)
		is.Equal(err, structs.ErrNutzerHatKeinenSessiontoken)
	})

	t.Run("Authenticate: username='anton@tobi.com' token='test' (abgelaufener Token)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "anton@tobi.com"
		var token = "test"
		authMap[username] = session{
			Sessiontoken: token,
			GenTime:      time.Date(2020, 12, 20, 11, 12, 12, 54, time.UTC),
		}
		err := Authenticate(username, token)
		is.Equal(err, structs.ErrAbgelaufenerSessiontoken)
	})

	t.Run("Authenticate: username='anton@tobi.com token='falscherToken' (falscher Token)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var username = "anton@tobi.com"
		var token = "falscherToken"

		GeneriereSessionToken(username)
		err := Authenticate(username, token)
		is.Equal(err, structs.ErrFalscherSessiontoken)
	})
}
