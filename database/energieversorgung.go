package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"runtime/debug"
)

// EnergieversorgungFind liefert einen Energieversorgung struct mit idEnergieversorgung gleich dem Parameter.
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
		log.Println(err)
		debug.PrintStack()
		return structs.Energieversorgung{}, err
	}

	return data, nil
}

// EnergieversorgungAddFaktor updated ein Dokument in der Datenbank, um den CO2-Faktor {jahr, wert},
// falls das Dokument vorhanden, aber das Jahr noch nicht vorhanden ist.
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
			log.Println(structs.ErrJahrVorhanden)
			debug.PrintStack()
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
		log.Println(err)
		debug.PrintStack()
		return err
	}

	return nil
}
