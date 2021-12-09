package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
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
		return structs.Nutzerdaten{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return structs.Nutzerdaten{}, err
	}

	return data, nil
}
