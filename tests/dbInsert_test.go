package tests

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/mongo"
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
	t.Run("TestZaehlerInsert", TestZaehlerInsert)
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

func TestZaehlerInsert(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("ZaehlerInsert: Waermezaehler, ID = 190", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			PKEnergie:           19,
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{1101},
			IDEnergieversorgung: 1,
		}

		err := database.ZaehlerInsert(data)
		is.NoErr(err)

		neuerZaehler, err := database.ZaehlerFind(data.PKEnergie, data.IDEnergieversorgung)
		is.NoErr(err)
		is.Equal(neuerZaehler, structs.Zaehler{
			Zaehlertyp:   "Waerme",
			PKEnergie:    19,
			Bezeichnung:  "Testzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "kWh",
			Spezialfall:  1,
			Revision:     1,
			GebaeudeRef:  []int32{1101},
		})

		updatedGebaeude, err := database.GebaeudeFind(1101)
		is.NoErr(err)
		is.Equal(updatedGebaeude, structs.Gebaeude{
			Nr:          1101,
			Bezeichnung: "Universitaetszentrum, karo 5, Audimax",
			Flaeche: structs.GebaeudeFlaeche{
				HNF:     6395.56,
				NNF:     3081.85,
				NGF:     15365.03,
				FF:      5539.21,
				VF:      5887.62,
				FreiF:   96.57,
				GesamtF: 21000.81,
			},
			Einheit:     "m^2",
			Spezialfall: 1,
			Revision:    1,
			KaelteRef:   []int32{},
			WaermeRef:   []int32{2084, 19},
			StromRef:    []int32{},
		})
	})

	// Errortests
	t.Run("ZaehlerInsert: keine Gebaeudereferenzen", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			PKEnergie:           190,
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{},
			IDEnergieversorgung: 1,
		}

		err := database.ZaehlerInsert(data)
		is.Equal(err, structs.ErrFehlendeGebaeuderef)
	})

	t.Run("ZaehlerInsert: Zaehler vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			PKEnergie:           2107,
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{190},
			IDEnergieversorgung: 1,
		}

		err := database.ZaehlerInsert(data)
		is.Equal(err, structs.ErrZaehlerVorhanden)
	})

	t.Run("ZaehlerInsert: ungueltige Gebaeudereferenz", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			PKEnergie:           14,
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{12},
			IDEnergieversorgung: 1,
		}

		err := database.ZaehlerInsert(data)
		is.Equal(err, mongo.ErrNoDocuments)
	})

	t.Run("ZaehlerInsert: IDEnergieversorgung = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			PKEnergie:           15,
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{12},
			IDEnergieversorgung: 0,
		}

		err := database.ZaehlerInsert(data)
		is.Equal(err, database.ErrIDEnergieversorgungNichtVorhanden)
	})
}
