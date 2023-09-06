package tests

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/server"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestDelete(t *testing.T) {
	is := is.NewRelaxed(t)

	err := database.ConnectDatabase("dev")
	is.NoErr(err)

	dir, err := database.CreateDump("TestDelete")
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

	t.Run("TestUmfrageDelete", TestUmfrageDelete)
	t.Run("TestUmfrageDeleteMitarbeiterUmfrage", TestUmfrageDeleteMitarbeiterUmfrage)
	t.Run("TestNutzerdatenDelete", TestNutzerdatenDelete)
}

func TestUmfrageDelete(t *testing.T) {
	is := is.NewRelaxed(t)

	//Normalfall
	t.Run("UmfrageDelete umfrageID vorhanden ohne MitarbeiterUmfragen", func(t *testing.T) {
		is := is.NewRelaxed(t)
		username := "anton@tobi.com"

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
			Auth: structs.AuthToken{
				Username:     username,
				Sessiontoken: server.GeneriereSessionToken(username),
			},
		}

		objectID, err := database.UmfrageInsert(data) // Neuen Eintrag erstellen
		is.NoErr(err)                                 // Kein Fehler im Normalfall

		err = database.UmfrageDelete(username, objectID) // Eintrag loeschen
		is.NoErr(err)                                    // Kein Fehler im Normalfall

		_, err = database.UmfrageFind(objectID) // Eintrag wird nicht mehr gefunden?
		is.Equal(err, mongo.ErrNoDocuments)     // Datenbank wirft ErrNoDocuments

	})

	t.Run("UmfrageDelete umfrageID vorhanden mit Mitarbeiterumfragen", func(t *testing.T) {
		is := is.NewRelaxed(t)
		username := "anton@tobi.com"

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
			Auth: structs.AuthToken{
				Username:     username,
				Sessiontoken: server.GeneriereSessionToken(username),
			},
		}

		objectID, err := database.UmfrageInsert(data) // Neuen Eintrag erstellen
		is.NoErr(err)                                 // Kein Fehler im Normalfall

		// Fuege zwei Mitarbeiterumfragen ein
		var mitarbeiterID [2]primitive.ObjectID
		mitarbeiter := structs.InsertMitarbeiterUmfrage{
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
			IDUmfrage: objectID,
		}
		mitarbeiterID[0], err = database.MitarbeiterUmfrageInsert(mitarbeiter)
		is.NoErr(err) // Kein Fehler im Normalfall

		mitarbeiter = structs.InsertMitarbeiterUmfrage{
			Pendelweg: []structs.UmfragePendelweg{
				{IDPendelweg: 2, Strecke: 20},
			},
			TageImBuero: 3,
			Dienstreise: []structs.UmfrageDienstreise{
				{IDDienstreise: 1, Strecke: 200},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 6, Anzahl: 31},
			},
			IDUmfrage: objectID,
		}
		mitarbeiterID[1], err = database.MitarbeiterUmfrageInsert(mitarbeiter)
		is.NoErr(err) // Kein Fehler im Normalfall

		err = database.UmfrageDelete(username, objectID) // Loesche Umfrage und Mitarbeiterumfragen
		is.NoErr(err)                                    // Kein Fehler im Normalfall

		_, err = database.UmfrageFind(objectID) // Eintraege koennen nicht mehr gefunden werden
		is.Equal(err, mongo.ErrNoDocuments)     // Datenbank wirft ErrNoDocument

		_, err = database.MitarbeiterUmfrageFind(mitarbeiterID[0])
		is.Equal(err, mongo.ErrNoDocuments)
		_, err = database.MitarbeiterUmfrageFind(mitarbeiterID[1])
		is.Equal(err, mongo.ErrNoDocuments)
	})

	//Fehlerfall
	t.Run("UmfrageDelete umfrageID nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idUmfrage primitive.ObjectID
		err := idUmfrage.UnmarshalText([]byte("61b23e9835aa64762baf76a9"))
		is.NoErr(err)

		err = database.UmfrageDelete("anton@tobi.com", idUmfrage)
		is.Equal(err, mongo.ErrNoDocuments) // Eintrag konnte nicht gefunden
	})

	t.Run("UmfrageDelete mitarbeiterUmfrageID nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idUmfrage primitive.ObjectID
		err := idUmfrage.UnmarshalText([]byte("61dc0c543e48998484eefaea"))
		is.NoErr(err)

		err = database.UmfrageDelete("anton@tobi.com", idUmfrage)
		is.Equal(err, structs.ErrObjectIDNichtGefunden) // Eintrag konnte nicht gefunden
	})
}

