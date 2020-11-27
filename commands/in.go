package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/logsquaredn/jenkins-job-resource"
	"github.com/yosida95/golang-jenkins"
)

// In runs the in script which checks stdin for a JSON object of the form of an InRequest
// fetches and writes the requested Version as well as Metadata about it to stdout and
// writes each of its output Artifacts to src/Artifact.DisplayPath/Artifact.FileName
func (c *Command) In() error {
	var req resource.InRequest

	decoder := json.NewDecoder(c.stdin)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	if len(c.args) < 2 {
		return fmt.Errorf("destination path not specified")
	}

	src := c.args[1]

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

	err = os.MkdirAll(src, 0755)
	if err != nil {
		return fmt.Errorf("unable to make directory %s", src)
	}

	for _, build := range builds {
		if build.Number == req.Version.Number {
			resp.Version = resource.ToVersion(&build)
			resp.Metadata = []resource.Metadata{
				{ Name: "description", Value: build.Description },
				{ Name: "displayName", Value: build.FullDisplayName },
				{ Name: "id", Value: build.Id },
				{ Name: "url", Value: build.Url },
				{ Name: "duration", Value: strconv.Itoa(build.Duration) },
				{ Name: "estimatedDuration", Value: strconv.Itoa(build.EstimatedDuration) },
				{ Name: "result", Value: build.Result },
			}

			for _, artifact := range build.Artifacts {
				b, err := jenkins.GetArtifact(build, artifact)
				if err != nil {
					ioutil.WriteFile(filepath.Join(src, artifact.DisplayPath, artifact.FileName), b, 0644)
				}
			}

			break;
		}
	}

	if resp.Version.Number != req.Version.Number {
		return fmt.Errorf("requested version of resource %d unavailable: %s", req.Version.Number, err)
	}

	err = json.NewEncoder(c.stdout).Encode(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}
