package tests

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	is := is.NewRelaxed(t)

	err := database.ConnectDatabase("dev")
	is.NoErr(err)

	dir, err := database.CreateDump("TestAdd")
	is.NoErr(err)

	fmt.Println(dir)

	defer func(dir string) {
		err := database.DisconnectDatabase()
		is.NoErr(err)
		err = database.RestoreDump(dir)
		is.NoErr(err)
		err = database.RemoveDump(dir)
		is.NoErr(err)
	}(dir)

	t.Run("TestEnergieversorgungAddFaktor", TestEnergieversorgungAddFaktor)
	t.Run("TestZaehlerAddZaehlerdaten", TestZaehlerAddZaehlerdaten)
	t.Run("TestZaehlerAddStandardZaehlerdaten", TestZaehlerAddStandardZaehlerdaten)
	t.Run("TestGebaeudeAddZaehlerref", TestGebaeudeAddZaehlerref)
	t.Run("TestGebaeudeAddVersorger", TestGebaeudeAddVersorger)
	t.Run("TestGebaeudeAddStandardVersorger", TestGebaeudeAddStandardVersorger)
	t.Run("TestNutzerdatenAddUmfrageref", TestNutzerdatenAddUmfrageref)
	t.Run("TestUmfrageAddMitarbeiterUmfrageRef", TestUmfrageAddMitarbeiterUmfrageRef)
}

func TestEnergieversorgungAddFaktor(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("EnergieversorgungAddFaktor: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddCO2Faktor{
			IDEnergieversorgung: 1,
			IDVertrag:           1,
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
			CO2Faktor: []structs.CO2Energie{
				{Jahr: 2020, Vertraege: []structs.CO2FaktorVetrag{
					{IDVertrag: 1, Wert: 144},
				}},
				{Jahr: 2021, Vertraege: []structs.CO2FaktorVetrag{
					{IDVertrag: 1, Wert: 125},
				}},
				{Jahr: 2030, Vertraege: []structs.CO2FaktorVetrag{
					{IDVertrag: 1, Wert: 400},
				}},
			},
		}) // Ueberpruefung des geaenderten Elementes
	})

	t.Run("EnergieversorgungAddFaktor: IDVertrag = 2", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddCO2Faktor{
			IDEnergieversorgung: 2,
			IDVertrag:           2,
			Wert:                400,
			Jahr:                2030,
		}

		err := database.EnergieversorgungAddFaktor(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.EnergieversorgungFind(data.IDEnergieversorgung)

		is.Equal(updatedDoc, structs.Energieversorgung{
			IDEnergieversorgung: 2,
			Kategorie:           "Strom",
			Einheit:             "g/kWh",
			Revision:            1,
			CO2Faktor: []structs.CO2Energie{
				{Jahr: 2020, Vertraege: []structs.CO2FaktorVetrag{
					{IDVertrag: 1, Wert: 285},
				}},
				{Jahr: 2021, Vertraege: []structs.CO2FaktorVetrag{
					{IDVertrag: 1, Wert: 357},
				}},
				{Jahr: 2030, Vertraege: []structs.CO2FaktorVetrag{
					{IDVertrag: 2, Wert: 400},
				}},
			},
		}) // Ueberpruefung des geaenderten Elementes
	})

	// Errortests
	t.Run("EnergieversorgungAddFaktor: ID nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddCO2Faktor{
			IDEnergieversorgung: 0,
			IDVertrag:           1,
			Wert:                400,
			Jahr:                2030,
		}

		err := database.EnergieversorgungAddFaktor(data)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})

	t.Run("EnergieversorgungAddFaktor: Jahr schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddCO2Faktor{
			IDEnergieversorgung: 2,
			IDVertrag:           2,
			Wert:                400,
			Jahr:                2030,
		}

		err := database.EnergieversorgungAddFaktor(data)
		is.Equal(err, structs.ErrJahrVorhanden) // Funktion wirft ErrJahrVorhanden
	})

	t.Run("EnergieversorgungAddFaktor: Jahr schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddCO2Faktor{
			IDEnergieversorgung: 1,
			IDVertrag:           1,
			Wert:                400,
			Jahr:                2020,
		}

		err := database.EnergieversorgungAddFaktor(data)
		is.Equal(err, structs.ErrJahrVorhanden) // Funktion wirft ErrJahrVorhanden
	})

	t.Run("EnergieversorgungAddFaktor: IDVertrag invalid", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddCO2Faktor{
			IDEnergieversorgung: 1,
			IDVertrag:           0,
			Wert:                400,
			Jahr:                2020,
		}

		err := database.EnergieversorgungAddFaktor(data)
		is.Equal(err, structs.ErrVertragNichtVorhanden) // Funktion wirft ErrJahrVorhanden
	})
}

