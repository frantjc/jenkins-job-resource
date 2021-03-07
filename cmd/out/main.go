package main

import (
	"log"
	"os"

	"github.com/logsquaredn/jenkins-job-resource/commands"
)

func main() {
	command := commands.NewJenkinsJobResource(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		os.Args,
	)

	err := command.Out()
	if err != nil {
		log.Fatal(err)
	}
}
