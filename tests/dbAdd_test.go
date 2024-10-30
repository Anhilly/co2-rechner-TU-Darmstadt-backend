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

		// hole altes Dokument
		compareDoc, err := database.EnergieversorgungFind(1)
		is.NoErr(err)

		compareDoc.CO2Faktor = append(compareDoc.CO2Faktor, structs.CO2Energie{
			Jahr: 2030,
			Vertraege: []structs.CO2FaktorVetrag{
				{IDVertrag: 1, Wert: 400},
			},
		})

		// update Dokument per Funktion
		data := structs.AddCO2Faktor{
			IDEnergieversorgung: 1,
			IDVertrag:           1,
			Wert:                400,
			Jahr:                2030,
		}

		err = database.EnergieversorgungAddFaktor(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.EnergieversorgungFind(data.IDEnergieversorgung)
		is.NoErr(err)
		is.Equal(updatedDoc, compareDoc) // Ueberpruefung des geaenderten Elementes
	})

	t.Run("EnergieversorgungAddFaktor: IDVertrag = 2", func(t *testing.T) {
		is := is.NewRelaxed(t)

		// hole altes Dokument
		compareDoc, err := database.EnergieversorgungFind(2)
		is.NoErr(err)

		compareDoc.CO2Faktor = append(compareDoc.CO2Faktor, structs.CO2Energie{
			Jahr: 2030,
			Vertraege: []structs.CO2FaktorVetrag{
				{IDVertrag: 2, Wert: 400},
			},
		})

		// update Dokument per Funktion
		data := structs.AddCO2Faktor{
			IDEnergieversorgung: 2,
			IDVertrag:           2,
			Wert:                400,
			Jahr:                2030,
		}

		err = database.EnergieversorgungAddFaktor(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.EnergieversorgungFind(data.IDEnergieversorgung)
		is.NoErr(err)
		is.Equal(updatedDoc, compareDoc) // Ueberpruefung des geaenderten Elementes
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
	t.Run("ZaehlerAddZaehlerdaten: Waermezaehler, DPName = B101XXXXXXHE000XXXXXXZ40CO00001", func(t *testing.T) {
		is := is.NewRelaxed(t)

		// hole altes Dokument
		compareDoc, err := database.ZaehlerFindDPName("B101XXXXXXHE000XXXXXXZ40CO00001", structs.IDEnergieversorgungWaerme)
		is.NoErr(err)

		location, _ := time.LoadLocation("Etc/GMT")
		compareDoc.Zaehlerdaten = append(compareDoc.Zaehlerdaten, structs.Zaehlerwerte{
			Wert:        1000.0,
			Zeitstempel: time.Date(3001, time.May, 01, 0, 0, 0, 0, location).UTC(),
		})

		// update Dokument per Funktion
		data := structs.AddZaehlerdaten{
			DPName:              "B101XXXXXXHE000XXXXXXZ40CO00001",
			Wert:                1000.0,
			Jahr:                3001,
			Monat:               5,
			IDEnergieversorgung: 1,
		}

		err = database.ZaehlerAddZaehlerdaten(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.ZaehlerFindDPName(data.DPName, data.IDEnergieversorgung)
		is.NoErr(err)                    // Datenbank wirft keinen Fehler
		is.Equal(updatedDoc, compareDoc) // Ueberpruefung des geaenderten Elementes
	})

	t.Run("ZaehlerAddZaehlerdaten: Stromzaehler, ID = B102xxxxxxNA000xxxxxxZ01ED11005", func(t *testing.T) {
		is := is.NewRelaxed(t)

		// hole altes Dokument
		compareDoc, err := database.ZaehlerFindDPName("B102xxxxxxNA000xxxxxxZ01ED11005", structs.IDEnergieversorgungStrom)
		is.NoErr(err)

		location, _ := time.LoadLocation("Etc/GMT")
		compareDoc.Zaehlerdaten = append(compareDoc.Zaehlerdaten, structs.Zaehlerwerte{
			Wert:        1000.0,
			Zeitstempel: time.Date(3001, time.November, 01, 0, 0, 0, 0, location).UTC(),
		})

		// update Dokument per Funktion
		data := structs.AddZaehlerdaten{
			DPName:              "B102xxxxxxNA000xxxxxxZ01ED11005",
			Wert:                1000.0,
			Jahr:                3001,
			Monat:               11,
			IDEnergieversorgung: 2,
		}

		err = database.ZaehlerAddZaehlerdaten(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.ZaehlerFindDPName(data.DPName, data.IDEnergieversorgung)
		is.NoErr(err)
		is.Equal(updatedDoc, compareDoc) // Ueberpruefung des geaenderten Elementes
	})

	t.Run("ZaehlerAddZaehlerdaten: Kaeltezaehler, DPName = L101XXXXXXKA000XXXXXXZ50CO00001", func(t *testing.T) {
		is := is.NewRelaxed(t)

		// hole altes Dokument
		compareDoc, err := database.ZaehlerFindDPName("L101XXXXXXKA000XXXXXXZ50CO00001", structs.IDEnergieversorgungKaelte)
		is.NoErr(err)

		location, _ := time.LoadLocation("Etc/GMT")
		compareDoc.Zaehlerdaten = append(compareDoc.Zaehlerdaten, structs.Zaehlerwerte{
			Wert:        1000.0,
			Zeitstempel: time.Date(3001, time.February, 01, 0, 0, 0, 0, location).UTC(),
		})

		// update Dokument per Funktion
		data := structs.AddZaehlerdaten{
			DPName:              "L101XXXXXXKA000XXXXXXZ50CO00001",
			Wert:                1000.0,
			Jahr:                3001,
			Monat:               2,
			IDEnergieversorgung: 3,
		}

		err = database.ZaehlerAddZaehlerdaten(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.ZaehlerFindDPName(data.DPName, data.IDEnergieversorgung)
		is.NoErr(err)
		is.Equal(updatedDoc, compareDoc) // Ueberpruefung des geaenderten Elementes
	})

	t.Run("ZaehlerAddZaehlerdaten: Ueberschreiben von Wert 0.0 bei Zaehler", func(t *testing.T) {
		is := is.NewRelaxed(t)

		// hole altes Dokument
		compareDoc, err := database.ZaehlerFindDPName("B101XXXXXXHE000XXXXXXZ40CO00001", structs.IDEnergieversorgungWaerme)
		is.NoErr(err)

		location, _ := time.LoadLocation("Etc/GMT")
		compareDoc.Zaehlerdaten = append(compareDoc.Zaehlerdaten, structs.Zaehlerwerte{
			Wert:        1030.0,
			Zeitstempel: time.Date(3005, time.May, 01, 0, 0, 0, 0, location).UTC(),
		})

		// update Dokument per Funktion
		data := structs.AddZaehlerdaten{
			DPName:              "B101XXXXXXHE000XXXXXXZ40CO00001",
			Wert:                0.0,
			Jahr:                3005,
			Monat:               5,
			IDEnergieversorgung: 1,
		}

		err = database.ZaehlerAddZaehlerdaten(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		data = structs.AddZaehlerdaten{
			DPName:              "B101XXXXXXHE000XXXXXXZ40CO00001",
			Wert:                1030.0,
			Jahr:                3005,
			Monat:               5,
			IDEnergieversorgung: 1,
		}

		err = database.ZaehlerAddZaehlerdaten(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, _ := database.ZaehlerFindDPName(data.DPName, data.IDEnergieversorgung)

		is.Equal(updatedDoc, compareDoc) // Ueberpruefung des geaenderten Elementes
	})

	// Errortests
	t.Run("ZaehlerAddZaehlerdaten: IDEnergieversorgung = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			DPName:              "B101XXXXXXHE000XXXXXXZ40CO00001",
			Wert:                550.0,
			Jahr:                3006,
			Monat:               5,
			IDEnergieversorgung: 0,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, structs.ErrIDEnergieversorgungNichtVorhanden) // Funktion wirft ErrIDEnergieversorgungNichtVorhanden
	})

	t.Run("ZaehlerAddZaehlerdaten: Waermezaehler, DPName nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			DPName:              "B101XXXXXXHE000XXXXXXZ40CO00",
			Wert:                550.0,
			Jahr:                3006,
			Monat:               5,
			IDEnergieversorgung: 1,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})

	t.Run("ZaehlerAddZaehlerdaten: Waermezaehler, Jahr und Monat schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			DPName:              "B101XXXXXXHE000XXXXXXZ40CO00001",
			Wert:                1000.0,
			Jahr:                3001,
			Monat:               5,
			IDEnergieversorgung: 1,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, structs.ErrJahrUndMonatVorhanden) // Funktion wirft ErrJahrVorhanden
	})

	t.Run("ZaehlerAddZaehlerdaten: Stromzaehler, DPName nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			DPName:              "B101XXXXXXHE000XXXXXXZ40CO00",
			Wert:                550.0,
			Jahr:                3006,
			Monat:               5,
			IDEnergieversorgung: 2,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})

	t.Run("ZaehlerAddZaehlerdaten: Stromzaehler, Jahr und Monat schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			DPName:              "B102xxxxxxNA000xxxxxxZ01ED11005",
			Wert:                1000.0,
			Jahr:                3001,
			Monat:               11,
			IDEnergieversorgung: 2,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, structs.ErrJahrUndMonatVorhanden) // Funktion wirft ErrJahrVorhanden
	})

	t.Run("ZaehlerAddZaehlerdaten: Kaeltezaehler, DPName nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			DPName:              "B101XXXXXXHE000XXXXXXZ40CO00",
			Wert:                550.0,
			Jahr:                3006,
			Monat:               5,
			IDEnergieversorgung: 3,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})

	t.Run("ZaehlerAddZaehlerdaten: Kaeltezaehler, Jahr und Monat schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.AddZaehlerdaten{
			DPName:              "L101XXXXXXKA000XXXXXXZ50CO00001",
			Wert:                1000.0,
			Jahr:                3001,
			Monat:               2,
			IDEnergieversorgung: 3,
		}

		err := database.ZaehlerAddZaehlerdaten(data)
		is.Equal(err, structs.ErrJahrUndMonatVorhanden) // Funktion wirft ErrJahrVorhanden
	})
}