func TestZaehlerAddZaehlerdaten(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("ZaehlerAddZaehlerdaten: Waermezaehler, ID = 2017", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie:           2107,
			Wert:                1000.0,
			Jahr:                3001,
			IDEnergieversorgung: 1,
		}
		location, _ := time.LoadLocation("Etc/GMT")

		err := database.ZaehlerAddZaehlerdaten(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.ZaehlerFind(data.PKEnergie, data.IDEnergieversorgung)

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
					Wert:        859.29,
					Zeitstempel: time.Date(2021, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        697.07,
					Zeitstempel: time.Date(2022, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        1000.0,
					Zeitstempel: time.Date(3001, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "MWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{2101, 2102, 2108},
		}) // Ueberpruefung des geaenderten Elementes
	})

	t.Run("ZaehlerAddZaehlerdaten: Stromzaehler, ID = 5967", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie:           5967,
			Wert:                1000.0,
			Jahr:                3001,
			IDEnergieversorgung: 2,
		}
		location, _ := time.LoadLocation("Etc/GMT")

		err := database.ZaehlerAddZaehlerdaten(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.ZaehlerFind(data.PKEnergie, data.IDEnergieversorgung)

		is.Equal(updatedDoc, structs.Zaehler{Zaehlertyp: "Strom",
			PKEnergie:   5967,
			Bezeichnung: "2201 Strom Hauptzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{
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
					Wert:        165440,
					Zeitstempel: time.Date(2021, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        197599.6,
					Zeitstempel: time.Date(2022, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        1000.0,
					Zeitstempel: time.Date(3001, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "kWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{2201},
		}) // Ueberpruefung des geaenderten Elementes
	})

	t.Run("ZaehlerAddZaehlerdaten: Kaeltezaehler, ID = 4023", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie:           4023,
			Wert:                1000.0,
			Jahr:                3001,
			IDEnergieversorgung: 3,
		}
		location, _ := time.LoadLocation("Etc/GMT")

		err := database.ZaehlerAddZaehlerdaten(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.ZaehlerFind(data.PKEnergie, data.IDEnergieversorgung)

		is.Equal(updatedDoc, structs.Zaehler{Zaehlertyp: "Kaelte",
			PKEnergie:   4023,
			Bezeichnung: "3101 Kaelte Hauptzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{
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
					Wert:        380.67,
					Zeitstempel: time.Date(2021, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        370.39,
					Zeitstempel: time.Date(2022, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        1000.0,
					Zeitstempel: time.Date(3001, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "MWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{3101},
		}) // Ueberpruefung des geaenderten Elementes
	})

	t.Run("ZaehlerAddZaehlerdaten: Ueberschreiben von Wert 0.0 bei Zaehler", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie:           5967,
			Wert:                5643,
			Jahr:                2018,
			IDEnergieversorgung: 2,
		}
		location, _ := time.LoadLocation("Etc/GMT")

		err := database.ZaehlerAddZaehlerdaten(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.ZaehlerFind(data.PKEnergie, data.IDEnergieversorgung)

		is.Equal(updatedDoc, structs.Zaehler{Zaehlertyp: "Strom",
			PKEnergie:   5967,
			Bezeichnung: "2201 Strom Hauptzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{
				{
					Wert:        0.0,
					Zeitstempel: time.Date(2020, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        0.0,
					Zeitstempel: time.Date(2019, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        5643.0,
					Zeitstempel: time.Date(2018, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        165440,
					Zeitstempel: time.Date(2021, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        197599.6,
					Zeitstempel: time.Date(2022, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        1000.0,
					Zeitstempel: time.Date(3001, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "kWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{2201},
		}) // Ueberpruefung des geaenderten Elementes
	})

	// Errortests
	t.Run("ZaehlerAddZaehlerdaten: IDEnergieversorgung = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie:           0,
			Wert:                1000.0,
			Jahr:                3000,
			IDEnergieversorgung: 0,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, structs.ErrIDEnergieversorgungNichtVorhanden) // Funktion wirft ErrIDEnergieversorgungNichtVorhanden
	})

	t.Run("ZaehlerAddZaehlerdaten: Waermezaehler, ID nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie:           0,
			Wert:                1000.0,
			Jahr:                3000,
			IDEnergieversorgung: 1,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})

	t.Run("ZaehlerAddZaehlerdaten: Waermezaehler, Jahr schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie:           2107,
			Wert:                1000.0,
			Jahr:                2020,
			IDEnergieversorgung: 1,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, structs.ErrJahrVorhanden) // Funktion wirft ErrJahrVorhanden
	})

	t.Run("ZaehlerAddZaehlerdaten: Stromzaehler, ID nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie:           0,
			Wert:                1000.0,
			Jahr:                3000,
			IDEnergieversorgung: 2,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})

	t.Run("ZaehlerAddZaehlerdaten: Stromzaehler, Jahr schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie:           5967,
			Wert:                1000.0,
			Jahr:                2021,
			IDEnergieversorgung: 2,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, structs.ErrJahrVorhanden) // Funktion wirft ErrJahrVorhanden
	})

	t.Run("ZaehlerAddZaehlerdaten: Kaeltezaehler, ID nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie:           0,
			Wert:                1000.0,
			Jahr:                3000,
			IDEnergieversorgung: 3,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})

	t.Run("ZaehlerAddZaehlerdaten: Kaeltezaehler, Jahr schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			PKEnergie:           4023,
			Wert:                1000.0,
			Jahr:                2020,
			IDEnergieversorgung: 3,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, structs.ErrJahrVorhanden) // Funktion wirft ErrJahrVorhanden
	})
}

