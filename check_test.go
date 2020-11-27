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

var _ = Describe("Check", func () {
	var (
		req             resource.CheckRequest
		resp		    resource.CheckResponse
		cmdErr			error
	)


	BeforeEach(func() {
		checkEnvConfigured()

		req.Source = resource.Source{}
		req.Version = nil

		resp = nil
	})

	JustBeforeEach(func() {
		cmd := exec.Command(bins.Check)

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

	Context("when called with no version", func() {
		BeforeEach(func() {
			req.Source = source
		})

		It("returns the current build number", func() {
			Expect(cmdErr).NotTo(HaveOccurred())
			Expect(len(resp)).To(Equal(1))
			Expect(resp[0].Number).To(BeNumerically(">", 0))
		})
	})

	Context("when called with a version", func() {
		BeforeEach(func() {
			req.Source = source

			req.Version = &resource.Version{
				Number: 1,
			}
		})

		It("returns all builds since the given version", func() {
			if cmdErr != nil {
				Expect(len(resp)).To(BeNumerically(">", 0))
				for _, build := range resp {
					Expect(build.Number).To(BeNumerically(">=", req.Version.Number))
				}
			} else {
				Skip("the specified $JENKINS_JOB must use a jenkinsfile like jenkins-job-resource/cicd/pipelines/jenkinsfile")
			}
		})
	})
})
