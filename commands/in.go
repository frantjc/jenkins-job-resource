package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/logsquaredn/jenkins-job-resource"
)

type In struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

func NewIn(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
) *In {
	return &In{
		stdin,
		stderr,
		stdout,
		args,
	}
}

func (i *In) Execute() error {
	var req resource.InRequest

	decoder := json.NewDecoder(i.stdin)
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

	builds, err := job.GetBuilds()
	if err != nil {
		return fmt.Errorf("unable to get builds for job %s: %s", req.Source.Job, err)
	}

	var resp resource.InResponse

	if len(builds) > 0 {
		build := builds[0]
		resp.Version = build.ToVersion()

		info, err := build.GetInfo()
		if err != nil {
			return fmt.Errorf("unable to get metadata for build %d: %s", build.Number, err)
		}
		
		resp.Metadata = []resource.Metadata{
			{ Name: "description", Value: info.Description },
			{ Name: "displayName", Value: info.DisplayName },
			{ Name: "id", Value: info.ID },
			{ Name: "url", Value: info.URL },
			{ Name: "duration", Value: strconv.Itoa(info.Duration) },
			{ Name: "estimatedDuration", Value: strconv.Itoa(info.EstimatedDuration) },
		}
	}

	err = json.NewEncoder(i.stdout).Encode(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}
