package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
)

/**
Die Funktion liefert einen ITGeraete struct mit idITGeraete gleich dem Parameter.
*/
func ITGeraeteFind(idITGeraete int32) (structs.ITGeraete, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.ITGeraeteCol)

	var data structs.ITGeraete
	err := collection.FindOne(
		ctx, bson.D{{"idITGeraete", idITGeraete}},
	).Decode(&data)
	if err != nil {
		return structs.ITGeraete{}, err
	}

	return data, nil
}
