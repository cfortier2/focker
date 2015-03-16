package focker

import (
	// "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// focker - short for Fake Docker will mimic a bare bones Docker API for testing

type Focker struct {
	Mux    *http.ServeMux
	Server *httptest.Server
}

// CreateFocker will create a test Docker API server
func CreateFocker() *Focker {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return &Focker{
		Mux:    mux,
		Server: server,
	}
}

// Close will close the underlying http server
func (f *Focker) Close() {
	f.Server.Close()
}

// Address returns the server address
func (f *Focker) Address() string {
	url, _ := url.Parse(f.Server.URL)
	return url.String()
}

// MockEndpoint contains all of the information necassary to mock an endpoint.
type MockEndpoint struct {
	// Endpoint is the endpoint to mock. ex: /images/json
	Endpoint string

	// Method is the HTTP method type. ex: GET, POST, DELETE, etc.
	Method string

	// ContentType is the value to return in the header for Content-Type
	ContentType string

	// Body is the value to return in the body of the response
	Body string
}

// HandleFunc will register a handler with the mux based on the MockEndpoint passed to it
func (f *Focker) HandleFunc(m *MockEndpoint) {
	f.Mux.HandleFunc(m.Endpoint,
		func(w http.ResponseWriter, r *http.Request) {

			// Check the request method
			if r.Method != m.Method {
				w.Header().Set("Allow", m.Method)
				w.WriteHeader(http.StatusMethodNotAllowed)
			}

			w.Header().Set("Content-Type", m.ContentType)
			fmt.Fprint(w, m.Body)
		})
}
