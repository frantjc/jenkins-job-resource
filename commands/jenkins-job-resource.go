package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	resource "github.com/logsquaredn/jenkins-job-resource"
	gojenkins "github.com/yosida95/golang-jenkins"
)

// JenkinsJobResource struct which has the Check, In, and Out methods on it which comprise
// the three scripts needed to implement a Concourse Resource Type
type JenkinsJobResource struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

// NewJenkinsJobResource creates a new JenkinsJobResource struct
func NewJenkinsJobResource(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
) *JenkinsJobResource {
	return &JenkinsJobResource{
		stdin:  stdin,
		stderr: stderr,
		stdout: stdout,
		args:   args,
	}
}

func (j *JenkinsJobResource) readInput(req interface{}) error {
	decoder := json.NewDecoder(j.stdin)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	return nil
}

func (j *JenkinsJobResource) writeOutput(resp interface{}) error {
	err := json.NewEncoder(j.stdout).Encode(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}

func (j *JenkinsJobResource) getMetadata(build *gojenkins.Build) []resource.Metadata {
	if build != nil {
		return []resource.Metadata{
			{Name: "description", Value: build.Description},
			{Name: "displayName", Value: build.FullDisplayName},
			{Name: "id", Value: build.Id},
			{Name: "url", Value: build.Url},
			{Name: "duration", Value: strconv.Itoa(build.Duration)},
			{Name: "estimatedDuration", Value: strconv.Itoa(build.EstimatedDuration)},
			{Name: "result", Value: build.Result},
		}
	}

	return []resource.Metadata{}
}

func (j *JenkinsJobResource) getSrc() (string, error) {
	if len(j.args) < 2 {
		return "", fmt.Errorf("destination path not specified")
	}
	return j.args[1], nil
}

func (j *JenkinsJobResource) newJenkins(s resource.Source) *gojenkins.Jenkins {
	return gojenkins.NewJenkins(
		&gojenkins.Auth{
			Username: s.Username,
			ApiToken: s.Login,
		},
		s.URL,
	)
}

func (j *JenkinsJobResource) getVersion(b *gojenkins.Build) resource.Version {
	return resource.Version{
		Number: b.Number,
	}
}