func TestZaehlerAddStandardZaehlerdaten(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("AddStandardZaehlerdaten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddStandardZaehlerdaten{
			Jahr: 3000,
		}
		location, _ := time.LoadLocation("Etc/GMT")

		err := database.ZaehlerAddStandardZaehlerdaten(data)
		is.NoErr(err) // Datenbank wirft keinen Fehler

		zaehler, err := database.ZaehlerFind(2107, 1)
		is.NoErr(err) // Datenbank wirft keinen Fehler
		is.Equal(zaehler, structs.Zaehler{Zaehlertyp: "Waerme",
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
					Wert:        859.29,
					Zeitstempel: time.Date(2021, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        697.07,
					Zeitstempel: time.Date(2022, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        1000.0,
					Zeitstempel: time.Date(3001, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
				{
					Wert:        0.0,
					Zeitstempel: time.Date(3000, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "MWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{2101, 2102, 2108},
		}) // Ueberpruefung des geaenderten Elementes)
	})
}

func TestGebaeudeAddZaehlerref(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("GebaeudeAddZaehlerref: ID = 1101, idEnergieversorgung = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var nr int32 = 1101
		var ref int32 = 999
		var idEnergieversorgung int32 = 1

		err := database.GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.GebaeudeFind(nr)
		is.Equal(updatedDoc, structs.Gebaeude{
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
			Stromversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
			},
			Waermeversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
			},
			Kaelteversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
			},
			KaelteRef: []int32{},
			WaermeRef: []int32{2084, 999},
			StromRef:  []int32{26024, 24799},
		}) // Ueberpruefung des zurueckgelieferten Elements
	})

	t.Run("GebaeudeAddZaehlerref: ID = 1102, idEnergieversorgung = 2", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var nr int32 = 1102
		var ref int32 = 999
		var idEnergieversorgung int32 = 2

		err := database.GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.GebaeudeFind(nr)
		is.Equal(updatedDoc, structs.Gebaeude{
			Nr:          1102,
			Bezeichnung: "Altes Hauptgebaeude (Westfluegel)",
			Flaeche: structs.GebaeudeFlaeche{
				HNF:     2632.27,
				NNF:     168.53,
				NGF:     4152.24,
				FF:      99,
				VF:      1351.44,
				FreiF:   0,
				GesamtF: 4251.24,
			},
			Einheit:     "m^2",
			Spezialfall: 1,
			Revision:    1,
			Stromversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
			},
			Waermeversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
			},
			Kaelteversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
			},
			KaelteRef: []int32{},
			WaermeRef: []int32{2348},
			StromRef:  []int32{28175, 999},
		}) // Ueberpruefung des zurueckgelieferten Elements
	})

	t.Run("GebaeudeAddZaehlerref: ID = 1103, idEnergieversorgung = 3", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var nr int32 = 1103
		var ref int32 = 999
		var idEnergieversorgung int32 = 3

		err := database.GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.GebaeudeFind(nr)
		is.Equal(updatedDoc, structs.Gebaeude{
			Nr:          1103,
			Bezeichnung: "Altes Hauptgebaeude",
			Flaeche: structs.GebaeudeFlaeche{
				HNF:     11745.22,
				NNF:     2191.86,
				NGF:     20297.78,
				FF:      2438.26,
				VF:      6360.7,
				FreiF:   0,
				GesamtF: 22736.04,
			},
			Einheit:     "m^2",
			Spezialfall: 1,
			Revision:    1,
			Stromversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
			},
			Waermeversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
			},
			Kaelteversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
			},
			KaelteRef: []int32{999},
			WaermeRef: []int32{2349, 2350, 2351, 2352, 2353, 2354},
			StromRef:  []int32{28175},
		}) // Ueberpruefung des zurueckgelieferten Elements
	})

	// dieser Fall sollte nicht auftreten
	t.Run("GebaeudeAddZaehlerref: Doppelte Referenz", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var nr int32 = 1108
		var ref int32 = 2250
		var idEnergieversorgung int32 = 1

		doc, _ := database.GebaeudeFind(nr)

		err := database.GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.GebaeudeFind(nr)
		is.Equal(updatedDoc, doc) // Ueberpruefung des zurueckgelieferten Elements
	})

	// Errortests
	t.Run("GebaeudeAddZaehlerref: ID = -12 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var nr int32 = -12
		var ref int32 = 999
		var idEnergieversorgung int32 = 3

		err := database.GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})

	t.Run("GebaeudeAddZaehlerref: IDEnergieversorgung = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var nr int32 = 1101
		var ref int32 = 999
		var idEnergieversorgung int32 = 0

		err := database.GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung)
		is.Equal(err, structs.ErrIDEnergieversorgungNichtVorhanden) // Datenbank wirft ErrNoDocuments
	})
}

