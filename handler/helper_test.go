package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
)

func setupHelper() {
	secretKey = "correct"
	readAllFunc = ioutil.ReadAll
	jsonUnmarshal = json.Unmarshal
	parseIP = net.ParseIP
}

func TestValidateKey(t *testing.T) {
	setupHelper()
	if !validateKey("correct") {
		t.Errorf("Failed Correct Validate Key")
	}

	if validateKey("wrong") {
		t.Error("Failed Wrong Validate Key")
	}
}

type readCloserMock struct {
	io.Reader
}

func (readCloserMock) Close() error {
	return nil
}

type testStruct struct {
	Test string `json:"test"`
}

func TestReadJSONBodyReadAllFuncFail(t *testing.T) {
	setupHelper()
	readAllFunc = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("test")
	}
	body := readCloserMock{bytes.NewBufferString(`{"test":"test"}`)}
	testThing := testStruct{}
	readJSONBody(body, &testThing)

	structText := ""
	if testThing.Test != structText {
		t.Errorf("function returned struct: got %v want %v",
			testThing.Test, structText)
	}
}

func TestReadJSONZeroFail(t *testing.T) {
	setupHelper()
	readAllFunc = func(r io.Reader) ([]byte, error) {
		return []byte{}, nil
	}
	body := readCloserMock{bytes.NewBufferString(`{"test":"test"}`)}
	testThing := testStruct{}
	readJSONBody(body, &testThing)

	structText := ""
	if testThing.Test != structText {
		t.Errorf("function returned struct: got %v want %v",
			testThing.Test, structText)
	}
}

func TestReadJSONUnmarshalFail(t *testing.T) {
	setupHelper()
	jsonUnmarshal = func(data []byte, v interface{}) error {
		return errors.New("test")
	}
	body := readCloserMock{bytes.NewBufferString(`{"test":"test"}`)}
	testThing := testStruct{}
	readJSONBody(body, &testThing)

	structText := ""
	if testThing.Test != structText {
		t.Errorf("function returned struct: got %v want %v",
			testThing.Test, structText)
	}
}

func TestGetIPPortFail(t *testing.T) {
	req := http.Request{}
	err := getIP(&req)

	errorText := "missing port in address"
	if err.Error() != errorText {
		t.Errorf("function returned struct: got %v want %v",
			err.Error(), errorText)
	}
	setupHelper()
}

func TestGetIPFail(t *testing.T) {
	req := http.Request{}
	req.RemoteAddr = `127.0.0.1:8080`
	parseIP = func(s string) net.IP { return nil }
	err := getIP(&req)

	errorText := "parseIP Error"
	if err.Error() != errorText {
		t.Errorf("function returned struct: got %v want %v",
			err.Error(), errorText)
	}
	setupHelper()
}

func TestGetIPSuccess(t *testing.T) {
	req := http.Request{}
	req.RemoteAddr = `127.0.0.1:8080`
	err := getIP(&req)

	errorText := "nil"
	if err != nil {
		t.Errorf("function returned struct: got %v want %v",
			err.Error(), errorText)
	}
	setupHelper()
}
