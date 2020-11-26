package resource

import (
	"fmt"
)

type CheckRequest struct {
	Source  Source `json:"source"`
	Version *Build `json:"version"`
}


type Version struct {
	Number string `json:"number,string"`
	URL    string `json:"url"`
}

func (b *Build) ToVersion() Version {
	return Version{
		Number: fmt.Sprint(b.Number),
		URL: b.URL,
	}
}

type CheckResponse []Version

type InRequest struct {
	Source  Source    `json:"source"`
	Params  GetParams `json:"params"`
	Version Build     `json:"version"`
}

type InResponse struct {
	Version  Build      `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type OutRequest struct {
	Source Source    `json:"source"`
	Params PutParams `json:"params"`
}

type OutResponse struct {
	Version  Build      `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type Source struct {
	URL string `json:"url"`
	Job string `json:"job"`

	BasicCredentials
	Token string `json:"token"`
}

type BasicCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type GetParams struct{}

type PutParams struct{}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
