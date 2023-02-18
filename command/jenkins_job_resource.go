package command

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	gojenkins "github.com/yosida95/golang-jenkins"

	resource "github.com/frantjc/jenkins-job-resource"
)

// JenkinsJobResource struct which has the Check, In, and Out methods on it which comprise
// the three scripts needed to implement a Concourse Resource Type.
type JenkinsJobResource struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

// NewJenkinsJobResource creates a new JenkinsJobResource struct.
func NewJenkinsJobResource(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
) *JenkinsJobResource {
	return &JenkinsJobResource{
		stdin,
		stderr,
		stdout,
		args,
	}
}

func (r *JenkinsJobResource) readInput(req interface{}) error {
	decoder := json.NewDecoder(r.stdin)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	return nil
}

func (r *JenkinsJobResource) writeOutput(resp interface{}) error {
	err := json.NewEncoder(r.stdout).Encode(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}

func (r *JenkinsJobResource) getMetadata(build *gojenkins.Build) []resource.Metadata {
	if build != nil {
		return []resource.Metadata{
			{Name: "name", Value: build.FullDisplayName},
			{Name: "id", Value: build.Id},
			{Name: "url", Value: build.Url},
			{Name: "duration", Value: strconv.Itoa(build.Duration)},
			{Name: "expectedDuration", Value: strconv.Itoa(build.EstimatedDuration)},
			{Name: "result", Value: build.Result},
		}
	}

	return []resource.Metadata{}
}

func (r *JenkinsJobResource) acceptResult(build *gojenkins.Build, acceptResults []string) error {
	if len(acceptResults) == 0 {
		return nil
	}

	accept := false
	result := build.Result
	for _, acceptable := range acceptResults {
		if result == acceptable {
			accept = true
			break
		}
	}

	if !accept {
		return fmt.Errorf("build resulted in %s", result)
	}

	return nil
}

func (r *JenkinsJobResource) getSrc() (string, error) {
	if len(r.args) < 2 {
		return "", fmt.Errorf("destination path not specified")
	}
	return r.args[1], nil
}

func (r *JenkinsJobResource) newJenkins(s resource.Source) *gojenkins.Jenkins {
	return gojenkins.NewJenkins(
		&gojenkins.Auth{
			Username: s.Username,
			ApiToken: s.Login,
		},
		s.URL,
	)
}

func (r *JenkinsJobResource) getVersion(b gojenkins.Build) resource.Version {
	return resource.Version{
		Build: b.Number,
	}
}

func (r *JenkinsJobResource) writeMetadata(mds []resource.Metadata) error {
	src, err := r.getSrc()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(src, ".metadata"), 0755)
	if err != nil {
		return fmt.Errorf("unable to make directory %s", filepath.Join(src, ".metadata"))
	}

	for _, md := range mds {
		if err = os.WriteFile(filepath.Join(src, ".metadata", md.Name), []byte(md.Value), 0600); err != nil {
			return err
		}
	}

	return nil
}

func (r *JenkinsJobResource) getCause(p *resource.PutParams) (string, error) {
	if p.CauseFile != "" {
		src, err := r.getSrc()
		if err != nil {
			return "", err
		}

		b, err := os.ReadFile(filepath.Join(src, p.CauseFile))
		if err != nil {
			return "", err
		}

		return r.expandEnv(string(b)), nil
	} else if p.Cause != "" {
		return r.expandEnv(p.Cause), nil
	}

	return r.expandEnv("caused by $ATC_EXTERNAL_URL/builds/$BUILD_ID"), nil
}

func (r *JenkinsJobResource) getDescription(p *resource.PutParams) (string, error) {
	if p.DescriptionFile != "" {
		src, err := r.getSrc()
		if err != nil {
			return "", err
		}

		b, err := os.ReadFile(filepath.Join(src, p.DescriptionFile))
		if err != nil {
			return "", err
		}

		return r.expandEnv(string(b)), nil
	} else if p.Description != "" {
		return r.expandEnv(p.Description), nil
	}

	return r.expandEnv("build triggered by $ATC_EXTERNAL_URL/builds/$BUILD_ID"), nil
}

func (r *JenkinsJobResource) expandEnv(s string) string {
	return os.Expand(s, func(v string) string {
		switch v {
		case "BUILD_ID", "BUILD_NAME", "BUILD_JOB_NAME", "BUILD_PIPELINE_NAME", "BUILD_TEAM_NAME", "ATC_EXTERNAL_URL":
			return os.Getenv(v)
		}
		return "$" + v
	})
}
