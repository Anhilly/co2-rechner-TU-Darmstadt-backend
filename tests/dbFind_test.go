package tests

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"io"
	"testing"
	"time"
)

func TestFind(t *testing.T) {
	is := is.NewRelaxed(t)

	err := database.ConnectDatabase()
	is.NoErr(err)
	defer func() {
		err := database.DisconnectDatabase()
		is.NoErr(err)
	}()

	t.Run("TestITGeraeteFind", TestITGeraeteFind)
	t.Run("TestEnergieversorgungFind", TestEnergieversorgungFind)
	t.Run("TestDienstreisenFind", TestDienstreisenFind)
	t.Run("TestGebaeudeFind", TestGebaeudeFind)
	t.Run("TestZaehlerFind", TestZaehlerFind)
}

func TestITGeraeteFind(t *testing.T) {
	is := is.NewRelaxed(t)

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

	t.Run("ITGeraeteFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ITGeraeteFind(0)

		is.Equal(err, io.EOF)               // Datenbank wirft EOF-Error
		is.Equal(data, structs.ITGeraete{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestDienstreisenFind(t *testing.T) {
	is := is.NewRelaxed(t)

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

	t.Run("DienstreisenFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.DienstreisenFind(0)

		is.Equal(err, io.EOF)                  // Datenbank wirft EOF-Error
		is.Equal(data, structs.Dienstreisen{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestEnergieversorgungFind(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("EnergieversorgungFind: ID = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.EnergieversorgungFind(1)

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, structs.Energieversorgung{
			IDEnergieversorgung: 1,
			Kategorie:           "Fernwaerme",
			Einheit:             "g/kWh",
			Revision:            1,
			CO2Faktor:           []structs.CO2Energie{{Wert: 144, Jahr: 2020}},
		}) // Überprüfung des zurückgelieferten Elements
	})

	t.Run("EnergieversorgungFind: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.EnergieversorgungFind(0)

		is.Equal(err, io.EOF)                       // Datenbank wirft EOF-Error
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
			KaelteRef:   []int32{},
			WaermeRef:   []int32{2084},
			StromRef:    []int32{},
		}) // Überprüfung des zurückgelieferten Elements
	})

	t.Run("GebaeudeFind: ID = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(0)

		is.Equal(err, io.EOF)              // Datenbank wirft EOF-Error
		is.Equal(data, structs.Gebaeude{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Justitzzentrum
	t.Run("GebaeudeFind: ID = 1473 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1437)

		is.Equal(err, io.EOF)              // Datenbank wirft EOF-Error
		is.Equal(data, structs.Gebaeude{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Landgericht Gebaeude A
	t.Run("GebaeudeFind: ID = 1475 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1475)

		is.Equal(err, io.EOF)              // Datenbank wirft EOF-Error
		is.Equal(data, structs.Gebaeude{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Landgericht Gebaeude B
	t.Run("GebaeudeFind: ID = 1476 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1476)

		is.Equal(err, io.EOF)              // Datenbank wirft EOF-Error
		is.Equal(data, structs.Gebaeude{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Regierungspraesidium
	t.Run("GebaeudeFind: ID = 1477 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1477)

		is.Equal(err, io.EOF)              // Datenbank wirft EOF-Error
		is.Equal(data, structs.Gebaeude{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Staatsbauamt
	t.Run("GebaeudeFind: ID = 1479 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1479)

		is.Equal(err, io.EOF)              // Datenbank wirft EOF-Error
		is.Equal(data, structs.Gebaeude{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Landesmuseum
	t.Run("GebaeudeFind: ID = 1480 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1480)

		is.Equal(err, io.EOF)              // Datenbank wirft EOF-Error
		is.Equal(data, structs.Gebaeude{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Staatsarchiv
	t.Run("GebaeudeFind: ID = 1481 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1481)

		is.Equal(err, io.EOF)              // Datenbank wirft EOF-Error
		is.Equal(data, structs.Gebaeude{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Gebaeude soll nicht beachtet werden, Frauenhofer Institut (LBF)
	t.Run("GebaeudeFind: ID = 1213 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.GebaeudeFind(1213)

		is.Equal(err, io.EOF)              // Datenbank wirft EOF-Error
		is.Equal(data, structs.Gebaeude{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
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

	// Errortests
	t.Run("ZaehlerFind: Kaeltezaehler, ID = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 0
		var idEnergieversorgung int32 = 3

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	t.Run("ZaehlerFind: Waermezaehler, ID = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 0
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1473 Waerme Hauptzaehler Justitzzentrum
	t.Run("ZaehlerFind: Waermezaehler, ID = 2014 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2014
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1475 Waerme Hauptzaehler Landgericht Gebaeude A
	t.Run("ZaehlerFind: Waermezaehler, ID = 2015 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2015
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1476 Waerme Hauptzaehler Landgericht Gebaeude B
	t.Run("ZaehlerFind: Waermezaehler, ID = 2016 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2016
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1477 Waerme Hauptzaehler Regierungspraesidium
	t.Run("ZaehlerFind: Waermezaehler, ID = 2256 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2256
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1479 Waerme Hauptzaehler Staatsbauamt
	t.Run("ZaehlerFind: Waermezaehler, ID = 3613 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 3613
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1480 Waerme Hauptzaehler Landesmuseum
	t.Run("ZaehlerFind: Waermezaehler, ID = 3614 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 3614
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1481 Waerme Hauptzaehler Staatsarchiv
	t.Run("ZaehlerFind: Waermezaehler, ID = 2102 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2012
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1213a Altbau Frauenhofer Institut (LBF) Waerme
	t.Run("ZaehlerFind: Waermezaehler, ID = 2377 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2377
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, 1213b Frauenhofer Institut (LBF) Neubau Waerme
	t.Run("ZaehlerFind: Waermezaehler, ID = 2378 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 2378
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, Geb_Rechenwerk_Neues_RP_Entega_ENERGIE
	t.Run("ZaehlerFind: Waermezaehler, ID = 4193 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 4193
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, Geb_Rechenwerk_Neues_RP_Steag_ENERGIE
	t.Run("ZaehlerFind: Waermezaehler, ID = 4194 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 4194
		var idEnergieversorgung int32 = 1

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	t.Run("ZaehlerFind: Stromzaehler, ID = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 0
		var idEnergieversorgung int32 = 2

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// dieser Zaehler wurde rausgenommen, weil die Einheit kW ist
	t.Run("ZaehlerFind: Stromzaehler, ID = 3576 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pkEnergie int32 = 3576
		var idEnergieversorgung int32 = 2

		data, err := database.ZaehlerFind(pkEnergie, idEnergieversorgung)

		is.Equal(err, io.EOF)             // Datenbank wirft EOF-Error
		is.Equal(data, structs.Zaehler{}) // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
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
