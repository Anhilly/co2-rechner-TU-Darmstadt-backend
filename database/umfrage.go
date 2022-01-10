package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/**
Die Funktion liefert einen Umfrage struct aus der Datenbank zurueck mit ObjectID gleich dem Parameter,
falls ein Document vorhanden ist.
*/
func UmfrageFind(id primitive.ObjectID) (structs.Umfrage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.UmfrageCol)

	var data structs.Umfrage
	err := collection.FindOne(
		ctx,
		bson.D{{"_id", id}},
	).Decode(&data)

	if err != nil {
		return structs.Umfrage{}, err
	}

	return data, nil
}

/**
Die Funktion fuegt eine Umfrage in die Datenbank ein und liefert die ObjectId der Umfrage zurueck.
*/
func UmfrageInsert(data structs.InsertUmfrage) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.UmfrageCol)

	insertedDoc, err := collection.InsertOne(
		ctx,
		structs.Umfrage{
			ID:                    primitive.NewObjectID(),
			Mitarbeiteranzahl:     data.Mitarbeiteranzahl,
			Jahr:                  data.Jahr,
			Gebaeude:              data.Gebaeude,
			ITGeraete:             data.ITGeraete,
			Revision:              1,
			MitarbeiterUmfrageRef: []primitive.ObjectID{},
		})
	if err != nil {
		return primitive.NilObjectID, err
	}

	id, ok := insertedDoc.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, structs.ErrObjectIDNichtKonvertierbar
	}

	// TODO needs to be commented out if not used with user authentification to work properly
	err = NutzerdatenAddUmfrageref(data.NutzerEmail, id)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id, nil
}

/**
Die Funktion fuegt eine Referenz an eine Umfrage an.
*/
func UmfrageAddMitarbeiterUmfrageRef(idUmfrage primitive.ObjectID, referenz primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.UmfrageCol)

	var updatedDoc structs.Umfrage
	err := collection.FindOneAndUpdate(
		ctx,
		bson.D{{"_id", idUmfrage}},
		bson.D{{"$addToSet",
			bson.D{{"mitarbeiterUmfrageRef", referenz}}}},
	).Decode(&updatedDoc)
	if err != nil {
		return err
	}

	return nil
}

/**
Die Funktion loescht eine Umfrage mit der ObjectID und alle assoziierten Mitarbeiterumfragen aus der Datenbank,
falls der Eintrag existiert, liefert Fehler oder nil zurueck
*/
func UmfrageDelete(umfrageID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.UmfrageCol)

	//Loesche assoziierte Mitarbeiterumfragen
	eintrag, err := UmfrageFind(umfrageID)
	if err != nil {
		return err
	}

	for _, mitarbeiterumfrage := range eintrag.MitarbeiterUmfrageRef {
		err = UmfrageDeleteMitarbeiterUmfrage(mitarbeiterumfrage)
		if err != nil {
			return err
		}
	}

	// Loesche Eintrag aus Umfragen
	anzahl, err := collection.DeleteOne(
		ctx,
		bson.M{"_id": umfrageID})

	if anzahl.DeletedCount == 0 {
		return structs.ErrObjectIDNichtGefunden
	}
	return err
}

/**
Die Funktion loescht eine Mitarbeiterumfrage mit der UmfrageID
*/
func UmfrageDeleteMitarbeiterUmfrage(umfrageID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.MitarbeiterUmfrageCol)

	anzahl, err := collection.DeleteOne(
		ctx,
		bson.M{"_id": umfrageID})

	if anzahl.DeletedCount == 0 {
		return structs.ErrObjectIDNichtGefunden
	}
	return err
}
