package resource

import (
	"github.com/yosida95/golang-jenkins"
)

// CheckRequest is the JSON object that Concourse passes to /opt/resource/check through stdin
type CheckRequest struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version"`
}

// Version is the JSON object that is passed to and from Concourse
type Version struct {
	Number int    `json:"number,string"`
	URL    string `json:"url"`
}

// ToBuild converts Version to Build to transfer information from Concourse to Jenkins
func ToBuild(v *Version) gojenkins.Build {
	return gojenkins.Build{
		Number: v.Number,
		Url: v.URL,
	}
}

// ToVersion converts Build to Version to transfer infromation from Jenkins to Concourse
func ToVersion(b *gojenkins.Build) Version {
	return Version{
		Number: b.Number,
		URL: b.Url,
	}
}

// CheckResponse is the JSON object that we pass back to Concourse through stdout from /opt/resource/check
type CheckResponse []Version

// InRequest is the JSON object that Concourse passes to /opt/resource/in through stdin
type InRequest struct {
	Source  Source    `json:"source"`
	Params  GetParams `json:"params"`
	Version Version   `json:"version"`
}

// InResponse is the JSON object that we pass back to Concourse through stdout from /opt/resource/in
type InResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

// OutRequest is the JSON object that Concourse passes to /opt/resource/out through stdin
type OutRequest struct {
	Source Source    `json:"source"`
	Params PutParams `json:"params"`
}

// OutResponse is the JSON object that we pass back to Concourse through stdout from /opt/resource/out
type OutResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

// Source is the JSON (yaml) object configured under the resources array in a Concourse pipeline
type Source struct {
	URL string `json:"url"`
	Job string `json:"job"`
	Username string `json:"username,omitempty"`
	APIToken string `json:"api_token,omitempty"`
	Token string `json:"token"`
}

// GetParams are additional parameters that can be passed to this Concourse Resource Type during a get step
type GetParams struct{}

// PutParams are additional parameters that can be passed to this Concourse Resource Type during a put step
type PutParams struct{
	Cause       string      `json:"cause,omitempty"`
	BuildParams interface{} `json:"build_params,omitempty"`
	Description string      `json:"description,omitempty"`
}

// Metadata is the object which is passed in array form to Concourse through stdout from /opt/resource/out and /opt/resource/in
// to provide additional information about the Version
type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
