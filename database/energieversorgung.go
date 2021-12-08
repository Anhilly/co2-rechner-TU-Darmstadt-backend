package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	energieversorgungCol = "energieversorgung"
)

/**
Die Funktion liefert einen Energieversorgung struct mit idEnergieversorgung gleich dem Parameter.
*/
func EnergieversorgungFind(idEnergieversorgung int32) (structs.Energieversorgung, error) {
	var data structs.Energieversorgung
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(energieversorgungCol)

	cursor, err := collection.Find(ctx, bson.D{{"idEnergieversorgung", idEnergieversorgung}}) //nolint:govet
	if err != nil {
		return structs.Energieversorgung{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return structs.Energieversorgung{}, err
	}

	return data, nil
}
