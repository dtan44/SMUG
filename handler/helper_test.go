package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"testing"
)

func setupHelper() {
	secretKey = "correct"
	readAllFunc = ioutil.ReadAll
	jsonUnmarshal = json.Unmarshal
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
