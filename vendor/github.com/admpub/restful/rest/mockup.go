package rest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
)

/*
func init() {
	flag.BoolVar(&mockUpEnv, "mock", false,
		"Use 'mock' flag to tell package rest that you would like to use mockups.")

	flag.Parse()
	startMockupServ()
}
*/

var mockUpEnv bool
var mockMap = make(map[string]*Mock)

var mockServer *httptest.Server
var mux *http.ServeMux

var mockServerURL *url.URL

// Mock serves the purpose of creating Mockups.
// All requests will be sent to the mockup server if mockup is activated.
// To activate the mockup *environment* you have two ways: using the flag -mock
//	go test -mock
//
// Or by programmatically starting the mockup server
// 	StartMockupServer()
type Mock struct {

	// Request URL
	URL string

	// Request HTTP Method (GET, POST, PUT, PATCH, HEAD, DELETE, OPTIONS)
	// As a good practice use the constants in http package (http.MethodGet, etc.)
	HTTPMethod string

	// Request array Headers
	ReqHeaders http.Header

	// Request Body, used with POST, PUT & PATCH
	ReqBody string

	// Response HTTP Code
	RespHTTPCode int

	// Response Array Headers
	RespHeaders http.Header

	// Response Body
	RespBody string
}

// StartMockupServer sets the environment to send all client requests
// to the mockup server.
func StartMockupServer() {

	mockUpEnv = true

	if mockServer == nil {
		startMockupServ()
	}
}

// StopMockupServer stop sending requests to the mockup server
func StopMockupServer() {

	mockUpEnv = false
	mockServer.Close()

	mockServer = nil
	mockServerURL = nil
	mux = nil
}

func startMockupServ() {

	if mockUpEnv {
		mux = http.NewServeMux()
		mockServer = httptest.NewServer(mux)
		mux.HandleFunc("/", mockupHandler)

		var err error
		if mockServerURL, err = url.Parse(mockServer.URL); err != nil {
			panic(err)
		}

	}
}

// AddMockups ...
func AddMockups(mocks ...*Mock) {
	for _, m := range mocks {
		mockMap[m.HTTPMethod+" "+m.URL] = m
	}
}

// FlushMockups ...
func FlushMockups() {
	mockMap = make(map[string]*Mock)
}

func mockupHandler(writer http.ResponseWriter, req *http.Request) {

	url := req.Header.Get("X-Original-URL")

	if m := mockMap[req.Method+" "+url]; m != nil {

		// Add headers
		for k, v := range m.RespHeaders {
			for _, vv := range v {
				writer.Header().Add(k, vv)
			}
		}

		writer.WriteHeader(m.RespHTTPCode)
		writer.Write([]byte(m.RespBody))
		return
	}

	writer.WriteHeader(http.StatusBadRequest)
	writer.Write([]byte("MockUp nil!"))
}
