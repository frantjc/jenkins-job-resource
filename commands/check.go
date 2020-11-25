package commands

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/logsquaredn/jenkins-job-resource"
)

type Check struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

func NewCheck(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
) *Check {
	return &Check{
		stdin,
		stderr,
		stdout,
		args,
	}
}

func (c *Check) Execute() error {
	var req resource.CheckRequest

	decoder := json.NewDecoder(c.stdin)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	// currently impossible to get error here
	jenk, _ := resource.NewJenkins(&resource.JenkinsInput{
		URL: req.Source.URL,
		BasicCredentials: resource.BasicCredentials{
			Username: req.Source.Username,
			Password: req.Source.Password,
		},
	})

	job, err := jenk.GetJob(req.Source.Job)
	if err != nil {
		return fmt.Errorf("unable to find job %s: %s", req.Source.Job, err)
	}

	info, err := job.GetInfo()
	if err != nil {
		return fmt.Errorf("unable to get info for job %s: %s", req.Source.Job, err)
	}

	var resp resource.CheckResponse

	if hasBuilds := len(info.Builds) > 0; hasBuilds {
		if req.Version != nil {
			resp = append(resp, info.Builds[len(info.Builds) - 1])
		} else {
			for _, build := range info.Builds {
				if build.Number >= req.Version.Number {
					resp = append(resp, build)
				}
			}
		}

		if foundResp := len(resp) > 0; !foundResp {
			resp = append(resp, info.Builds[len(info.Builds) - 1])
		}
	}

	err = json.NewEncoder(c.stdout).Encode(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}
