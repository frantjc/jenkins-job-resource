package resource_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"testing"

	resource "github.com/logsquaredn/jenkins-job-resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var bins struct {
	In    string `json:"in"`
	Out   string `json:"out"`
	Check string `json:"check"`
}

var jenkinsUrl = os.Getenv("JENKINS_URL")
var jenkinsJob = os.Getenv("JENKINS_JOB")
var authenticationToken = os.Getenv("JENKINS_JOB_TOKEN")
var jenkinsUsername = os.Getenv("JENKINS_USERNAME")
var apiToken = os.Getenv("JENKINS_API_TOKEN")
var source = resource.Source{
	URL: jenkinsUrl,
	Job: jenkinsJob,
	Token: authenticationToken,
	Username: jenkinsUsername,
	Login: apiToken,
}

func checkEnvConfigured() {
	if jenkinsUrl == "" || jenkinsJob == "" || authenticationToken == "" || jenkinsUsername == "" || apiToken == "" {
		Skip("must specify $JENKINS_URL, $JENKINS_JOB, $JENKINS_JOB_TOKEN, $JENKINS_USERNAME and $JENKINS_API_TOKEN")
	}
}

var _ = SynchronizedBeforeSuite(func() []byte {
	b := bins

	if _, err := os.Stat("/opt/resource/in"); err == nil {
		b.In = "/opt/resource/in"
	} else {
		b.In, err = gexec.Build("github.com/logsquaredn/jenkins-job-resource/cmd/in")
		Expect(err).ToNot(HaveOccurred())
	}

	if _, err := os.Stat("/opt/resource/out"); err == nil {
		b.Out = "/opt/resource/out"
	} else {
		b.Out, err = gexec.Build("github.com/logsquaredn/jenkins-job-resource/cmd/out")
		Expect(err).ToNot(HaveOccurred())
	}

	if _, err := os.Stat("/opt/resource/check"); err == nil {
		b.Check = "/opt/resource/check"
	} else {
		b.Check, err = gexec.Build("github.com/logsquaredn/jenkins-job-resource/cmd/check")
		Expect(err).ToNot(HaveOccurred())
	}

	j, err := json.Marshal(b)
	Expect(err).ToNot(HaveOccurred())

	return j
}, func(bp []byte) {
	err := json.Unmarshal(bp, &bins)
	Expect(err).ToNot(HaveOccurred())

	// make sure the job has at least 1 build
	var (
		req resource.OutRequest
		resp resource.OutResponse
	)

	req.Source = source

	cmd := exec.Command(bins.Check)

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
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	gexec.CleanupBuildArtifacts()
})

func TestJenkinsJobResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JenkinsJobResource Suite")
}
