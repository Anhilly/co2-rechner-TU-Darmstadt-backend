package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	kaeltezaehlerCol = "kaeltezaehler"
)

/**
Die Funktion liefert einen Zaehler struct für den Kaeltezaehler mit pkEnergie gleich dem Parameter.
*/
func KaeltezaehlerFind(pkEnergie int32) (structs.Zaehler, error) {
	var data structs.Zaehler
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(kaeltezaehlerCol)

	cursor, err := collection.Find(ctx, bson.D{{"pkEnergie", pkEnergie}}) //nolint:govet
	if err != nil {
		return structs.Zaehler{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return structs.Zaehler{}, err
	}

	data.Zaehlertyp = "Kaelte"

	return data, nil
}
