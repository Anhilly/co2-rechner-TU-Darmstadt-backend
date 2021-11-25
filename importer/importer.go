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
