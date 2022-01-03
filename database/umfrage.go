package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/**
Funktion gibt alle Umfragen in der Datenbank zurueck.
*/
func AlleUmfragen() ([]structs.Umfrage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.UmfrageCol)

	cursor, err := collection.Find(
		ctx,
		bson.D{},
	)
	if err != nil {
		return nil, err
	}

	var results []structs.Umfrage

	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

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

// UmfrageUpdate Updates a umfrage with value given in data and returns the ID of the updated Umfrage
func UmfrageUpdate(data structs.UpdateUmfrageReq) (primitive.ObjectID, error) {
	// TODO Tests
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.UmfrageCol)

	var updatedDoc structs.Umfrage
	var umfrageID, err = primitive.ObjectIDFromHex(data.UmfrageID)
	if err != nil {
		return primitive.NilObjectID, err
	}

	err = collection.FindOneAndUpdate(
		ctx,
		bson.D{{"_id", umfrageID}},
		bson.D{{"$set",
			bson.D{
				{"mitarbeiteranzahl", data.Mitarbeiteranzahl},
				{"jahr", data.Jahr},
				{"gebaeude", data.Gebaeude},
				{"itGeraete", data.ITGeraete},
				// TODO also update "revision"-field?
			},
		}},
	).Decode(&updatedDoc)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return updatedDoc.ID, nil
}

// UmfrageInsert Die Funktion fuegt eine Umfrage in die Datenbank ein und liefert die ObjectId der Umfrage zurueck.
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
	err = NutzerdatenAddUmfrageref(data.Hauptverantwortlicher.Username, id)
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
