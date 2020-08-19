package pvoutput

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	apiBaseURL                = "https://pvoutput.org/service/r2/"
	apiAddOutputEndpoint      = "addoutput.jsp"
	apiAddBatchOutputEndpoint = "addbatchoutput.jsp"
	apiAddStatusEndpoint      = "addstatus.jsp"
	apiAddBatchStatusEndpoint = "addbatchstatus.jsp"
)

var (
	ErrRateExceeded = errors.New("rate limit exceeded")
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

// AddOutput implements PVOutput's /addoutput.jsp service
// TODO: this function and equivalent should be refactored
// instead of duplicated
func (a API) AddOutput(o Output) error {
	req, err := a.getPOSTRequest(apiAddOutputEndpoint, o)
	if err != nil {
		return err
	}

	resp, err := a.client.Do(req)
	if err != nil {

	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != "OK 200: Added Output" {
		return errors.New(string(body))
	}

	return nil
}

// AddBatchOutput implements PVOutput's /addbatchoutput.jsp service
func (a API) AddBatchOutput(b BatchOutput) error {
	req, err := a.getPOSTRequest(apiAddBatchOutputEndpoint, b)
	if err != nil {
		return err
	}

	resp, err := a.client.Do(req)
	if err != nil {

	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != "OK 200: Added Batch" {
		return errors.New(string(body))
	}

	return nil
}

// AddStatus implements PVOutput's /addstatus.jsp service
func (a API) AddStatus(s Status) error {
	req, err := a.getPOSTRequest(apiAddStatusEndpoint, s)
	if err != nil {
		return err
	}

	resp, err := a.client.Do(req)
	if err != nil {

	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != "OK 200: Added Status" {
		return errors.New(string(body))
	}

	return nil
}

// AddBatchStatus implements PVOutput's /addbatchstatus.jsp service
func (a API) AddBatchStatus(b BatchStatus) error {
	req, err := a.getPOSTRequest(apiAddBatchStatusEndpoint, b)
	if err != nil {
		return err
	}

	resp, err := a.client.Do(req)
	if err != nil {

	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if match, err := regexp.MatchString(`^(\d+,[\d:]+,[01];?)+$`, string(body)); !match || err != nil {
		return errors.New(string(body))
	}

	return nil
}
