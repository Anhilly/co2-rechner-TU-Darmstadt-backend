package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"runtime/debug"
)

// ITGeraeteFind liefert einen ITGeraete struct mit idITGeraete gleich dem Parameter.
func ITGeraeteFind(idITGeraete int32) (structs.ITGeraete, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.ITGeraeteCol)

	var data structs.ITGeraete
	err := collection.FindOne(
		ctx, bson.D{{"idITGeraete", idITGeraete}},
	).Decode(&data)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return structs.ITGeraete{}, err
	}

	return data, nil
}

// ITGeraeteFindAll liefert einen Slice aller IT-Geraete structs.
func ITGeraeteFindAll() ([]structs.ITGeraete, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.ITGeraeteCol)

	var data []structs.ITGeraete
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return nil, err
	}
	err = cursor.All(ctx, &data)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return nil, err
	}

	return data, nil
}
