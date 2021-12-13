package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
)

/**
Die Funktion liefert einen Pendelweg struct mit idPendelweg gleich dem Parameter.
*/
func PendelwegFind(idPendelweg int32) (structs.Pendelweg, error) {
	var data structs.Pendelweg
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.PendelwegCol)

	cursor, err := collection.Find(ctx, bson.D{{"idPendelweg", idPendelweg}}) //nolint:govet
	if err != nil {
		return structs.Pendelweg{}, err
	}
	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return structs.Pendelweg{}, err
	}

	return data, nil
}
