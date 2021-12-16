
package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/**
Die Funktion liefert einen Nutzerdaten struct zurueck, der die uebergegebene E-Mail hat, falls ein solches Dokument in der
Datenbank vorhanden ist.
*/
func NutzerdatenFind(email string) (structs.Nutzerdaten, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.NutzerdatenCol)

	var data structs.Nutzerdaten
	err := collection.FindOne(
		ctx,
		bson.D{{"email", email}},
	).Decode(&data)
	if err != nil {
		return structs.Nutzerdaten{}, err
	}

	return data, nil
}

/**
Die Funktion fuegt einem Nutzer eine ObjectID einer Umfrage hinzu, falls der Nutzer vorhanden sind.
*/
func NutzerdatenAddUmfrageref(email string, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.NutzerdatenCol)

	var updatedDoc structs.Nutzerdaten
	err := collection.FindOneAndUpdate(
		ctx,
		bson.D{{"email", email}},
		bson.D{{"$addToSet",
			bson.D{{"umfrageRef", id}}}},
	).Decode(&updatedDoc)
	if err != nil {
		return err
	}

	return nil
}

/**
Fuegt einen Datenbankeintrag in Form des Nutzerdaten structs ein, dabei wird das Passwort gehashed
*/
func NutzerdatenInsert(anmeldedaten structs.AuthReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

    collection := client.Database(dbName).Collection(structs.NutzerdatenCol)
	// Pruefe ob bereits ein Eintrag mit dieser Email existiert
	_, err := NutzerdatenFind(anmeldedaten.Username)

	if err == nil {
		// Eintrag mit dieser Email existiert bereits
		return structs.ErrInsertExistingAccount
	}
	// Kein Eintrag vorhanden

	passwordhash, err := bcrypt.GenerateFromPassword([]byte(anmeldedaten.Passwort), bcrypt.DefaultCost)
		return err // Bcrypt hashing error
	if err != nil {
	}
	_, err = collection.InsertOne(ctx, structs.Nutzerdaten{

		Email:    anmeldedaten.Username,
		Passwort: string(passwordhash),
		Revision: 1,
	})
	if err != nil {
		return err // DB Error
	}

	return nil
}

