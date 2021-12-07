package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	itGeraeteCol = "itGeraete"
)

/**
Die Funktion liefert einen ITGeraete struct mit idITGeraete gleich dem Parameter.
*/
func ITGeraeteFind(idITGeraete int32) (structs.ITGeraete, error) {
	var data structs.ITGeraete
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(itGeraeteCol)

	cursor, err := collection.Find(ctx, bson.D{{"idITGeraete", idITGeraete}}) //nolint:govet
	if err != nil {
		return structs.ITGeraete{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return structs.ITGeraete{}, err
	}

	return data, nil
}
