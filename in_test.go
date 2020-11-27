package resource_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	resource "github.com/logsquaredn/jenkins-job-resource"
)

var _ = Describe("In", func () {
	var (
		src    string
		req    resource.InRequest
		resp   resource.InResponse
		cmdErr error
	)


	BeforeEach(func() {
		checkEnvConfigured()

		var err error
		src, err = ioutil.TempDir("", "in-jenkins-job-resource")
		Expect(err).ToNot(HaveOccurred())

		req.Source = resource.Source{}
		req.Params = resource.GetParams{}
		req.Version = resource.Version{}

		resp.Version = resource.Version{}
		resp.Metadata = nil
	})

	JustBeforeEach(func() {
		cmd := exec.Command(bins.In, src)

		payload, err := json.Marshal(req)
		Expect(err).ToNot(HaveOccurred())

		outBuf := new(bytes.Buffer)

		cmd.Stdin = bytes.NewBuffer(payload)
		cmd.Stdout = outBuf
		cmd.Stderr = GinkgoWriter

		cmdErr = cmd.Run()

		if cmdErr == nil {
			err = json.Unmarshal(outBuf.Bytes(), &resp)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	AfterEach(func() {
		Expect(os.RemoveAll(src)).To(Succeed())
	})

	Context("when called with a version that does not exists", func() {
		BeforeEach(func() {
			req.Source = source

			req.Version = resource.Version{
				Number: 0,
			}
		})

		It("errors", func() {
			Expect(cmdErr).To(HaveOccurred())
		})
	})

	Context("when called with a version that exists", func() {
		BeforeEach(func() {
			req.Source = source

			req.Version = resource.Version{
				Number: 1,
			}
		})

		It("captures metadata", func() {
			if cmdErr != nil {
				Expect(cmdErr).NotTo(HaveOccurred())
				Expect(len(resp.Metadata)).To(BeNumerically(">", 0))
			} else {
				Skip("the specified $JENKINS_JOB must use a jenkinsfile like jenkins-job-resource/cicd/pipelines/jenkinsfile")
			}
		})

		It("gets the requested version", func() {
			if cmdErr != nil {
				Expect(cmdErr).NotTo(HaveOccurred())
				Expect(resp.Version.Number).To(Equal(req.Version.Number))
			} else {
				Skip("the specified $JENKINS_JOB must use a jenkinsfile like jenkins-job-resource/cicd/pipelines/jenkinsfile")
			}
		})

		It("gets the version's artifacts", func() {
			if cmdErr != nil {
				Expect(cmdErr).NotTo(HaveOccurred())
				_, err := os.Stat(filepath.Join(src, "output.txt"))
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("the specified $JENKINS_JOB must use a jenkinsfile like jenkins-job-resource/cicd/pipelines/jenkinsfile")
			}
		})
	})
})