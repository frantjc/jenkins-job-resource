package commands

import (
	"fmt"
	"time"

	goqs "github.com/google/go-querystring/query"
	resource "github.com/logsquaredn/jenkins-job-resource"
)

const defaultCause = "Triggered by Concourse"
const defaultDescription = "Build triggered by Concourse"

// Out runs the in script which checks stdin for a JSON object of the form of an OutRequest
// triggers a new build and then fetches and writes it as well as Metadata about it to stdout
func (j *JenkinsJobResource) Out() error {
	var (
		req resource.OutRequest
	    resp resource.OutResponse
	)

	err := j.readInput(&req)
	if err != nil {
		return err
	}

	jenkins := j.newJenkins(req.Source)

	job, err := jenkins.GetJob(req.Source.Job)
	if err != nil {
		return fmt.Errorf("unable to find job %s: %s", req.Source.Job, err)
	}

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

	err = jenkins.Build(job, params)
	if err != nil {
		return fmt.Errorf("unable to trigger build payload: %s", err)
	}

	for {
		if build, err := jenkins.GetBuild(job, job.LastCompletedBuild.Number + 1); err == nil {
			// currently don't care if there is an error here
			jenkins.SetBuildDescription(build, req.Params.Description)

			resp.Version = j.getVersion(&build)
			resp.Metadata = j.getMetadata(&build)
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	j.writeOutput(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}
