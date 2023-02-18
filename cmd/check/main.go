package main

import (
	"log"
	"os"

	"github.com/frantjc/jenkins-job-resource/command"
)

func main() {
	if err := command.NewJenkinsJobResource(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		os.Args,
	).Check(); err != nil {
		log.Fatal(err)
	}
}
