package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

/**
Die Funktion liefert einen Nutzerdaten struct mit email gleich dem Parameter.
*/
func NutzerdatenFind(emailNutzer string) (structs.Nutzerdaten, error) {
	var data structs.Nutzerdaten
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.NutzerdatenCol)

	cursor, err := collection.Find(ctx, bson.D{{"email", emailNutzer}}) //nolint:govet
	if err != nil {
		//Problem mit Datenbankanbindung
		return structs.Nutzerdaten{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		//Kein Dateneintrag gefunden
		return structs.Nutzerdaten{}, err
	}

	return data, nil
}

/**
Fügt einen Datenbankeintrag in Form des Nutzerdaten structs ein, dabei wird das Passwort gehashed
*/
func NutzerdatenInsert(anmeldedaten structs.AuthReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.NutzerdatenCol)
	//Prüfe ob bereits ein Eintrag mit dieser Email existiert
	_, err := NutzerdatenFind(anmeldedaten.Email)
	if err != nil {
		//Kein Eintrag vorhanden
		passwordhash, err := bcrypt.GenerateFromPassword([]byte(anmeldedaten.Passwort), bcrypt.DefaultCost)
		if err != nil {
			return err //Bcrypt hashing error
		}
		_, err = collection.InsertOne(ctx, structs.Nutzerdaten{
			Email:    anmeldedaten.Email,
			Passwort: string(passwordhash),
			Revision: 1,
		})
		if err != nil {
			return err //DB Error
		}
		return nil
	}
	//Eintrag mit dieser Email existiert bereits
	return structs.ErrInsertExistingAccount
}
