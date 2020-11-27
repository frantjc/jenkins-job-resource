package main

import (
	"log"
	"os"

	"github.com/logsquaredn/jenkins-job-resource/commands"
)

func main() {
	command := commands.NewCommand(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		os.Args,
	)

	err := command.In()
	if err != nil {
		log.Fatal(err)
	}
}
