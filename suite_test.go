package resource_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"testing"

	resource "github.com/frantjc/jenkins-job-resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var bins struct {
	In    string `json:"in"`
	Out   string `json:"out"`
	Check string `json:"check"`
}

var (
	jenkinsURL          = os.Getenv("JENKINS_URL")
	jenkinsJob          = os.Getenv("JENKINS_JOB")
	authenticationToken = os.Getenv("JENKINS_JOB_TOKEN")
	jenkinsUsername     = os.Getenv("JENKINS_USERNAME")
	apiToken            = os.Getenv("JENKINS_API_TOKEN")
	jobArtifact         = os.Getenv("JENKINS_JOB_ARTIFACT")
	jobResult           = os.Getenv("JENKINS_JOB_RESULT")
	source              = resource.Source{
		URL:      jenkinsURL,
		Job:      jenkinsJob,
		Token:    authenticationToken,
		Username: jenkinsUsername,
		Login:    apiToken,
	}
)

func checkJenkinsConfigured() {
	if jenkinsURL == "" || jenkinsJob == "" || authenticationToken == "" || jenkinsUsername == "" || apiToken == "" {
		Skip("must specify $JENKINS_URL, $JENKINS_JOB, $JENKINS_JOB_TOKEN, $JENKINS_USERNAME and $JENKINS_API_TOKEN")
	}
}

func checkJenkinsArtifactConfigured() {
	if jobArtifact == "" {
		Skip("must specify $JENKINS_JOB_ARTIFACT")
	}
}

var _ = SynchronizedBeforeSuite(func() []byte {
	b := bins

	if _, err := os.Stat("/opt/resource/in"); err == nil {
		b.In = "/opt/resource/in"
	} else {
		b.In, err = gexec.Build("github.com/frantjc/jenkins-job-resource/cmd/in")
		Expect(err).ToNot(HaveOccurred())
	}

	if _, err := os.Stat("/opt/resource/out"); err == nil {
		b.Out = "/opt/resource/out"
	} else {
		b.Out, err = gexec.Build("github.com/frantjc/jenkins-job-resource/cmd/out")
		Expect(err).ToNot(HaveOccurred())
	}

	if _, err := os.Stat("/opt/resource/check"); err == nil {
		b.Check = "/opt/resource/check"
	} else {
		b.Check, err = gexec.Build("github.com/frantjc/jenkins-job-resource/cmd/check")
		Expect(err).ToNot(HaveOccurred())
	}

	j, err := json.Marshal(b)
	Expect(err).ToNot(HaveOccurred())

	return j
}, func(bp []byte) {
	err := json.Unmarshal(bp, &bins)
	Expect(err).ToNot(HaveOccurred())

	if !(jenkinsURL == "" || jenkinsJob == "" || authenticationToken == "" || jenkinsUsername == "" || apiToken == "") {
		// make sure the job has at least 1 build
		var (
			req  resource.OutRequest
			resp resource.OutResponse
		)

		req.Source = source

		cmd := exec.Command(bins.Out) //nolint:gosec

		payload, err := json.Marshal(req)
		Expect(err).ToNot(HaveOccurred())

		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		cmd.Stdin = bytes.NewBuffer(payload)
		cmd.Stdout = outBuf
		cmd.Stderr = io.MultiWriter(GinkgoWriter, errBuf)

		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred())

		err = json.Unmarshal(outBuf.Bytes(), &resp)
		Expect(err).ToNot(HaveOccurred())
	}
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	gexec.CleanupBuildArtifacts()
})

func TestJenkinsJobResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JenkinsJobResource Suite")
}
