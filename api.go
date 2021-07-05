package pvoutput

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	apiBaseURL                = "https://pvoutput.org/service/r2/"
	apiAddOutputEndpoint      = "addoutput.jsp"
	apiAddBatchOutputEndpoint = "addbatchoutput.jsp"
	apiAddStatusEndpoint      = "addstatus.jsp"
	apiAddBatchStatusEndpoint = "addbatchstatus.jsp"
)

// API is a struct holding relevant session data
type API struct {
	Key      string
	SystemID string
	client   http.Client
}

// NewAPI returns a new API object for given systemID and API key
func NewAPI(key, systemID string) API {
	return API{
		SystemID: systemID,
		Key:      key,
		client:   http.Client{},
	}
}

func (a API) getPOSTRequest(path string, enc PVEncodable) (*http.Request, error) {
	req, err := a.getRequest(http.MethodPost, path)
	if err != nil {
		return nil, err
	}

	data, err := enc.Encode()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", len(data)))
	req.Body = ioutil.NopCloser(strings.NewReader(data))

	return req, nil
}

func (a API) getRequest(method, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", apiBaseURL, path), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Pvoutput-Apikey", a.Key)
	req.Header.Add("X-Pvoutput-SystemId", a.SystemID)

	return req, nil
}

func (a API) handleRequest(req *http.Request) error {
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	return nil
}

// AddOutput implements PVOutput's /addoutput.jsp service
func (a API) AddOutput(o Output) error {
	req, err := a.getPOSTRequest(apiAddOutputEndpoint, o)
	if err != nil {
		return err
	}

	return a.handleRequest(req)
}

// AddBatchOutput implements PVOutput's /addbatchoutput.jsp service
func (a API) AddBatchOutput(b BatchOutput) error {
	req, err := a.getPOSTRequest(apiAddBatchOutputEndpoint, b)
	if err != nil {
		return err
	}

	return a.handleRequest(req)
}

// AddStatus implements PVOutput's /addstatus.jsp service
func (a API) AddStatus(s Status) error {
	req, err := a.getPOSTRequest(apiAddStatusEndpoint, s)
	if err != nil {
		return err
	}

	return a.handleRequest(req)
}

// AddBatchStatus implements PVOutput's /addbatchstatus.jsp service
func (a API) AddBatchStatus(b BatchStatus) error {
	req, err := a.getPOSTRequest(apiAddBatchStatusEndpoint, b)
	if err != nil {
		return err
	}

	return a.handleRequest(req)
}
