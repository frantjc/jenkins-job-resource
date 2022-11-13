package resource_test

import (
	"encoding/json"
	"testing"

	resource "github.com/frantjc/jenkins-job-resource"
)

func TestVersion(t *testing.T) {
	var (
		raw     = []byte(`{"build":"11"}`)
		version resource.Version
		err     = json.Unmarshal(raw, &version)
	)
	if err != nil {
		t.FailNow()
	}

	if version.Build != 11 {
		t.FailNow()
	}

	marshalled, err := json.Marshal(version)
	if err != nil {
		t.FailNow()
	}

	if string(marshalled) != string(raw) {
		t.FailNow()
	}
}
