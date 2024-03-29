package database

import (
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

	cmd := exec.Command("docker", "exec", "-i", dbContainer, "/usr/bin/mongodump",
		"--username", dbUsername, "--password", dbPassword, "--authenticationDatabase", "admin",
		"--db", dbName, "--out", structs.DumpPath+dirTimestamp)

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
	cmd := exec.Command("docker", "exec", "-i", dbContainer, "/usr/bin/mongorestore",
		"--username", dbUsername, "--password", dbPassword, "--authenticationDatabase", "admin",
		"--drop", "--db", dbName, structs.DumpPath+directoryName+"/"+dbName)

	err := cmd.Run()
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	return nil
}

func RemoveDump(directoryName string) error {
	cmd := exec.Command("docker", "exec", "-i", dbContainer, "rm", "-rf", structs.DumpPath+directoryName)

	err := cmd.Run()
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		return err
	}

	return nil
}
