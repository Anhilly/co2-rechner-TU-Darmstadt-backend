package importer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

func ImportEnergieversorgung() {
	in, err := os.Open("importer\\Energieversorgung.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	reader := csv.NewReader(in)
	reader.FieldsPerRecord = -1
	reader.Comma = ';'

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var energieversorgung database.Energieversorgung
	var energieversorgungArray []database.Energieversorgung

	for _, record := range rawCSVdata {
		energieversorgung.Kategorie = record[0]

		temp, _ := strconv.ParseFloat(record[1], 64)
		temp2, _ := strconv.ParseInt(record[2], 10, 32)
		energieversorgung.CO2Faktor = []database.CO2Energie{{Wert: temp, Jahr: int32(temp2)}}

		energieversorgung.Einheit = record[3]
		energieversorgung.Revision = 1

		fmt.Println(energieversorgung)

		b, _ := json.Marshal(energieversorgung)
		fmt.Println(string(b))

		energieversorgungArray = append(energieversorgungArray, energieversorgung)
	}

	file, _ := json.MarshalIndent(energieversorgungArray, "", " ")
	_ = ioutil.WriteFile("importer\\energieversorgung.json", file, 0644)

	fmt.Println(energieversorgungArray)
}

func ImportGebaeude() {
	in, err := os.Open("importer\\Gebaeudedaten_erweitert.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	reader := csv.NewReader(in)
	reader.FieldsPerRecord = -1
	reader.Comma = ';'

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var gebaeudeArray []database.Gebaeude

	var data []byte
	for _, b1 := range []byte("[\n") {
		data = append(data, b1)
	}

	for _, record := range rawCSVdata {
		var gebaeude database.Gebaeude
		if record[0] == "" {
			continue
		}

		temp, err := strconv.ParseInt(record[0], 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		gebaeude.Nr = int32(temp)

		gebaeude.Bezeichnung = record[1]

		hnf, _ := strconv.ParseFloat(record[2], 64)
		nnf, _ := strconv.ParseFloat(record[3], 64)
		ngf, _ := strconv.ParseFloat(record[4], 64)
		ff, _ := strconv.ParseFloat(record[5], 64)
		vf, _ := strconv.ParseFloat(record[6], 64)
		freif, _ := strconv.ParseFloat(record[7], 64)
		gesamtf, _ := strconv.ParseFloat(record[8], 64)
		gebaeude.Flaeche = database.GebaeudeFlaeche{HNF: hnf, NNF: nnf, NGF: ngf, FF: ff, VF: vf, FreiF: freif, GesamtF: gesamtf}

		gebaeude.Einheit = "m^2"
		gebaeude.Revision = 1

		specialcase, _ := strconv.ParseFloat(record[9], 64)
		gebaeude.Spezialfall = int32(specialcase)

		gebaeude.StromRef = []int32{}
		for i := 0; i < 10; i++ {
			if record[i+10] == "" {
				continue
			}

			temp, _ := strconv.ParseInt(record[i+10], 10, 32)
			gebaeude.StromRef = append(gebaeude.StromRef, int32(temp))
		}

		gebaeude.KaelteRef = []int32{}
		for i := 0; i < 10; i++ {
			if record[i+20] == "" {
				continue
			}

			temp, _ := strconv.ParseInt(record[i+20], 10, 32)
			gebaeude.KaelteRef = append(gebaeude.KaelteRef, int32(temp))
		}

		gebaeude.WaermeRef = []int32{}
		for i := 0; i < 10; i++ {
			if record[i+30] == "" {
				continue
			}

			temp, _ := strconv.ParseInt(record[i+30], 10, 32)
			gebaeude.WaermeRef = append(gebaeude.WaermeRef, int32(temp))
		}

		fmt.Println(gebaeude)

		b, _ := bson.MarshalExtJSONIndent(gebaeude, false, false, "", " ")
		fmt.Println(string(b))

		for _, b1 := range b {
			data = append(data, b1)
		}
		for _, b1 := range []byte(",\n") {
			data = append(data, b1)
		}

		gebaeudeArray = append(gebaeudeArray, gebaeude)
	}
	data = data[0 : len(data)-2]
	for _, b1 := range []byte("\n]") {
		data = append(data, b1)
	}

	fmt.Println(string(data))
	_ = ioutil.WriteFile("importer\\Gebaeudedaten.json", data, 0644)
}

func ImportStromzaehler() {
	in, err := os.Open("importer\\Strom_erweitert.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	reader := csv.NewReader(in)
	reader.FieldsPerRecord = -1
	reader.Comma = ';'

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var stromArray []database.Stromzaehler

	var data []byte
	for _, b1 := range []byte("[\n") {
		data = append(data, b1)
	}

	for _, record := range rawCSVdata {
		var strom database.Stromzaehler
		if record[0] == "" {
			continue
		}

		pke, err := strconv.ParseInt(record[1], 10, 32)
		if err != nil {
			fmt.Println(err)
		}

		strom.Bezeichnung = record[0]
		strom.PKEnergie = int32(pke)

		d21, _ := strconv.ParseFloat(record[2], 64)
		d20, _ := strconv.ParseFloat(record[3], 64)
		d19, _ := strconv.ParseFloat(record[4], 64)
		d18, _ := strconv.ParseFloat(record[5], 64)
		strom.Einheit = record[6]
		strom.Revision = 1

		location, err := time.LoadLocation("Etc/GMT")
		if err != nil {
			log.Fatal(err)
		}

		strom.Zaehlerdaten = []database.Zaehlerwerte{
			database.Zaehlerwerte{Wert: d21, Zeitstempel: time.Date(2021, time.January, 01, 0, 0, 0, 0, location)},
			database.Zaehlerwerte{Wert: d20, Zeitstempel: time.Date(2020, time.January, 01, 0, 0, 0, 0, location)},
			database.Zaehlerwerte{Wert: d19, Zeitstempel: time.Date(2019, time.January, 01, 0, 0, 0, 0, location)},
			database.Zaehlerwerte{Wert: d18, Zeitstempel: time.Date(2018, time.January, 01, 0, 0, 0, 0, location)}}

		specialcase, _ := strconv.ParseFloat(record[7], 64)
		strom.Spezialfall = int32(specialcase)

		strom.GebaeudeRef = []int32{}

		for i := 0; i < 10; i++ {
			if record[i+8] == "" {
				continue
			}

			temp, _ := strconv.ParseInt(record[i+8], 10, 32)
			strom.GebaeudeRef = append(strom.GebaeudeRef, int32(temp))
		}

		fmt.Println(strom)

		b, _ := bson.MarshalExtJSONIndent(strom, false, false, "", " ")
		fmt.Println(string(b))

		for _, b1 := range b {
			data = append(data, b1)
		}
		for _, b1 := range []byte(",\n") {
			data = append(data, b1)
		}

		stromArray = append(stromArray, strom)
	}

	data = data[0 : len(data)-2]
	for _, b1 := range []byte("\n]") {
		data = append(data, b1)
	}

	fmt.Println(string(data))
	_ = ioutil.WriteFile("importer\\Stromdaten.json", data, 0644)
}

func ImportWaermedaten() {
	in, err := os.Open("importer\\Waerme_erweitert.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	reader := csv.NewReader(in)
	reader.FieldsPerRecord = -1
	reader.Comma = ';'

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var waermeArray []database.Waermezaehler

	var data []byte
	for _, b1 := range []byte("[\n") {
		data = append(data, b1)
	}

	for _, record := range rawCSVdata {
		var waerme database.Waermezaehler
		if record[0] == "" {
			continue
		}

		waerme.ExtSystemID = record[0]
		waerme.Bezeichnung = record[1]
		waerme.Einheit = record[2]

		waerme.Revision = 1

		pke, err := strconv.ParseInt(record[6], 10, 32)
		if err != nil {
			fmt.Println(err)
		}

		waerme.PKEnergie = int32(pke)

		d20, _ := strconv.ParseFloat(record[3], 64)
		d19, _ := strconv.ParseFloat(record[4], 64)
		d18, _ := strconv.ParseFloat(record[5], 64)

		location, err := time.LoadLocation("Etc/GMT")
		if err != nil {
			log.Fatal(err)
		}

		waerme.Zaehlerdaten = []database.Zaehlerwerte{
			database.Zaehlerwerte{Wert: d20, Zeitstempel: time.Date(2020, time.January, 01, 0, 0, 0, 0, location)},
			database.Zaehlerwerte{Wert: d19, Zeitstempel: time.Date(2019, time.January, 01, 0, 0, 0, 0, location)},
			database.Zaehlerwerte{Wert: d18, Zeitstempel: time.Date(2018, time.January, 01, 0, 0, 0, 0, location)}}

		specialcase, _ := strconv.ParseFloat(record[7], 64)
		waerme.Spezialfall = int32(specialcase)

		waerme.GebaeudeRef = []int32{}

		for i := 0; i < 10; i++ {
			if record[i+8] == "" {
				continue
			}

			temp, _ := strconv.ParseInt(record[i+8], 10, 32)
			waerme.GebaeudeRef = append(waerme.GebaeudeRef, int32(temp))
		}

		fmt.Println(waerme)

		b, _ := bson.MarshalExtJSONIndent(waerme, false, false, "", " ")
		fmt.Println(string(b))

		for _, b1 := range b {
			data = append(data, b1)
		}
		for _, b1 := range []byte(",\n") {
			data = append(data, b1)
		}

		waermeArray = append(waermeArray, waerme)
	}

	data = data[0 : len(data)-2]
	for _, b1 := range []byte("\n]") {
		data = append(data, b1)
	}

	fmt.Println(string(data))
	_ = ioutil.WriteFile("importer\\Waermedaten.json", data, 0644)
}

func ImportKaeltedaten() {
	in, err := os.Open("importer\\Kaelte_erweitert.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	reader := csv.NewReader(in)
	reader.FieldsPerRecord = -1
	reader.Comma = ';'

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	//var kaelte database.Kaeltezaehler
	var kaelteArray []database.Kaeltezaehler

	var data []byte

	for _, b1 := range []byte("[\n") {
		data = append(data, b1)
	}

	for _, record := range rawCSVdata {
		var kaelte database.Kaeltezaehler

		if record[0] == "" {
			continue
		}

		kaelte.Bezeichnung = record[0]
		kaelte.ExtSystemID = record[3]
		kaelte.Einheit = record[2]

		kaelte.Revision = 1

		pke, err := strconv.ParseInt(record[1], 10, 32)
		if err != nil {
			fmt.Println(err)
		}

		kaelte.PKEnergie = int32(pke)

		d21, _ := strconv.ParseFloat(record[4], 64)
		d20, _ := strconv.ParseFloat(record[5], 64)
		d19, _ := strconv.ParseFloat(record[6], 64)
		d18, _ := strconv.ParseFloat(record[7], 64)

		location, err := time.LoadLocation("Etc/GMT")
		if err != nil {
			log.Fatal(err)
		}

		kaelte.Zaehlerdaten = []database.Zaehlerwerte{
			database.Zaehlerwerte{Wert: d21, Zeitstempel: time.Date(2021, time.January, 01, 0, 0, 0, 0, location)},
			database.Zaehlerwerte{Wert: d20, Zeitstempel: time.Date(2020, time.January, 01, 0, 0, 0, 0, location)},
			database.Zaehlerwerte{Wert: d19, Zeitstempel: time.Date(2019, time.January, 01, 0, 0, 0, 0, location)},
			database.Zaehlerwerte{Wert: d18, Zeitstempel: time.Date(2018, time.January, 01, 0, 0, 0, 0, location)}}

		specialcase, _ := strconv.ParseFloat(record[8], 64)
		kaelte.Spezialfall = int32(specialcase)

		kaelte.GebaeudeRef = []int32{}

		for i := 0; i < 10; i++ {
			if record[i+9] == "" {
				continue
			}

			temp, _ := strconv.ParseInt(record[i+9], 10, 32)
			kaelte.GebaeudeRef = append(kaelte.GebaeudeRef, int32(temp))
		}

		fmt.Println(kaelte)

		b, _ := bson.MarshalExtJSONIndent(kaelte, false, false, "", " ")
		fmt.Println(string(b))

		for _, b1 := range b {
			data = append(data, b1)
		}
		for _, b1 := range []byte(",\n") {
			data = append(data, b1)
		}
		kaelteArray = append(kaelteArray, kaelte)
	}

	data = data[0 : len(data)-2]
	for _, b1 := range []byte("\n]") {
		data = append(data, b1)
	}

	fmt.Println(kaelteArray)
	fmt.Println(string(data))
	_ = ioutil.WriteFile("importer\\Kaeltedaten.json", data, 0644)
}
