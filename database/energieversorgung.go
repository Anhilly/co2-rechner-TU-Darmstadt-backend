package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
)

/**
Die Funktion liefert einen Energieversorgung struct mit idEnergieversorgung gleich dem Parameter.
*/
func EnergieversorgungFind(idEnergieversorgung int32) (structs.Energieversorgung, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.EnergieversorgungCol)

	var data structs.Energieversorgung
	err := collection.FindOne(
		ctx,
		bson.D{{"idEnergieversorgung", idEnergieversorgung}},
	).Decode(&data)
	if err != nil {
		return structs.Energieversorgung{}, err
	}

	return data, nil
}

/**
Funktion updated ein Dokument in der Datenbank, um den CO2-Faktor {jahr, wert}, falls Dokument vorhanden
und Jahr noch nicht vorhanden.
*/
func EnergieversorgungAddFaktor(data structs.AddCO2Faktor) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	collection := client.Database(dbName).Collection(structs.EnergieversorgungCol)

	// Ueberpruefung, ob ID in Datenbank vorhanden
	currentDoc, err := EnergieversorgungFind(data.IDEnergieversorgung)
	if err != nil {
		return err
	}

	// Ueberpruefung, ob schon Wert zu angegebenen Jahr existiert
	for _, co2Faktor := range currentDoc.CO2Faktor {
		if co2Faktor.Jahr == data.Jahr {
			return structs.ErrJahrVorhanden
		}
	}

	// Update des Eintrages
	_, err = collection.UpdateOne(
		ctx,
		bson.D{{"idEnergieversorgung", data.IDEnergieversorgung}},
		bson.D{{"$push",
			bson.D{{"CO2Faktor",
				bson.D{{"wert", data.Wert}, {"jahr", data.Jahr}}}}}},
	)
	if err != nil {
		return err
	}

	return nil
}
