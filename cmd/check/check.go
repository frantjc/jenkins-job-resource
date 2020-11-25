package main

import (
	"log"
	"os"

	"github.com/logsquaredn/jenkins-job-resource/commands"
)

func main() {
	command := commands.NewCheck(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		os.Args,
	)

	err := command.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
