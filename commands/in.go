package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	resource "github.com/logsquaredn/jenkins-job-resource"
)

// In runs the in script which checks stdin for a JSON object of the form of an InRequest
// fetches and writes the requested Version as well as Metadata about it to stdout and
// writes each of its output Artifacts to src/Artifact.FileName
func (j *JenkinsJobResource) In() error {
	var (
		req resource.InRequest
		resp resource.InResponse
	)

	err := j.readInput(req)
	if err != nil {
		return err
	}

	src, err := j.getSrc()
	if err != nil {
		return err
	}

	err = os.MkdirAll(src, 0755)
	if err != nil {
		return fmt.Errorf("unable to make directory %s", src)
	}

	jenkins := j.newJenkins(req.Source)

	job, err := jenkins.GetJob(req.Source.Job)
	if err != nil {
		return fmt.Errorf("unable to find job %s: %s", req.Source.Job, err)
	}

	build, err := jenkins.GetBuild(job, req.Version.Number)
	if err != nil {
		return fmt.Errorf("requested version of resource %d unavailable: %s", req.Version.Number, err)
	}

	resp.Version = j.getVersion(&build)
	resp.Metadata = j.getMetadata(&build)

	for _, artifact := range build.Artifacts {
		data, err := jenkins.GetArtifact(build, artifact)
		if err == nil {
			ioutil.WriteFile(filepath.Join(src, artifact.FileName), data, 0644)
		}
	}

	j.writeOutput(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}
