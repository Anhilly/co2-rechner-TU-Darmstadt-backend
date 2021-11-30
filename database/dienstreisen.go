package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	dienstreisenCol = "dienstreisen"
)

/**
Die Funktion liefert einen Dienstreisen struct mit idDienstreisen gleich dem Parameter.
*/
func DienstreisenFind(idDienstreisen int32) (Dienstreisen, error) {
	var data Dienstreisen
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	collection := client.Database(dbName).Collection(dienstreisenCol)

	cursor, err := collection.Find(ctx, bson.D{{"idDienstreisen", idDienstreisen}})
	if err != nil {
		return Dienstreisen{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return Dienstreisen{}, err
	}

	return data, nil
}
