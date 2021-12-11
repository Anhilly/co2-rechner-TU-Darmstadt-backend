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

	t.Run("TestGebaeudeInsert", TestGebaeudeInsert)
}

func TestGebaeudeInsert(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("GebaeudeInsert: Nr = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertGebaeude{
			Nr:          0,
			Bezeichnung: "Testgebaeude",
			Flaeche: structs.GebaeudeFlaeche{
				HNF:     1000.0,
				NNF:     1000.0,
				NGF:     1000.0,
				FF:      1000.0,
				VF:      1000.0,
				FreiF:   1000.0,
				GesamtF: 1000.0,
			},
		}

		err := database.GebaeudeInsert(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		insertedDoc, err := database.GebaeudeFind(data.Nr)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(insertedDoc, structs.Gebaeude{
			Nr:          0,
			Bezeichnung: "Testgebaeude",
			Flaeche: structs.GebaeudeFlaeche{
				HNF:     1000.0,
				NNF:     1000.0,
				NGF:     1000.0,
				FF:      1000.0,
				VF:      1000.0,
				FreiF:   1000.0,
				GesamtF: 1000.0,
			},
			Einheit:     "m^2",
			Spezialfall: 1,
			Revision:    1,
			KaelteRef:   []int32{},
			WaermeRef:   []int32{},
			StromRef:    []int32{},
		}) // Ueberpruefung des geaenderten Elementes
	})

	// Errortest
	t.Run("GebaeudeInsert: Nr = 1101 schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertGebaeude{
			Nr:          0,
			Bezeichnung: "Testgebaeude",
			Flaeche: structs.GebaeudeFlaeche{
				HNF:     1000.0,
				NNF:     1000.0,
				NGF:     1000.0,
				FF:      1000.0,
				VF:      1000.0,
				FreiF:   1000.0,
				GesamtF: 1000.0,
			},
		}

		err := database.GebaeudeInsert(data)
		is.Equal(err, database.ErrGebaeudeVorhanden) // Funktion wirft ErrGebaeudeVorhanden
	})
}
