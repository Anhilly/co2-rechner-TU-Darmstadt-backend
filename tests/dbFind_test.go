package tests

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
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
	t.Run("TestGebaeudeFindOID", TestGebaeudeFindOID)
	t.Run("TestZaehlerFindDPName", TestZaehlerFindDPName)
	t.Run("TestZaehlerFindOID", TestZaehlerFindOID)
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
		is.Equal(len(data), 14)
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

		// entferne jahresabhängige Daten
		data.Waermeversorger = []structs.Versoger{}
		data.Stromversorger = []structs.Versoger{}
		data.Kaelteversorger = []structs.Versoger{}

		waermeZaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff656")
		gebaeudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a286")
		compareDoc := structs.Gebaeude{
			GebaeudeID:  gebaeudeOID,
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
			Einheit:         "m^2",
			Spezialfall:     1,
			Revision:        2,
			Stromversorger:  []structs.Versoger{},
			Waermeversorger: []structs.Versoger{},
			Kaelteversorger: []structs.Versoger{},
			KaelteRef:       []primitive.ObjectID{},
			WaermeRef:       []primitive.ObjectID{waermeZaehlerOID},
			StromRef:        []primitive.ObjectID{},
		}

		is.NoErr(err)              // kein Error seitens der Datenbank
		is.Equal(data, compareDoc) // Überprüfung des zurückgelieferten Elements
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

