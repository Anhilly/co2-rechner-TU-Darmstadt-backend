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
	t.Run("TestEnergieversorgungFind", TestEnergieversorgungFind)
	t.Run("TestDienstreisenFind", TestDienstreisenFind)
	t.Run("TestGebaeudeFind", TestGebaeudeFind)
}

func TestITGeraeteFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("ITGeraeteFind: Test ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ITGeraeteFind(1)

		is.NoErr(err)                                                                                                                                               // Error seitens der Datenbank
		is.Equal(data, database.ITGeraete{IDITGerate: 1, Kategorie: "Notebooks", CO2FaktorGesamt: 588000, CO2FaktorJahr: 147000, Einheit: "g/Stueck", Revision: 1}) // Überprüfung des zurrückgelieferten Elements
	})

}

func TestDienstreisenFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("DienstreisenFind: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.DienstreisenFind(1)

		is.NoErr(err)                                                                                                                                             // Error seitens der Datenbank
		is.Equal(data, database.Dienstreisen{IDDienstreisen: 1, Medium: "Bahn", Einheit: "g/Pkm", Revision: 1, CO2Faktor: []database.CO2Dienstreisen{{Wert: 8}}}) // Überprüfung des zurrückgelieferten Elements
	})
}

func TestEnergieversorgungFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("EnergieversorgungFind: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.EnergieversorgungFind(1)

		is.NoErr(err)                                                                                                                                                                         // Error seitens der Datenbank
		is.Equal(data, database.Energieversorgung{
			IDEnergieversorgung: 1,
			Kategorie: "Fernwaerme",
			Einheit: "g/kWh",
			Revision: 1,
			CO2Faktor: []database.CO2Energie{{Wert: 144, Jahr: 2020}}}) // Überprüfung des zurrückgelieferten Elements
	})
}

func TestGebaeudeFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("GebaeudeFind: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1102)

		is.NoErr(err) // Error seitens der Datenbank
		is.Equal(data, database.Gebaeude{
			Nr:          1101,
			Bezeichnung: "Universitaetszentrum, karo 5, Audimax",
			Flaeche: database.GebaeudeFlaeche{
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
			WaermeRef:   []int32{2084},
			StromRef:    []int32{}}) // Überprüfung des zurrückgelieferten Elements
	})
}
