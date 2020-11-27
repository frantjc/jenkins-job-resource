package resource_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestJenkinsJobResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JenkinsJobResource Suite")
}
