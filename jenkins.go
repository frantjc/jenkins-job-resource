package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
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
	Description           string         `json:"description"`
	DisplayName           string         `json:"displayName"`
	DisplayNameOrNull     string         `json:"displayNameOrNull,omitempty"`
	FullDisplayName       string         `json:"fullDisplayName"`
	FullName              string         `json:"fullName"`
	Name                  string         `json:"name"`
	URL                   string         `json:"url"`
	Buildable             bool           `json:"buildable"`
	Builds                []Build        `json:"builds"`
	Color                 string         `json:"color"`
	FirstBuild            Build          `json:"firstBuild,omitempty"`
	HealthReport          []HealthReport `json:"healthReport,omitempty"`
	InQueue               bool           `json:"inQueue"`
	KeepDependencies      bool           `json:"keepDependencies"`
	LastBuild             Build          `json:"lastBuild,omitempty"`
	LastCompletedBuild    Build          `json:"lastCompletedBuild,omitempty"`
	LastFailedBuild       Build          `json:"lastFailedBuild,omitempty"`
	LastStableBuild       Build          `json:"lastStableBuild,omitempty"`
	LastSuccessfulBuild   Build          `json:"lastSuccessfulBuild,omitempty"`
	LastUnstableBuild     Build          `json:"lastUnstableBuild,omitempty"`
	LastUnsuccessfulBuild Build          `json:"lastUnsuccessfulBuild,omitempty"`
	NextBuildNumber       int            `json:"nextBuildNumber"`
	ConcurrentBuild       bool           `json:"concurrentBuild"`
	ResumeBlocked         bool           `json:"resumeBlocked"`
}

type Build struct {
	Number  int    `json:"number"`
	URL     string `json:"url"`
	BasicCredentials
}

type HealthReport struct {
	Description   string `json:"description"`
	IconClassName string `json:"iconClassName"`
	IconURL       string `json:"iconUrl"`
	Score         int    `json:"score"`
}

func (j *Job) Build(token string, cause string, params *interface{}) (*Build, error) {
	info, err := j.GetInfo()

	url := j.URL
	if params != nil {
		url += "buildWithParameters"
	} else {
		url += "build"
	}

	qs, err := query.Values(*params)
	if err != nil {
		return nil, err
	}

	if token != "" {
		qs.Set("token", token)
	}

	if cause != "" {
		qs.Set("cause", cause)
	}

	url += fmt.Sprintf("?%s", qs.Encode())

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	c := http.Client{}
	_, err = c.Do(req)
	if err != nil {
		return nil, err
	}

	var resp Build

	for newInfo, err := j.GetInfo(); err != nil; newInfo, err = j.GetInfo() {
		if newInfo.LastBuild.Number != info.LastBuild.Number {
			resp = newInfo.LastBuild
			resp.BasicCredentials = j.BasicCredentials
			break;
		}
		
		time.Sleep(5 * time.Second)
	}

	return &resp, nil
}

func (j *Job) GetInfo() (*JobResponse, error) {
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

	return &jResp, nil
}

func (j *Job) GetBuilds() ([]Build, error) {
	info, err := j.GetInfo()
	if err != nil {
		return nil, err
	}

	for _, build := range info.Builds {
		build.BasicCredentials = j.BasicCredentials
	}

	return info.Builds, nil
}

type BuildResponse struct {
	// Artifacts `json:"artifacts"` idk what this looks like
	Building bool `json:"building"`
	Description string `json:"Description"`
	DisplayName string `json:"displayName"`
	Duration int `json:"duration"`
	EstimatedDuration int `json:"estimatedDuration"`
	// Executor `json:"executor"` idk what this looks like
	FullDisplayName string `json:"fullDisplayName"`
	ID string `json:"id"`
	KeepLog bool `json:"keepLog"`
	Number int `json:"number"`
	QueueID int `json:"queueId"`
	Result string `json:"result"`
	Timestamp int `json:"timestamp"`
	URL string `json:"url"`
	// ChangeSets `json:"changeSets"` idk what this looks like
	// Culprits `json:"culprits"` idk what this looks like
	NextBuild Build `json:"nextBuild"`
	PreviousBuild Build `json:"previousBuild"`
}

func (b *Build) GetInfo() (*BuildResponse, error) {
	req, err := http.NewRequest("GET", b.URL + "api/json/", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(b.BasicCredentials.Username, b.BasicCredentials.Password)

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var bResp BuildResponse
	err = json.NewDecoder(resp.Body).Decode(&bResp)
	if err != nil {
		return nil, err
	}

	return &bResp, nil
}
