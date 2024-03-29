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

func TestFind(t *testing.T) {
	is := is.NewRelaxed(t)

	err := database.ConnectDatabase("dev")
	is.NoErr(err)

	dir, err := database.CreateDump("TestFind")
	is.NoErr(err)

	fmt.Println(dir)

	defer func() {
		err := database.DisconnectDatabase()
		is.NoErr(err)
		err = database.RestoreDump(dir)
		is.NoErr(err)
		err = database.RemoveDump(dir)
		is.NoErr(err)
	}()

	t.Run("TestITGeraeteFind", TestITGeraeteFind)
	t.Run("TestITGeraeteFindAll", TestITGeraeteFindAll)
	t.Run("TestEnergieversorgungFind", TestEnergieversorgungFind)
	t.Run("TestDienstreisenFind", TestDienstreisenFind)
	t.Run("TestDienstreisenFindAll", TestDienstreisenFindAll)
	t.Run("TestPendelwegFind", TestPendelwegFind)
	t.Run("TestPendelwegFindAll", TestPendelwegFindAll)
	t.Run("TestGebaeudeFind", TestGebaeudeFind)
	t.Run("TestZaehlerFind", TestZaehlerFind)
	t.Run("TestTestUmfrageFind", TestUmfrageFind)
	t.Run("TestMitarbeiterUmfrageFind", TestMitarbeiterUmfrageFind)
	t.Run("TestNutzerdatenFind", TestNutzerdatenFind)
	t.Run("TestNutzerdatenFindByEMail", TestNutzerdatenFindByEMail)
	t.Run("TestGebaeudeAlleNr", TestGebaeudeAlleNr)
	t.Run("TestGebaeudeAlleNrUndZaehlerRef", TestGebaeudeAlleNrUndZaehlerRef)
	t.Run("TestZaehlerAlleZaehlerUndDaten", TestZaehlerAlleZaehlerUndDaten)
	t.Run("TestMitarbeiterUmfrageFindMany", TestMitarbeiterUmfrageFindMany)
	t.Run("TestMitarbeiterUmfageForUmfrage", TestMitarbeiterUmfageForUmfrage)
	t.Run("TestAlleUmfragen", TestAlleUmfragen)
	t.Run("TestAlleUmfragenForUser", TestAlleUmfragenForUser)
}

func TestITGeraeteFind(t *testing.T) {
	is := is.NewRelaxed(t)

	//Normalfall
	t.Run("ITGeraeteFind: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ITGeraeteFind(1)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data,
			structs.ITGeraete{
				IDITGerate:      1,
				Kategorie:       "Notebooks",
				CO2FaktorGesamt: 588000,
				CO2FaktorJahr:   147000,
				Einheit:         "g/Stueck",
				Revision:        1,
			}) // Überprüfung des zurückgelieferten Elements
	})

	// Errorfall
	t.Run("ITGeraeteFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ITGeraeteFind(0)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.ITGeraete{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestITGeraeteFindAll(t *testing.T) {
	is := is.NewRelaxed(t)

	//Normalfall
	t.Run("TestITGeraeteFindAll: length", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ITGeraeteFindAll()

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(len(data), 10)
	})
}

func TestDienstreisenFind(t *testing.T) {
	is := is.NewRelaxed(t)

	//Normalfall
	t.Run("DienstreisenFind: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.DienstreisenFind(1)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data,
			structs.Dienstreisen{
				IDDienstreisen: 1,
				Medium:         "Bahn",
				Einheit:        "g/Pkm",
				Revision:       1,
				CO2Faktor:      []structs.CO2Dienstreisen{{Wert: 8}},
			}) // Überprüfung des zurückgelieferten Elements
	})

	// Errorfall
	t.Run("DienstreisenFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.DienstreisenFind(0)

		is.Equal(err, mongo.ErrNoDocuments)    // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Dienstreisen{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestDienstreisenFindAll(t *testing.T) {
	is := is.NewRelaxed(t)

	//Normalfall
	t.Run("TestDienstreisenFindAll: length", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.DienstreisenFindAll()

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(len(data), 3)
	})
}

