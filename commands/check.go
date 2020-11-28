package commands

import (
	"fmt"

	resource "github.com/logsquaredn/jenkins-job-resource"
)

// Check runs the in script which checks stdin for a JSON object of the form of a CheckRequest
// fetches and writes the all Versions that are newer than the provided Version to stdout
func (r *JenkinsJobResource) Check() error {
	var (
		req resource.CheckRequest
	    resp resource.CheckResponse
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

	builds := job.Builds

	if len(builds) > 0 {
		if req.Version != nil {
			for _, build := range builds {
				version := r.getVersion(&build)
				if version.Build >= req.Version.Build {
					// prepend
					resp = append([]resource.Version{version}, resp...)
				}
			}
		}

		if len(resp) <= 0 {
			resp = append(resp, r.getVersion(&builds[0]))
		}
	}

	r.writeOutput(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}
