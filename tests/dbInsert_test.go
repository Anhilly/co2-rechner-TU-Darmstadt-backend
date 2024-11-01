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

func TestInsert(t *testing.T) {
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

	t.Run("TestGebaeudeInsert", TestGebaeudeInsert)
	t.Run("TestZaehlerInsert", TestZaehlerInsert)
	t.Run("TestUmfrageInsert", TestUmfrageInsert)
	t.Run("TestMitarbeiterUmfrageInsert", TestMitarbeiterUmfrageInsert)
	t.Run("TestNutzerdatenInsert", TestNutzerdatenInsert)
}

func TestGebaeudeInsert(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("GebaeudeInsert: Nr = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertGebaeude{
			Nr:          0,
			Bezeichnung: "Testgebaeude",
			Flaeche: structs.GebaeudeFlaeche{
				HNF:     1000.0,
				NNF:     1000.0,
				NGF:     1000.0,
				FF:      1000.0,
				VF:      1000.0,
				FreiF:   1000.0,
				GesamtF: 1000.0,
			},
		}

		newGebaeudeOID, err := database.GebaeudeInsert(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		aktuellesJahr := int32(time.Now().Year())
		var versoger []structs.Versoger
		for i := structs.ErstesJahr; i <= aktuellesJahr; i++ {
			versoger = append(versoger, structs.Versoger{
				Jahr:      int32(i),
				IDVertrag: structs.IDVertragTU,
			})
		}

		insertedDoc, err := database.GebaeudeFind(data.Nr)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(insertedDoc, structs.Gebaeude{
			GebaeudeID:  newGebaeudeOID,
			Nr:          0,
			Bezeichnung: "Testgebaeude",
			Flaeche: structs.GebaeudeFlaeche{
				HNF:     1000.0,
				NNF:     1000.0,
				NGF:     1000.0,
				FF:      1000.0,
				VF:      1000.0,
				FreiF:   1000.0,
				GesamtF: 1000.0,
			},
			Einheit:         "m^2",
			Spezialfall:     1,
			Revision:        2,
			Stromversorger:  versoger,
			Waermeversorger: versoger,
			Kaelteversorger: versoger,
			KaelteRef:       []primitive.ObjectID{},
			WaermeRef:       []primitive.ObjectID{},
			StromRef:        []primitive.ObjectID{},
		}) // Ueberpruefung des geaenderten Elementes
	})

	t.Run("GebaeudeInsert: Nr = 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertGebaeude{
			Nr:          1,
			Bezeichnung: "Testgebaeude",
			Flaeche: structs.GebaeudeFlaeche{
				HNF:     1000.0,
				NNF:     1000.0,
				NGF:     1000.0,
				FF:      1000.0,
				VF:      1000.0,
				FreiF:   1000.0,
				GesamtF: 1000.0,
			},
			WaermeVersorgerJahre: []int32{2020, 2019, 2018},
			KaelteVersorgerJahre: []int32{2021},
			StromVersorgerJahre:  []int32{2023},
		}

		newGebaeudeOID, err := database.GebaeudeInsert(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		aktuellesJahr := int32(time.Now().Year())
		var stromversoger []structs.Versoger
		var waermeversoger []structs.Versoger
		var kaelteversoger []structs.Versoger
		for i := structs.ErstesJahr; i <= aktuellesJahr; i++ {
			stromversoger = append(stromversoger, structs.Versoger{
				Jahr:      int32(i),
				IDVertrag: structs.IDVertragTU,
			})

			waermeversoger = append(waermeversoger, structs.Versoger{
				Jahr:      int32(i),
				IDVertrag: structs.IDVertragTU,
			})

			kaelteversoger = append(kaelteversoger, structs.Versoger{
				Jahr:      int32(i),
				IDVertrag: structs.IDVertragTU,
			})
		}

		compareDoc := structs.Gebaeude{
			GebaeudeID:  newGebaeudeOID,
			Nr:          1,
			Bezeichnung: "Testgebaeude",
			Flaeche: structs.GebaeudeFlaeche{
				HNF:     1000.0,
				NNF:     1000.0,
				NGF:     1000.0,
				FF:      1000.0,
				VF:      1000.0,
				FreiF:   1000.0,
				GesamtF: 1000.0,
			},
			Einheit:         "m^2",
			Spezialfall:     1,
			Revision:        2,
			Stromversorger:  stromversoger,
			Waermeversorger: waermeversoger,
			Kaelteversorger: kaelteversoger,
			KaelteRef:       []primitive.ObjectID{},
			WaermeRef:       []primitive.ObjectID{},
			StromRef:        []primitive.ObjectID{},
		}

		compareDoc.Waermeversorger[0].IDVertrag = 2
		compareDoc.Kaelteversorger[1].IDVertrag = 2
		compareDoc.Stromversorger[3].IDVertrag = 2

		insertedDoc, err := database.GebaeudeFind(data.Nr)
		is.NoErr(err)                     // kein Error seitens der Datenbank
		is.Equal(insertedDoc, compareDoc) // Ueberpruefung des geaenderten Elementes
	})

	// Errortest
	t.Run("GebaeudeInsert: Nr = 1101 schon vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertGebaeude{
			Nr:          1101,
			Bezeichnung: "Testgebaeude",
			Flaeche: structs.GebaeudeFlaeche{
				HNF:     1000.0,
				NNF:     1000.0,
				NGF:     1000.0,
				FF:      1000.0,
				VF:      1000.0,
				FreiF:   1000.0,
				GesamtF: 1000.0,
			},
		}

		oid, err := database.GebaeudeInsert(data)
		is.Equal(err, structs.ErrGebaeudeVorhanden) // Funktion wirft ErrGebaeudeVorhanden
		is.Equal(oid, primitive.NilObjectID)
	})
}

