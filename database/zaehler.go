package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"time"
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

/**
Funktion updated einen Zaehler in der Datenbank um den Zaehlerwert {jahr, wert}, falls Zaehler vorhanden
und Jahr noch nicht vorhanden.
*/
func ZaehlerAddZaehlerdaten(data structs.AddZaehlerdaten) error {
	var collectionname string

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	switch data.IDEnergieversorgung { //TODO: Ersetzte Zahlen mit Konstanten
	case 1: // Waerme
		collectionname = waermezaehlerCol
	case 2: // Strom
		collectionname = stromzaehlerCol
	case 3: // Kaelte
		collectionname = kaeltezaehlerCol
	default:
		return ErrIDEnergieversorgungNichtVorhanden
	}

	collection := client.Database(dbName).Collection(collectionname)

	// Ueberpruefung, ob PK in Datenbank vorhanden
	currentDoc, err := ZaehlerFind(data.PKEnergie, data.IDEnergieversorgung)
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
