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
	t.Run("TestZaehlerAddZaehlerdaten", TestZaehlerAddZaehlerdaten)
	t.Run("TestGebaeudeAddZaehlerref", TestGebaeudeAddZaehlerref)
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
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})

	t.Run("EnergieversorgungAddFaktor: Jahr schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddCO2Faktor{
			IDEnergieversorgung: 1,
			Wert:                400,
			Jahr:                2020,
		}

		err := database.EnergieversorgungAddFaktor(data)
		is.Equal(err, structs.ErrJahrVorhanden) // Funktion wirft ErrJahrVorhanden
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
					Zeitstempel: time.Date(3001, time.January, 01, 0, 0, 0, 0, location).UTC(),
				},
			},
			Einheit:     "MWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{3101},
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
			Jahr:                2020,
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
			KaelteRef:   []int32{},
			WaermeRef:   []int32{2084, 999},
			StromRef:    []int32{},
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
			KaelteRef:   []int32{},
			WaermeRef:   []int32{2348},
			StromRef:    []int32{999},
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
			KaelteRef:   []int32{999},
			WaermeRef:   []int32{2349, 2350, 2351, 2352, 2353, 2354},
			StromRef:    []int32{},
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

func TestNutzerdatenAddUmfrageref(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("NutzerdatenAddUmfrageref: Email = anton@tobi.com", func(t *testing.T) {
		is.NewRelaxed(t)

		email := "anton@tobi.com"
		id := primitive.NewObjectID()

		var idVorhanden primitive.ObjectID
		err := idVorhanden.UnmarshalText([]byte("61b23e9855aa64762baf76d7"))
		is.NoErr(err)

		err = database.NutzerdatenAddUmfrageref(email, id)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.NutzerdatenFind(email)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(updatedDoc, structs.Nutzerdaten{
			Email:      "anton@tobi.com",
			Passwort:   "test_pw",
			Revision:   1,
			UmfrageRef: []primitive.ObjectID{idVorhanden, id},
		}) // Überprüfung des zurückgelieferten Elements
	})

	// Errortests
	t.Run("NutzerdatenAddUmfrageref: Email = 0 nicht vorhanden", func(t *testing.T) {
		is.NewRelaxed(t)

		email := "0"
		id := primitive.NewObjectID()

		err := database.NutzerdatenAddUmfrageref(email, id)
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
			Mitarbeiteranzahl: 1,
			Jahr:              2020,
			Gebaeude: []structs.UmfrageGebaeude{
				{GebaeudeNr: 1101, Nutzflaeche: 100},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 5, Anzahl: 10},
			},
			Revision:              1,
			MitarbeiterUmfrageRef: []primitive.ObjectID{idVorhanden, referenz},
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
