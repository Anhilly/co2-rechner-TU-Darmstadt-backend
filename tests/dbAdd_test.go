package tests

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"io"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
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

	t.Run("TestEnergieversorgungAddFaktor", TestEnergieversorgungAddFaktor)
	t.Run("TestWaermezaehlerAddZaehlerdaten", TestWaermezaehlerAddZaehlerdaten)
	t.Run("TestStromzaehlerAddZaehlerdaten", TestStromzaehlerAddZaehlerdaten)
	t.Run("TestKaeltezaehlerAddZaehlerdaten", TestKaeltezaehlerAddZaehlerdaten)
}

func TestEnergieversorgungAddFaktor(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("EnergieversorgungAddFaktor: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddCO2Faktor{
			IDEnergieversorgung: 1,
			Wert:                400,
			Jahr:                2030,
		}

		err := database.EnergieversorgungAddFaktor(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.EnergieversorgungFind(data.IDEnergieversorgung)

		is.Equal(updatedDoc, structs.Energieversorgung{
			IDEnergieversorgung: 1,
			Kategorie:           "Fernwaerme",
			Einheit:             "g/kWh",
			Revision:            1,
			CO2Faktor:           []structs.CO2Energie{{Wert: 144, Jahr: 2020}, {Wert: 400, Jahr: 2030}},
		}) // Ueberpruefung des geaenderten Elementes
	})

	// Errortests
	t.Run("EnergieversorgungAddFaktor: ID nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddCO2Faktor{
			IDEnergieversorgung: 0,
			Wert:                400,
			Jahr:                2030,
		}

		err := database.EnergieversorgungAddFaktor(data)
		is.Equal(err, io.EOF) // Datenbank wirft EOF-Error
	})

	t.Run("EnergieversorgungAddFaktor: Jahr schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddCO2Faktor{
			IDEnergieversorgung: 1,
			Wert:                400,
			Jahr:                2020,
		}

		err := database.EnergieversorgungAddFaktor(data)
		is.Equal(err, database.ErrJahrVorhanden) // Funktion wirft ErrJahrVorhanden
	})
}

func TestWaermezaehlerAddZaehlerdaten(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("WaermezaehlerAddZaehlerdaten: ID = 2017", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie: 2107,
			Wert:      1000.0,
			Jahr:      3000,
		}
		location, _ := time.LoadLocation("Etc/GMT")

		err := database.WaermezaehlerAddZaehlerdaten(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.WaermezaehlerFind(data.PKEnergie)

		is.Equal(updatedDoc, structs.Zaehler{Zaehlertyp: "Waerme",
			PKEnergie:   2107,
			Bezeichnung: " 2101,2102,2108 Waerme Gruppenzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{
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
				{
					Wert:        1000.0,
					Zeitstempel: time.Date(3000, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "MWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{2101, 2102, 2108},
		}) // Ueberpruefung des geaenderten Elementes
	})

	// Errortests
	t.Run("WaermezaehlerAddZaehlerdaten: ID nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie: 0,
			Wert:      1000.0,
			Jahr:      3000,
		}

		err := database.WaermezaehlerAddZaehlerdaten(data)
		is.Equal(err, io.EOF) // Datenbank wirft EOF-Error
	})

	t.Run("WaermezaehlerAddZaehlerdaten: Jahr schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie: 2107,
			Wert:      1000.0,
			Jahr:      2020,
		}

		err := database.WaermezaehlerAddZaehlerdaten(data)
		is.Equal(err, database.ErrJahrVorhanden) // Funktion wirft ErrJahrVorhanden
	})
}

func TestStromzaehlerAddZaehlerdaten(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("StromzaehlerAddZaehlerdaten: ID = 5967", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie: 5967,
			Wert:      1000.0,
			Jahr:      3000,
		}
		location, _ := time.LoadLocation("Etc/GMT")

		err := database.StromzaehlerAddZaehlerdaten(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.StromzaehlerFind(data.PKEnergie)

		is.Equal(updatedDoc, structs.Zaehler{Zaehlertyp: "Strom",
			PKEnergie:   5967,
			Bezeichnung: "2201 Strom Hauptzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{
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
				{
					Wert:        1000.0,
					Zeitstempel: time.Date(3000, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "kWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{2201},
		}) // Ueberpruefung des geaenderten Elementes
	})

	// Errortests
	t.Run("StromzaehlerAddZaehlerdaten: ID nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie: 0,
			Wert:      1000.0,
			Jahr:      3000,
		}

		err := database.StromzaehlerAddZaehlerdaten(data)
		is.Equal(err, io.EOF) // Datenbank wirft EOF-Error
	})

	t.Run("StromzaehlerAddZaehlerdaten: Jahr schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie: 5967,
			Wert:      1000.0,
			Jahr:      2020,
		}

		err := database.StromzaehlerAddZaehlerdaten(data)
		is.Equal(err, database.ErrJahrVorhanden) // Funktion wirft ErrJahrVorhanden
	})
}

func TestKaeltezaehlerAddZaehlerdaten(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("KaeltezaehlerAddZaehlerdaten: ID = 4023", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie: 4023,
			Wert:      1000.0,
			Jahr:      3000,
		}
		location, _ := time.LoadLocation("Etc/GMT")

		err := database.KaeltezaehlerAddZaehlerdaten(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.KaeltezaehlerFind(data.PKEnergie)

		is.Equal(updatedDoc, structs.Zaehler{Zaehlertyp: "Kaelte",
			PKEnergie:   4023,
			Bezeichnung: "3101 Kaelte Hauptzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{
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
				{
					Wert:        1000.0,
					Zeitstempel: time.Date(3000, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "MWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{3101},
		}) // Ueberpruefung des geaenderten Elementes
	})

	// Errortests
	t.Run("KaeltezaehlerAddZaehlerdaten: ID nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie: 0,
			Wert:      1000.0,
			Jahr:      3000,
		}

		err := database.KaeltezaehlerAddZaehlerdaten(data)
		is.Equal(err, io.EOF) // Datenbank wirft EOF-Error
	})

	t.Run("KaeltezaehlerAddZaehlerdaten: Jahr schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie: 4023,
			Wert:      1000.0,
			Jahr:      2020,
		}

		err := database.KaeltezaehlerAddZaehlerdaten(data)
		is.Equal(err, database.ErrJahrVorhanden) // Funktion wirft ErrJahrVorhanden
	})
}
