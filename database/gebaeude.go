package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	gebaeudeCol = "gebaeude"
)

/**
Die Funktion liefert einen Gebaeude struct mit nr gleich dem Parameter.
*/
func GebaeudeFind(nr int32) (structs.Gebaeude, error) {
	var data structs.Gebaeude
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(gebaeudeCol)

	cursor, err := collection.Find(ctx, bson.D{{"nr", nr}}) //nolint:govet
	if err != nil {
		return structs.Gebaeude{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return structs.Gebaeude{}, err
	}

	return data, nil
}
