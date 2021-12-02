package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	stromzaehlerCol = "stromzaehler"
)

/**
Die Funktion liefert einen Zaehler struct für den Stromzaehler mit pkEnergie gleich dem Parameter.
*/
func StromzaehlerFind(pkEnergie int32) (Zaehler, error) {
	var data Zaehler
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(dbName).Collection(stromzaehlerCol)

	cursor, err := collection.Find(ctx, bson.D{{"pkEnergie", pkEnergie}})
	if err != nil {
		return Zaehler{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return Zaehler{}, err
	}

	data.Zaehlertyp = "Strom"

	return data, nil
}
