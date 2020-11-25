package commands

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/logsquaredn/jenkins-job-resource"
)

type In struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

func NewIn(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
) *In {
	return &In{
		stdin,
		stderr,
		stdout,
		args,
	}
}

func (i *In) Execute() error {
	var req resource.InRequest

	decoder := json.NewDecoder(i.stdin)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	return nil
}
