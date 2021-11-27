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
Die Methode liefert einen Splice an ITGeraete Dokumenten mit Kategorie gleich den Parameter.
*/
func ITGeraeteFind(kategorie string) ([]ITGeraete, error) {
	var data []ITGeraete

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	collection := client.Database(dbName).Collection(itGeraeteCol)

	cursor, err := collection.Find(ctx, bson.D{{"kategorie", kategorie}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var document ITGeraete

		err := cursor.Decode(&document)
		if err != nil {
			return nil, err
		}
		data = append(data, document)
	}

	return data, nil
}
