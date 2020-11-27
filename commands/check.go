package commands

import (
	"encoding/json"
	"fmt"

	"github.com/logsquaredn/jenkins-job-resource"
	"github.com/yosida95/golang-jenkins"
)

func (c *Command) Check() error {
	var req resource.CheckRequest

	decoder := json.NewDecoder(c.stdin)
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

	var resp resource.CheckResponse

	if len(builds) > 0 {
		if req.Version != nil {
			for _, build := range builds {
				version := resource.ToVersion(&build)
				if version.Number >= req.Version.Number {
					// prepend
					resp = append([]resource.Version{version}, resp...)
				}
			}
		}

		if len(resp) <= 0 {
			resp = append(resp, resource.ToVersion(&builds[0]))
		}
	}

	err = json.NewEncoder(c.stdout).Encode(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}
