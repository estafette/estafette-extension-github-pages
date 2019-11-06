package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sethgrid/pester"
)

// GithubAPIClient allows to communicate with the Github api
type GithubAPIClient interface {
	RequestPageBuild(repoOwner, repoName string) (err error)
}

type githubAPIClientImpl struct {
	accessToken string
}

func newGithubAPIClient(accessToken string) GithubAPIClient {
	return &githubAPIClientImpl{
		accessToken: accessToken,
	}
}

func (gh *githubAPIClientImpl) RequestPageBuild(repoOwner, repoName string) (err error) {

	// https://developer.github.com/v3/repos/pages/#request-a-page-build
	log.Printf("\nRequesting github pages build...")

	_, err = gh.callGithubAPI("POST", fmt.Sprintf("https://api.github.com/repos/%v/%v/pages/builds", repoOwner, repoName), []int{http.StatusOK,http.StatusCreated}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (gh *githubAPIClientImpl) callGithubAPI(method, url string, validStatusCodes []int, params interface{}) (body []byte, err error) {

	// convert params to json if they're present
	var requestBody io.Reader
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return body, err
		}
		requestBody = bytes.NewReader(data)
	}

	// create client, in order to add headers
	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	request, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return
	}

	// add headers
	request.Header.Add("Authorization", fmt.Sprintf("%v %v", "token", gh.accessToken))
	request.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	hasValidStatusCode := false
	for _, sc := range validStatusCodes {
		if response.StatusCode == sc {
			hasValidStatusCode = true
		}
	}
	if !hasValidStatusCode {
		return body, fmt.Errorf("Status code %v for '%v %v' is not one of the valid status codes %v for this request. Body: %v", response.StatusCode, method, url, validStatusCodes, string(body))
	}

	if string(body) == "" {
		log.Printf("Received successful response without body for '%v %v' with status code %v", method, url, response.StatusCode)
		return
	}

	// unmarshal json body
	var b interface{}
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Printf("Deserializing response for '%v' Github api call failed. Body: %v. Error: %v", url, string(body), err)
		return
	}

	return
}
