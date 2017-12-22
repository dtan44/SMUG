package service

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Global variables
var defaultClient clientInterface
var dumpRequestFunc func(req *http.Request, body bool) ([]byte, error)
var requestFunc func(method, url string, body io.Reader) (*http.Request, error)
var readAllFunc func(r io.Reader) ([]byte, error)

func init() {
	defaultClient = &http.Client{}
	dumpRequestFunc = httputil.DumpRequest
	requestFunc = http.NewRequest
	readAllFunc = ioutil.ReadAll
}

//ClientInterface
type clientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

//sendRequest: sends generic http request
//return: response body or error
func sendRequest(url, httpMethod string, headers map[string]string,
	body string, client clientInterface) (*http.Response, []byte, error) {

	// Default client
	if client == nil {
		client = defaultClient
	}

	// Initialize request
	var req *http.Request
	var err error
	if body == "" {
		req, err = requestFunc(httpMethod, url, nil)
	} else {
		req, err = requestFunc(httpMethod, url, bytes.NewBufferString(body))
	}
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Set(key, val)
		}
	}

	// Log request
	requestDump, err := dumpRequestFunc(req, true)
	if err != nil {
		log.Error(err)
	}
	log.Info(strings.TrimRight(string(requestDump), "\r\n"))

	// Send request
	rs, err := client.Do(req)
	if err != nil {
		log.Error("Error with request: " + err.Error())
		return nil, nil, err
	}
	defer rs.Body.Close()

	// Process Response Body
	rsBody, err := readAllFunc(rs.Body)
	if err != nil {
		log.Error("Error reading response body")
		return nil, nil, err
	}
	return rs, rsBody, nil
}
