package server

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
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

	t.Run("TestcheckValidSessionToken", TestcheckValidSessionToken)
	t.Run("TestloescheSessionToken", TestloescheSessionToken)
	t.Run("TestgeneriereSessionToken", TestgeneriereSessionToken)
	t.Run("TestAuthenticate", TestAuthenticate)
}

func TestcheckValidSessionToken(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("checkValidSessionToken: email='test1'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "test1"
		GeneriereSessionToken(email)
		err := checkValidSessionToken(email)
		is.NoErr(err) // Normalfall wirft keine Errors
	})

	t.Run("checkValidSessionToken: email='name@tu-darmstadt.de'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "name@tu-darmstadt.de"
		GeneriereSessionToken(email)
		err := checkValidSessionToken(email)
		is.NoErr(err) // Normalfall wirft keine Errors
	})

	// Errortests
	t.Run("checkValidSessionToken: email='abcdefg' (nicht vorhanden)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "abcdefg"
		err := checkValidSessionToken(email)
		is.Equal(err, structs.ErrNutzerHatKeinenSessiontoken) // Nicht vorhandener Eintrag
	})

	t.Run("checkValidSessionToken: email='test@stud.tu-darmstadt.de' (abgelaufen)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "test@stud.tu-darmstadt.de"
		authMap[email] = session{
			Sessiontoken: "test",
			GenTime:      time.Date(2020, 12, 20, 11, 12, 12, 54, time.UTC),
		}
		err := checkValidSessionToken(email)
		is.Equal(err, structs.ErrAbgelaufenerSessiontoken)
	})
}

func TestloescheSessionToken(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("loescheSessionToken: email='name@stud.tu-darmstadt.de'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "name@stud.tu-darmstadt.de"
		GeneriereSessionToken(email)
		err := loescheSessionToken(email)
		is.NoErr(err) // Nomalfall wirft keinen Fehler
	})

	t.Run("loescheSessionToken: email='test@stud.tu-darmstadt.de' (abgelaufener)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "test@stud.tu-darmstadt.de"
		authMap[email] = session{
			Sessiontoken: "test",
			GenTime:      time.Date(2020, 12, 20, 11, 12, 12, 54, time.UTC),
		}
		err := loescheSessionToken(email)
		is.NoErr(err) // Normalfall wirft keinen Fehler
	})

	// Fehlerfall
	t.Run("loescheSessionToken: email='nichtvorhanden' (nicht vorhanden)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "nichtvorhanden"
		err := loescheSessionToken(email)
		is.Equal(err, structs.ErrNutzerHatKeinenSessiontoken)
	})
}

func TestgeneriereSessionToken(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("GeneriereSessionToken: email='felix@stud.tu-darmstadt.de'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "felix@stud.tu-darmstadt.de"
		token := GeneriereSessionToken(email)
		is.Equal(token, authMap[email].Sessiontoken)
	})

	t.Run("GeneriereSessionToken: email='repeatedEntry'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "repeatedEntry"
		GeneriereSessionToken(email)
		token := GeneriereSessionToken(email) // Ueberschreibe alten Token
		is.Equal(token, authMap[email].Sessiontoken)
	})
}

func TestAuthenticate(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("Authenticate: email='anton@tobi.com'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "anton@tobi.com"
		token := GeneriereSessionToken(email)
		err := Authenticate(email, token)

		is.NoErr(err) // Im Normalfall wird kein Fehler geworfen
	})

	// Errorfall
	t.Run("Authenticate: email='test123' token='keinToken' (kein Token)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "test123"
		var token = "keinToken"
		err := Authenticate(email, token)
		is.Equal(err, structs.ErrNutzerHatKeinenSessiontoken)
	})

	t.Run("Authenticate: email='anton@tobi.com' token='test' (abgelaufener Token)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "anton@tobi.com"
		var token = "test"
		authMap[email] = session{
			Sessiontoken: token,
			GenTime:      time.Date(2020, 12, 20, 11, 12, 12, 54, time.UTC),
		}
		err := Authenticate(email, token)
		is.Equal(err, structs.ErrAbgelaufenerSessiontoken)
	})

	t.Run("Authenticate: email='anton@tobi.com token='falscherToken' (falscher Token)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "anton@tobi.com"
		var token = "falscherToken"

		GeneriereSessionToken(email)
		err := Authenticate(email, token)
		is.Equal(err, structs.ErrFalscherSessiontoken)
	})
}
