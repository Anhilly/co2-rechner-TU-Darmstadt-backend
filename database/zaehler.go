package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	kaeltezaehlerCol = "kaeltezaehler"
	stromzaehlerCol  = "stromzaehler"
	waermezaehlerCol = "waermezaehler"
)

/**
Die Funktion liefert einen Zaehler struct fuer den Zaehler mit pkEnergie und uebergebenen Energieform.
*/
func ZaehlerFind(pkEnergie, idEnergieversorgung int32) (structs.Zaehler, error) {
	var data structs.Zaehler
	var collectionname string
	var zaehlertyp string

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	switch idEnergieversorgung { //TODO: Ersetzte Zahlen mit Konstanten
	case 1: // Waerme
		collectionname = waermezaehlerCol
		zaehlertyp = "Waerme"
	case 2: // Strom
		collectionname = stromzaehlerCol
		zaehlertyp = "Strom"
	case 3: // Kaelte
		collectionname = kaeltezaehlerCol
		zaehlertyp = "Kaelte"
	default:
		return structs.Zaehler{}, ErrIDEnergieversorgungNichtVorhanden
	}

	collection := client.Database(dbName).Collection(collectionname)

	cursor, err := collection.Find(ctx, bson.D{{"pkEnergie", pkEnergie}}) //nolint:govet
	if err != nil {
		return structs.Zaehler{}, err
	}

	cursor.Next(ctx)
	err = cursor.Decode(&data)
	if err != nil {
		return structs.Zaehler{}, err
	}

	data.Zaehlertyp = zaehlertyp

	return data, nil
}
