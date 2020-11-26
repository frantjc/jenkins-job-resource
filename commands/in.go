package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/logsquaredn/jenkins-job-resource"
	"github.com/yosida95/golang-jenkins"
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

	jenkins := gojenkins.NewJenkins(
		&gojenkins.Auth{
			Username: req.Source.Username,
			ApiToken: req.Source.APIToken,
		},
		req.Source.URL,
	)

	job, err := jenkins.GetJob(req.Source.Job)
	if err != nil {
		return fmt.Errorf("unable to find job %s: %s", req.Source.Job, err)
	}

	builds := job.Builds

	var resp resource.InResponse

	if len(builds) > 0 {
		build := builds[0]
		resp.Version = resource.ToVersion(&build)
		
		resp.Metadata = []resource.Metadata{
			{ Name: "description", Value: build.Description },
			{ Name: "displayName", Value: build.FullDisplayName },
			{ Name: "id", Value: build.Id },
			{ Name: "url", Value: build.Url },
			{ Name: "duration", Value: strconv.Itoa(build.Duration) },
			{ Name: "estimatedDuration", Value: strconv.Itoa(build.EstimatedDuration) },
		}
	}

	err = json.NewEncoder(i.stdout).Encode(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}
