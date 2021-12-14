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
}

func TestcheckValidSessionToken(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	//Normalfall
	t.Run("checkValidSessionToken: email='test1'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "test1"
		generiereSessionToken(email)
		err := checkValidSessionToken(email)
		is.NoErr(err) // Normalfall wirft keine Errors
	})

	t.Run("checkValidSessionToken: email='name@tu-darmstadt.de'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "name@tu-darmstadt.de"
		generiereSessionToken(email)
		err := checkValidSessionToken(email)
		is.NoErr(err) // Normalfall wirft keine Errors
	})

	//Errortests
	t.Run("checkValidSessionToken: email='abcdefg' (nicht vorhanden)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "abcdefg"
		err := checkValidSessionToken(email)
		is.Equal(err, structs.ErrNutzerHatKeinenSessiontoken) // Nicht vorhandener Eintrag
	})

	t.Run("checkValidSessionToken: email='test@stud.tu-darmstadt.de' (abgelaufen)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "test@stud.tu-darmstadt.de"
		AuthMap[email] = Session{
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
		generiereSessionToken(email)
		err := loescheSessionToken(email)
		is.NoErr(err) // Nomalfall wirft keinen Fehler
	})

	t.Run("loescheSessionToken: email='test@stud.tu-darmstadt.de' (abgelaufener)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "test@stud.tu-darmstadt.de"
		AuthMap[email] = Session{
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

	//Normalfall
	t.Run("generiereSessionToken: email='felix@stud.tu-darmstadt.de'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "felix@stud.tu-darmstadt.de"
		token := generiereSessionToken(email)
		is.Equal(token, AuthMap[email].Sessiontoken)
	})

	t.Run("generiereSessionToken: email='repeatedEntry'", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var email = "repeatedEntry"
		generiereSessionToken(email)
		token := generiereSessionToken(email) //Überschreibe alten Token
		is.Equal(token, AuthMap[email].Sessiontoken)
	})
}
