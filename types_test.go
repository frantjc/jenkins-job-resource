package resource_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	resource "github.com/logsquaredn/jenkins-job-resource"
)

var _ = Describe("Version", func () {
	It("Should marshal build string into int", func() {
		var version resource.Version
		raw := []byte(`{ "build": "11" }`)

		err := json.Unmarshal(raw, &version)
		Expect(err).ToNot(HaveOccurred())
		Expect(version.Build).To(Equal(11))
	})

	It("Should unmarshal build int into string", func() {
		version := resource.Version{Build: 11}

		json, err := json.Marshal(version)
		Expect(err).ToNot(HaveOccurred())
		Expect(json).To(MatchJSON(`{"build":"11"}`))
	})
})
