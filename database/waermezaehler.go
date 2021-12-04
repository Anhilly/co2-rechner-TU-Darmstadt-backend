package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	waermezaehlerCol = "waermezaehler"
)

/**
Die Funktion liefert einen Zaehler struct für den Waermezaehler mit pkEnergie gleich dem Parameter.
*/
func WaermezaehlerFind(pkEnergie int32) (structs.Zaehler, error) {
	var data structs.Zaehler
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(dbName).Collection(waermezaehlerCol)

	cursor, err := collection.Find(ctx, bson.D{{"pkEnergie", pkEnergie}})
	if err != nil {
		return structs.Zaehler{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return structs.Zaehler{}, err
	}

	data.Zaehlertyp = "Waerme"

	return data, nil
}
