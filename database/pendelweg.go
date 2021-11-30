package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)
const (
	pendelwegCol = "pendelweg"
)

/**
Die Funktion liefert einen Pendelweg struct mit idPendelweg gleich dem Parameter.
*/
func PendelwegFind(idPendelweg int32) (Pendelweg, error) {
	var data Pendelweg
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	collection := client.Database(dbName).Collection(pendelwegCol)

	cursor, err := collection.Find(ctx, bson.D{{"idPendelweg", idPendelweg}})
	if err != nil {
		return Pendelweg{}, err
	}
	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return Pendelweg{}, err
	}

	return data, nil
}
