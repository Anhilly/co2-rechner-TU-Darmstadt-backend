package tests

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/matryer/is"
	"testing"
)

func TestFind(t *testing.T) {
	//is := is.NewRelaxed(t)

	database.ConnectDatabase()
	defer database.DisconnectDatabase()

	t.Run("TestITGeareteFind", TestITGeraeteFind)

}

func TestITGeraeteFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("Test ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ITGeraeteFind(1)

		is.NoErr(err)                                                                                                                                               // Error seitens der Datenbank
		is.Equal(data, database.ITGeraete{IDITGerate: 1, Kategorie: "Notebooks", CO2FaktorGesamt: 588000, CO2FaktorJahr: 147000, Einheit: "g/Stueck", Revision: 1}) // Überprüfung des zurrückgelieferten Elements
	})

}
