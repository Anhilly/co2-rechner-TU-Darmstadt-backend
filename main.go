package main

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/server"
)

func main() {
	//importer.ImportEnergieversorgung()
	//importer.ImportGebaeude()
	//importer.ImportStromzaehler()
	//importer.ImportWaermedaten()
	//importer.ImportKaeltedaten()

	database.ConnectDatabase()

	server.StartServer()

	//database.DisconnectDatabase()
}
