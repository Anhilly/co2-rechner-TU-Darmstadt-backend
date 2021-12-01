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
Die Funktion liefert einen Zaehler struct für den Waermezaehler mit pkEnergie gleich dem Parameter.
*/
func WaermezaehlerFind(pkEnergie int32) (Zaehler, error) {
	var data Zaehler
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	collection := client.Database(dbName).Collection(waermezaehlerCol)

	cursor, err := collection.Find(ctx, bson.D{{"pkEnergie", pkEnergie}})
	if err != nil {
		return Zaehler{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return Zaehler{}, err
	}

	data.Zaehlertyp = "Waerme"

	return data, nil
}
