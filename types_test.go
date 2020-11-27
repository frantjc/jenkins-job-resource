package resource_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	resource "github.com/logsquaredn/jenkins-job-resource"
)

var _ = Describe("Version", func () {
	It("Should marshal number string into int", func() {
		var version resource.Version
		raw := []byte(`{ "number": "11" }`)

		err := json.Unmarshal(raw, &version)
		Expect(err).ToNot(HaveOccurred())
		Expect(version.Build).To(Equal(11))
	})

	It("Should unmarshal number int into string", func() {
		version := resource.Version{Build: 11}

		json, err := json.Marshal(version)
		Expect(err).ToNot(HaveOccurred())
		Expect(json).To(MatchJSON(`{"number":"11"}`))
	})
})

var _ = Describe("OutRequest", func () {
	It("Should return the defaultDescription if no Description is provided", func() {
		request := resource.OutRequest{
			Params: resource.PutParams{
				Description: "",
			},
		}

		description := request.Description()
		Expect(description).To(Equal("Build triggered by Concourse"))
	})


	It("Should return the defaultCause if no Cause is provided", func() {
		request := resource.OutRequest{
			Params: resource.PutParams{
				Cause: "",
			},
		}

		cause := request.Cause()
		Expect(cause).To(Equal("Triggered by Concourse"))
	})

	It("Should return the provided description if a Description is provided", func() {
		request := resource.OutRequest{
			Params: resource.PutParams{
				Description: "description",
			},
		}

		description := request.Description()
		Expect(description).To(Equal("description"))
	})


	It("Should return the provided cause if no a Cause is provided", func() {
		request := resource.OutRequest{
			Params: resource.PutParams{
				Cause: "cause",
			},
		}

		cause := request.Cause()
		Expect(cause).To(Equal("cause"))
	})
})
