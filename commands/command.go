package commands

import (
	"io"
)

type Command struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

func NewCommand(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
) *Command {
	return &Command{
		stdin:  stdin,
		stderr: stderr,
		stdout: stdout,
		args:   args,
	}
}
