package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
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
func NutzerdatenInsert(anmeldedaten structs.AuthReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.NutzerdatenCol)
	// Pruefe, ob bereits ein Eintrag mit diesem Nutzernamen existiert
	_, err := NutzerdatenFind(anmeldedaten.Username)

	if err == nil {
		// Eintrag mit diesem Nutzernamen existiert bereits
		return structs.ErrInsertExistingAccount
	}
	// Kein Eintrag vorhanden

	passwordhash, err := bcrypt.GenerateFromPassword([]byte(anmeldedaten.Passwort), bcrypt.DefaultCost)
	if err != nil {
		return err // Bcrypt hashing error
	}
	_, err = collection.InsertOne(ctx, structs.Nutzerdaten{
		Nutzername: anmeldedaten.Username,
		Passwort:   string(passwordhash),
		Rolle:      structs.IDRolleNutzer,
		Revision:   1,
		UmfrageRef: []primitive.ObjectID{},
	})
	if err != nil {
		return err // DB Error
	}

	return nil
}

// NutzerdatenDelete loescht einen Nutzer mit dem gegebenen username und alle assoziierten Umfragen aus der Datenbank.
// falls der Eintrag nicht existiert, wird ein Fehler bzw nil zurückgeliefert
func NutzerdatenDelete(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.NutzerdatenCol)

	// finde nutzerdaten
	nutzerdaten, err := NutzerdatenFind(username)
	if err != nil {
		return err
	}

	// Loesche assoziierte Umfragen
	for _, umfrageID := range nutzerdaten.UmfrageRef {

		umfrage, err := UmfrageFind(umfrageID)
		if err != nil {
			return err
		}

		// Loesche assoziierte Mitarbeiterumfragen pro Umfrage
		for _, mitarbeiterumfrage := range umfrage.MitarbeiterUmfrageRef {
			err = UmfrageDeleteMitarbeiterUmfrage(mitarbeiterumfrage)
			if err != nil {
				return err
			}
		}

		// loesche Umfrage nun selbst
		err = UmfrageDelete(username, umfrageID)
		if err != nil {
			return err
		}
	}

	// Loesche Nutzerdaten
	anzahl, err := collection.DeleteOne(
		ctx,
		bson.M{"nutzername": username})

	if err != nil {
		return err
	}

	if anzahl.DeletedCount == 0 {
		return structs.ErrUsernameNichtGefunden
	}

	return err
}
