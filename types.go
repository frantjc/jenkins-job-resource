package resource

import (
	"strconv"
)

type CheckRequest struct {
	Source  Source `json:"source"`
	Version *Version `json:"version"`
}


type Version struct {
	Number string `json:"number,string"`
	URL    string `json:"url"`
}

func (v *Version) ToBuild() Build {
	n, err := strconv.Atoi(v.Number)
	if err != nil {
		return Build{
			Number: 0,
			URL: v.URL,
		}
	}

	return Build{
		Number: n,
		URL: v.URL,
	}
}

func (b *Build) ToVersion() Version {
	return Version{
		Number: strconv.Itoa(b.Number),
		URL: b.URL,
	}
}

type CheckResponse []Version

type InRequest struct {
	Source  Source    `json:"source"`
	Params  GetParams `json:"params"`
	Version Version   `json:"version"`
}

type InResponse struct {
	Version  Version      `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type OutRequest struct {
	Source Source    `json:"source"`
	Params PutParams `json:"params"`
}

type OutResponse struct {
	Version  Version      `json:"version"`
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