func TestPendelwegFind(t *testing.T) {
	is := is.NewRelaxed(t)

	//Normalfall
	t.Run("TestPendelwegFind: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.PendelwegFind(1)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data,
			structs.Pendelweg{
				IDPendelweg: 1,
				Medium:      "Fahrrad",
				CO2Faktor:   4,
				Einheit:     "g/Pkm",
				Revision:    1,
			}) // Überprüfung des zurückgelieferten Elements
	})

	// Errorfall
	t.Run("TestPendelwegFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.PendelwegFind(0)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Pendelweg{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestPendelwegFindAll(t *testing.T) {
	is := is.NewRelaxed(t)

	//Normalfall
	t.Run("TestPendelwegFindAll: length", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.PendelwegFindAll()

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(len(data), 12)
	})
}

func TestEnergieversorgungFind(t *testing.T) {
	is := is.NewRelaxed(t)

	//Normalfall
	t.Run("EnergieversorgungFind: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.EnergieversorgungFind(1)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, structs.Energieversorgung{
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
			},
		}) // Überprüfung des zurückgelieferten Elements
	})

	t.Run("EnergieversorgungFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.EnergieversorgungFind(0)

		is.Equal(err, mongo.ErrNoDocuments)         // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Energieversorgung{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestGebaeudeFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("GebaeudeFind: ID = 1101", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1101)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, structs.Gebaeude{
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
			WaermeRef: []int32{2084},
			StromRef:  []int32{26024, 24799},
		}) // Überprüfung des zurückgelieferten Elements
	})

	t.Run("GebaeudeFind: ID = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(0)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Gebaeude{})  // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Justitzzentrum
	t.Run("GebaeudeFind: ID = 1473 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1437)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Gebaeude{})  // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Landgericht Gebaeude A
	t.Run("GebaeudeFind: ID = 1475 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1475)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Gebaeude{})  // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Landgericht Gebaeude B
	t.Run("GebaeudeFind: ID = 1476 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1476)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Gebaeude{})  // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Regierungspraesidium
	t.Run("GebaeudeFind: ID = 1477 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1477)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Gebaeude{})  // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Staatsbauamt
	t.Run("GebaeudeFind: ID = 1479 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1479)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Gebaeude{})  // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Landesmuseum
	t.Run("GebaeudeFind: ID = 1480 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1480)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Gebaeude{})  // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Staatsarchiv
	t.Run("GebaeudeFind: ID = 1481 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1481)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Gebaeude{})  // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Frauenhofer Institut (LBF)
	t.Run("GebaeudeFind: ID = 1213 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1213)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Gebaeude{})  // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

}

func TestZaehlerFind(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("ZaehlerFind: Kaelterzaehler, ID = 4023", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 4023
		var idEnergieversorgung int32 = 3

		location, _ := time.LoadLocation("Etc/GMT")
		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, structs.Zaehler{Zaehlertyp: "Kaelte",
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
			},
			Einheit:     "MWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{3101},
		})
	})

	t.Run("ZaehlerFind: Waermezaehler, ID = 2107", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2107
		var idEnergieversorgung int32 = 1

		location, _ := time.LoadLocation("Etc/GMT")
		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, structs.Zaehler{Zaehlertyp: "Waerme",
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
			},
			Einheit:     "MWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{2101, 2102, 2108},
		})
	})

	t.Run("ZaehlerFind: Stromzaehler, ID = 5967", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 5967
		var idEnergieversorgung int32 = 2

		location, _ := time.LoadLocation("Etc/GMT")
		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, structs.Zaehler{Zaehlertyp: "Strom",
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
			},
			Einheit:     "kWh",
			Spezialfall: 1,
			Revision:    1,
			GebaeudeRef: []int32{2201},
		})
	})

	// Errortests
	t.Run("ZaehlerFind: Kaeltezaehler, ID = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 0
		var idEnergieversorgung int32 = 3

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	t.Run("ZaehlerFind: Waermezaehler, ID = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 0
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1473 Waerme Hauptzaehler Justitzzentrum
	t.Run("ZaehlerFind: Waermezaehler, ID = 2104 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2104
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1475 Waerme Hauptzaehler Landgericht Gebaeude A
	t.Run("ZaehlerFind: Waermezaehler, ID = 2105 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2105
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1476 Waerme Hauptzaehler Landgericht Gebaeude B
	t.Run("ZaehlerFind: Waermezaehler, ID = 2106 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2106
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1477 Waerme Hauptzaehler Regierungspraesidium
	t.Run("ZaehlerFind: Waermezaehler, ID = 2256 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2256
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1479 Waerme Hauptzaehler Staatsbauamt
	t.Run("ZaehlerFind: Waermezaehler, ID = 3613 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 3613
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1480 Waerme Hauptzaehler Landesmuseum
	t.Run("ZaehlerFind: Waermezaehler, ID = 3614 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 3614
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1481 Waerme Hauptzaehler Staatsarchiv
	t.Run("ZaehlerFind: Waermezaehler, ID = 2102 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2012
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1213a Altbau Frauenhofer Institut (LBF) Waerme
	t.Run("ZaehlerFind: Waermezaehler, ID = 2377 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2377
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1213b Frauenhofer Institut (LBF) Neubau Waerme
	t.Run("ZaehlerFind: Waermezaehler, ID = 2378 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2378
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, Geb_Rechenwerk_Neues_RP_Entega_ENERGIE
	t.Run("ZaehlerFind: Waermezaehler, ID = 4193 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 4193
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, Geb_Rechenwerk_Neues_RP_Steag_ENERGIE
	t.Run("ZaehlerFind: Waermezaehler, ID = 4194 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 4194
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, Heizgruppe Süd Schloss gesamt + E-Technik gesamt
	t.Run("ZaehlerFind: Waermezaehler, ID = 3960 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 3960
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, Netz Darmstadt Nord
	t.Run("ZaehlerFind: Waermezaehler, ID = 6697 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 6697
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, Wärmezähler Lichtwiese Forschungsprojekt
	t.Run("ZaehlerFind: Waermezaehler, ID = 3789 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 3789
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, Wärme Hauptzähler Tagesverbrauch
	t.Run("ZaehlerFind: Waermezaehler, ID = 2558 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2558
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, Wärme Hauptzähler Wochenverbrauch
	t.Run("ZaehlerFind: Waermezaehler, ID = 2560 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2560
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	t.Run("ZaehlerFind: Stromzaehler, ID = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 0
		var idEnergieversorgung int32 = 2

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// dieser Zaehler wurde rausgenommen, weil die Einheit kW ist
	t.Run("ZaehlerFind: Stromzaehler, ID = 3576 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 3576
		var idEnergieversorgung int32 = 2

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	t.Run("ZaehlerFind: IDEnergieversorgung = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 1
		var idEnergieversorgung int32 = 0

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, structs.ErrIDEnergieversorgungNichtVorhanden) // Funktion wirft ErrIDEnergieversorgungNichtVorhanden
		is.Equal(data, structs.Zaehler{})                           // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestUmfrageFind(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("UmfrageFind: ID = 61b23e9855aa64762baf76d7", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var id primitive.ObjectID
		err := id.UnmarshalText([]byte("61b23e9855aa64762baf76d7"))
		is.NoErr(err)

		var idMitarbeiterumfrage primitive.ObjectID
		err = idMitarbeiterumfrage.UnmarshalText([]byte("61b34f9324756df01eee5ff4"))
		is.NoErr(err)

		data, err := database.UmfrageFind(id)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data,
			structs.Umfrage{
				ID:                id,
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
				MitarbeiterumfrageRef: []primitive.ObjectID{idMitarbeiterumfrage},
			}) // Überprüfung des zurückgelieferten Elements
	})

	// Errortests
	t.Run("ITGeraeteFind: ID aus aktuellem Zeitstempel nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		id := primitive.NewObjectID()

		data, err := database.UmfrageFind(id)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Umfrage{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestMitarbeiterUmfrageFind(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("MitarbeiterUmfrageFind: ID = 61b34f9324756df01eee5ff4", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var id primitive.ObjectID
		err := id.UnmarshalText([]byte("61b34f9324756df01eee5ff4"))
		is.NoErr(err)

		data, err := database.MitarbeiterUmfrageFind(id)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data,
			structs.MitarbeiterUmfrage{
				ID: id,
				Pendelweg: []structs.UmfragePendelweg{
					{IDPendelweg: 1, Strecke: 123, Personenanzahl: 1},
				},
				TageImBuero: 7,
				Dienstreise: []structs.UmfrageDienstreise{
					{IDDienstreise: 3, Streckentyp: "Langstrecke", Strecke: 321},
				},
				ITGeraete: []structs.UmfrageITGeraete{
					{IDITGeraete: 3, Anzahl: 45},
				},
				Revision: 1,
			}) // Überprüfung des zurückgelieferten Elements
	})

	// Errortests
	t.Run("MitarbeiterUmfrageFind: ID aus aktuellem Zeitstempel nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		id := primitive.NewObjectID()

		data, err := database.MitarbeiterUmfrageFind(id)

		is.Equal(err, mongo.ErrNoDocuments)          // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.MitarbeiterUmfrage{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestMitarbeiterUmfrageFindMany(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("MitarbeiterumfrageFindMany: einzelne ID", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var id primitive.ObjectID
		err := id.UnmarshalText([]byte("61b34f9324756df01eee5ff4"))
		is.NoErr(err)

		ids := []primitive.ObjectID{id}

		data, err := database.MitarbeiterumfrageFindMany(ids)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data,
			[]structs.MitarbeiterUmfrage{{
				ID: id,
				Pendelweg: []structs.UmfragePendelweg{
					{IDPendelweg: 1, Strecke: 123, Personenanzahl: 1},
				},
				TageImBuero: 7,
				Dienstreise: []structs.UmfrageDienstreise{
					{IDDienstreise: 3, Streckentyp: "Langstrecke", Strecke: 321},
				},
				ITGeraete: []structs.UmfrageITGeraete{
					{IDITGeraete: 3, Anzahl: 45},
				},
				Revision: 1,
			}}) // Überprüfung des zurückgelieferten Elements
	})

	t.Run("MitarbeiterumfrageFindMany: mehrere IDs", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var id primitive.ObjectID
		err := id.UnmarshalText([]byte("61b34f9324756df01eee5ff4"))
		is.NoErr(err)

		var idUmfrage primitive.ObjectID
		err = idUmfrage.UnmarshalText([]byte("61b23e9855aa64762baf76d7"))
		is.NoErr(err)

		id2, err := database.MitarbeiterUmfrageInsert(structs.InsertMitarbeiterUmfrage{IDUmfrage: idUmfrage})
		is.NoErr(err)

		ids := []primitive.ObjectID{id, id2}

		data, err := database.MitarbeiterumfrageFindMany(ids)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, []structs.MitarbeiterUmfrage{{
			ID: id,
			Pendelweg: []structs.UmfragePendelweg{
				{IDPendelweg: 1, Strecke: 123, Personenanzahl: 1},
			},
			TageImBuero: 7,
			Dienstreise: []structs.UmfrageDienstreise{
				{IDDienstreise: 3, Streckentyp: "Langstrecke", Strecke: 321},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 3, Anzahl: 45},
			},
			Revision: 1,
		},
			{
				ID:          id2,
				Pendelweg:   nil,
				TageImBuero: 0,
				Dienstreise: nil,
				ITGeraete:   nil,
				Revision:    1,
			},
		}) // Überprüfung des zurückgelieferten Elements
	})

	// Errortests
	t.Run("MitarbeiterumfrageFindMany: zu wenige Dokumente", func(t *testing.T) {
		is := is.NewRelaxed(t)

		ids := []primitive.ObjectID{primitive.NewObjectID()}

		data, err := database.MitarbeiterumfrageFindMany(ids)

		is.Equal(err, structs.ErrDokumenteNichtGefunden) // Funktion wirft ErrDokumenteNichtGefunden
		is.Equal(data, nil)                              // Bei einem Fehler soll nil zurückgeliefert werden
	})
}

func TestNutzerdatenFind(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("NutzerdatenFind: username = anton", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "anton"
		var idUmfrage primitive.ObjectID
		err := idUmfrage.UnmarshalText([]byte("61b23e9855aa64762baf76d7"))
		is.NoErr(err)
		objID, _ := primitive.ObjectIDFromHex("61b1ceb3dfb93b34b1305b70")
		data, err := database.NutzerdatenFind(username)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data,
			structs.Nutzerdaten{
				NutzerID:   objID,
				EMail:      "anton@tobi.com",
				Nutzername: "anton",
				Rolle:      0,
				Revision:   2,
				UmfrageRef: []primitive.ObjectID{idUmfrage},
			}) // Überprüfung des zurückgelieferten Elements
	})

	// Errortests
	t.Run("MitarbeiterUmfrageFind: username = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "0"

		data, err := database.NutzerdatenFind(username)

		is.Equal(err, mongo.ErrNoDocuments)   // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Nutzerdaten{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestNutzerdatenFindByEMail(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("NutzerdatenFind: email = anton@tobi.com", func(t *testing.T) {
		is := is.NewRelaxed(t)

		email := "anton@tobi.com"
		var idUmfrage primitive.ObjectID
		err := idUmfrage.UnmarshalText([]byte("61b23e9855aa64762baf76d7"))
		is.NoErr(err)
		objID, _ := primitive.ObjectIDFromHex("61b1ceb3dfb93b34b1305b70")

		data, err := database.NutzerdatenFindByEMail(email)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data,
			structs.Nutzerdaten{
				NutzerID:   objID,
				EMail:      "anton@tobi.com",
				Nutzername: "anton",
				Rolle:      0,
				Revision:   2,
				UmfrageRef: []primitive.ObjectID{idUmfrage},
			}) // Überprüfung des zurückgelieferten Elements
	})

	// Errortests
	t.Run("MitarbeiterUmfrageFind: username = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		email := "0"

		data, err := database.NutzerdatenFindByEMail(email)

		is.Equal(err, mongo.ErrNoDocuments)   // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Nutzerdaten{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestGebaeudeAlleNr(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("GebaeudeAlleNr: liefert Slice zurueck", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudenummer, err := database.GebaeudeAlleNr()
		is.NoErr(err)                    // kein Error seitens der Datenbank
		is.True(len(gebaeudenummer) > 0) // Slice ist nicht leer
	})
}

func TestGebaeudeAlleNrUndZaehlerRef(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("GebaeudeAlleNrUndZaehlerRef: liefert Slice zurueck", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, err := database.GebaeudeAlleNrUndZaehlerRef()
		is.NoErr(err)              // kein Error seitens der Datenbank
		is.True(len(gebaeude) > 0) // Slice ist nicht leer
	})
}

func TestZaehlerAlleZaehlerUndDaten(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("ZaehlerAlleZaehlerUndDaten: liefert Slice zurueck", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeude, err := database.ZaehlerAlleZaehlerUndDaten()

		is.NoErr(err)              // kein Error seitens der Datenbank
		is.True(len(gebaeude) > 0) // Slice ist nicht leer
	})
}

func TestMitarbeiterUmfageForUmfrage(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("MitarbeiterUmfrageFindForUmfrage: liefert MitarbeiterUmfrageIDs zurück", func(t *testing.T) {
		is := is.NewRelaxed(t)

		umfrageID, err := primitive.ObjectIDFromHex("61b23e9855aa64762baf76d7")
		is.NoErr(err)

		mitarbeiterUmfragen, err := database.MitarbeiterUmfrageFindForUmfrage(umfrageID)
		is.NoErr(err)                         // kein Error seitens der Datenbank
		is.True(len(mitarbeiterUmfragen) > 0) // Slice ist nicht leer

		correctMitarbeiterUmfrageID, err := primitive.ObjectIDFromHex("61b34f9324756df01eee5ff4")
		is.NoErr(err)

		// after clarification
		is.Equal(mitarbeiterUmfragen[0], structs.MitarbeiterUmfrage{
			ID: correctMitarbeiterUmfrageID,
			Pendelweg: []structs.UmfragePendelweg{
				{
					IDPendelweg:    1,
					Strecke:        123,
					Personenanzahl: 1,
				},
			},
			TageImBuero: 7,
			Dienstreise: []structs.UmfrageDienstreise{
				{
					IDDienstreise: 3,
					Streckentyp:   "Langstrecke",
					Strecke:       321,
				},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{
					IDITGeraete: 3,
					Anzahl:      45,
				},
			},
			Revision: 1})
	})

	// Normalfall
	t.Run("MitarbeiterUmfrageFindForUmfrage: liefert keine MitarbeiterUmfrageRefs", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var umfrageID primitive.ObjectID
		err := umfrageID.UnmarshalText([]byte("61dc0c543e48998484eefae9"))
		is.NoErr(err)

		mitarbeiterUmfragen, err := database.MitarbeiterUmfrageFindForUmfrage(umfrageID)
		is.NoErr(err)                         // kein Error seitens der Datenbank
		is.Equal(len(mitarbeiterUmfragen), 0) // Slice ist nicht leer
	})

	// Errorfaelle
	t.Run("MitarbeiterUmfrageFindForUmfrage: umfrageID existiert nicht", func(t *testing.T) {
		is := is.NewRelaxed(t)

		// assumes there will be no zero objectID
		umfrageID, err := primitive.ObjectIDFromHex("000000000000000000000000")
		is.NoErr(err)

		mitarbeiterUmfragen, err := database.MitarbeiterUmfrageFindForUmfrage(umfrageID)
		is.Equal(mitarbeiterUmfragen, nil)  // leerer Array
		is.Equal(err, mongo.ErrNoDocuments) // Error raised
	})
}

func TestAlleUmfragen(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("AlleUmfragen: liefert alle existenten Umfragen zurück", func(t *testing.T) {
		is := is.NewRelaxed(t)

		alleUmfragen, err := database.AlleUmfragen()
		is.NoErr(err)                  // kein Error seitens der Datenbank
		is.Equal(len(alleUmfragen), 5) // Slice ist nicht leer
	})
}

func TestAlleUmfragenForUser(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("AlleUmfragenForUser: liefert alle existenten Umfragen zurueck, die mit gegebenem User assoziiert sind.", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "anton"
		alleUmfragen, err := database.AlleUmfragenForUser(username)
		is.NoErr(err)                  // kein Error seitens der Datenbank
		is.True(len(alleUmfragen) > 0) // Slice ist nicht leer

		correctRefID, err := primitive.ObjectIDFromHex("61b23e9855aa64762baf76d7")
		is.NoErr(err)
		is.Equal(alleUmfragen[0].ID, correctRefID)
	})

	t.Run("AlleUmfragenForUser: Keine Umfragen mit User assoziiert", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "lorem"
		alleUmfragen, err := database.AlleUmfragenForUser(username)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(alleUmfragen, []structs.Umfrage{})
	})

	// Errorfaelle
	t.Run("AlleUmfragenForUser: User existiert nicht", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "KaeptnBlaubaer"
		alleUmfragen, err := database.AlleUmfragenForUser(username)
		is.Equal(alleUmfragen, nil)
		is.Equal(err, mongo.ErrNoDocuments)
	})
}