func TestZaehlerAddStandardZaehlerdaten(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("AddStandardZaehlerdaten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		// hole altes Dokument
		compareDoc, err := database.ZaehlerFindDPName("B101XXXXXXHE000XXXXXXZ40CO00001", structs.IDEnergieversorgungWaerme)
		is.NoErr(err)

		location, _ := time.LoadLocation("Etc/GMT")
		compareDoc.Zaehlerdaten = append(compareDoc.Zaehlerdaten, structs.Zaehlerwerte{
			Wert:        0.0,
			Zeitstempel: time.Date(3000, time.August, 01, 0, 0, 0, 0, location).UTC(),
		})

		// update Dokument per Funktion
		data := structs.AddStandardZaehlerdaten{
			Jahr:  3000,
			Monat: 8,
		}

		err = database.ZaehlerAddStandardZaehlerdaten(data)
		is.NoErr(err) // Datenbank wirft keinen Fehler

		updatedDoc, err := database.ZaehlerFindDPName("B101XXXXXXHE000XXXXXXZ40CO00001", structs.IDEnergieversorgungWaerme)
		is.NoErr(err)                    // Datenbank wirft keinen Fehler
		is.Equal(updatedDoc, compareDoc) // Ueberpruefung des geaenderten Elementes)
	})
}

