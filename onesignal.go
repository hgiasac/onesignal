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

// AuthKeyType specifies the token used to authenticate the requests
// https://documentation.onesignal.com/docs/accounts-and-keys
type AuthKeyType uint

const (
	APP AuthKeyType = iota
	USER
)

// ClientOptions represent required OneSignal client options
type ClientOptions struct {
	BaseURL string
	// Private key used for most API calls like sending push notifications and updating users.
	// https://documentation.onesignal.com/docs/accounts-and-keys#rest-api-key
	ApiKey string
	// Another type of REST API key used for viewing Apps and related updates.
	UserKey string
	Client  *http.Client
	IsDebug bool
	Logger  func(...interface{})
}

// A Client manages communication with the OneSignal API.
type Client struct {
	baseURL *url.URL
	apiKey  string
	userKey string
	client  *http.Client
	logger  func(...interface{})

	Apps          *AppsService
	Players       *PlayersService
	Notifications *NotificationsService
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

// New returns a new OneSignal API client.
func New(options ClientOptions) (*Client, error) {

	if options.ApiKey == "" && options.UserKey == "" {
		return nil, errors.New("require ApiKey or UserKey")
	}

	sBaseUrl := options.BaseURL
	if sBaseUrl == "" {
		sBaseUrl = defaultBaseURL
	}

	baseURL, err := url.Parse(sBaseUrl)
	if err != nil {
		return nil, err
	}

	httpClient := options.Client
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		baseURL: baseURL,
		client:  httpClient,
		apiKey:  options.ApiKey,
		userKey: options.UserKey,
		logger:  options.Logger,
	}

	c.Apps = &AppsService{client: c}
	c.Players = &PlayersService{client: c}
	c.Notifications = &NotificationsService{client: c}

	return c, err
}

// NewRequest creates an API request.
// path is a relative URL, like "/apps".
// The value pointed to by body is JSON encoded and included as the request body.
// The AuthKeyType will determine which authorization token (APP or USER) is
// used for the request.
func (c *Client) NewRequest(method, path string, body interface{}, authKeyType AuthKeyType) (*http.Request, error) {
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

	var token string
	if authKeyType == APP {
		token = fmt.Sprintf("Basic %s", c.apiKey)
	} else {
		token = fmt.Sprintf("Basic %s", c.userKey)
	}
	c.printDebug("[OneSignal] Authorization:", token)
	req.Header.Add("Authorization", token)

	return req, nil
}

// Sends an API request and returns the API response.
// Return JSON decoded and stored in the value pointed to by v,
// or an error if an API error has occurred.
func (c *Client) Do(r *http.Request, v interface{}) (*http.Response, error) {
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

func (c *Client) printDebug(args ...interface{}) {
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
