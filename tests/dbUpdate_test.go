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

func TestUpdate(t *testing.T) {
	is := is.NewRelaxed(t)

	dir, err := database.CreateDump("TestUpdate")
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

	t.Run("TestUmfrageUpdate", TestUmfrageUpdate)
}

func TestUmfrageUpdate(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("UmfrageUpdate: Update von Umfragewerten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "anton@tobi.com"
		token := server.GeneriereSessionToken(username)

		data := structs.InsertUmfrage{
			Bezeichnung:       "TestUmfrageUpdated",
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
				Sessiontoken: token,
			},
		}

		id, err := database.UmfrageInsert(data)
		is.NoErr(err) // kein Error seitens der Datenbank

		updateData := structs.UpdateUmfrage{
			UmfrageID:         id,
			Bezeichnung:       "neuer Name",
			Mitarbeiteranzahl: 12,
			Jahr:              2077,
			Gebaeude: []structs.UmfrageGebaeude{
				{GebaeudeNr: 1102, Nutzflaeche: 100},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 6, Anzahl: 30},
				{IDITGeraete: 4, Anzahl: 12},
			},
		}

		idOfUpdatedUmfrage, err := database.UmfrageUpdate(updateData)
		is.NoErr(err) // kein Error seitens der Datenbank
		is.Equal(idOfUpdatedUmfrage, id)

		updatedUmfrage, err := database.UmfrageFind(id)
		is.NoErr(err) // kein Error seitens der Datenbank

		is.Equal(updatedUmfrage, structs.Umfrage{
			ID:                    id,
			Bezeichnung:           "neuer Name",
			Mitarbeiteranzahl:     updateData.Mitarbeiteranzahl,
			Jahr:                  updateData.Jahr,
			Gebaeude:              updateData.Gebaeude,
			ITGeraete:             updateData.ITGeraete,
			Revision:              1,
			MitarbeiterUmfrageRef: []primitive.ObjectID{},
		}) // Ueberpruefung des geaenderten Elementes

		// check that reference from user to umfrage is still correct
		user, err := database.NutzerdatenFind(data.Auth.Username)
		is.NoErr(err) // kein Error seitens der Datenbank
		idStillInUserRefs := false

		for _, b := range user.UmfrageRef {
			if b == id {
				idStillInUserRefs = true
			}
		}

		is.Equal(idStillInUserRefs, true)
	})

	// Errorfall
	t.Run("UmfrageUpdate: Umfrage ID unbekannt", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var id primitive.ObjectID
		err := id.UnmarshalText([]byte("aaaaaaaaaaaaaaaaaaaaaaaa"))

		updateData := structs.UpdateUmfrage{
			UmfrageID:         id,
			Bezeichnung:       "neuer Name",
			Mitarbeiteranzahl: 12,
			Jahr:              2077,
			Gebaeude: []structs.UmfrageGebaeude{
				{GebaeudeNr: 1102, Nutzflaeche: 100},
			},
			ITGeraete: []structs.UmfrageITGeraete{
				{IDITGeraete: 6, Anzahl: 30},
				{IDITGeraete: 4, Anzahl: 12},
			},
		}

		idOfUpdatedUmfrage, err := database.UmfrageUpdate(updateData)
		is.Equal(err, mongo.ErrNoDocuments) // ErrNoDocuments error
		is.Equal(idOfUpdatedUmfrage, primitive.NilObjectID)
	})
}
