package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	waermezaehlerCol = "waermezaehler"
)


/**
Die Funktion liefert einen Waermezaehler struct mit pkEnergie gleich dem Parameter.
*/
func WaermezaehlerFind(pkEnergie int32) (Waermezaehler, error) {
	var data Waermezaehler
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	collection := client.Database(dbName).Collection(waermezaehlerCol)

	cursor, err := collection.Find(ctx, bson.D{{"pkEnergie", pkEnergie}})
	if err != nil {
		return Waermezaehler{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return Waermezaehler{}, err
	}

	return data, nil
}
