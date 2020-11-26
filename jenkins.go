package resource

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func NewJenkins(i *JenkinsInput) (*Jenkins, error) {
	if !strings.HasSuffix(i.URL, "/") {
		i.URL = i.URL + "/"
	}

	j := Jenkins{
		URL: i.URL,
		BasicCredentials: BasicCredentials{
			Username: i.Username,
			Password: i.Password,
		},
	}

	return &j, nil
}

func (j *Jenkins) getJobs() ([]Job, error) {
	req, err := http.NewRequest("GET", j.URL + "api/json/", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(j.BasicCredentials.Username, j.BasicCredentials.Password)

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var jResp JenkinsResponse
	err = json.NewDecoder(resp.Body).Decode(&jResp)
	if err != nil {
		return nil, err
	}

	if jResp.Jobs == nil {
		return nil, errors.New("no jobs found")
	}

	return jResp.Jobs, nil
}

func (j *Jenkins) GetJob(name string) (*Job, error) {
	jobs, err := j.getJobs()
	if err != nil {
		return nil, err
	}

	for _, job := range jobs {
		if strings.EqualFold(job.Name, name) {
			job.BasicCredentials = j.BasicCredentials
			return &job, nil
		}
	}

	return nil, errors.New("no job found matching the given name")
}

type Jenkins struct {
	URL string
	BasicCredentials
}

type JenkinsResponse struct {
	Jobs []Job `json:"jobs"`
}

type JenkinsInput struct {
	URL string
	BasicCredentials
}

type Job struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Color string `json:"color"`
	BasicCredentials
}

type JobResponse struct {
	Description          string `json:"description"`
	DisplayName          string `json:"displayName"`
	DisplayNameOrNull    string `json:"displayNameOrNull"`
	FullDisplayName      string `json:"fullDisplayName"`
	FullName             string `json:"fullName"`
	Name                 string `json:"name"`
	URL                  string `json:"url"`
	Buildable              bool `json:"buildable"`
	Builds              []Build `json:"builds"`
	Color                string `json:"color"`
	FirstBuild            Build `json:"firstBuild"`
	HealthReport []HealthReport `json:"healthReport"`
	InQueue                bool `json:"inQueue"`
	KeepDependencies       bool `json:"keepDependencies"`
	LastBuild             Build `json:"lastBuild"`
	LastCompletedBuild    Build `json:"lastCompletedBuild"`
	LastFailedBuild       Build `json:"lastFailedBuild"`
	LastStableBuild       Build `json:"lastStableBuild"`
	LastSuccessfulBuild   Build `json:"lastSuccessfulBuild"`
	LastUnstableBuild     Build `json:"lastUnstableBuild"`
	LastUnsuccessfulBuild Build `json:"lastUnsuccessfulBuild"`
	NextBuildNumber       int64 `json:"nextBuildNumber"`
	ConcurrentBuild        bool `json:"concurrentBuild"`
	ResumeBlocked          bool `json:"resumeBlocked"`
}

type Build struct {
	Number int64 `json:"number"`
	URL   string `json:"url"`
}

type HealthReport struct {
	Description   string `json:"description"`
	IconClassName string `json:"iconClassName"`
	IconURL       string `json:"iconUrl"`
	Score          int64 `json:"score"`
}

func (j *Job) Build(token string, cause string, params *interface{}) error {
	var body io.Reader
	url := j.URL
	if params != nil {
		buf, err := json.Marshal(params)
		if err != nil {
			return err
		}

		body = bytes.NewBuffer(buf)
		url += "buildWithParameters"
	} else {
		body = nil
		url += "build"
	}

	url += fmt.Sprintf("?token=%s", token)

	if cause != "" {
		url += fmt.Sprintf("&cause=%s", cause)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	c := http.Client{}
	_, err = c.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (j *Job) GetBuilds() ([]Build, error) {
	req, err := http.NewRequest("GET", j.URL + "api/json/", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(j.BasicCredentials.Username, j.BasicCredentials.Password)

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var jResp JobResponse
	err = json.NewDecoder(resp.Body).Decode(&jResp)
	if err != nil {
		return nil, err
	}

	return jResp.Builds, nil
}
