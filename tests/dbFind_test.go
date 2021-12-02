package tests

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/matryer/is"
	"testing"
	"time"
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
		is.Equal(data, database.ITGeraete{IDITGerate: 1, Kategorie: "Notebooks", CO2FaktorGesamt: 588000, CO2FaktorJahr: 147000, Einheit: "g/Stueck", Revision: 1}) // Überprüfung des zurückgelieferten Elements
	})

}

func TestDienstreisenFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("DienstreisenFind: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.DienstreisenFind(1)

		is.NoErr(err)                                                                                                                                             // Error seitens der Datenbank
		is.Equal(data, database.Dienstreisen{IDDienstreisen: 1, Medium: "Bahn", Einheit: "g/Pkm", Revision: 1, CO2Faktor: []database.CO2Dienstreisen{{Wert: 8}}}) // Überprüfung des zurückgelieferten Elements
	})
}

func TestEnergieversorgungFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("EnergieversorgungFind: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.EnergieversorgungFind(1)

		is.NoErr(err) // Error seitens der Datenbank
		is.Equal(data, database.Energieversorgung{
			IDEnergieversorgung: 1,
			Kategorie:           "Fernwaerme",
			Einheit:             "g/kWh",
			Revision:            1,
			CO2Faktor:           []database.CO2Energie{{Wert: 144, Jahr: 2020}}}) // Überprüfung des zurückgelieferten Elements
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
			StromRef:    []int32{}}) // Überprüfung des zurückgelieferten Elements
	})
}

func TestKaeltezaehlerFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("KaelteFind: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.KaeltezaehlerFind(4023)
		location, err := time.LoadLocation("Etc/GMT")

		is.NoErr(err) // Error seitens der Datenbank
		is.Equal(data, database.Zaehler{Zaehlertyp: "Kaelte",
			PKEnergie:   4023,
			Bezeichnung: "3101 Kaelte Hauptzaehler",
			Zaehlerdaten: []database.Zaehlerwerte{
				{
					Wert:        311.06,
					Zeitstempel: time.Date(2021, time.January, 01, 0, 0, 0, 0, location),
				},
				{
					Wert:        414.61,
					Zeitstempel: time.Date(2020, time.January, 01, 0, 0, 0, 0, location),
				},
				{
					Wert:        555.3,
					Zeitstempel: time.Date(2019, time.January, 01, 0, 0, 0, 0, location),
				},
				{
					Wert:        169.59,
					Zeitstempel: time.Date(2018, time.January, 01, 0, 0, 0, 0, location),
				},
			}})
	})
}