func TestGebaeudeFindOID(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("GebaeudeFindOID: OID = 61b712f66a1a52dea358a286", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a286")

		data, err := database.GebaeudeFindOID(gebaeudeOID)

		// entferne jahresabhängige Daten
		data.Waermeversorger = []structs.Versoger{}
		data.Stromversorger = []structs.Versoger{}
		data.Kaelteversorger = []structs.Versoger{}

		waermeZaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff656")
		compareDoc := structs.Gebaeude{
			GebaeudeID:  gebaeudeOID,
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
			Einheit:         "m^2",
			Spezialfall:     1,
			Revision:        2,
			Stromversorger:  []structs.Versoger{},
			Waermeversorger: []structs.Versoger{},
			Kaelteversorger: []structs.Versoger{},
			KaelteRef:       []primitive.ObjectID{},
			WaermeRef:       []primitive.ObjectID{waermeZaehlerOID},
			StromRef:        []primitive.ObjectID{},
		}

		is.NoErr(err)              // kein Error seitens der Datenbank
		is.Equal(data, compareDoc) // Überprüfung des zurückgelieferten Elements
	})

	t.Run("GebaeudeFindOID: OID = 000000000000000000000000 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		gebaeudeOID := primitive.NilObjectID

		data, err := database.GebaeudeFindOID(gebaeudeOID)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Gebaeude{})  // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestZaehlerFindDPName(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("ZaehlerFindDPName: Kaelterzaehler, DPName = L101XXXXXXKA000XXXXXXZ50CO00001", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("L101XXXXXXKA000XXXXXXZ50CO00001", structs.IDEnergieversorgungKaelte)
		data.Zaehlerdaten = []structs.Zaehlerwerte{}

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff67d")
		gebaudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a2e7")
		compareDoc := structs.Zaehler{
			ZaehlerID:    zaehlerOID,
			Zaehlertyp:   "Kaelte",
			DPName:       "L101XXXXXXKA000XXXXXXZ50CO00001",
			Bezeichnung:  "3101 Kälte Hauptzähler",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "MWh",
			Spezialfall:  1,
			Revision:     2,
			GebaeudeRef:  []primitive.ObjectID{gebaudeOID},
		}

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, compareDoc)
	})

	t.Run("ZaehlerFindDPName: Waermezaehler, DPName = B101XXXXXXHE000XXXXXXZ40CO00001", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("B101XXXXXXHE000XXXXXXZ40CO00001", structs.IDEnergieversorgungWaerme)
		data.Zaehlerdaten = []structs.Zaehlerwerte{}

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1c93a47b613426ff62b")
		gebaudeOID2101, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a2cf")
		gebaudeOID2102, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a2d0")
		gebaudeOID2108, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a2d5")
		compareDoc := structs.Zaehler{
			ZaehlerID:    zaehlerOID,
			Zaehlertyp:   "Waerme",
			DPName:       "B101XXXXXXHE000XXXXXXZ40CO00001",
			Bezeichnung:  "2101 Zoologie,2102 Botanik (Altbau), 2108 Mobi-Office 1 Wärme Gruppenzähler",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "MWh",
			Spezialfall:  1,
			Revision:     2,
			GebaeudeRef:  []primitive.ObjectID{gebaudeOID2101, gebaudeOID2102, gebaudeOID2108},
		}

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, compareDoc)
	})

	t.Run("ZaehlerFindDPName: Stromzaehler, DPName = B102xxxxxxNA000xxxxxxZ01ED11005", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("B102xxxxxxNA000xxxxxxZ01ED11005", structs.IDEnergieversorgungStrom)
		data.Zaehlerdaten = []structs.Zaehlerwerte{}

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff68e")
		gebaeudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a2d0")
		compareDoc := structs.Zaehler{
			ZaehlerID:    zaehlerOID,
			Zaehlertyp:   "Strom",
			DPName:       "B102xxxxxxNA000xxxxxxZ01ED11005",
			Bezeichnung:  "2102 Elektro HZ Blechverteiler Altbau",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "kWh",
			Spezialfall:  1,
			Revision:     2,
			GebaeudeRef:  []primitive.ObjectID{gebaeudeOID},
		}

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, compareDoc)
	})

	// Errortests
	t.Run("ZaehlerFindDPName: Kaeltezaehler, DPName = xxxx nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("xxxx", structs.IDEnergieversorgungKaelte)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	t.Run("ZaehlerFindDPName: Waermezaehler, DPName = xxxx nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("xxxx", structs.IDEnergieversorgungWaerme)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	t.Run("ZaehlerFindDPName: Stromzaehler, DPName = xxxx nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("xxxx", structs.IDEnergieversorgungStrom)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	t.Run("ZaehlerFindDPName: IDEnergieversorgung = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("B102xxxxxxNA000xxxxxxZ01ED11005", 0)

		is.Equal(err, structs.ErrIDEnergieversorgungNichtVorhanden) // Funktion wirft ErrIDEnergieversorgungNichtVorhanden
		is.Equal(data, structs.Zaehler{})                           // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, S473XXXXXXHE000XXXXXXZ40CO00001 Waerme Hauptzaehler Justitzzentrum
	t.Run("ZaehlerFindDPName: Waermezaehler, DPName = S473XXXXXXHE000XXXXXXZ40CO00001 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("S473XXXXXXHE000XXXXXXZ40CO00001", structs.IDEnergieversorgungWaerme)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, S475XXXXXXHE000XXXXXXZ40CO00001 Waerme Hauptzaehler Landgericht Gebaeude A
	t.Run("ZaehlerFindDPName: Waermezaehler, DPName = S475XXXXXXHE000XXXXXXZ40CO00001 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("S475XXXXXXHE000XXXXXXZ40CO00001", structs.IDEnergieversorgungWaerme)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, S476XXXXXXHE000XXXXXXZ40CO00001 Waerme Hauptzaehler Landgericht Gebaeude B
	t.Run("ZaehlerFindDPName: Waermezaehler, DPName = S476XXXXXXHE000XXXXXXZ40CO00001 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("S476XXXXXXHE000XXXXXXZ40CO00001", structs.IDEnergieversorgungWaerme)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, S477XXXXXXHE000XXXXXXZ40CO00001 Waerme Hauptzaehler Regierungspraesidium
	t.Run("ZaehlerFindDPName: Waermezaehler, DPName = S477XXXXXXHE000XXXXXXZ40CO00001 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("S477XXXXXXHE000XXXXXXZ40CO00001", structs.IDEnergieversorgungWaerme)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, S479XXXXXXHE000XXXXXXZ40CO00001 Waerme Hauptzaehler Staatsbauamt
	t.Run("ZaehlerFindDPName: Waermezaehler, DPName = S479XXXXXXHE000XXXXXXZ40CO00001 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("S479XXXXXXHE000XXXXXXZ40CO00001", structs.IDEnergieversorgungWaerme)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, S480XXXXXXHE000XXXXXXZ40CO00001 Waerme Hauptzaehler Landesmuseum
	t.Run("ZaehlerFindDPName: Waermezaehler, DPName = S480XXXXXXHE000XXXXXXZ40CO00001 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("S480XXXXXXHE000XXXXXXZ40CO00001", structs.IDEnergieversorgungWaerme)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	// Zaehler soll nicht beachtet werden, S481XXXXXXHE000XXXXXXZ40CO00001 Waerme Hauptzaehler Staatsarchiv
	t.Run("ZaehlerFindDPName: Waermezaehler, DPName = S481XXXXXXHE000XXXXXXZ40CO00001 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindDPName("S481XXXXXXHE000XXXXXXZ40CO00001", structs.IDEnergieversorgungWaerme)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})
}

func TestZaehlerFindOID(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("TestZaehlerFindOID: Kaelterzaehler, OID = 6710a1ca3a47b613426ff67d", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff67d")

		data, err := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungKaelte)
		data.Zaehlerdaten = []structs.Zaehlerwerte{}

		gebaudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a2e7")
		compareDoc := structs.Zaehler{
			ZaehlerID:    zaehlerOID,
			Zaehlertyp:   "Kaelte",
			DPName:       "L101XXXXXXKA000XXXXXXZ50CO00001",
			Bezeichnung:  "3101 Kälte Hauptzähler",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "MWh",
			Spezialfall:  1,
			Revision:     2,
			GebaeudeRef:  []primitive.ObjectID{gebaudeOID},
		}

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, compareDoc)
	})

	t.Run("TestZaehlerFindOID: Waermezaehler, OID = 6710a1c93a47b613426ff62b", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1c93a47b613426ff62b")

		data, err := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungWaerme)
		data.Zaehlerdaten = []structs.Zaehlerwerte{}

		gebaudeOID2101, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a2cf")
		gebaudeOID2102, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a2d0")
		gebaudeOID2108, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a2d5")
		compareDoc := structs.Zaehler{
			ZaehlerID:    zaehlerOID,
			Zaehlertyp:   "Waerme",
			DPName:       "B101XXXXXXHE000XXXXXXZ40CO00001",
			Bezeichnung:  "2101 Zoologie,2102 Botanik (Altbau), 2108 Mobi-Office 1 Wärme Gruppenzähler",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "MWh",
			Spezialfall:  1,
			Revision:     2,
			GebaeudeRef:  []primitive.ObjectID{gebaudeOID2101, gebaudeOID2102, gebaudeOID2108},
		}

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, compareDoc)
	})

	t.Run("TestZaehlerFindOID: Stromzaehler, OID = 6710a1ca3a47b613426ff68e", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff68e")

		data, err := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungStrom)
		data.Zaehlerdaten = []structs.Zaehlerwerte{}

		gebaeudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a2d0")
		compareDoc := structs.Zaehler{
			ZaehlerID:    zaehlerOID,
			Zaehlertyp:   "Strom",
			DPName:       "B102xxxxxxNA000xxxxxxZ01ED11005",
			Bezeichnung:  "2102 Elektro HZ Blechverteiler Altbau",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "kWh",
			Spezialfall:  1,
			Revision:     2,
			GebaeudeRef:  []primitive.ObjectID{gebaeudeOID},
		}

		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(data, compareDoc)
	})

	// Errortests
	t.Run("TestZaehlerFindOID: Kaeltezaehler, OID = 000000000000000000000000 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindOID(primitive.NilObjectID, structs.IDEnergieversorgungKaelte)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	t.Run("TestZaehlerFindOID: Waermezaehler, OID = 000000000000000000000000 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindOID(primitive.NilObjectID, structs.IDEnergieversorgungWaerme)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	t.Run("TestZaehlerFindOID: Stromzaehler, OID = 000000000000000000000000 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data, err := database.ZaehlerFindOID(primitive.NilObjectID, structs.IDEnergieversorgungStrom)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(data, structs.Zaehler{})   // Bei einem Fehler soll ein leer Struct zurückgeliefert werden
	})

	t.Run("TestZaehlerFindOID: IDEnergieversorgung = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		zaehlerOID, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff68e")

		data, err := database.ZaehlerFindOID(zaehlerOID, 0)

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
					{IDDienstreise: 3, Streckentyp: "Langstrecke", Strecke: 321, Klasse: "average"},
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
					{IDDienstreise: 3, Streckentyp: "Langstrecke", Strecke: 321, Klasse: "average"},
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
				{IDDienstreise: 3, Streckentyp: "Langstrecke", Strecke: 321, Klasse: "average"},
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
					Klasse:        "average",
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
