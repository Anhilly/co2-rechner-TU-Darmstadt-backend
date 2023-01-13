package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"runtime/debug"
	"time"
)

// ZaehlerFind liefert einen Zaehler struct fuer den Zaehler mit pkEnergie und uebergebenen Energieform.
func ZaehlerFind(pkEnergie, idEnergieversorgung int32) (structs.Zaehler, error) {
	var collectionname string
	var zaehlertyp string

	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	switch idEnergieversorgung {
	case structs.IDEnergieversorgungWaerme: // Waerme
		collectionname = structs.WaermezaehlerCol
		zaehlertyp = structs.ZaehlertypWaerme
	case structs.IDEnergieversorgungStrom: // Strom
		collectionname = structs.StromzaehlerCol
		zaehlertyp = structs.ZaehlertypStrom
	case structs.IDEnergieversorgungKaelte: // Kaelte
		collectionname = structs.KaeltezaehlerCol
		zaehlertyp = structs.ZaehlertypKaelte
	default:
		log.Println(structs.ErrIDEnergieversorgungNichtVorhanden)
		log.Println(string(debug.Stack()))
		return structs.Zaehler{}, structs.ErrIDEnergieversorgungNichtVorhanden
	}

	collection := client.Database(dbName).Collection(collectionname)

	var data structs.Zaehler
	err := collection.FindOne(
		ctx,
		bson.D{{"pkEnergie", pkEnergie}},
	).Decode(&data)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return structs.Zaehler{}, err
	}

	data.Zaehlertyp = zaehlertyp

	return data, nil
}

// ZaehlerAddZaehlerdaten updated einen Zaehler in der Datenbank um den Zaehlerwert {jahr, wert},
// falls Zaehler vorhanden und Jahr noch nicht vorhanden.
func ZaehlerAddZaehlerdaten(data structs.AddZaehlerdaten) error {
	var collectionname string

	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	switch data.IDEnergieversorgung {
	case structs.IDEnergieversorgungWaerme: // Waerme
		collectionname = structs.WaermezaehlerCol
	case structs.IDEnergieversorgungStrom: // Strom
		collectionname = structs.StromzaehlerCol
	case structs.IDEnergieversorgungKaelte: // Kaelte
		collectionname = structs.KaeltezaehlerCol
	default:
		log.Println(structs.ErrIDEnergieversorgungNichtVorhanden)
		log.Println(string(debug.Stack()))
		return structs.ErrIDEnergieversorgungNichtVorhanden
	}

	collection := client.Database(dbName).Collection(collectionname)

	// Ueberpruefung, ob PK in Datenbank vorhanden
	currentDoc, err := ZaehlerFind(data.PKEnergie, data.IDEnergieversorgung)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	// Ueberpruefung, ob schon Wert zu angegebenen Jahr existiert
	for _, zaehlerdatum := range currentDoc.Zaehlerdaten {
		if int32(zaehlerdatum.Zeitstempel.Year()) == data.Jahr {
			log.Println(structs.ErrJahrVorhanden)
			log.Println(string(debug.Stack()))
			return structs.ErrJahrVorhanden
		}
	}

	// Update des Eintrages
	location, _ := time.LoadLocation("Etc/GMT")
	zeitstempel := time.Date(int(data.Jahr), time.January, 01, 0, 0, 0, 0, location).UTC()

	_, err = collection.UpdateOne(
		ctx,
		bson.D{{"pkEnergie", data.PKEnergie}},
		bson.D{{"$push",
			bson.D{{"zaehlerdaten",
				bson.D{{"wert", data.Wert}, {"zeitstempel", zeitstempel}}}}}},
	)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	return nil
}

// ZaehlerAddStandardZaehlerdaten updated alle Zaehler in der Datenbank um den Zaehlerwert {jahr, 0.0},
// falls Jahr noch nicht vorhanden.
func ZaehlerAddStandardZaehlerdaten(data structs.AddStandardZaehlerdaten) error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	// Update aller Zaehler ohne Wert für Zeitstempel
	location, _ := time.LoadLocation("Etc/GMT")
	zeitstempel := time.Date(int(data.Jahr), time.January, 01, 0, 0, 0, 0, location).UTC()

	for _, tmp := range []string{structs.WaermezaehlerCol, structs.StromzaehlerCol, structs.KaeltezaehlerCol} {
		collection := client.Database(dbName).Collection(tmp)

		_, err := collection.UpdateMany(
			ctx,
			bson.D{{"zaehlerdaten",
				bson.D{{"$not",
					bson.D{{"$elemMatch",
						bson.D{{"zeitstempel", zeitstempel}}}}}}}},
			bson.D{{"$push",
				bson.D{{"zaehlerdaten",
					bson.D{{"wert", 0.0}, {"zeitstempel", zeitstempel}}}}}},
		)
		if err != nil {
			log.Println(err)
			log.Println(string(debug.Stack()))
			return err
		}
	}

	return nil
}

// ZaehlerInsert fuegt einen Zaehler in die Datenbank ein, falls PK noch nicht vergeben.
// Außerdem werden die referenzierten Gebaeude um eine Referenz auf diesen Zaehler erweitert.
// Sollte die Funktion durch einen Fehler beendet werden, kann es zu inkonsistenten Daten in der Datenbank fuehren!
func ZaehlerInsert(data structs.InsertZaehler) error {
	var collectionname string

	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	switch data.IDEnergieversorgung {
	case structs.IDEnergieversorgungWaerme: // Waerme
		collectionname = structs.WaermezaehlerCol
	case structs.IDEnergieversorgungStrom: // Strom
		collectionname = structs.StromzaehlerCol
	case structs.IDEnergieversorgungKaelte: // Kaelte
		collectionname = structs.KaeltezaehlerCol
	default:
		log.Println(structs.ErrIDEnergieversorgungNichtVorhanden)
		log.Println(string(debug.Stack()))
		return structs.ErrIDEnergieversorgungNichtVorhanden
	}
	collection := client.Database(dbName).Collection(collectionname)

	if len(data.GebaeudeRef) == 0 {
		log.Println(structs.ErrFehlendeGebaeuderef)
		log.Println(string(debug.Stack()))
		return structs.ErrFehlendeGebaeuderef
	}

	_, err := ZaehlerFind(data.PKEnergie, data.IDEnergieversorgung)
	if err == nil { // kein Error = Nr schon vorhanden
		log.Println(err)
		log.Println(string(debug.Stack()))
		return structs.ErrZaehlerVorhanden
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
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	for _, referenz := range data.GebaeudeRef {
		err := GebaeudeAddZaehlerref(referenz, data.PKEnergie, data.IDEnergieversorgung)
		if err != nil {
			log.Println(err)
			log.Println(string(debug.Stack()))
			return err
		}
	}

	return nil
}
