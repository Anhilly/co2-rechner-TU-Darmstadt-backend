package main

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/server"
)

func main() {
	err := database.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	server.StartServer()
}