func TestGebaeudeAddVersorger(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("AddVersorger: Gebaeude Nr 1101", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddVersorger{
			Nr:                  1102,
			Jahr:                3005,
			IDEnergieversorgung: 1,
			IDVertrag:           2,
		}

		err := database.GebaeudeAddVersorger(data)
		is.NoErr(err) // Datenbank wirft keinen Fehler

		gebaeude, err := database.GebaeudeFind(1102)
		is.NoErr(err)
		is.Equal(gebaeude, structs.Gebaeude{
			Nr:          1102,
			Bezeichnung: "Altes Hauptgebaeude (Westfluegel)",
			Flaeche: structs.GebaeudeFlaeche{
				HNF:     2632.27,
				NNF:     168.53,
				NGF:     4152.24,
				FF:      99,
				VF:      1351.44,
				FreiF:   0,
				GesamtF: 4251.24,
			},
			Einheit:     "m^2",
			Spezialfall: 1,
			Revision:    1,
			Stromversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
			},
			Waermeversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
				{Jahr: 3005, IDVertrag: 2},
			},
			Kaelteversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
			},
			KaelteRef: []int32{},
			WaermeRef: []int32{2348},
			StromRef:  []int32{28175, 999},
		}) // Ueberpruefung des zurueckgelieferten Elements
	})

	// Errortests
	t.Run("AddVersorger: Gebaeude nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddVersorger{
			Nr:                  0,
			Jahr:                3005,
			IDEnergieversorgung: 1,
			IDVertrag:           2,
		}

		err := database.GebaeudeAddVersorger(data)
		is.Equal(err, structs.ErrGebaeudeNichtVorhanden)
	})

	t.Run("AddVersorger: Versorger nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddVersorger{
			Nr:                  1101,
			Jahr:                3005,
			IDEnergieversorgung: 1,
			IDVertrag:           3,
		}

		err := database.GebaeudeAddVersorger(data)
		is.Equal(err, structs.ErrVertragNichtVorhanden)
	})

	t.Run("AddVersorger: IDEnergieversorgung nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddVersorger{
			Nr:                  1101,
			Jahr:                3005,
			IDEnergieversorgung: 0,
			IDVertrag:           1,
		}

		err := database.GebaeudeAddVersorger(data)
		is.Equal(err, structs.ErrIDEnergieversorgungNichtVorhanden)
	})

	t.Run("AddVersorger: Jahr schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddVersorger{
			Nr:                  1101,
			Jahr:                2018,
			IDEnergieversorgung: 1,
			IDVertrag:           1,
		}

		err := database.GebaeudeAddVersorger(data)
		is.Equal(err, structs.ErrJahrVorhanden)
	})
}

