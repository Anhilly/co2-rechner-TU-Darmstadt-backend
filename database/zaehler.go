package database

import (
	"context"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"runtime/debug"
	"time"
)

// ZaehlerFindOID liefert einen Zaehler struct fuer den Zaehler mit ObjectID und uebergebenen Energieform.
func ZaehlerFindOID(oid primitive.ObjectID, idEnergieversorgung int32) (structs.Zaehler, error) {
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
		bson.D{{"_id", oid}},
	).Decode(&data)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return structs.Zaehler{}, err
	}

	data.Zaehlertyp = zaehlertyp

	return data, nil
}

// ZaehlerFindDPName liefert einen Zaehler struct fuer den Zaehler mit DPName und uebergebenen Energieform.
func ZaehlerFindDPName(dpName string, idEnergieversorgung int32) (structs.Zaehler, error) {
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
		bson.D{{"dpName", dpName}},
	).Decode(&data)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return structs.Zaehler{}, err
	}

	data.Zaehlertyp = zaehlertyp

	return data, nil
}

// ZaehlerAlleZaehlerUndDaten
func ZaehlerAlleZaehlerUndDaten() ([]structs.ZaehlerUndZaehlerdaten, error) {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	var results []structs.ZaehlerUndZaehlerdaten
	var temp []structs.ZaehlerUndZaehlerdaten

	collectionWaermezaehler := client.Database(dbName).Collection(structs.WaermezaehlerCol)
	cursorWaermezaehler, err := collectionWaermezaehler.Find(
		ctx,
		bson.D{},
		options.Find().SetProjection(bson.M{"_id": 1, "zaehlerdaten": 1}),
	)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return nil, err
	}
	err = cursorWaermezaehler.All(ctx, &temp)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return nil, err
	}
	results = append(results, temp...)

	collectionKaeltezaehler := client.Database(dbName).Collection(structs.KaeltezaehlerCol)
	cursorKaeltezaehler, err := collectionKaeltezaehler.Find(
		ctx,
		bson.D{},
		options.Find().SetProjection(bson.M{"_id": 1, "zaehlerdaten": 1}),
	)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return nil, err
	}
	err = cursorKaeltezaehler.All(ctx, &temp)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return nil, err
	}
	results = append(results, temp...)

	collectionStromzaehler := client.Database(dbName).Collection(structs.StromzaehlerCol)
	cursorStromzaehler, err := collectionStromzaehler.Find(
		ctx,
		bson.D{},
		options.Find().SetProjection(bson.M{"_id": 1, "zaehlerdaten": 1}),
	)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return nil, err
	}
	err = cursorStromzaehler.All(ctx, &temp)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return nil, err
	}
	results = append(results, temp...)

	return results, nil
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

	// Ueberpruefung, ob DPName in Datenbank vorhanden
	currentDoc, err := ZaehlerFindDPName(data.DPName, data.IDEnergieversorgung)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	// Ueberpruefung, ob schon Wert zu angegebenen Jahr existiert
	wert_ersetzten := false

	for _, zaehlerdatum := range currentDoc.Zaehlerdaten {
		if int32(zaehlerdatum.Zeitstempel.Year()) == data.Jahr && int32(zaehlerdatum.Zeitstempel.Month()) == data.Monat {
			if zaehlerdatum.Wert == 0.0 {
				wert_ersetzten = true
				break
			}
			log.Println(structs.ErrJahrUndMonatVorhanden)
			log.Println(string(debug.Stack()))
			return structs.ErrJahrUndMonatVorhanden
		}
	}

	// Update des Eintrages
	location, _ := time.LoadLocation("Etc/GMT")
	zeitstempel := time.Date(int(data.Jahr), time.Month(data.Monat), 01, 0, 0, 0, 0, location).UTC()

	if wert_ersetzten {
		_, err = collection.UpdateOne(
			ctx,
			bson.D{{"dpName", data.DPName}},
			bson.D{{"$set",
				bson.D{{"zaehlerdaten.$[element]",
					bson.D{{"wert", data.Wert}, {"zeitstempel", zeitstempel}}}}}},
			options.Update().SetArrayFilters(options.ArrayFilters{
				Filters: []interface{}{bson.M{"element.zeitstempel": zeitstempel}},
			}),
		)

	} else {
		_, err = collection.UpdateOne(
			ctx,
			bson.D{{"dpName", data.DPName}},
			bson.D{{"$push",
				bson.D{{"zaehlerdaten",
					bson.D{{"wert", data.Wert}, {"zeitstempel", zeitstempel}}}}}},
		)
	}
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
	zeitstempel := time.Date(int(data.Jahr), time.Month(data.Monat), 01, 0, 0, 0, 0, location).UTC()

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
func ZaehlerInsert(data structs.InsertZaehler) (primitive.ObjectID, error) {
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
		return primitive.NilObjectID, structs.ErrIDEnergieversorgungNichtVorhanden
	}
	collection := client.Database(dbName).Collection(collectionname)

	if len(data.GebaeudeRef) == 0 {
		log.Println(structs.ErrFehlendeGebaeuderef)
		log.Println(string(debug.Stack()))
		return primitive.NilObjectID, structs.ErrFehlendeGebaeuderef
	}

	_, err := ZaehlerFindDPName(data.DPName, data.IDEnergieversorgung)
	if err == nil { // kein Error = DPName schon vorhanden
		log.Println(err)
		log.Println(string(debug.Stack()))
		return primitive.NilObjectID, structs.ErrZaehlerVorhanden
	}

	location, _ := time.LoadLocation("Etc/GMT")
	aktuellesJahr := int32(time.Now().Year())
	var zaehlerdaten []structs.Zaehlerwerte
	for i := structs.ErstesJahr; i <= aktuellesJahr; i++ {
		for j := 1; j <= 12; j++ {
			zaehlerdaten = append(zaehlerdaten, structs.Zaehlerwerte{
				Wert:        0.0,
				Zeitstempel: time.Date(int(i), time.Month(j), 01, 0, 0, 0, 0, location).UTC(),
			})
		}
	}

	// konvertiert die Gebaeudenummern in ObjectIDs aus der Datenbank
	var gebaeudeRefOID []primitive.ObjectID
	for _, referenz := range data.GebaeudeRef {
		gebaeude, err := GebaeudeFind(referenz)
		if err != nil {
			log.Println(err)
			log.Println(string(debug.Stack()))
			return primitive.NilObjectID, structs.ErrGebaeudeNichtVorhanden
		}
		gebaeudeRefOID = append(gebaeudeRefOID, gebaeude.GebaeudeID)
	}

	zaehlerID := primitive.NewObjectID()

	_, err = collection.InsertOne(
		ctx,
		structs.Zaehler{
			ZaehlerID:    zaehlerID,
			DPName:       data.DPName,
			Bezeichnung:  data.Bezeichnung,
			Einheit:      data.Einheit,
			Zaehlerdaten: zaehlerdaten,
			Spezialfall:  1,
			Revision:     2,
			GebaeudeRef:  gebaeudeRefOID,
		},
	)
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return primitive.NilObjectID, err
	}

	for _, referenz := range data.GebaeudeRef {
		err := GebaeudeAddZaehlerref(referenz, zaehlerID, data.IDEnergieversorgung)
		if err != nil {
			log.Println(err)
			log.Println(string(debug.Stack()))
			return primitive.NilObjectID, err
		}
	}

	return zaehlerID, nil
}
