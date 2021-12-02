package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	gebaeudeCol = "gebaeude"
)

/**
Die Funktion liefert einen Gebaeude struct mit nr gleich dem Parameter.
*/
func GebaeudeFind(nr int32) (Gebaeude, error) {
	var data Gebaeude
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(dbName).Collection(gebaeudeCol)

	cursor, err := collection.Find(ctx, bson.D{{"nr", nr}})
	if err != nil {
		return Gebaeude{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return Gebaeude{}, err
	}

	return data, nil
}
