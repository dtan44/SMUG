package service

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

func setupServiceDiscovery() {
	serviceMap = make(map[string]string)
	readAllFunc = ioutil.ReadAll
	request = sendRequest
}

func TestDiscoveryList(t *testing.T) {
	setupServiceDiscovery()
	var ds DiscoveryService
	serviceMap["test"] = "test"

	// invalid URL
	res := ds.List()
	if len(res) != 1 {
		t.Errorf("service returned unexpected error: got %v want %v",
			strconv.Itoa(len(res)), strconv.Itoa(0))
	}
}

func TestDiscoveryRouteServiceNameFail(t *testing.T) {
	setupServiceDiscovery()
	var ds DiscoveryService
	req := http.Request{}
	url, _ := url.ParseRequestURI("http://www.test.com")
	req.URL = url
	req.Body = readCloserMock{bytes.NewBufferString("test")}
	req.Method = http.MethodPost
	req.Header = http.Header{}

	errorText := "Invalid Service Name"

	// invalid URL
	_, _, err := ds.Route(&req)
	if err.Error() != errorText {
		t.Errorf("service returned unexpected error: got %v want %v",
			err.Error(), errorText)
	}
}

func TestDiscoveryRouteReadAllFuncFail(t *testing.T) {
	setupServiceDiscovery()
	var ds DiscoveryService

	// setup helper
	serviceMap["test"] = "http://www.test.com"
	readAllFunc = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("test")
	}

	// set up request
	req := http.Request{}
	url, _ := url.ParseRequestURI("http://www.test.com/service/test/check")
	req.URL = url
	req.Body = readCloserMock{bytes.NewBufferString("test")}
	req.Method = http.MethodPost
	req.Header = http.Header{}

	errorText := "test"

	// invalid URL
	_, _, err := ds.Route(&req)
	if err.Error() != errorText {
		t.Errorf("service returned unexpected error: got %v want %v",
			err.Error(), errorText)
	}
}

func TestDiscoveryRouteRequestFail(t *testing.T) {
	setupServiceDiscovery()
	var ds DiscoveryService

	// setup helper
	serviceMap["test"] = "http://www.test.com"
	request = func(url, httpMethod string,
		headers map[string]string, body string,
		client clientInterface) (*http.Response, []byte, error) {
		return nil, nil, errors.New("test")
	}

	// set up request
	req := http.Request{}
	url, _ := url.ParseRequestURI("http://www.test.com/service/test/check")
	req.URL = url
	req.Body = readCloserMock{bytes.NewBufferString("test")}
	req.Method = http.MethodPost
	req.Header = http.Header{}
	req.Header.Set("test", "test,test")

	errorText := "test"

	// invalid URL
	_, _, err := ds.Route(&req)
	if err.Error() != errorText {
		t.Errorf("service returned unexpected error: got %v want %v",
			err.Error(), errorText)
	}
}
