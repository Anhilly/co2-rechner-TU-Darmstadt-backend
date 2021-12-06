package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	dienstreisenCol = "dienstreisen"
)

/**
Die Funktion liefert einen Dienstreisen struct mit idDienstreisen gleich dem Parameter.
*/
func DienstreisenFind(idDienstreisen int32) (structs.Dienstreisen, error) {
	var data structs.Dienstreisen
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(dbName).Collection(dienstreisenCol)

	cursor, err := collection.Find(ctx, bson.D{{"idDienstreisen", idDienstreisen}})
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
