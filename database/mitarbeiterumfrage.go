package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MitarbeiterUmfrageUpdate Updates a umfrage with value given in data and returns the ID of the updated Umfrage
func MitarbeiterUmfrageUpdate(data structs.UpdateMitarbeiterUmfrage) (primitive.ObjectID, error) {
	// TODO Tests
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.MitarbeiterUmfrageCol)

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
				{"pendelweg", data.Pendelweg},
				{"tageImBuero", data.TageImBuero},
				{"dienstreise", data.Dienstreise},
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

/**
Die Funktion liefert einen Array aus Umfrage structs aus der Datenbank zurueck, die mit der gegebenen Umfrage(ID) assoziiert sind.
*/
func MitarbeiterUmfrageFindForUmfrage(umfrageID primitive.ObjectID) ([]structs.MitarbeiterUmfrage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	// find umfrage for given id
	umfrage, err := UmfrageFind(umfrageID)
	if err != nil {
		return nil, err
	}

	// get ref ids from umfrage
	umfrageRefs := umfrage.MitarbeiterUmfrageRef

	collection := client.Database(dbName).Collection(structs.MitarbeiterUmfrageCol)

	cursor, err := collection.Find(
		ctx,
		bson.D{{"_id", bson.M{"$in": umfrageRefs}}},
	)
	if err != nil {
		return nil, err
	}

	var results []structs.MitarbeiterUmfrage

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
func MitarbeiterUmfrageFind(id primitive.ObjectID) (structs.MitarbeiterUmfrage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.MitarbeiterUmfrageCol)

	var data structs.MitarbeiterUmfrage
	err := collection.FindOne(
		ctx,
		bson.D{{"_id", id}},
	).Decode(&data)
	if err != nil {
		return structs.MitarbeiterUmfrage{}, err
	}

	return data, nil
}

func MitarbeiterUmfrageFindMany(ids []primitive.ObjectID) ([]structs.MitarbeiterUmfrage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.MitarbeiterUmfrageCol)

	cursor, err := collection.Find(
		ctx,
		bson.D{{"_id",
			bson.D{{"$in", ids}}}},
	)
	if err != nil {
		return nil, err
	}

	var data []structs.MitarbeiterUmfrage
	err = cursor.All(ctx, &data)
	if err != nil {
		return nil, err
	}

	if len(ids) != len(data) {
		return nil, structs.ErrDokumenteNichtGefunden
	}

	return data, nil
}

/**
Die Funktion fügt eine neue Mitarbeiterumfrage in die Datenbank ein und liefert die ObjectId mit.
*/
func MitarbeiterUmfrageInsert(data structs.InsertMitarbeiterUmfrage) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.MitarbeiterUmfrageCol)

	insertedDoc, err := collection.InsertOne(
		ctx,
		structs.MitarbeiterUmfrage{
			ID:          primitive.NewObjectID(),
			Pendelweg:   data.Pendelweg,
			TageImBuero: data.TageImBuero,
			Dienstreise: data.Dienstreise,
			ITGeraete:   data.ITGeraete,
			Revision:    1,
		},
	)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id, ok := insertedDoc.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, structs.ErrObjectIDNichtKonvertierbar
	}

	err = UmfrageAddMitarbeiterUmfrageRef(data.IDUmfrage, id)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id, nil
}
