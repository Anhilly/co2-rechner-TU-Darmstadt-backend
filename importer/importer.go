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