func TestZaehlerInsert(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("ZaehlerInsert: Waermezaehler, DPName = 190acd", func(t *testing.T) {
		is := is.NewRelaxed(t)

		// hole unveraendertes Gebaeude
		gebaeude, err := database.GebaeudeFind(1101)
		is.NoErr(err)

		// fuege neuen Zaehler hinzu
		data := structs.InsertZaehler{
			DPName:              "190acd",
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{1101},
			IDEnergieversorgung: 1,
		}

		newZaehlerOID, err := database.ZaehlerInsert(data)
		is.NoErr(err)

		// Vergleichszaehler wird dynamisch erstellt, weil eingetragene Zaehlerdaten vom aktuellen Jahr abhaengen
		location, _ := time.LoadLocation("Etc/GMT")
		aktuellesJahr := int32(time.Now().Year())
		gebaeudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a286")

		vergleichszaehler := structs.Zaehler{
			ZaehlerID:    newZaehlerOID,
			Zaehlertyp:   "Waerme",
			DPName:       "190acd",
			Bezeichnung:  "Testzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "kWh",
			Spezialfall:  1,
			Revision:     2,
			GebaeudeRef:  []primitive.ObjectID{gebaeudeOID},
		}
		for i := structs.ErstesJahr; i <= aktuellesJahr; i++ {
			for j := 1; j <= 12; j++ {
				vergleichszaehler.Zaehlerdaten = append(vergleichszaehler.Zaehlerdaten, structs.Zaehlerwerte{
					Wert:        0.0,
					Zeitstempel: time.Date(int(i), time.Month(j), 01, 0, 0, 0, 0, location).UTC(),
				})
			}
		}

		neuerZaehler, err := database.ZaehlerFindDPName(data.DPName, data.IDEnergieversorgung)
		is.NoErr(err)
		is.Equal(neuerZaehler, vergleichszaehler)

		gebaeude.WaermeRef = append(gebaeude.WaermeRef, newZaehlerOID)

		updatedGebaeude, err := database.GebaeudeFind(1101)
		is.NoErr(err)
		is.Equal(updatedGebaeude, gebaeude)
	})

	t.Run("ZaehlerInsert: Stromzaehler, DPName = 191acd", func(t *testing.T) {
		is := is.NewRelaxed(t)

		// hole unveraendertes Gebaeude
		gebaeude, err := database.GebaeudeFind(1102)
		is.NoErr(err)

		// fuege neuen Zaehler hinzu
		data := structs.InsertZaehler{
			DPName:              "191acd",
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{1102},
			IDEnergieversorgung: 2,
		}

		newZaehlerOID, err := database.ZaehlerInsert(data)
		is.NoErr(err)

		// Vergleichszaehler wird dynamisch erstellt, weil eingetragene Zaehlerdaten vom aktuellen Jahr abhaengen
		location, _ := time.LoadLocation("Etc/GMT")
		aktuellesJahr := int32(time.Now().Year())
		gebaeudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a287")

		vergleichszaehler := structs.Zaehler{
			ZaehlerID:    newZaehlerOID,
			Zaehlertyp:   "Strom",
			DPName:       "191acd",
			Bezeichnung:  "Testzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "kWh",
			Spezialfall:  1,
			Revision:     2,
			GebaeudeRef:  []primitive.ObjectID{gebaeudeOID},
		}
		for i := structs.ErstesJahr; i <= aktuellesJahr; i++ {
			for j := 1; j <= 12; j++ {
				vergleichszaehler.Zaehlerdaten = append(vergleichszaehler.Zaehlerdaten, structs.Zaehlerwerte{
					Wert:        0.0,
					Zeitstempel: time.Date(int(i), time.Month(j), 01, 0, 0, 0, 0, location).UTC(),
				})
			}
		}

		neuerZaehler, err := database.ZaehlerFindDPName(data.DPName, data.IDEnergieversorgung)
		is.NoErr(err)
		is.Equal(neuerZaehler, vergleichszaehler)

		gebaeude.StromRef = append(gebaeude.StromRef, newZaehlerOID)

		updatedGebaeude, err := database.GebaeudeFind(1102)
		is.NoErr(err)
		is.Equal(updatedGebaeude, gebaeude)
	})

	t.Run("ZaehlerInsert: Kaeltezaehler, DPName = 192acd", func(t *testing.T) {
		is := is.NewRelaxed(t)

		// hole unveraendertes Gebaeude
		gebaeude, err := database.GebaeudeFind(1103)
		is.NoErr(err)

		// fuege neuen Zaehler hinzu
		data := structs.InsertZaehler{
			DPName:              "192acd",
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{1103},
			IDEnergieversorgung: 3,
		}

		newZaehlerOID, err := database.ZaehlerInsert(data)
		is.NoErr(err)

		// Vergleichszaehler wird dynamisch erstellt, weil eingetragene Zaehlerdaten vom aktuellen Jahr abhaengen
		location, _ := time.LoadLocation("Etc/GMT")
		aktuellesJahr := int32(time.Now().Year())
		gebaeudeOID, _ := primitive.ObjectIDFromHex("61b712f66a1a52dea358a288")

		vergleichszaehler := structs.Zaehler{
			ZaehlerID:    newZaehlerOID,
			Zaehlertyp:   "Kaelte",
			DPName:       "192acd",
			Bezeichnung:  "Testzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "kWh",
			Spezialfall:  1,
			Revision:     2,
			GebaeudeRef:  []primitive.ObjectID{gebaeudeOID},
		}
		for i := structs.ErstesJahr; i <= aktuellesJahr; i++ {
			for j := 1; j <= 12; j++ {
				vergleichszaehler.Zaehlerdaten = append(vergleichszaehler.Zaehlerdaten, structs.Zaehlerwerte{
					Wert:        0.0,
					Zeitstempel: time.Date(int(i), time.Month(j), 01, 0, 0, 0, 0, location).UTC(),
				})
			}
		}

		neuerZaehler, err := database.ZaehlerFindDPName(data.DPName, data.IDEnergieversorgung)
		is.NoErr(err)
		is.Equal(neuerZaehler, vergleichszaehler)

		gebaeude.KaelteRef = append(gebaeude.KaelteRef, newZaehlerOID)

		updatedGebaeude, err := database.GebaeudeFind(1103)
		is.NoErr(err)
		is.Equal(updatedGebaeude, gebaeude)
	})

	// Errortests
	t.Run("ZaehlerInsert: keine Gebaeudereferenzen", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			DPName:              "190",
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{},
			IDEnergieversorgung: 1,
		}

		oid, err := database.ZaehlerInsert(data)
		is.Equal(err, structs.ErrFehlendeGebaeuderef) // Funktion wirft ErrFehlendeGebaeudered
		is.Equal(oid, primitive.NilObjectID)
	})

	t.Run("ZaehlerInsert: Zaehler vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			DPName:              "B101XXXXXXHE000XXXXXXZ40CO00001",
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{190},
			IDEnergieversorgung: 1,
		}

		oid, err := database.ZaehlerInsert(data)
		is.Equal(err, structs.ErrZaehlerVorhanden) // Funktion wirft ErrZaehlerVorhanden
		is.Equal(oid, primitive.NilObjectID)
	})

	t.Run("ZaehlerInsert: ungueltige Gebaeudereferenz", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			DPName:              "78",
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{12},
			IDEnergieversorgung: 1,
		}

		oid, err := database.ZaehlerInsert(data)
		is.Equal(err, structs.ErrGebaeudeNichtVorhanden) // Funktion wirft ErrGebaeudeNichtVorhanden
		is.Equal(oid, primitive.NilObjectID)
	})

	t.Run("ZaehlerInsert: IDEnergieversorgung = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			DPName:              "780",
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{12},
			IDEnergieversorgung: 0,
		}

		oid, err := database.ZaehlerInsert(data)
		is.Equal(err, structs.ErrIDEnergieversorgungNichtVorhanden) // Funktion wirft ErrIDEnergieversorgungNichtVorhanden
		is.Equal(oid, primitive.NilObjectID)
	})
}

