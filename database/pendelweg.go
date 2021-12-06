package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	pendelwegCol = "pendelweg"
)

/**
Die Funktion liefert einen Pendelweg struct mit idPendelweg gleich dem Parameter.
*/
func PendelwegFind(idPendelweg int32) (structs.Pendelweg, error) {
	var data structs.Pendelweg
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(dbName).Collection(pendelwegCol)

	cursor, err := collection.Find(ctx, bson.D{{"idPendelweg", idPendelweg}})
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
