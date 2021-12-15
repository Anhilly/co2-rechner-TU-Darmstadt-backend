package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
)

/**
Die Funktion liefert einen Gebaeude struct mit nr gleich dem Parameter.
*/
func GebaeudeFind(nr int32) (structs.Gebaeude, error) {
	var data structs.Gebaeude
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.GebaeudeCol)

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

/**
Die Funktion fuegt ein Gebaeude in die Datenbank ein, falls die Nr noch nicht vorhanden ist.
*/
func GebaeudeInsert(data structs.InsertGebaeude) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.GebaeudeCol)

	_, err := GebaeudeFind(data.Nr)
	if err == nil { // kein Error = Nr schon vorhanden
		return structs.ErrGebaeudeVorhanden
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

/**
Die Funktion fuegt einem Gebaeude eine Zaehlereferenz hinzu, falls diese noch nicht vorhanden ist.
*/
func GebaeudeAddZaehlerref(nr, ref, idEnergieversorgung int32) error {
	var referenzname string

	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.GebaeudeCol)

	switch idEnergieversorgung {
	case structs.IDEnergieversorgungWaerme: // Waerme
		referenzname = "waermeRef"
	case structs.IDEnergieversorgungStrom: // Strom
		referenzname = "stromRef"
	case structs.IDEnergieversorgungKaelte: // Kaelte
		referenzname = "kaelteRef"
	default:
		return structs.ErrIDEnergieversorgungNichtVorhanden
	}

	var updatedDoc bson.M
	err := collection.FindOneAndUpdate(
		ctx,
		bson.D{{"nr", nr}},
		bson.D{{"$addToSet", // $addToSet verhindert, dass eine Referenz doppelt im Array steht (sollte nicht vorkommen)
			bson.D{{referenzname, ref}}}},
	).Decode(&updatedDoc)

	if err != nil {
		return err
	}

	return nil
}
