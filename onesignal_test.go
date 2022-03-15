package onesignal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func setup(t *testing.T) (*httptest.Server, *http.ServeMux, *Client) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client := setupClient(t)
	client.SetBaseURL(server.URL)
	return server, mux, client
}

func setupClient(t *testing.T) *Client {

	c, err := NewClient("fake-app-id", "mock-api-key")
	if err != nil {
		t.Fatal(err)
	}

	return c
}

func setupUserClient(t *testing.T) (*httptest.Server, *http.ServeMux, *UserClient) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client, err := NewUserClient("mock-user-key")
	if err != nil {
		t.Fatal(err)
	}
	client.SetBaseURL(server.URL)

	return server, mux, client
}

func teardown(server *httptest.Server) {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	if got := r.Header.Get(header); got != want {
		t.Errorf("NewRequest() %s header is %v, want %v", header, got, want)
	}
}

func testBody(t *testing.T, r *http.Request, body interface{}, want interface{}) {
	json.NewDecoder(r.Body).Decode(body)
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Request body: %+v, want %+v", body, want)
	}
}

func TestNewClient(t *testing.T) {
	_, err := NewClient("", "")
	if err == nil {
		t.Error("expected error, not nil")
	}

	c := setupClient(t)

	if got, want := c.baseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}

	if got, want := c.client, http.DefaultClient; got != want {
		t.Errorf("NewClient Client is %v, want %v", got, want)
	}

	if got, want := c.Players.client, c; got != want {
		t.Errorf("NewClient.PlayersService.client is %v, want %v", got, want)
	}

	if got, want := c.Notifications.client, c; got != want {
		t.Errorf("NewClient.NotificationsService.client is %v, want %v", got, want)
	}
}

func TestCustomHTTPClient(t *testing.T) {
	httpClient := &http.Client{}

	c := setupClient(t)
	c.SetHTTPClient(httpClient)

	if got, want := c.client, httpClient; got != want {
		t.Errorf("NewClient Client is %v, want %v", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	apiKey := "mock-api-key"
	c := setupClient(t)

	method := "GET"
	inURL, outURL := "foo", defaultBaseURL+"foo"
	inBody := struct{ Foo string }{Foo: "Bar"}
	outBody := `{"Foo":"Bar"}` + "\n"
	req, _ := c.NewRequest(method, inURL, inBody)

	// test the HTTP method
	if got, want := req.Method, method; got != want {
		t.Errorf("NewRequest(%q) Method is %v, want %v", method, got, want)
	}

	// test the URL
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	testHeader(t, req, "Content-Type", "application/json")
	testHeader(t, req, "Accept", "application/json")
	testHeader(t, req, "Authorization", fmt.Sprintf("Basic %s", apiKey))
}

func TestNewRequest_userKeyType(t *testing.T) {
	c := setupClient(t)

	req, _ := c.NewRequest("GET", "foo", nil)

	testHeader(t, req, "Authorization", "Basic mock-api-key")
}

func TestNewRequest_emptyBody(t *testing.T) {
	c := setupClient(t)

	req, err := c.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if req.Body != nil {
		t.Fatalf("Request contains a non-nil Body: %v", req.Body)
	}
}

func TestDo(t *testing.T) {
	server, mux, client := setup(t)
	defer teardown(server)

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)
	client.Do(req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_httpError(t *testing.T) {
	server, mux, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	_, ok := err.(*ErrorResponse)
	if ok {
		t.Errorf("Error should be `couldn't decode response body JSON` but got %v: %+v", reflect.TypeOf(err), err)
	}
}

func TestCheckResponse_ok(t *testing.T) {
	r := &http.Response{
		StatusCode: http.StatusOK,
	}

	err := checkErrorResponse(r)
	if err != nil {
		t.Fatalf("checkErrorResponse shouldn't return an error, but returned: %+v", err)
	}
}

func TestCheckResponse_badRequest(t *testing.T) {
	r := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(`{
			"errors":
			[
				"Invalid or missing authentication token"
			]
		}`)),
	}

	err, ok := checkErrorResponse(r).(*ErrorResponse)
	if !ok {
		t.Errorf("checkErrorResponse return value should be of type ErrorResponse but is %v: %+v", reflect.TypeOf(err), err)
	}

	if err == nil {
		t.Fatalf("checkErrorResponse should return an error")
	}

	if len(err.Messages) == 0 {
		t.Fatalf("checkErrorResponse ErrorResponse should contain messages")
	}

	want := "Invalid or missing authentication token"
	if got := err.Messages[0]; want != got {
		t.Errorf("Error message: %v, want %v", got, want)
	}
}

func TestCheckResponse_noBody(t *testing.T) {
	r := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}

	err := checkErrorResponse(r)
	if err == nil {
		t.Fatalf("checkErrorResponse should return an error")
	}

	if err.Error() != "couldn't decode response body JSON: EOF" {
		t.Errorf("checkErrorResponse return value should be of type ErrorResponse but is %v: %+v", reflect.TypeOf(err), err)
	}

}
