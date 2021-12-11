package database

import (
	"context"
	"errors"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	waermezaehlerCol = "waermezaehler"
)

var (
	ErrZaehlerVorhanden    = errors.New("Es ist schon ein Zaehler mit dem PK vorhanden!")
	ErrFehlendeGebaeuderef = errors.New("Neuer Zaehler hat keine Referenzen auf Gebaeude!")
)

/**
Die Funktion liefert einen Zaehler struct für den Waermezaehler mit pkEnergie gleich dem Parameter.
*/
func WaermezaehlerFind(pkEnergie int32) (structs.Zaehler, error) {
	var data structs.Zaehler
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(waermezaehlerCol)

	cursor, err := collection.Find(ctx, bson.D{{"pkEnergie", pkEnergie}}) //nolint:govet
	if err != nil {
		return structs.Zaehler{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return structs.Zaehler{}, err
	}

	data.Zaehlertyp = "Waerme"

	return data, nil
}

/**
Funktion updated ein Dokument in der Datenbank, um den Zaehlerwert {jahr, wert}, falls Dokument vorhanden
und Jahr noch nicht vorhanden.
*/
func WaermezaehlerAddZaehlerdaten(data structs.AddZaehlerdaten) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(waermezaehlerCol)

	// Ueberpruefung, ob PK in Datenbank vorhanden
	currentDoc, err := WaermezaehlerFind(data.PKEnergie)
	if err != nil {
		return err
	}

	// Ueberpruefung, ob schon Wert zu angegebenen Jahr existiert
	for _, zaehlerdatum := range currentDoc.Zaehlerdaten {
		if int32(zaehlerdatum.Zeitstempel.Year()) == data.Jahr {
			return ErrJahrVorhanden
		}
	}

	// Update des Eintrages
	location, _ := time.LoadLocation("Etc/GMT")
	zeitstemple := time.Date(int(data.Jahr), time.January, 01, 0, 0, 0, 0, location).UTC()

	_, err = collection.UpdateOne(
		ctx,
		bson.D{{"pkEnergie", data.PKEnergie}},
		bson.D{{"$push",
			bson.D{{"zaehlerdaten",
				bson.D{{"wert", data.Wert}, {"zeitstempel", zeitstemple}}}}}},
	)
	if err != nil {
		return err
	}

	return nil
}

func WaermezaehlerInsert(data structs.InsertZaehler) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(waermezaehlerCol)

	if len(data.GebaeudeRef) == 0 {
		return ErrFehlendeGebaeuderef
	}

	_, err := WaermezaehlerFind(data.PKEnergie)
	if err == nil { // kein Error = Nr schon vorhanden
		return ErrZaehlerVorhanden
	}

	_, err = collection.InsertOne(
		ctx,
		structs.Zaehler{
			PKEnergie:    data.PKEnergie,
			Bezeichnung:  data.Bezeichnung,
			Einheit:      data.Einheit,
			Zaehlerdaten: []structs.Zaehlerwerte{},
			Spezialfall:  1,
			Revision:     1,
			GebaeudeRef:  data.GebaeudeRef,
		},
	)
	if err != nil {
		return err
	}

	for _, referenz := range data.GebaeudeRef {
		err := GebaeudeAddZaehlerref(referenz, data.PKEnergie, data.IDEnergieversorgung)
		if err != nil {
			return err
		}
	}

	return nil
}
