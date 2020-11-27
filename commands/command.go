package commands

import (
	"io"
)

// Command struct which has the Check, In, and Out methods on it which comprise
// the three scripts needed to implement a Concourse Resource Type
type Command struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

// NewCommand creates a new Command struct
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
