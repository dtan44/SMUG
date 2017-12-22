package service

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"testing"
)

func setup() {
	defaultClient = &http.Client{}
	dumpRequestFunc = httputil.DumpRequest
	requestFunc = http.NewRequest
	readAllFunc = ioutil.ReadAll
}

// Test nil response with error from request
type clientError struct {
}

func (clientError) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("test")
}

func TestSendRequestFail(t *testing.T) {
	setup()
	var client clientError

	_, _, err := sendRequest("", "", nil, "test", client)
	if err.Error() != errors.New("test").Error() {
		t.Fail()
	}
}

func TestSendRequestClientNil(t *testing.T) {
	setup()
	defaultClient = clientError{}
	_, _, err := sendRequest("", "", nil, "", nil)
	if err.Error() != errors.New("test").Error() {
		t.Fail()
	}
}

func TestSendRequestDumpRequestFail(t *testing.T) {
	setup()
	var client clientError
	dumpRequestFunc = func(req *http.Request, body bool) ([]byte, error) {
		return nil, errors.New("test")
	}
	_, _, err := sendRequest("", "", nil, "", client)
	if err.Error() != errors.New("test").Error() {
		t.Fail()
	}
}

func TestSendRequestFuncFail(t *testing.T) {
	setup()
	var client clientError
	requestFunc = func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("test")
	}
	_, _, err := sendRequest("", "", nil, "", client)
	if err.Error() != errors.New("test").Error() {
		t.Fail()
	}
}

func TestSendRequestReadAllFail(t *testing.T) {
	setup()
	var client clientStatusSuccess
	readAllFunc = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("test")
	}
	_, _, err := sendRequest("", "", nil, "", client)
	if err.Error() != errors.New("test").Error() {
		t.Fail()
	}
}

// Test status error from request
// type clientStatusError struct {
// }

type readCloserMock struct {
	io.Reader
}

func (readCloserMock) Close() error {
	return nil
}

// func (clientStatusError) Do(req *http.Request) (*http.Response, error) {
// 	rs := http.Response{}
// 	rs.StatusCode = http.StatusBadRequest
// 	rs.Status = "status error"
// 	rs.Body = readCloserMock{bytes.NewBufferString("")}
// 	return &rs, nil
// }

// Test success request
type clientStatusSuccess struct {
}

func (clientStatusSuccess) Do(req *http.Request) (*http.Response, error) {
	rs := http.Response{}
	rs.StatusCode = http.StatusOK
	rs.Body = readCloserMock{bytes.NewBufferString("SUCCESS")}
	return &rs, nil
}

func TestSendRequestSuccess(t *testing.T) {
	setup()
	var client clientStatusSuccess

	_, rs, _ := sendRequest("", "", map[string]string{"test": "test"}, "", client)
	if string(rs) != "SUCCESS" {
		t.Fail()
	}
}
