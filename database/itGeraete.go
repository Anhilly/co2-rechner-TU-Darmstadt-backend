package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	itGeraeteCol = "itGeraete"
)

/**
Die Funktion liefert einen ITGeraete struct mit idITGeraete gleich dem Parameter.
*/
func ITGeraeteFind(idITGeraete int32) (ITGeraete, error) {
	var data ITGeraete

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	collection := client.Database(dbName).Collection(itGeraeteCol)

	cursor, err := collection.Find(ctx, bson.D{{"idITGeraete", idITGeraete}})
	if err != nil {
		return ITGeraete{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return ITGeraete{}, err
	}

	return data, nil
}
