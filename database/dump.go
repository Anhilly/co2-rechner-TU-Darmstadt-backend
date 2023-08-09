package database

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/config"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"log"
	"os/exec"
	"runtime/debug"
	"time"
)

// CreateDump erstellt ein Dump der Abbildung mit mongodump im Verzeichnis "DumpPath + timestamp + directoryName".
// Zurueckgeliefert wird der Ordnername mit Timestamp.
// Beim Ausfuehren unter Linux Systemen muss Docker per default Sudo Rechte besitzen, da der Befehl Sudo Rechte benoetigt.
func CreateDump(directoryName string) (string, error) {
	dirTimestamp := time.Now().Format("20060102150405") + directoryName // Format: yyyyMMddHHmmss

	cmd := exec.Command("docker", "exec", "-i", config.ContainerName, "/usr/bin/mongodump",
		"--username", config.Username, "--password", config.Password, "--authenticationDatabase", "admin",
		"--db", config.DBName, "--out", structs.DumpPath+dirTimestamp)

	err := cmd.Run()
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return "", err
	}

	return dirTimestamp, nil
}

// RestoreDump spielt einen Dump, der in "DumpPath + directoryName" liegt, wieder in die Datenbank ein mittels mongorestore.
func RestoreDump(directoryName string) error {
	cmd := exec.Command("docker", "exec", "-i", config.ContainerName, "/usr/bin/mongorestore",
		"--username", config.Username, "--password", config.Password, "--authenticationDatabase", "admin",
		"--drop", "--db", config.DBName, structs.DumpPath+directoryName+"/"+config.DBName)

	err := cmd.Run()
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	return nil
}

func RemoveDump(directoryName string) error {
	cmd := exec.Command("docker", "exec", "-i", config.ContainerName, "rm", "-rf", structs.DumpPath+directoryName)

	err := cmd.Run()
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	return nil
}
