package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"runtime/debug"
)

// AlleUmfragen gibt alle Umfragen in der Datenbank zurueck.
func AlleUmfragen() ([]structs.Umfrage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.UmfrageCol)

	cursor, err := collection.Find(
		ctx,
		bson.D{},
	)
	if err != nil {
		log.Println(err)
		debug.PrintStack()
		return nil, err
	}

	var results []structs.Umfrage

	err = cursor.All(ctx, &results)
	if err != nil {
		log.Println(err)
		debug.PrintStack()
		return nil, err
	}

	return results, nil
}

// AlleUmfragenForUser gibt alle Umfragen in der Datenbank zurueck, die mit gegebenem User assoziiert sind.
func AlleUmfragenForUser(username string) ([]structs.Umfrage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	// Finde Nutzerdaten fuer den gegeben Usernamen
	nutzerdaten, err := NutzerdatenFind(username)
	if err != nil {
		return nil, err
	}

	// Hole mit Nutzer assoziierten UmfrageIDs
	umfrageRefs := nutzerdaten.UmfrageRef

	// liefere leere Liste zurueck, falls keine assoziierten Umfragen gefunden wurden
	if umfrageRefs == nil || len(umfrageRefs) == 0 {
		return []structs.Umfrage{}, nil
	}

	collection := client.Database(dbName).Collection(structs.UmfrageCol)

	cursor, err := collection.Find(
		ctx,
		bson.D{{"_id", bson.M{"$in": umfrageRefs}}},
	)
	if err != nil {
		log.Println(err)
		debug.PrintStack()
		return nil, err
	}

	var results []structs.Umfrage

	err = cursor.All(ctx, &results)
	if err != nil {
		log.Println(err)
		debug.PrintStack()
		return nil, err
	}

	return results, nil
}

// UmfrageFind liefert einen Umfrage struct aus der Datenbank zurueck mit ObjectID gleich dem Parameter,
// falls ein Document vorhanden ist.
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
		log.Println(err)
		debug.PrintStack()
		return structs.Umfrage{}, err
	}

	return data, nil
}

// UmfrageUpdate Updates a umfrage with value given in data and returns the ID of the updated Umfrage
func UmfrageUpdate(data structs.UpdateUmfrage) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.UmfrageCol)

	var updatedDoc structs.Umfrage

	err := collection.FindOneAndUpdate(
		ctx,
		bson.D{{"_id", data.UmfrageID}},
		bson.D{{"$set",
			bson.D{
				{"mitarbeiteranzahl", data.Mitarbeiteranzahl},
				{"bezeichnung", data.Bezeichnung},
				{"jahr", data.Jahr},
				{"gebaeude", data.Gebaeude},
				{"itGeraete", data.ITGeraete},
			},
		}},
	).Decode(&updatedDoc)

	if err != nil {
		log.Println(err)
		debug.PrintStack()
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
			Bezeichnung:           data.Bezeichnung,
			Mitarbeiteranzahl:     data.Mitarbeiteranzahl,
			Jahr:                  data.Jahr,
			Gebaeude:              data.Gebaeude,
			ITGeraete:             data.ITGeraete,
			Revision:              1,
			MitarbeiterUmfrageRef: []primitive.ObjectID{},
			AuswertungFreigegeben: 0,
		})
	if err != nil {
		log.Println(err)
		debug.PrintStack()
		return primitive.NilObjectID, err
	}

	id, ok := insertedDoc.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println(structs.ErrObjectIDNichtKonvertierbar)
		debug.PrintStack()
		return primitive.NilObjectID, structs.ErrObjectIDNichtKonvertierbar
	}

	err = NutzerdatenAddUmfrageref(data.Auth.Username, id)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id, nil
}

// UmfrageAddMitarbeiterUmfrageRef haengt eine Mitarbeiterumfrage Referenz an eine Umfrage an.
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
		log.Println(err)
		debug.PrintStack()
		return err
	}

	return nil
}

// UmfrageDelete loescht eine Umfrage mit der ObjectID und alle assoziierten Mitarbeiterumfragen aus der Datenbank,
// falls der Eintrag existiert, liefert Fehler oder nil zurueck
func UmfrageDelete(username string, umfrageID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.UmfrageCol)

	// Loesche assoziierte Mitarbeiterumfragen
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
		bson.M{"_id": umfrageID},
	)
	if err != nil {
		log.Println(err)
		debug.PrintStack()
		return err
	}

	if anzahl.DeletedCount == 0 {
		log.Println(structs.ErrObjectIDNichtGefunden)
		debug.PrintStack()
		return structs.ErrObjectIDNichtGefunden
	}

	// Loesche Umfrage aus RefListe des Nutzers
	collection = client.Database(dbName).Collection(structs.NutzerdatenCol)

	var updatedDoc structs.Nutzerdaten
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"nutzername": username},
		bson.D{{"$pull",
			bson.D{{"umfrageRef", umfrageID}}}},
	).Decode(&updatedDoc)
	if err != nil {
		log.Println(err)
		debug.PrintStack()
		return err
	}

	return nil
}

// UmfrageUpdateLinkShare setzt den auswertungFreigegeben Wert der Umfrage mit der gegebenen umfrageID
// auf den uebergebenen setValue Wert. Dieser ist entweder 0 oder 1.
// Der Wert steuert ob die Auswertung der Umfrage geteilt werden darf.
func UmfrageUpdateLinkShare(setValue int32, umfrageID primitive.ObjectID) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.UmfrageCol)

	var updatedDoc structs.Umfrage

	err := collection.FindOneAndUpdate(
		ctx,
		bson.D{{"_id", umfrageID}},
		bson.D{{"$set",
			bson.D{
				{"auswertungFreigegeben", setValue},
			},
		}},
	).Decode(&updatedDoc)
	if err != nil {
		log.Println(err)
		debug.PrintStack()
		return primitive.NilObjectID, err
	}

	return updatedDoc.ID, nil
}
