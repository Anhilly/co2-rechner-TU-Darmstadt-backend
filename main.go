package main

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/config"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/keycloak"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/server"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)

var mode = "dev" // changed using symbol substitution at link time

func main() {
	// setting up logger
	var filename string
	if mode == "prod" {
		println("prod mode")
		filename = config.ProdLogFilename
	} else if mode == "dev" {
		println("dev mode")
		filename = config.DevLogFilename
	} else if mode == "test" {
		println("test mode")
		filename = config.TestLogFilename
	} else {
		panic("MODE not set")
	}

	println("cgo:", cgoEnabled)

	logger := lumberjack.Logger{
		Filename:  filename,
		MaxSize:   100,
		LocalTime: true,
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetOutput(&logger)

	// setting up database
	err := database.ConnectDatabase(mode)
	if err != nil {
		log.Fatalln(err)
	}

	// setting up keycloak
	err = keycloak.SetupKeycloakClient(mode)
	if err != nil {
		log.Fatalln(err)
	}

	// starting server
	server.StartServer(&logger, mode)
}
