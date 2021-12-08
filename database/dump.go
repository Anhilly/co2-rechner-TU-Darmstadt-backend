package database

import (
	"log"
	"os/exec"
)

func CreateDump() {
	cmd := exec.Command("docker", "exec", "-i", "mongodb", "/usr/bin/mongodump",
		"--username", username, "--password", password, "--authenticationDatabase", "admin",
		"--db", dbName, "--out", "/testDump")

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("docker", "cp", "mongodb:/testDump", "testDump")

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
