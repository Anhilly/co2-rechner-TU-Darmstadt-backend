package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/config"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"runtime/debug"
)

// PendelwegFind liefert einen Pendelweg struct mit idPendelweg gleich dem Parameter.
func PendelwegFind(idPendelweg int32) (structs.Pendelweg, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(config.DBName).Collection(structs.PendelwegCol)

	var data structs.Pendelweg
	err := collection.FindOne(
		ctx,
		bson.D{{"idPendelweg", idPendelweg}},
	).Decode(&data)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return structs.Pendelweg{}, err
	}

	return data, nil
}
