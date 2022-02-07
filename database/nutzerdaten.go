package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// NutzerdatenFind liefert einen Nutzerdaten struct zurueck, der die uebergegebene E-Mail hat,
// falls ein solches Dokument in der Datenbank vorhanden ist.
func NutzerdatenFind(username string) (structs.Nutzerdaten, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.NutzerdatenCol)

	var data structs.Nutzerdaten
	err := collection.FindOne(
		ctx,
		bson.D{{"nutzername", username}},
	).Decode(&data)
	if err != nil {
		return structs.Nutzerdaten{}, err
	}

	return data, nil
}

func NutzerdatenUpdate(nutzer structs.Nutzerdaten) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.NutzerdatenCol)

	// Update des Eintrages
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": nutzer.NutzerID},
		bson.D{
			{"$set",
				bson.D{
					{"nutzername", nutzer.Nutzername},
					{"passwort", nutzer.Passwort},
					{"rolle", nutzer.Rolle},
					{"emailBestaetigt", nutzer.EmailBestaetigt},
				}},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// NutzerdatenUpdateMailBestaetigung updatet den emailBestaetigt Eintrag von angegeben Nutzer in der Datenbank
func NutzerdatenUpdateMailBestaetigung(id primitive.ObjectID, value int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.NutzerdatenCol)

	// Update des Eintrages
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"emailBestaetigt", value}}},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// NutzerdatenAddUmfrageref fuegt einem Nutzer eine ObjectID einer Umfrage hinzu, falls der Nutzer vorhanden sind.
func NutzerdatenAddUmfrageref(username string, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.NutzerdatenCol)

	var updatedDoc structs.Nutzerdaten
	err := collection.FindOneAndUpdate(
		ctx,
		bson.D{{"nutzername", username}},
		bson.D{{"$addToSet",
			bson.D{{"umfrageRef", id}}}},
	).Decode(&updatedDoc)
	if err != nil {
		return err
	}

	return nil
}

// NutzerdatenInsert fuegt einen Datenbankeintrag in Form des Nutzerdaten structs ein, dabei wird das Passwort gehashed.
func NutzerdatenInsert(anmeldedaten structs.AuthReq) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.NutzerdatenCol)
	// Pruefe, ob bereits ein Eintrag mit diesem Nutzernamen existiert
	_, err := NutzerdatenFind(anmeldedaten.Username)

	if err == nil {
		// Eintrag mit diesem Nutzernamen existiert bereits
		return primitive.NilObjectID, structs.ErrInsertExistingAccount
	}
	// Kein Eintrag vorhanden

	passwordhash, err := bcrypt.GenerateFromPassword([]byte(anmeldedaten.Passwort), bcrypt.DefaultCost)
	if err != nil {
		return primitive.NilObjectID, err // Bcrypt hashing error
	}
	result, err := collection.InsertOne(ctx, structs.Nutzerdaten{
		NutzerID:        primitive.NewObjectID(),
		Nutzername:      anmeldedaten.Username,
		Passwort:        string(passwordhash),
		Rolle:           structs.IDRolleNutzer,
		EmailBestaetigt: structs.IDEmailNichtBestaetigt,
		Revision:        1,
		UmfrageRef:      []primitive.ObjectID{},
	})
	if err != nil {
		return primitive.NilObjectID, err // DB Error
	}

	id, ok := result.InsertedID.(primitive.ObjectID)

	if !ok {
		return primitive.NilObjectID, errors.New("Casten der ObjectID ist fehlgeschlagen")
	}

	return id, nil
}
