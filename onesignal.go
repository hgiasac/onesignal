package onesignal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL = "https://onesignal.com/api/v1"
)

// Client manages communication with the OneSignal application API.
type Client struct {
	*httpClient
	appID string

	Players       *PlayersService
	Notifications *NotificationsService
}

// UserClient manages OneSignal applications.
type UserClient struct {
	*httpClient

	Apps *AppsService
}

// NewUserClient returns a UserClient
func NewUserClient(userKey string) (*UserClient, error) {

	if userKey == "" {
		return nil, errors.New("user auth key is required")
	}

	c := &UserClient{
		httpClient: newHTTPClient(userKey),
	}

	c.Apps = &AppsService{client: c}

	return c, nil
}

// SuccessResponse wraps the standard http.Response for several API methods
// that just return a Success flag.
type SuccessResponse struct {
	Success bool `json:"success"`
}

// ErrorResponse reports one or more errors caused by an API request.
type ErrorResponse struct {
	Messages []string `json:"errors"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("OneSignal errors:\n - %s", strings.Join(e.Messages, "\n - "))
}

// NewClient returns a new OneSignal API client.
func NewClient(appID string, apiKey string) (*Client, error) {

	if appID == "" {
		return nil, errors.New("app ID is required")
	}

	if apiKey == "" {
		return nil, errors.New("api key is required")
	}

	c := &Client{
		appID:      appID,
		httpClient: newHTTPClient(apiKey),
	}

	c.Players = &PlayersService{client: c}
	c.Notifications = &NotificationsService{client: c}

	return c, nil
}

// GetAppID returns the application ID
func (c *Client) GetAppID() string {
	return c.appID
}

type httpClient struct {
	baseURL *url.URL
	apiKey  string
	client  *http.Client
	logger  func(...interface{})
}

func newHTTPClient(apiKey string) *httpClient {
	baseURL, _ := url.Parse(defaultBaseURL)
	return &httpClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  http.DefaultClient,
	}
}

// SetBaseURL change the default OneSignal base URL
func (c *httpClient) SetBaseURL(baseURL string) error {
	sBaseURL, err := url.Parse(baseURL)
	if err != nil {
		panic(fmt.Sprintf("incorrect base url format: %s", baseURL))
	}

	c.baseURL = sBaseURL
	return nil
}

// SetHTTPClient set custom http client
func (c *httpClient) SetHTTPClient(client *http.Client) {
	c.client = client
}

// SetHTTPClient set custom debug logger
func (c *httpClient) SetLogger(logger func(args ...interface{})) {
	c.logger = logger
}

// NewRequest creates an API request.
// path is a relative URL, like "/apps".
// The value pointed to by body is JSON encoded and included as the request body.
// The AuthKeyType will determine which authorization token (APP or USER) is
// used for the request.
func (c *httpClient) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(c.baseURL.String() + path)
	if err != nil {
		return nil, err
	}

	c.printDebug(fmt.Sprintf("[OneSignal] requesting url: %s", u.String()))

	var buf io.ReadWriter
	if body != nil {
		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(body)
		if err != nil {
			return nil, err
		}
		buf = b

		if c.logger != nil {
			c.logger("[OneSignal] request body: " + b.String())
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// set header and access token
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	token := fmt.Sprintf("Basic %s", c.apiKey)
	c.printDebug("[OneSignal] Authorization:", token)
	req.Header.Add("Authorization", token)

	return req, nil
}

// Sends an API request and returns the API response.
// Return JSON decoded and stored in the value pointed to by v,
// or an error if an API error has occurred.
func (c *httpClient) Do(r *http.Request, v interface{}) (*http.Response, error) {
	// send the request
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkErrorResponse(resp)
	if err != nil {
		return resp, err
	}

	if c.logger != nil {
		var b bytes.Buffer
		b.ReadFrom(resp.Body)
		c.printDebug("response body: ", b.String())
		err = json.Unmarshal(b.Bytes(), &v)
	} else {
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&v)
	}

	// it returns EOF if http status 204 (no content available)
	if err != nil && err != io.EOF {
		return resp, err
	}
	return resp, nil
}

func (c *httpClient) printDebug(args ...interface{}) {
	if c.logger != nil {
		c.logger(args...)
	}
}

// checkErrorResponse checks the API response for errors, by http status code
// and returns them if present
func checkErrorResponse(r *http.Response) error {
	switch r.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusInternalServerError:
		return errors.New("internal server error")
	default:
		errResp := new(ErrorResponse)
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&errResp)
		if err != nil {
			return fmt.Errorf("couldn't decode response body JSON: %v", err)
		}
		return errResp
	}
}
