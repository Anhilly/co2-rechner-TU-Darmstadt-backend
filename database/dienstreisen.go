package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
)

/**
Die Funktion liefert einen Dienstreisen struct mit idDienstreisen gleich dem Parameter.
*/
func DienstreisenFind(idDienstreisen int32) (structs.Dienstreisen, error) {
	var data structs.Dienstreisen
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.DienstreisenCol)

	cursor, err := collection.Find(ctx, bson.D{{"idDienstreisen", idDienstreisen}}) //nolint:govet
	if err != nil {
		return structs.Dienstreisen{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return structs.Dienstreisen{}, err
	}

	return data, nil
}
