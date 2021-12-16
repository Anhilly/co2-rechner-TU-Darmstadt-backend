package tests

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestInsert(t *testing.T) {
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

		err := database.GebaeudeInsert(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		insertedDoc, err := database.GebaeudeFind(data.Nr)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(insertedDoc, structs.Gebaeude{
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
			Einheit:     "m^2",
			Spezialfall: 1,
			Revision:    1,
			KaelteRef:   []int32{},
			WaermeRef:   []int32{},
			StromRef:    []int32{},
		}) // Ueberpruefung des geaenderten Elementes
	})

	// Errortest
	t.Run("GebaeudeInsert: Nr = 1101 schon vorhanden", func(t *testing.T) {
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

		err := database.GebaeudeInsert(data)
		is.Equal(err, structs.ErrGebaeudeVorhanden) // Funktion wirft ErrGebaeudeVorhanden
	})
}

func TestZaehlerInsert(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("ZaehlerInsert: Waermezaehler, ID = 190", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			PKEnergie:           190,
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{1101},
			IDEnergieversorgung: 1,
		}

		err := database.ZaehlerInsert(data)
		is.NoErr(err)

		neuerZaehler, err := database.ZaehlerFind(data.PKEnergie, data.IDEnergieversorgung)
		is.NoErr(err)
		is.Equal(neuerZaehler, structs.Zaehler{
			Zaehlertyp:   "Waerme",
			PKEnergie:    190,
			Bezeichnung:  "Testzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "kWh",
			Spezialfall:  1,
			Revision:     1,
			GebaeudeRef:  []int32{1101},
		})

		updatedGebaeude, err := database.GebaeudeFind(1101)
		is.NoErr(err)
		is.Equal(updatedGebaeude, structs.Gebaeude{
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
			WaermeRef:   []int32{2084, 190},
			StromRef:    []int32{},
		})
	})

	t.Run("ZaehlerInsert: Stromzaehler, ID = 191", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			PKEnergie:           191,
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{1101},
			IDEnergieversorgung: 2,
		}

		err := database.ZaehlerInsert(data)
		is.NoErr(err)

		neuerZaehler, err := database.ZaehlerFind(data.PKEnergie, data.IDEnergieversorgung)
		is.NoErr(err)
		is.Equal(neuerZaehler, structs.Zaehler{
			Zaehlertyp:   "Strom",
			PKEnergie:    191,
			Bezeichnung:  "Testzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "kWh",
			Spezialfall:  1,
			Revision:     1,
			GebaeudeRef:  []int32{1101},
		})

		updatedGebaeude, err := database.GebaeudeFind(1101)
		is.NoErr(err)
		is.Equal(updatedGebaeude, structs.Gebaeude{
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
			WaermeRef:   []int32{2084, 190},
			StromRef:    []int32{191},
		})
	})

	t.Run("ZaehlerInsert: Kaeltezaehler, ID = 192", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			PKEnergie:           192,
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{1101},
			IDEnergieversorgung: 3,
		}

		err := database.ZaehlerInsert(data)
		is.NoErr(err)

		neuerZaehler, err := database.ZaehlerFind(data.PKEnergie, data.IDEnergieversorgung)
		is.NoErr(err)
		is.Equal(neuerZaehler, structs.Zaehler{
			Zaehlertyp:   "Kaelte",
			PKEnergie:    192,
			Bezeichnung:  "Testzaehler",
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Einheit:      "kWh",
			Spezialfall:  1,
			Revision:     1,
			GebaeudeRef:  []int32{1101},
		})

		updatedGebaeude, err := database.GebaeudeFind(1101)
		is.NoErr(err)
		is.Equal(updatedGebaeude, structs.Gebaeude{
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
			KaelteRef:   []int32{192},
			WaermeRef:   []int32{2084, 190},
			StromRef:    []int32{191},
		})
	})

	// Errortests
	t.Run("ZaehlerInsert: keine Gebaeudereferenzen", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			PKEnergie:           190,
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{},
			IDEnergieversorgung: 1,
		}

		err := database.ZaehlerInsert(data)
		is.Equal(err, structs.ErrFehlendeGebaeuderef) // Funktion wirft ErrFehlendeGebaeudered
	})

	t.Run("ZaehlerInsert: Zaehler vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			PKEnergie:           2107,
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{190},
			IDEnergieversorgung: 1,
		}

		err := database.ZaehlerInsert(data)
		is.Equal(err, structs.ErrZaehlerVorhanden) // Funktion wirft ErrZaehlerVorhanden
	})

	t.Run("ZaehlerInsert: ungueltige Gebaeudereferenz", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			PKEnergie:           14,
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{12},
			IDEnergieversorgung: 1,
		}

		err := database.ZaehlerInsert(data)
		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
	})

	t.Run("ZaehlerInsert: IDEnergieversorgung = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertZaehler{
			PKEnergie:           15,
			Bezeichnung:         "Testzaehler",
			Einheit:             "kWh",
			GebaeudeRef:         []int32{12},
			IDEnergieversorgung: 0,
		}

		err := database.ZaehlerInsert(data)
		is.Equal(err, structs.ErrIDEnergieversorgungNichtVorhanden) // Funktion wirft ErrIDEnergieversorgungNichtVorhanden
	})
}

