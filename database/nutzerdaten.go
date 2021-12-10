package database

import (
	"context"
	"errors"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var (
	//Nutzer will Account mit bestehender Email registrieren
	ErrInsertExistingAccount = errors.New("Account with this Email already exists")
)

const (
	nutzerdatenCol = "nutzerdaten"
)

/**
Die Funktion liefert einen Nutzerdaten struct mit email gleich dem Parameter.
*/
func NutzerdatenFind(emailNutzer string) (structs.Nutzerdaten, error) {
	var data structs.Nutzerdaten
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(nutzerdatenCol)

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

func NutzerdatenInsert(anmeldedaten structs.AnmeldungReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(nutzerdatenCol)

	//Check if entry already exists
	_, err := NutzerdatenFind(anmeldedaten.Email)
	if err != nil {
		//No entry found
		passwordhash, _ := bcrypt.GenerateFromPassword([]byte(anmeldedaten.Passwort), bcrypt.DefaultCost)
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
	//Entry with given email already exists
	return ErrInsertExistingAccount
}
