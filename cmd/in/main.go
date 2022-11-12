package main

import (
	"log"
	"os"

	"github.com/frantjc/jenkins-job-resource/pkg/command"
)

func main() {
	if err := command.NewJenkinsJobResource(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		os.Args,
	).In(); err != nil {
		log.Fatal(err)
	}
}
