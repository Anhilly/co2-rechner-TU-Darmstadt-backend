package main

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
)

func main() {
	//importer.ImportEnergieversorgung()
	//importer.ImportGebaeude()
	//importer.ImportStromzaehler()
	//importer.ImportWaermedaten()
	//importer.ImportKaeltedaten()

	database.ConnectDatabase()

	data, err := database.ITGeraeteFind(100)

	fmt.Println(err)
	fmt.Println(data)

	database.DisconnectDatabase()
}
