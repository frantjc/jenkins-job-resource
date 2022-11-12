package command

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	resource "github.com/frantjc/jenkins-job-resource"
)

// In runs the in script which checks stdin for a JSON object of the form of an InRequest
// fetches and writes the requested Version as well as Metadata about it to stdout and
// writes each of its output Artifacts to src/Artifact.FileName.
func (r *JenkinsJobResource) In() error {
	var (
		req  resource.InRequest
		resp resource.InResponse
	)

	err := r.readInput(&req)
	if err != nil {
		return err
	}

	src, err := r.getSrc()
	if err != nil {
		return err
	}

	err = os.MkdirAll(src, 0755)
	if err != nil {
		return fmt.Errorf("unable to make directory %s", src)
	}

	jenkins := r.newJenkins(req.Source)

	job, err := jenkins.GetJob(req.Source.Job)
	if err != nil {
		return fmt.Errorf("unable to find job %s: %s", req.Source.Job, err)
	}

	build, err := jenkins.GetBuild(job, req.Version.Build)
	if err != nil {
		return fmt.Errorf("requested version of resource %d unavailable: %s", req.Version.Build, err)
	}

	resp.Version = r.getVersion(build)
	resp.Metadata = r.getMetadata(&build)

	err = r.acceptResult(&build, req.Params.AcceptResults)
	if err != nil {
		return fmt.Errorf("unaccepted result: %s", err)
	}

	if !req.Params.SkipDownload {
		for _, artifact := range build.Artifacts {
			if !strings.HasPrefix(artifact.FileName, ".") {
				match := false
				if req.Params.Regexp != nil {
					for _, expr := range req.Params.Regexp {
						if re, err := regexp.Compile(expr); err == nil && re.Match([]byte(artifact.FileName)) {
							match = true
							break
						}
					}
				} else {
					match = true
				}

				if match {
					data, err := jenkins.GetArtifact(build, artifact)
					if err == nil {
						if err = os.WriteFile(filepath.Join(src, artifact.FileName), data, 0600); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	if err = r.writeMetadata(resp.Metadata); err != nil {
		return err
	}

	return r.writeOutput(resp)
}
