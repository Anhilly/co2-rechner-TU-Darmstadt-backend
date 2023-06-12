package main

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/server"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)

var Mode = "dev" // changed using symbol substitution at link time

func main() {
	var filename string
	if Mode == "prod" {
		print("prod mode")
		filename = prod_log_filename
	} else if Mode == "dev" {
		print("dev mode")
		filename = dev_log_filename
	} else {
		panic("MODE not set")
	}

	logger := lumberjack.Logger{
		Filename:  filename,
		MaxSize:   100,
		LocalTime: true,
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetOutput(&logger)

	err := database.ConnectDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	server.StartServer(&logger)
}
