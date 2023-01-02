package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"runtime/debug"
	"time"
)

// GebaeudeFind liefert einen Gebaeude struct mit nr gleich dem Parameter.
func GebaeudeFind(nr int32) (structs.Gebaeude, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.GebaeudeCol)

	var data structs.Gebaeude
	err := collection.FindOne(
		ctx,
		bson.D{{"nr", nr}},
	).Decode(&data)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return structs.Gebaeude{}, err
	}

	return data, nil
}

// GebaeudeInsert fuegt ein Gebaeude in die Datenbank ein, falls die Nr noch nicht vorhanden ist.
func GebaeudeInsert(data structs.InsertGebaeude) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.GebaeudeCol)

	_, err := GebaeudeFind(data.Nr)
	if err == nil { // kein Error = Nr schon vorhanden
		return structs.ErrGebaeudeVorhanden
	}

	aktuellesJahr := time.Now().Year()
	var versorger []structs.Versoger

	for i := structs.ErstesJahr; i < aktuellesJahr; i++ {
		versorger = append(versorger, structs.Versoger{
			Jahr:      int32(i),
			IDVertrag: structs.IDVertragTU,
		})
	}

	_, err = collection.InsertOne(
		ctx,
		structs.Gebaeude{
			Nr:              data.Nr,
			Bezeichnung:     data.Bezeichnung,
			Flaeche:         data.Flaeche,
			Einheit:         "m^2",
			Spezialfall:     1,
			Revision:        1,
			KaelteRef:       []int32{},
			WaermeRef:       []int32{},
			StromRef:        []int32{},
			Kaelteversorger: versorger,
			Waermeversorger: versorger,
			Stromversorger:  versorger,
		},
	)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	return nil
}

// GebaeudeAddZaehlerref fuegt einem Gebaeude eine Zaehlereferenz hinzu, falls diese noch nicht vorhanden ist.
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
		log.Println(structs.ErrIDEnergieversorgungNichtVorhanden)
		log.Println(string(debug.Stack()))
		return structs.ErrIDEnergieversorgungNichtVorhanden
	}

	var updatedDoc structs.Gebaeude
	err := collection.FindOneAndUpdate(
		ctx,
		bson.D{{"nr", nr}},
		bson.D{{"$addToSet", // $addToSet verhindert, dass eine Referenz doppelt im Array steht (sollte nicht vorkommen)
			bson.D{{referenzname, ref}}}},
	).Decode(&updatedDoc)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	return nil
}

// GebaeudeAlleNr gibt alle Nummern von Gebaeuden in der Datenbank zurueck.
func GebaeudeAlleNr() ([]int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.GebaeudeCol)

	cursor, err := collection.Find(
		ctx,
		bson.D{},
		options.Find().SetProjection(bson.M{"_id": 0, "nr": 1}),
	)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return nil, err
	}

	var results []struct {
		Nr int32 `bson:"nr"`
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return nil, err
	}

	var gebaeudenummern []int32
	for _, elem := range results {
		gebaeudenummern = append(gebaeudenummern, elem.Nr)
	}

	return gebaeudenummern, nil
}