func TestUmfrageDeleteMitarbeiterUmfrage(t *testing.T) {
	is := is.NewRelaxed(t)

	//Normalfall
	t.Run("UmfrageDeleteMitarbeiterUmfrage umfrageID vorhanden", func(t *testing.T) {
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

		err = database.UmfrageDeleteMitarbeiterUmfrage(idMitarbeiterumfrage)
		is.NoErr(err) // kein Error im Normallfall

		_, err = database.MitarbeiterUmfrageFind(idMitarbeiterumfrage) // Umfrage kann nicht mehr gefunden werden
		is.Equal(err, mongo.ErrNoDocuments)

		// pruefe, dass MitarbeiterUmfrage nicht mehr in Umfrage
		umfrageAfterDeletion, err := database.UmfrageFind(idUmfrage)
		is.NoErr(err)

		for _, mitarbeiterRef := range umfrageAfterDeletion.MitarbeiterUmfrageRef {
			is.Equal(mitarbeiterRef == idMitarbeiterumfrage, false)
		}
	})

	//Fehlerfall
	t.Run("UmfrageDeleteMitarbeiterUmfrage umfrageID nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idUmfrage primitive.ObjectID
		err := idUmfrage.UnmarshalText([]byte("aaaaaaaaaaaaaaaaaaaaaaaa"))
		is.NoErr(err)

		err = database.UmfrageDeleteMitarbeiterUmfrage(idUmfrage)
		is.Equal(err, structs.ErrObjectIDNichtGefunden)
	})
}

func TestNutzerdatenDelete(t *testing.T) {
	is := is.NewRelaxed(t)

	//Normalfall
	t.Run("NutzerdatenDelete username vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "lorem_ipsum_mustermann"

		err := database.NutzerdatenDelete(username)
		is.NoErr(err) // kein Error im Normalfall

		_, err = database.NutzerdatenFind(username) // Nutzer kann nicht mehr gefunden werden
		is.Equal(err, mongo.ErrNoDocuments)
	})

	t.Run("NutzerdatenDelete username vorhanden und Umfragen geloescht und Mitarbeiterumfragen der Umfragen geloescht", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "testuser"

		nutzerdaten, err := database.NutzerdatenFind(username)
		is.NoErr(err)

		umfragenBefore := make([]structs.Umfrage, len(nutzerdaten.UmfrageRef))
		for i, umfrageID := range nutzerdaten.UmfrageRef {
			umfragenBefore[i], err = database.UmfrageFind(umfrageID)
			is.NoErr(err)
		}

		err = database.NutzerdatenDelete(username)
		is.NoErr(err) // kein Error im Normalfall

		_, err = database.NutzerdatenFind(username) // Nutzer kann nicht mehr gefunden werden
		is.Equal(err, mongo.ErrNoDocuments)

		// Umfragen koennen nicht mehr gefunden werden
		for _, umfrage := range umfragenBefore {
			_, err = database.UmfrageFind(umfrage.ID)
			is.Equal(err, mongo.ErrNoDocuments)

			// mitarbeiterUmfragen koennen nicht mehr gefunden werden
			for _, mitarbeiterUmfrageID := range umfrage.MitarbeiterUmfrageRef {
				_, err = database.MitarbeiterUmfrageFind(mitarbeiterUmfrageID)
				is.Equal(err, mongo.ErrNoDocuments)
			}
		}
	})

	//Fehlerfall
	t.Run("NutzerdatenDelete username nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "bielefeld"
		err := database.NutzerdatenDelete(username)
		is.Equal(err, mongo.ErrNoDocuments)
	})
}
