package main

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/server"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)

func main() {
	logger := lumberjack.Logger{
		//Filename:  "/app/backend/logs/go.log",
		Filename:  "/home/tobias/Desktop/go.log",
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
