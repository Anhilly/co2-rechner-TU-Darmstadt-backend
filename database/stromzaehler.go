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
Die Funktion liefert einen Stromzaehler struct mit pkEnergie gleich dem Parameter.
*/
func StromzaehlerFind(pkEnergie int32) (Stromzaehler, error) {
	var data Stromzaehler
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	collection := client.Database(dbName).Collection(stromzaehlerCol)

	cursor, err := collection.Find(ctx, bson.D{{"pkEnergie", pkEnergie}})
	if err != nil {
		return Stromzaehler{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return Stromzaehler{}, err
	}

	return data, nil
}
