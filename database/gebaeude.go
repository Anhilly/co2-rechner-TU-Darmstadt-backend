package database

import (
	"context"
	"errors"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	gebaeudeCol = "gebaeude"
)

var (
	ErrGebaeudeVorhanden = errors.New("Ein Gebaeude mit der angegeben Nummer existiert schon in der Datenbank")
)

/**
Die Funktion liefert einen Gebaeude struct mit nr gleich dem Parameter.
*/
func GebaeudeFind(nr int32) (structs.Gebaeude, error) {
	var data structs.Gebaeude
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(gebaeudeCol)

	cursor, err := collection.Find(ctx, bson.D{{"nr", nr}}) //nolint:govet
	if err != nil {
		return structs.Gebaeude{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return structs.Gebaeude{}, err
	}

	return data, nil
}

func GebaeudeInsert(data structs.InsertGebaeude) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(gebaeudeCol)

	_, err := GebaeudeFind(data.Nr)
	if err == nil { // kein Error = Nr schon vorhanden
		return ErrGebaeudeVorhanden
	}

	_, err = collection.InsertOne(
		ctx,
		structs.Gebaeude{
			Nr:          data.Nr,
			Bezeichnung: data.Bezeichnung,
			Flaeche:     data.Flaeche,
			Einheit:     "m^2",
			Spezialfall: 1,
			Revision:    1,
			KaelteRef:   []int32{},
			WaermeRef:   []int32{},
			StromRef:    []int32{},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung int32) error {
	var referenzname string

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(gebaeudeCol)

	switch idEnergieversorgung { //TODO: Ersetzte Zahlen mit Konstanten
	case 1: // Waerme
		referenzname = "waermeRef"
	case 2: // Strom
		referenzname = "stromRef"
	case 3: // Kaelte
		referenzname = "kaelteRef"
	}

	var updatedDoc bson.M
	err := collection.FindOneAndUpdate(
		ctx,
		bson.D{{"nr", nr}},
		// $addToSet verhindert, dass eine Referenz doppelt im Array steht (sollte nicht vorkommen)
		bson.D{{"$addToSet", bson.D{{referenzname, ref}}}},
	).Decode(&updatedDoc)

	if err != nil {
		return err
	}

	return nil
}
