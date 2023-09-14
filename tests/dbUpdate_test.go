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

func TestUpdate(t *testing.T) {
	is := is.NewRelaxed(t)

	err := database.ConnectDatabase("dev")
	is.NoErr(err)

	dir, err := database.CreateDump("TestUpdate")
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

	t.Run("TestUmfrageUpdate", TestUmfrageUpdate)
	t.Run("TestUmfrageShareUpdate", TestUmfrageShareUpdate)
}

func TestUmfrageUpdate(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("UmfrageUpdate: Update von Umfragewerten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "anton"

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
		}

		id, err := database.UmfrageInsert(data, username)
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
		user, err := database.NutzerdatenFind(username)
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

func TestUmfrageShareUpdate(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("LinkShareUpdate: Update von LinkSharewerten", func(t *testing.T) {
		is := is.NewRelaxed(t)

		username := "anton"

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
		}

		oid, err := database.UmfrageInsert(data, username)
		is.NoErr(err)

		refOid, err := database.UmfrageUpdateLinkShare(0, oid)
		is.Equal(refOid, oid)
		is.NoErr(err)

		umfrage, err := database.UmfrageFind(refOid)
		is.NoErr(err)

		is.Equal(umfrage.AuswertungFreigegeben, int32(0))

		//Pruefe ob alle anderen Felder gleich geblieben sind
		is.Equal(umfrage.Bezeichnung, data.Bezeichnung)
		is.Equal(umfrage.Mitarbeiteranzahl, data.Mitarbeiteranzahl)
		is.Equal(umfrage.Jahr, data.Jahr)
		is.Equal(umfrage.Gebaeude, data.Gebaeude)
		is.Equal(umfrage.ITGeraete, data.ITGeraete)
	})

	//Errorfall
	t.Run("LinkShareUpdate: ungueltige UmfrageID", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var id primitive.ObjectID
		err := id.UnmarshalText([]byte("aaaaaaaaaaaaaaaaaaaaaaaa"))
		is.NoErr(err)

		oid, err := database.UmfrageUpdateLinkShare(1, id)
		is.Equal(err, mongo.ErrNoDocuments)
		is.Equal(oid, primitive.NilObjectID)
	})
}
