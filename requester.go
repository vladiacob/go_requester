package go_requester

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// ServiceInterface for that service
type ServiceInterface interface {
	Make(string, string, interface{}, interface{}) (*Response, error)
	SetAuthentication(string, string)
}

type Response struct {
	Status int
	Body   []byte
}

// Requester is structure for that package
type Requester struct {
	client    *http.Client
	userAgent string
	username  string
	password  string
}

// New return a Requester structure
func New(client *http.Client) *Requester {
	return &Requester{
		client: client,
	}
}

// SerUserAgent set user agent for the request
func (r *Requester) SerUserAgent(userAgent string) {
	r.userAgent = userAgent
}

// SetAuthentication is adding HTTP basic authentication for all the requests which
// are making with that requester
func (r *Requester) SetAuthentication(username string, password string) {
	r.username = username
	r.password = password
}

// Make method is using to do http requests,
// it will return error and Response struct
func (r *Requester) Make(method string, URL string, params interface{}, responseStruct interface{}) (*Response, error) {
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return &Response{}, errors.New(err.Error())
	}

	request, err := http.NewRequest(method, URL, bytes.NewBuffer(paramsBytes))
	if err != nil {
		return &Response{}, errors.New(err.Error())
	}

	if r.username != "" && r.password != "" {
		request.SetBasicAuth(r.username, r.password)
	}

	if r.userAgent != "" {
		request.Header.Set("User-Agent", r.userAgent)
	}

	response, err := r.client.Do(request)
	if err != nil {
		return &Response{}, errors.New(err.Error())
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &Response{}, errors.New(err.Error())
	}

	switch responseStruct.(type) {
	case *string:
		*responseStruct.(*string) = string(contents)
	default:
		err := json.Unmarshal(contents, responseStruct)
		if err != nil {
			return &Response{Status: response.StatusCode, Body: contents}, errors.New(err.Error())

		}
	}

	return &Response{Status: response.StatusCode, Body: contents}, nil
}