func TestGebaeudeAddStandardVersorger(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("AddStandardVersorger", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddStandardVersorger{
			Jahr: 3000,
		}

		err := database.GebaeudeAddStandardVersorger(data)
		is.NoErr(err) // Datenbank wirft keinen Fehler

		gebaeude, err := database.GebaeudeFind(1101)
		is.NoErr(err)
		is.Equal(gebaeude, structs.Gebaeude{
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
			Stromversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
				{Jahr: 3000, IDVertrag: 1},
			},
			Waermeversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
				{Jahr: 3000, IDVertrag: 1},
			},
			Kaelteversorger: []structs.Versoger{
				{Jahr: 2018, IDVertrag: 1},
				{Jahr: 2019, IDVertrag: 1},
				{Jahr: 2020, IDVertrag: 1},
				{Jahr: 2021, IDVertrag: 1},
				{Jahr: 2022, IDVertrag: 1},
				{Jahr: 3000, IDVertrag: 1},
			},
			KaelteRef: []int32{},
			WaermeRef: []int32{2084, 999},
			StromRef:  []int32{26024, 24799},
		}) // Überprüfung des zurückgelieferten Elements
	})
}

func TestNutzerdatenAddUmfrageref(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("NutzerdatenAddUmfrageref: username = anton", func(t *testing.T) {
		is.NewRelaxed(t)

		username := "anton"
		id := primitive.NewObjectID()
		objID, _ := primitive.ObjectIDFromHex("61b1ceb3dfb93b34b1305b70")

		var idVorhanden primitive.ObjectID
		err := idVorhanden.UnmarshalText([]byte("61b23e9855aa64762baf76d7"))
		is.NoErr(err)

		err = database.NutzerdatenAddUmfrageref(username, id)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.NutzerdatenFind(username)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(updatedDoc, structs.Nutzerdaten{
			NutzerID:   objID,
			EMail:      "anton@tobi.com",
			Nutzername: "anton",
			Rolle:      0,
			Revision:   2,
			UmfrageRef: []primitive.ObjectID{idVorhanden, id},
		}) // Überprüfung des zurückgelieferten Elements
	})

	// Errortests
	t.Run("NutzerdatenAddUmfrageref: username = 0 nicht vorhanden", func(t *testing.T) {
		is.NewRelaxed(t)

		username := "0"
		id := primitive.NewObjectID()

		err := database.NutzerdatenAddUmfrageref(username, id)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})
}

func TestUmfrageAddMitarbeiterUmfrageRef(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("UmfrageAddMitarbeiterUmfrageRef: ID = 61b23e9855aa64762baf76d7", func(t *testing.T) {
		is.NewRelaxed(t)

		var idUmfrage primitive.ObjectID
		err := idUmfrage.UnmarshalText([]byte("61b23e9855aa64762baf76d7"))
		is.NoErr(err)
		var idVorhanden primitive.ObjectID
		err = idVorhanden.UnmarshalText([]byte("61b34f9324756df01eee5ff4"))
		is.NoErr(err)

		referenz := primitive.NewObjectID()

		err = database.UmfrageAddMitarbeiterUmfrageRef(idUmfrage, referenz)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.UmfrageFind(idUmfrage)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(updatedDoc, structs.Umfrage{
			ID:                idUmfrage,
			Bezeichnung:       "testumfrage1",
			Mitarbeiteranzahl: 10,
			Jahr:              2020,
			Gebaeude: []structs.UmfrageGebaeude{
				{GebaeudeNr: 1101, Nutzflaeche: 100},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 5, Anzahl: 10},
			},
			AuswertungFreigegeben: 0,
			Revision:              1,
			MitarbeiterumfrageRef: []primitive.ObjectID{idVorhanden, referenz},
		}) // Überprüfung des zurückgelieferten Elements
	})

	// Errortests
	t.Run("UmfrageAddMitarbeiterUmfrageRef: ID nicht vorhanden", func(t *testing.T) {
		is.NewRelaxed(t)

		idUmfrage := primitive.NewObjectID()
		referenz := primitive.NewObjectID()

		err := database.UmfrageAddMitarbeiterUmfrageRef(idUmfrage, referenz)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})
}