func TestUmfrageInsert(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("UmfrageInsert: ID nach aktueller Zeitstempel", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertUmfrage{
			Mitarbeiteranzahl: 42,
			Jahr:              3442,
			Gebaeude: []structs.UmfrageGebaeude{
				{GebaeudeNr: 1103, Nutzflaeche: 200},
				{GebaeudeNr: 1105, Nutzflaeche: 200},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 6, Anzahl: 30},
			},
			NutzerEmail: "anton@tobi.com",
		}

		id, err := database.UmfrageInsert(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		insertedDoc, err := database.UmfrageFind(id)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(insertedDoc, structs.Umfrage{
			ID:                id,
			Mitarbeiteranzahl: 42,
			Jahr:              3442,
			Gebaeude: []structs.UmfrageGebaeude{
				{GebaeudeNr: 1103, Nutzflaeche: 200},
				{GebaeudeNr: 1105, Nutzflaeche: 200},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 6, Anzahl: 30},
			},
			Revision:              1,
			MitarbeiterUmfrageRef: []primitive.ObjectID{},
		}) // Ueberpruefung des geaenderten Elementes

		var idVorhanden primitive.ObjectID
		err = idVorhanden.UnmarshalText([]byte("61b23e9855aa64762baf76d7"))
		is.NoErr(err)

		updatedDoc, err := database.NutzerdatenFind(data.NutzerEmail)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(updatedDoc, structs.Nutzerdaten{
			Email:      "anton@tobi.com",
			Passwort:   "test_pw",
			Revision:   1,
			UmfrageRef: []primitive.ObjectID{idVorhanden, id},
		}) // Ueberpruefung des zurueckgelieferten Elements
	})

	// Errortest
	t.Run("UmfrageInsert: ungueltige NutzerEmail", func(t *testing.T) {
		is := is.NewRelaxed(t)

		data := structs.InsertUmfrage{
			Mitarbeiteranzahl: 42,
			Jahr:              3442,
			Gebaeude:          []structs.UmfrageGebaeude{},
			ITGeraete:         []structs.UmfrageITGeraete{},
			NutzerEmail:       "0",
		}

		id, err := database.UmfrageInsert(data)
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
			Mitarbeiteranzahl: 1,
			Jahr:              2020,
			Gebaeude: []structs.UmfrageGebaeude{
				{GebaeudeNr: 1101, Nutzflaeche: 100},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 5, Anzahl: 10},
			},
			Revision:              1,
			MitarbeiterUmfrageRef: []primitive.ObjectID{idVorhanden, idMitarbeiterumfrage},
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
}

func TestNutzerdatenInsert(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("NutzerdatenInsert: {email = 'testingEmailPlsDontUse' password='verysecurepassword'} (nicht vorhanden)", func(t *testing.T) {
		is := is.NewRelaxed(t)
		var email = "testingEmailPlsDontUse"
		testData := structs.AuthReq{
			Username: email,
			Passwort: "verysecurepassword",
		}
		err := database.NutzerdatenInsert(testData)
		is.NoErr(err) // Kein Fehler wird geworfen

		daten, err := database.NutzerdatenFind(email)
		is.NoErr(err) // Kein Fehler seitens der Datenbank
		// Eintrag wurde korrekt hinzugefuegt
		is.Equal(daten.Email, email)
		is.NoErr(bcrypt.CompareHashAndPassword([]byte(daten.Passwort), []byte(testData.Passwort)))
	})

	// Errorfall
	t.Run("NutzerdatenInsert: {email = 'anton@tobi.com' password='verysecurepassword'} (vorhanden)", func(t *testing.T) {
		is := is.NewRelaxed(t)
		testData := structs.AuthReq{
			Username: "anton@tobi.com",
			Passwort: "verysecurepassword",
		}
		err := database.NutzerdatenInsert(testData)
		is.Equal(err, structs.ErrInsertExistingAccount) // Dateneintrag existiert bereits
	})
}