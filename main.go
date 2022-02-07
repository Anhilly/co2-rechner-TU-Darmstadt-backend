package main

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/server"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	err := database.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	server.StartServer()
}
