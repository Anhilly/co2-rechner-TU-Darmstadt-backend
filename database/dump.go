package database

import (
	"os/exec"
	"time"
)

const containerDir = "/autoDump/"

/**
Funktion erstellt ein Dump der Abbildung mit mongodump im Verzeichnis "containerDir + directoryName + timestamp".
Zurueckgeliefert wird der Ordnername mit Timestamp.
*/
func CreateDump(directoryName string) (string, error) {
	dirTimestamp := directoryName + time.Now().Format("20060102150405") // Format: yyyyMMddHHmmss

	cmd := exec.Command("docker", "exec", "-i", "mongodb", "/usr/bin/mongodump",
		"--username", username, "--password", password, "--authenticationDatabase", "admin",
		"--db", dbName, "--out", containerDir+dirTimestamp)

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
Funktion spielt einen Dump, der in "containerDir + directoryName" liegt, wieder in die Datenbank ein mittels mongorestore.
*/
func RestoreDump(directoryName string) error {
	cmd := exec.Command("docker", "exec", "-i", "mongodb", "/usr/bin/mongorestore",
		"--username", username, "--password", password, "--authenticationDatabase", "admin",
		"--drop", "--db", dbName, containerDir+directoryName+"/"+dbName)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