func TestUmfrageInsert(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("UmfrageInsert: ID nach aktueller Zeitstempel", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "anton"

		data := structs.InsertUmfrage{
			Bezeichnung:       "TestUmfrageInsert",
			Mitarbeiteranzahl: 42,
			Jahr:              3442,
			Gebaeude: []structs.UmfrageGebaeude{
				{GebaeudeNr: 1103, Nutzflaeche: 200},
				{GebaeudeNr: 1105, Nutzflaeche: 200},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 6, Anzahl: 30},
			},
		}

		id, err := database.UmfrageInsert(data, username)
		is.NoErr(err) // kein Error seitens der Datenbank

		insertedDoc, err := database.UmfrageFind(id)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(insertedDoc, structs.Umfrage{
			ID:                id,
			Bezeichnung:       "TestUmfrageInsert",
			Mitarbeiteranzahl: 42,
			Jahr:              3442,
			Gebaeude: []structs.UmfrageGebaeude{
				{GebaeudeNr: 1103, Nutzflaeche: 200},
				{GebaeudeNr: 1105, Nutzflaeche: 200},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 6, Anzahl: 30},
			},
			AuswertungFreigegeben: 0,
			Revision:              1,
			MitarbeiterumfrageRef: []primitive.ObjectID{},
		}) // Ueberpruefung des geaenderten Elementes

		var idVorhanden primitive.ObjectID
		err = idVorhanden.UnmarshalText([]byte("61b23e9855aa64762baf76d7"))
		is.NoErr(err)

		objID, _ := primitive.ObjectIDFromHex("61b1ceb3dfb93b34b1305b70")

		updatedDoc, err := database.NutzerdatenFind(username)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(updatedDoc, structs.Nutzerdaten{
			NutzerID:   objID,
			Nutzername: "anton",
			EMail:      "anton@tobi.com",
			Rolle:      0,
			Revision:   2,
			UmfrageRef: []primitive.ObjectID{idVorhanden, id},
		}) // Ueberpruefung des zurueckgelieferten Elements
	})

	// Errortest
	t.Run("UmfrageInsert: ungueltiger Username", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "o123"

		data := structs.InsertUmfrage{
			Mitarbeiteranzahl: 42,
			Jahr:              3442,
			Gebaeude:          []structs.UmfrageGebaeude{},
			ITGeraete:         []structs.UmfrageITGeraete{},
		}

		id, err := database.UmfrageInsert(data, username)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(id, primitive.NilObjectID) // im Fehlerfall wird NilObjectID zurueckgegeben
	})
}

