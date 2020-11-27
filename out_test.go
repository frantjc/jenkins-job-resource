package resource_test

import (
	"io"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"

	"github.com/logsquaredn/jenkins-job-resource/commands"
	// resource "github.com/logsquaredn/jenkins-job-resource"
)

var _ = Describe("Out", func () {
	var (
		stdin io.Reader
		stderr io.Writer
		stdout io.Writer
		args []string
		jenkinsJobResource *commands.JenkinsJobResource
	)


	BeforeEach(func() {
		// initialize mock stdin, stderr, stdout and args
		jenkinsJobResource = commands.NewJenkinsJobResource(
			stdin,
			stderr,
			stdout,
			args,
		)
	})

	It("", func() {
		_ = jenkinsJobResource
	})
})
