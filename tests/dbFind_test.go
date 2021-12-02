package tests

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/matryer/is"
	"io"
	"testing"
	"time"
)

func TestFind(t *testing.T) {
	database.ConnectDatabase()
	defer database.DisconnectDatabase()

	t.Run("TestITGeareteFind", TestITGeraeteFind)
	t.Run("TestEnergieversorgungFind", TestEnergieversorgungFind)
	t.Run("TestDienstreisenFind", TestDienstreisenFind)
	t.Run("TestGebaeudeFind", TestGebaeudeFind)
	t.Run("TestKaelteFind", TestKaeltezaehlerFind)
	t.Run("TestWaermeFind", TestWaermezaehlerFind)
	t.Run("TestStromFind", TestStromzaehlerFind)
}

func TestITGeraeteFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("ITGeraeteFind: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ITGeraeteFind(1)

		is.NoErr(err)                                                                                                                                               // Error seitens der Datenbank
		is.Equal(data, database.ITGeraete{IDITGerate: 1, Kategorie: "Notebooks", CO2FaktorGesamt: 588000, CO2FaktorJahr: 147000, Einheit: "g/Stueck", Revision: 1}) // Überprüfung des zurückgelieferten Elements
	})

	t.Run("ITGeraeteFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ITGeraeteFind(0)

		is.Equal(err, io.EOF)                // Datenbank gibt EOF error
		is.Equal(data, database.ITGeraete{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
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

	t.Run("DienstreisenFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.DienstreisenFind(0)

		is.Equal(err, io.EOF)                   // Datenbank gibt EOF error
		is.Equal(data, database.Dienstreisen{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
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

	t.Run("EnergieversorgungFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.EnergieversorgungFind(0)

		is.Equal(err, io.EOF)                        // Datenbank gibt EOF error
		is.Equal(data, database.Energieversorgung{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestGebaeudeFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("GebaeudeFind: ID = 1101", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1101)

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

	t.Run("GebaeudeFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(0)

		is.Equal(err, io.EOF)               // Datenbank gibt EOF error
		is.Equal(data, database.Gebaeude{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestKaeltezaehlerFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("KaeltezaehlerFind: ID = 4023", func(t *testing.T) {
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
					Zeitstempel: time.Date(2021, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
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
			},
			Einheit:     "MWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{3101},
		})
	})

	t.Run("KaeltezaehlerFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.KaeltezaehlerFind(0)

		is.Equal(err, io.EOF)              // Datenbank gibt EOF error
		is.Equal(data, database.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestWaermezaehlerFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("WearmeFind: ID = 2107", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.WaermezaehlerFind(2107)
		location, err := time.LoadLocation("Etc/GMT")

		is.NoErr(err) // Error seitens der Datenbank
		is.Equal(data, database.Zaehler{Zaehlertyp: "Waerme",
			PKEnergie:   2107,
			Bezeichnung: " 2101,2102,2108 Waerme Gruppenzaehler",
			Zaehlerdaten: []database.Zaehlerwerte{
				{
					Wert:        788.66,
					Zeitstempel: time.Date(2020, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        794.8,
					Zeitstempel: time.Date(2019, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        736.9,
					Zeitstempel: time.Date(2018, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "MWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{2101, 2102, 2108},
		})
	})

	t.Run("WaermezaehlerFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.WaermezaehlerFind(0)

		is.Equal(err, io.EOF)              // Datenbank gibt EOF error
		is.Equal(data, database.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestStromzaehlerFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("WearmeFind: ID = 5967", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.StromzaehlerFind(5967)
		location, err := time.LoadLocation("Etc/GMT")

		is.NoErr(err) // Error seitens der Datenbank
		is.Equal(data, database.Zaehler{Zaehlertyp: "Strom",
			PKEnergie:   5967,
			Bezeichnung: "2201 Strom Hauptzaehler",
			Zaehlerdaten: []database.Zaehlerwerte{
				{
					Wert:        126048.9,
					Zeitstempel: time.Date(2021, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        0.0,
					Zeitstempel: time.Date(2020, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        0.0,
					Zeitstempel: time.Date(2019, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        0.0,
					Zeitstempel: time.Date(2018, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "kWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{2201},
		})
	})

	t.Run("StromzaehlerFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.StromzaehlerFind(0)

		is.Equal(err, io.EOF)              // Datenbank gibt EOF error
		is.Equal(data, database.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}
