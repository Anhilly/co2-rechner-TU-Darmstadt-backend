package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	kaeltezaehlerCol = "gebaeude"
)

/**
Die Funktion liefert einen Kaeltezaehler struct mit pkEnergie gleich dem Parameter.
*/
func KaeltezaehlerFind(pkEnergie int32) (Kaeltezaehler, error) {
	var data Kaeltezaehler
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	collection := client.Database(dbName).Collection(kaeltezaehlerCol)

	cursor, err := collection.Find(ctx, bson.D{{"pkEnergie", pkEnergie}})
	if err != nil {
		return Kaeltezaehler{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return Kaeltezaehler{}, err
	}

	return data, nil
}
