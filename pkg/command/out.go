package command

import (
	"fmt"
	"time"

	resource "github.com/frantjc/jenkins-job-resource"
	goqs "github.com/google/go-querystring/query"
)

// Out runs the in script which checks stdin for a JSON object of the form of an OutRequest
// triggers a new build and then fetches and writes it as well as Metadata about it to stdout.
func (r *JenkinsJobResource) Out() error {
	var (
		req  resource.OutRequest
		resp resource.OutResponse
	)

	err := r.readInput(&req)
	if err != nil {
		return err
	}

	jenkins := r.newJenkins(req.Source)

	job, err := jenkins.GetJob(req.Source.Job)
	if err != nil {
		return fmt.Errorf("unable to find job %s: %s", req.Source.Job, err)
	}

	params, err := goqs.Values(req.Params.BuildParams)
	if err != nil {
		return fmt.Errorf("unable to turn build_params into a query string: %s", err)
	}

	cause, err := r.getCause(&req.Params)
	if err != nil {
		return fmt.Errorf("unable to get cause: %s", err)
	}
	params.Set("cause", cause)

	if req.Source.Token != "" {
		params.Set("token", req.Source.Token)
	} else {
		return fmt.Errorf("no token supplied to source")
	}

	if err = jenkins.Build(job, params); err != nil {
		return fmt.Errorf("unable to trigger build payload: %s", err)
	}

	for {
		if build, err := jenkins.GetBuild(job, job.LastCompletedBuild.Number+1); err == nil {
			// TODO: do we care if there is an error here?
			description, _ := r.getDescription(&req.Params)
			_ = jenkins.SetBuildDescription(build, description)

			resp.Version = r.getVersion(build)
			resp.Metadata = r.getMetadata(&build)

			if err = r.acceptResult(&build, req.Params.AcceptResults); err != nil {
				return fmt.Errorf("unaccepted result: %s", err)
			}

			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	return r.writeOutput(resp)
}
