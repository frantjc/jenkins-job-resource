package commands

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	goqs "github.com/google/go-querystring/query"
	"github.com/logsquaredn/jenkins-job-resource"
	"github.com/yosida95/golang-jenkins"
)

const defaultCause = "Triggered by Concourse"
const defaultDescription = "Build triggered by Concourse"

// Out runs the in script which checks stdin for a JSON object of the form of an OutRequest
// triggers a new build and then fetches and writes it as well as Metadata about it to stdout
func (c *Command) Out() error {
	var req resource.OutRequest

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

	params, err := goqs.Values(req.Params.BuildParams)
	if err != nil {
		return fmt.Errorf("unable to turn build_params into a query string: %s", err)
	}

	if req.Params.Cause != "" {
		params.Set("cause", req.Params.Cause)
	} else {
		params.Set("cause", defaultCause)
	}

	if req.Source.Token != "" {
		params.Set("token", req.Source.Token)
	} else {
		return fmt.Errorf("no token supplied to source")
	}

	job, err := jenkins.GetJob(req.Source.Job)
	if err != nil {
		return fmt.Errorf("unable to find job %s: %s", req.Source.Job, err)
	}

	err = jenkins.Build(job, params)
	if err != nil {
		return fmt.Errorf("unable to trigger build payload: %s", err)
	}

	var resp resource.OutResponse

	for jobAfterBuild, err := jenkins.GetJob(req.Source.Job); resp.Version.Number == 0; jobAfterBuild, err = jenkins.GetJob(req.Source.Job) {
		if err != nil {
			return fmt.Errorf("unable to find job %s after triggering build: %s", req.Source.Job, err)
		}

		lastBuild := jobAfterBuild.LastCompletedBuild
		if lastBuild.Number > job.LastCompletedBuild.Number {
			err = jenkins.SetBuildDescription(lastBuild, defaultDescription)
			if err != nil {
				// Do I care?
			}

			resp.Version = resource.ToVersion(&lastBuild)
			resp.Metadata = []resource.Metadata{
				{ Name: "description", Value: lastBuild.Description },
				{ Name: "displayName", Value: lastBuild.FullDisplayName },
				{ Name: "id", Value: lastBuild.Id },
				{ Name: "url", Value: lastBuild.Url },
				{ Name: "duration", Value: strconv.Itoa(lastBuild.Duration) },
				{ Name: "estimatedDuration", Value: strconv.Itoa(lastBuild.EstimatedDuration) },
				{ Name: "result", Value: lastBuild.Result },
			}

			if lastBuild.Result != "SUCCESS" {
				return fmt.Errorf("%s %s resulted in %s", req.Source.Job, lastBuild.Id, lastBuild.Result)
			}

			break;
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	err = json.NewEncoder(c.stdout).Encode(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}
