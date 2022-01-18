package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
)

// DienstreisenFind liefert einen Dienstreisen struct mit idDienstreisen gleich dem Parameter.
func DienstreisenFind(idDienstreisen int32) (structs.Dienstreisen, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.DienstreisenCol)

	var data structs.Dienstreisen
	err := collection.FindOne(
		ctx,
		bson.D{{"idDienstreisen", idDienstreisen}},
	).Decode(&data)
	if err != nil {
		return structs.Dienstreisen{}, err
	}

	return data, nil
}
