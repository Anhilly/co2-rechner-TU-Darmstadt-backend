package database

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"os/exec"
	"time"
)

/**
Funktion erstellt ein Dump der Abbildung mit mongodump im Verzeichnis "DumpPath + directoryName + timestamp".
Zurueckgeliefert wird der Ordnername mit Timestamp.
*/
func CreateDump(directoryName string) (string, error) {
	dirTimestamp := directoryName + time.Now().Format("20060102150405") // Format: yyyyMMddHHmmss

	cmd := exec.Command("docker", "exec", "-i", "mongodb", "/usr/bin/mongodump",
		"--username", username, "--password", password, "--authenticationDatabase", "admin",
		"--db", dbName, "--out", structs.DumpPath+dirTimestamp)

	err := cmd.Run()
	if err != nil {
		return "", err
	}
	/*
		cmd = exec.Command("docker", "cp", "mongodb:/testDump", "testDump")

		err = cmd.Run()
		if err != nil {
			return "", err
		}*/

	return dirTimestamp, nil
}

/**
Funktion spielt einen Dump, der in "DumpPath + directoryName" liegt, wieder in die Datenbank ein mittels mongorestore.
*/
func RestoreDump(directoryName string) error {
	cmd := exec.Command("docker", "exec", "-i", "mongodb", "/usr/bin/mongorestore",
		"--username", username, "--password", password, "--authenticationDatabase", "admin",
		"--drop", "--db", dbName, structs.DumpPath+directoryName+"/"+dbName)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
