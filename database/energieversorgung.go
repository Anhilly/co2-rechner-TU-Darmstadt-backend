package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	energieversorgungCol = "energieversorgung"
)

/**
Die Funktion liefert einen Energieversorgung struct mit idEnergieversorgung gleich dem Parameter.
*/
func EnergieversorgungFind(idEnergieversorgung int32) (Energieversorgung, error) {
	var data Energieversorgung
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(dbName).Collection(energieversorgungCol)

	cursor, err := collection.Find(ctx, bson.D{{"idEnergieversorgung", idEnergieversorgung}})
	if err != nil {
		return Energieversorgung{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return Energieversorgung{}, err
	}

	return data, nil
}