func TestGebaeudeAddZaehlerref(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("GebaeudeAddZaehlerref: ID = 1101, idEnergieversorgung = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var nr int32 = 1101
		ref, _ := primitive.ObjectIDFromHex("B203XXXXXXHE000XXXXXXZ40CO00002")
		var idEnergieversorgung int32 = 1

		// hole altes Dokument
		compareDoc, err := database.GebaeudeFind(nr)
		is.NoErr(err)

		compareDoc.WaermeRef = append(compareDoc.WaermeRef, ref)

		// update Dokument per Funktion
		err = database.GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.GebaeudeFind(nr)
		is.NoErr(err)
		is.Equal(updatedDoc, compareDoc) // Ueberpruefung des zurueckgelieferten Elements
	})

	t.Run("GebaeudeAddZaehlerref: ID = 1102, idEnergieversorgung = 2", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var nr int32 = 1102
		ref, _ := primitive.ObjectIDFromHex("B203XXXXXXHE000XXXXXXZ40CO00002")
		var idEnergieversorgung int32 = 2

		// hole altes Dokument
		compareDoc, err := database.GebaeudeFind(nr)
		is.NoErr(err)

		compareDoc.StromRef = append(compareDoc.StromRef, ref)

		// update Dokument per Funktion
		err = database.GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.GebaeudeFind(nr)
		is.NoErr(err)
		is.Equal(updatedDoc, compareDoc) // Ueberpruefung des zurueckgelieferten Elements
	})

	t.Run("GebaeudeAddZaehlerref: ID = 1103, idEnergieversorgung = 3", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var nr int32 = 1103
		ref, _ := primitive.ObjectIDFromHex("B203XXXXXXHE000XXXXXXZ40CO00002")
		var idEnergieversorgung int32 = 3

		// hole altes Dokument
		compareDoc, err := database.GebaeudeFind(nr)
		is.NoErr(err)

		compareDoc.KaelteRef = append(compareDoc.KaelteRef, ref)

		// update Dokument per Funktion
		err = database.GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.GebaeudeFind(nr)
		is.NoErr(err)
		is.Equal(updatedDoc, compareDoc) // Ueberpruefung des zurueckgelieferten Elements
	})

	// dieser Fall sollte nicht auftreten
	t.Run("GebaeudeAddZaehlerref: Doppelte Referenz", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var nr int32 = 1108
		ref, _ := primitive.ObjectIDFromHex("6710a1ca3a47b613426ff65a")
		var idEnergieversorgung int32 = 1

		// hole altes Dokument
		compareDoc, _ := database.GebaeudeFind(nr)

		// update Dokument per Funktion
		err := database.GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.GebaeudeFind(nr)
		is.NoErr(err)
		is.Equal(updatedDoc, compareDoc) // Ueberpruefung des zurueckgelieferten Elements
	})

	// Errortests
	t.Run("GebaeudeAddZaehlerref: ID = -12 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var nr int32 = -12
		ref := primitive.NewObjectID()
		var idEnergieversorgung int32 = 3

		err := database.GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})

	t.Run("GebaeudeAddZaehlerref: IDEnergieversorgung = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var nr int32 = 1101
		ref := primitive.NewObjectID()
		var idEnergieversorgung int32 = 0

		err := database.GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung)
		is.Equal(err, structs.ErrIDEnergieversorgungNichtVorhanden) // Datenbank wirft ErrNoDocuments
	})
}

