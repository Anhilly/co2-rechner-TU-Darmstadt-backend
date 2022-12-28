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
		log.Println(string(debug.Stack()))
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

	jahrVorhanden := false

	// Ueberpruefung, ob schon Wert zu angegebenen Jahr und Vertrag existiert
	for _, co2Faktor := range currentDoc.CO2Faktor {
		if co2Faktor.Jahr == data.Jahr {
			jahrVorhanden = true

			for _, vertrag := range co2Faktor.Vertraege {
				if vertrag.IDVertrag == data.IDVertrag {
					log.Println(structs.ErrJahrVorhanden)
					log.Println(string(debug.Stack()))
					return structs.ErrJahrVorhanden
				}
			}
		}
	}

	// Update des Eintrages
	if jahrVorhanden {
		_, err = collection.UpdateOne(
			ctx,
			bson.D{{"$and", []bson.D{bson.D{{"idEnergieversorgung", data.IDEnergieversorgung}}, bson.D{{"CO2Faktor.jahr", data.Jahr}}}}},
			bson.D{{"$push",
				bson.D{{"CO2Faktor.$.vertraege", bson.D{{"idVertrag", data.IDVertrag}, {"wert", data.Wert}}}}}},
		)
	} else {
		_, err = collection.UpdateOne(
			ctx,
			bson.D{{"idEnergieversorgung", data.IDEnergieversorgung}},
			bson.D{{"$push",
				bson.D{{"CO2Faktor",
					bson.D{{"jahr", data.Jahr},
						{"vertraege", []bson.D{bson.D{{"idVertrag", data.IDVertrag}, {"wert", data.Wert}}}}}}}}},
		)
	}
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	return nil
}
