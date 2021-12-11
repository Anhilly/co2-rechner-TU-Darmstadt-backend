package tests

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"testing"
)

func TestInsert(t *testing.T) {
	is := is.NewRelaxed(t)

	dir, err := database.CreateDump("TestAdd")
	is.NoErr(err)

	fmt.Println(dir)

	err = database.ConnectDatabase()
	is.NoErr(err)

	defer func(dir string) {
		err := database.DisconnectDatabase()
		is.NoErr(err)
		err = database.RestoreDump(dir)
		is.NoErr(err)
	}(dir)

	t.Run("TestNutzerdatenInsert", TestNutzerdatenInsert)
}

func TestNutzerdatenInsert(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("NutzerdatenInsert: {email = 'testingEmailPlsDontUse' password='verysecurepassword'} (nicht vorhanden)", func(t *testing.T) {
		is := is.NewRelaxed(t)
		testData := structs.AnmeldungReq{
			Email:    "testingEmailPlsDontUse",
			Passwort: "verysecurepassword",
		}
		err := database.NutzerdatenInsert(testData)
		is.NoErr(err) //Kein Fehler wird geworfen
	})

	//Errorfall
	t.Run("NutzerdatenInsert: {email = 'anton@tobi.com' password='verysecurepassword'} (vorhanden)", func(t *testing.T) {
		is := is.NewRelaxed(t)
		testData := structs.AnmeldungReq{
			Email:    "anton@tobi.com",
			Passwort: "verysecurepassword",
		}
		err := database.NutzerdatenInsert(testData)
		is.Equal(err, database.ErrInsertExistingAccount) //Dateneintrag existiert bereits
	})
}
