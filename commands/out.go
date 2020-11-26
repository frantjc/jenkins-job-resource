package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/logsquaredn/jenkins-job-resource"
)

type Out struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

func NewOut(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
) *Out {
	return &Out{
		stdin:  stdin,
		stderr: stderr,
		stdout: stdout,
		args:   args,
	}
}

func (o *Out) Execute() error {
	var req resource.OutRequest

	decoder := json.NewDecoder(o.stdin)
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

	var resp resource.OutResponse

	build, err := job.Build(req.Source.Token, req.Params.Cause, &req.Params.BuildParams)
	if err != nil {
		return fmt.Errorf("unable to build job %s: %s", req.Source.Job, err)
	}

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

	err = json.NewEncoder(o.stdout).Encode(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}
