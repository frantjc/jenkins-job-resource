package resource

import (
	"github.com/yosida95/golang-jenkins"
)

type CheckRequest struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version"`
}

type Version struct {
	Number int    `json:"number,string"`
	URL    string `json:"url"`
}

func ToBuild(v *Version) gojenkins.Build {
	return gojenkins.Build{
		Number: v.Number,
		Url: v.URL,
	}
}

func ToVersion(b *gojenkins.Build) Version {
	return Version{
		Number: b.Number,
		URL: b.Url,
	}
}

type CheckResponse []Version

type InRequest struct {
	Source  Source    `json:"source"`
	Params  GetParams `json:"params"`
	Version Version   `json:"version"`
}

type InResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type OutRequest struct {
	Source Source    `json:"source"`
	Params PutParams `json:"params"`
}

type OutResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type Source struct {
	URL string `json:"url"`
	Job string `json:"job"`
	Username string `json:"username,omitempty"`
	APIToken string `json:"api_token,omitempty"`
	Token string `json:"token"`
}

type GetParams struct{}

type PutParams struct{
	Cause       string      `json:"cause,omitempty"`
	BuildParams interface{} `json:"build_params,omitempty"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