func TestGebaeudeAddVersorger(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("AddVersorger: Gebaeude Nr 1102", func(t *testing.T) {
		is := is.NewRelaxed(t)

		// hole altes Dokument
		compareDoc, err := database.GebaeudeFind(1102)
		is.NoErr(err)

		compareDoc.Waermeversorger = append(compareDoc.Waermeversorger, structs.Versoger{
			Jahr:      3005,
			IDVertrag: 2,
		})

		// update Dokument per Funktion
		data := structs.AddVersorger{
			Nr:                  1102,
			Jahr:                3005,
			IDEnergieversorgung: 1,
			IDVertrag:           2,
		}

		err = database.GebaeudeAddVersorger(data)
		is.NoErr(err) // Datenbank wirft keinen Fehler

		gebaeude, err := database.GebaeudeFind(1102)
		is.NoErr(err)
		is.Equal(gebaeude, compareDoc) // Ueberpruefung des zurueckgelieferten Elements
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
			Nr:                  1102,
			Jahr:                3005,
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

		// hole altes Dokument
		compareDoc, err := database.GebaeudeFind(1101)
		is.NoErr(err)

		compareDoc.Stromversorger = append(compareDoc.Stromversorger, structs.Versoger{
			Jahr:      3000,
			IDVertrag: 1,
		})
		compareDoc.Waermeversorger = append(compareDoc.Waermeversorger, structs.Versoger{
			Jahr:      3000,
			IDVertrag: 1,
		})
		compareDoc.Kaelteversorger = append(compareDoc.Kaelteversorger, structs.Versoger{
			Jahr:      3000,
			IDVertrag: 1,
		})

		// update Dokument per Funktion
		data := structs.AddStandardVersorger{
			Jahr: 3000,
		}

		err = database.GebaeudeAddStandardVersorger(data)
		is.NoErr(err) // Datenbank wirft keinen Fehler

		updatedDoc, err := database.GebaeudeFind(1101)
		is.NoErr(err)
		is.Equal(updatedDoc, compareDoc) // Überprüfung des zurückgelieferten Elements
	})
}

func TestNutzerdatenAddUmfrageref(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("NutzerdatenAddUmfrageref: username = anton", func(t *testing.T) {
		is.NewRelaxed(t)

		username := "anton"
		id := primitive.NewObjectID()

		// hole altes Dokument
		compareDoc, err := database.NutzerdatenFind(username)
		is.NoErr(err)

		compareDoc.UmfrageRef = append(compareDoc.UmfrageRef, id)

		// update Dokument per Funktion
		err = database.NutzerdatenAddUmfrageref(username, id)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.NutzerdatenFind(username)
		is.NoErr(err)                    // kein Error seitens der Datenbank
		is.Equal(updatedDoc, compareDoc) // Überprüfung des zurückgelieferten Elements
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

		idUmfrage, _ := primitive.ObjectIDFromHex("61b23e9855aa64762baf76d7")
		referenz := primitive.NewObjectID()

		// hole altes Dokument
		compareDoc, err := database.UmfrageFind(idUmfrage)
		is.NoErr(err)

		compareDoc.MitarbeiterumfrageRef = append(compareDoc.MitarbeiterumfrageRef, referenz)

		// update Dokument per Funktion
		err = database.UmfrageAddMitarbeiterUmfrageRef(idUmfrage, referenz)
		is.NoErr(err) // kein Error seitens der Datenbank

		updatedDoc, err := database.UmfrageFind(idUmfrage)
		is.NoErr(err)                    // kein Error seitens der Datenbank
		is.Equal(updatedDoc, compareDoc) // Überprüfung des zurückgelieferten Elements
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
