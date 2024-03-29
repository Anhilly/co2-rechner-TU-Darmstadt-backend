package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"runtime/debug"
)

// MitarbeiterUmfrageFind liefert einen Mitarbeiterumfrage struct aus der Datenbank zurueck mit ObjectID gleich dem Parameter,
// falls ein Document vorhanden ist.
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
		log.Println(err)
		log.Println(string(debug.Stack()))
		return structs.MitarbeiterUmfrage{}, err
	}

	return data, nil
}

//// MitarbeiterUmfrageUpdate updated eine Mitarbeiterumfrage mit den in data uebergebenen Werten und
//// gibt die ID der aktualisierten Umfrage zurueck.
//func MitarbeiterUmfrageUpdate(data structs.UpdateMitarbeiterUmfrage) (primitive.ObjectID, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
//	defer cancel()
//
//	collection := client.Database(dbName).Collection(structs.MitarbeiterUmfrageCol)
//
//	var updatedDoc structs.Umfrage
//
//	err := collection.FindOneAndUpdate(
//		ctx,
//		bson.D{{"_id", data.UmfrageID}},
//		bson.D{{"$set",
//			bson.D{
//				{"pendelweg", data.Pendelweg},
//				{"tageImBuero", data.TageImBuero},
//				{"dienstreise", data.Dienstreise},
//				{"itGeraete", data.ITGeraete},
//			},
//		}},
//	).Decode(&updatedDoc)
//
//	if err != nil {
//		log.Println(err)
//		log.Println(string(debug.Stack()))
//		return primitive.NilObjectID, err
//	}
//
//	return updatedDoc.ID, nil
//}

// MitarbeiterUmfrageFindForUmfrage liefert einen Array aus Mitarbeiterumfrage structs aus der Datenbank zurueck,
// die mit der gegebenen Umfrage(ID) assoziiert sind.
func MitarbeiterUmfrageFindForUmfrage(umfrageID primitive.ObjectID) ([]structs.MitarbeiterUmfrage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	// suche Umfrage mit der gegebenen umfrageID
	umfrage, err := UmfrageFind(umfrageID)
	if err != nil {
		return nil, err
	}

	// hole die IDs der assoziierten Mitarbeiterumfragen
	umfrageRefs := umfrage.MitarbeiterumfrageRef

	collection := client.Database(dbName).Collection(structs.MitarbeiterUmfrageCol)

	cursor, err := collection.Find(
		ctx,
		bson.D{{"_id", bson.M{"$in": umfrageRefs}}},
	)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return nil, err
	}

	var results []structs.MitarbeiterUmfrage

	err = cursor.All(ctx, &results)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return nil, err
	}

	return results, nil
}

// MitarbeiterumfrageFindMany liefert ein Array aus allen Mitarbeiterumfragen zurueck, deren ID in ids liegt.
// Wenn nicht alle IDs in ids in der DB gefunden wurden, wird der Fehler structs.ErrDokumenteNichtGefunden zurueckgegeben.
func MitarbeiterumfrageFindMany(ids []primitive.ObjectID) ([]structs.MitarbeiterUmfrage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.MitarbeiterUmfrageCol)

	cursor, err := collection.Find(
		ctx,
		bson.D{{"_id",
			bson.D{{"$in", ids}}}},
	)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return []structs.MitarbeiterUmfrage{}, err
	}

	var data []structs.MitarbeiterUmfrage
	err = cursor.All(ctx, &data)
	if err != nil {
		log.Println(err)
		log.Println(debug.Stack())
		return []structs.MitarbeiterUmfrage{}, err
	}

	if len(ids) != len(data) {
		log.Println(structs.ErrDokumenteNichtGefunden)
		log.Println(string(debug.Stack()))
		return nil, structs.ErrDokumenteNichtGefunden
	}

	return data, nil
}

// MitarbeiterUmfrageInsert fügt eine neue Mitarbeiterumfrage in die Datenbank ein und liefert die ObjectId zurueck.
func MitarbeiterUmfrageInsert(data structs.InsertMitarbeiterUmfrage) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.MitarbeiterUmfrageCol)

	// only insert if umfrage is not complete yet
	umfrageID := data.IDUmfrage

	// hole Umfrage aus der Datenbank
	umfrage, err := UmfrageFind(umfrageID)
	if err != nil {
		return primitive.NilObjectID, err
	}

	mitarbeiterumfragen, err := MitarbeiterumfrageFindMany(umfrage.MitarbeiterumfrageRef)
	if err != nil {
		return primitive.NilObjectID, err
	}

	mitarbeiterMax := umfrage.Mitarbeiteranzahl
	umfragenFilled := int32(len(mitarbeiterumfragen))

	if umfragenFilled >= mitarbeiterMax {
		log.Println(structs.ErrUmfrageVollstaendig)
		log.Println(string(debug.Stack()))
		return primitive.NilObjectID, structs.ErrUmfrageVollstaendig
	}

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
		log.Println(err)
		log.Println(string(debug.Stack()))
		return primitive.NilObjectID, err
	}

	id, ok := insertedDoc.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println(structs.ErrObjectIDNichtKonvertierbar)
		log.Println(string(debug.Stack()))
		return primitive.NilObjectID, structs.ErrObjectIDNichtKonvertierbar
	}

	err = UmfrageAddMitarbeiterUmfrageRef(umfrageID, id)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id, nil
}

// UmfrageDeleteMitarbeiterUmfrage loescht eine Mitarbeiterumfrage mit der gegebenen UmfrageID
func UmfrageDeleteMitarbeiterUmfrage(mitarbeiterUmfrageID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.MitarbeiterUmfrageCol)

	anzahl, err := collection.DeleteOne(
		ctx,
		bson.M{"_id": mitarbeiterUmfrageID})
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	if anzahl.DeletedCount == 0 {
		log.Println(structs.ErrObjectIDNichtGefunden)
		log.Println(string(debug.Stack()))
		return structs.ErrObjectIDNichtGefunden
	}

	// remove MitarbeiterUmfrage from Refs in Umfrage
	umfrageCollection := client.Database(dbName).Collection(structs.UmfrageCol)

	var updatedDocument structs.Umfrage

	err = umfrageCollection.FindOneAndUpdate(
		ctx,
		bson.M{"mitarbeiterUmfrageRef": mitarbeiterUmfrageID},
		bson.D{{"$pull",
			bson.D{{"mitarbeiterUmfrageRef", mitarbeiterUmfrageID}}}},
	).Decode(&updatedDocument)

	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	return err
}
