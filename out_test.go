package resource_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	resource "github.com/logsquaredn/jenkins-job-resource"
)

var _ = Describe("Out", func () {
	var (
		req resource.OutRequest
		resp resource.OutResponse
		cmdErr error
	)


	BeforeEach(func() {
		checkEnvConfigured()

		req.Source = resource.Source{}
		req.Params = resource.PutParams{}

		resp.Version = resource.Version{}
		resp.Metadata = nil
	})

	JustBeforeEach(func() {
		cmd := exec.Command(bins.Out)

		payload, err := json.Marshal(req)
		Expect(err).ToNot(HaveOccurred())

		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)

		cmd.Stdin = bytes.NewBuffer(payload)
		cmd.Stdout = outBuf
		cmd.Stderr = io.MultiWriter(GinkgoWriter, errBuf)

		cmdErr = cmd.Run()

		if cmdErr == nil {
			err = json.Unmarshal(outBuf.Bytes(), &resp)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	Context("with no params", func() {
		BeforeEach(func() {
			req.Source = source
		})

		It("captures metadata", func() {
			if cmdErr == nil {
				Expect(cmdErr).NotTo(HaveOccurred())
				Expect(len(resp.Metadata)).To(BeNumerically(">", 0))
			} else {
				Skip("the specified $JENKINS_JOB must use a jenkinsfile like jenkins-job-resource/cicd/pipelines/jenkinsfile")
			}
		})

		It("gets a version", func() {
			if cmdErr == nil {
				Expect(cmdErr).NotTo(HaveOccurred())
				Expect(resp.Version.Build).To(BeNumerically(">", 0))
			} else {
				Skip("the specified $JENKINS_JOB must use a jenkinsfile like jenkins-job-resource/cicd/pipelines/jenkinsfile")
			}
		})
	})
})
