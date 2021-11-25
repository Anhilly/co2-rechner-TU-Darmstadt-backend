package importer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"io/ioutil"
	"os"
	"strconv"
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
	in, err := os.Open("importer\\Gebaeudedaten.csv")
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

	var gebaeude database.Gebaeude
	var gebaeudeArray []database.Gebaeude

	for _, record := range rawCSVdata {
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

		fmt.Println(gebaeude)

		b, _ := json.Marshal(gebaeude)
		fmt.Println(string(b))

		gebaeudeArray = append(gebaeudeArray, gebaeude)
	}

	file, _ := json.MarshalIndent(gebaeudeArray, "", " ")
	_ = ioutil.WriteFile("importer\\gebaeudedaten.json", file, 0644)

	//fmt.Println(energieversorgungArray)
}

func ImportStromzaehler() {
	in, err := os.Open("importer/Stromdaten_aufbereitet.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	reader := csv.NewReader(in)
	reader.FieldsPerRecord = -1
	//reader.Comma = ';'

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var strom database.Stromzaehler
	var stromArray []database.Stromzaehler

	for _, record := range rawCSVdata {
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

		strom.Zaehlerdaten = []database.Zaehlerwerte{
			database.Zaehlerwerte{Wert: d21, Zeitstempel: "2021-01-01T00:00:00Z"},
			database.Zaehlerwerte{Wert: d20, Zeitstempel: "2020-01-01T00:00:00Z"},
			database.Zaehlerwerte{Wert: d19, Zeitstempel: "2019-01-01T00:00:00Z"},
			database.Zaehlerwerte{Wert: d18, Zeitstempel: "2018-01-01T00:00:00Z"}}

		fmt.Println(strom)

		b, _ := json.Marshal(strom)
		fmt.Println(string(b))

		stromArray = append(stromArray, strom)
	}

	file, _ := json.MarshalIndent(stromArray, "", " ")
	_ = ioutil.WriteFile("importer/Stromdaten.json", file, 0644)

	//fmt.Println(energieversorgungArray)
}

func ImportWaermedaten(){
	in, err := os.Open("importer/waermedaten_aufbereitet.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	reader := csv.NewReader(in)
	reader.FieldsPerRecord = -1
	//reader.Comma = ';'

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var waerme database.Waermezaehler
	var waermeArray []database.Waermezaehler

	for _, record := range rawCSVdata {
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

		waerme.Zaehlerdaten = []database.Zaehlerwerte{
			database.Zaehlerwerte{Wert: d20, Zeitstempel: "2020-01-01T00:00:00Z"},
			database.Zaehlerwerte{Wert: d19, Zeitstempel: "2019-01-01T00:00:00Z"},
			database.Zaehlerwerte{Wert: d18, Zeitstempel: "2018-01-01T00:00:00Z"}}

		fmt.Println(waerme)

		b, _ := json.Marshal(waerme)
		fmt.Println(string(b))

		waermeArray = append(waermeArray, waerme)
	}

	file, _ := json.MarshalIndent(waermeArray, "", " ")
	_ = ioutil.WriteFile("importer/waermedaten.json", file, 0644)

	//fmt.Println(energieversorgungArray)
}