func TestMitarbeiterUmfrageInsert(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("MitarbeiterUmfrageInsert: ID nach aktueller Zeitstempel", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idUmfrage primitive.ObjectID
		err := idUmfrage.UnmarshalText([]byte("61b23e9855aa64762baf76d7"))
		is.NoErr(err)

		data := structs.InsertMitarbeiterUmfrage{
			Pendelweg: []structs.UmfragePendelweg{
				{IDPendelweg: 2, Strecke: 20},
			},
			TageImBuero: 4,
			Dienstreise: []structs.UmfrageDienstreise{
				{IDDienstreise: 1, Strecke: 100},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 6, Anzahl: 30},
			},
			IDUmfrage: idUmfrage,
		}

		idMitarbeiterumfrage, err := database.MitarbeiterUmfrageInsert(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		insertedDoc, err := database.MitarbeiterUmfrageFind(idMitarbeiterumfrage)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(insertedDoc, structs.MitarbeiterUmfrage{
			ID: idMitarbeiterumfrage,
			Pendelweg: []structs.UmfragePendelweg{
				{IDPendelweg: 2, Strecke: 20},
			},
			TageImBuero: 4,
			Dienstreise: []structs.UmfrageDienstreise{
				{IDDienstreise: 1, Strecke: 100},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 6, Anzahl: 30},
			},
			Revision: 1,
		}) // Ueberpruefung des zurueckgelieferten Elementes

		var idVorhanden primitive.ObjectID
		err = idVorhanden.UnmarshalText([]byte("61b34f9324756df01eee5ff4"))
		is.NoErr(err)

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
			MitarbeiterumfrageRef: []primitive.ObjectID{idVorhanden, idMitarbeiterumfrage},
		}) // Ueberpruefung des zurueckgelieferten Elements
	})

	// Errortest
	t.Run("MitarbeiterUmfrageInsert: ungueltige UmfrageID", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertMitarbeiterUmfrage{
			Pendelweg: []structs.UmfragePendelweg{
				{IDPendelweg: 2, Strecke: 20},
			},
			TageImBuero: 4,
			Dienstreise: []structs.UmfrageDienstreise{},
			ITGeraete:   []structs.UmfrageITGeraete{},
			IDUmfrage:   primitive.NewObjectID(),
		}

		id, err := database.MitarbeiterUmfrageInsert(data)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(id, primitive.NilObjectID) // im Fehlerfall wird NilObjectID zurueckgegeben
	})

	t.Run("MitarbeiterUmfrageInsert: Umfrage vollstaendig", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idUmfrage primitive.ObjectID
		err := idUmfrage.UnmarshalText([]byte("61dc0c543e48998484eefaeb"))

		data := structs.InsertMitarbeiterUmfrage{
			Pendelweg: []structs.UmfragePendelweg{
				{IDPendelweg: 2, Strecke: 20},
			},
			TageImBuero: 4,
			Dienstreise: []structs.UmfrageDienstreise{},
			ITGeraete:   []structs.UmfrageITGeraete{},
			IDUmfrage:   idUmfrage,
		}

		id, err := database.MitarbeiterUmfrageInsert(data)
		is.Equal(err, structs.ErrUmfrageVollstaendig) // Datenbank wirft ErrUmfrageVollstaendig
		is.Equal(id, primitive.NilObjectID)           // im Fehlerfall wird NilObjectID zurueckgegeben
	})
}

func TestNutzerdatenInsert(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("NutzerdatenInsert: {username = 'testingUserPlsDontUse' password='verysecurepassword'} (nicht vorhanden)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "testingUserPlsDontUse"
		email := "dont-reply@test.com"

		id, err := database.NutzerdatenInsert(username, email)
		is.NoErr(err) // Kein Fehler wird geworfen

		daten, err := database.NutzerdatenFind(username)
		is.NoErr(err) // Kein Fehler seitens der Datenbank
		// Eintrag wurde korrekt hinzugefuegt
		is.Equal(daten, structs.Nutzerdaten{
			NutzerID:   id,
			EMail:      email,
			Nutzername: username,
			Rolle:      0,
			Revision:   2,
			UmfrageRef: []primitive.ObjectID{},
		})
	})

	// Errorfall
	t.Run("NutzerdatenInsert: {username = 'anton@tobi.com' password='verysecurepassword'} (vorhanden)", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "anton"
		email := "anton@tobi.com"

		_, err := database.NutzerdatenInsert(username, email)
		is.Equal(err, structs.ErrInsertExistingAccount) // Dateneintrag existiert bereits
	})
}
