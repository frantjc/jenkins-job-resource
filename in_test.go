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

var _ = Describe("In", func() {
	var (
		src    string
		req    resource.InRequest
		resp   resource.InResponse
		cmdErr error
	)

	BeforeEach(func() {
		checkJenkinsConfigured()

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

	Context("when called with a version that does not exist", func() {
		BeforeEach(func() {
			req.Source = source

			req.Version = resource.Version{
				Build: 0,
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
				Build: 1,
			}
		})

		It("captures metadata", func() {
			Expect(cmdErr).NotTo(HaveOccurred())
			Expect(len(resp.Metadata)).To(BeNumerically(">", 0))
		})

		It("gets the requested version", func() {
			Expect(cmdErr).NotTo(HaveOccurred())
			Expect(resp.Version.Build).To(Equal(req.Version.Build))
		})

		It("gets the version's artifacts", func() {
			checkJenkinsArtifactConfigured()
			Expect(cmdErr).NotTo(HaveOccurred())
			_, err := os.Stat(filepath.Join(src, jobArtifact))
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when skip_download is true", func() {
			BeforeEach(func() {
				req.Params.SkipDownload = true
			})

			AfterEach(func() {
				req.Params.SkipDownload = false
			})
	
			It("doesn't get the version's artifacts", func() {
				checkJenkinsArtifactConfigured()
				Expect(cmdErr).NotTo(HaveOccurred())
				_, err := os.Stat(filepath.Join(src, jobArtifact))
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the artifact matches the regexp", func() {
			BeforeEach(func() {
				req.Params.Regexp = []string{jobArtifact}
			})

			AfterEach(func() {
				req.Params.Regexp = nil
			})
	
			It("gets the version's artifacts", func() {
				checkJenkinsArtifactConfigured()
				Expect(cmdErr).NotTo(HaveOccurred())
				_, err := os.Stat(filepath.Join(src, jobArtifact))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the artifact doesn't match the regexp", func() {
			BeforeEach(func() {
				req.Params.Regexp = []string{"messup" + jobArtifact}
			})

			AfterEach(func() {
				req.Params.Regexp = nil
			})
	
			It("doesn't get the version's artifacts", func() {
				checkJenkinsArtifactConfigured()
				Expect(cmdErr).NotTo(HaveOccurred())
				_, err := os.Stat(filepath.Join(src, jobArtifact))
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the result isn't in accept_results", func() {
			BeforeEach(func() {
				req.Params.AcceptResults = []string{"messup" + jobResult}
			})

			AfterEach(func() {
				req.Params.AcceptResults = nil
			})
	
			It("errors", func() {
				checkJenkinsArtifactConfigured()
				Expect(cmdErr).To(HaveOccurred())
			})
		})

		Context("when the result is in accept_results", func() {
			BeforeEach(func() {
				req.Params.AcceptResults = []string{jobResult}
			})

			AfterEach(func() {
				req.Params.AcceptResults = nil
			})
	
			It("doesn't error", func() {
				checkJenkinsArtifactConfigured()
				Expect(cmdErr).NotTo(HaveOccurred())
			})
		})
	})
})
