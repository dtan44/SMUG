package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupServiceHandler() {
	jsonMarshal = json.Marshal
	secretKey = "correct"
}

var (
	ReturnError = func() error {
		return errors.New("")
	}
	ReturnNoError = func() error {
		return nil
	}
)

type RegisterMock struct {
	Work           func() error
	CreateUserWork func() (string, error)
}

func (rs RegisterMock) Register(serviceName, serviceURL string) error {
	return rs.Work()
}

func (rs RegisterMock) Deregister(serviceName string) error {
	return rs.Work()
}

func TestHandleRegisterFail(t *testing.T) {
	setupServiceHandler()
	var sh ServiceHandler
	sh.Registration = RegisterMock{Work: ReturnError}

	// Create empty request for handler.
	req, err := http.NewRequest("PUT", "/register/test", strings.NewReader(`{"URL" : ""}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("secret-key", "correct")

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sh.HandleRegister)

	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body.
	expected := `{"result":"failure"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandleRegisterJSONFail(t *testing.T) {
	setupServiceHandler()
	var sh ServiceHandler
	sh.Registration = RegisterMock{Work: ReturnError}
	jsonMarshal = func(v interface{}) ([]byte, error) {
		return nil, errors.New("test")
	}

	// Create empty request for handler.
	req, err := http.NewRequest("PUT", "/register/test", strings.NewReader(`{"URL" : ""}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("secret-key", "correct")

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sh.HandleRegister)

	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	// Check the response body.
	expected := "test\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandleRegisterKeyFail(t *testing.T) {
	setupServiceHandler()
	var sh ServiceHandler
	sh.Registration = RegisterMock{Work: ReturnError}

	// Create empty request for handler.
	req, err := http.NewRequest("PUT", "/register/test", strings.NewReader(`{"URL" : ""}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("secret-key", "wrong")

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sh.HandleRegister)

	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	// Check the response body.
	expected := `{"result":"failure","reason":"Incorrect Key"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandleRegisterSuccess(t *testing.T) {
	setupServiceHandler()
	var sh ServiceHandler
	sh.Registration = RegisterMock{Work: ReturnNoError}

	// Create empty request for handler.
	req, err := http.NewRequest("PUT", "/register/test", strings.NewReader(`{"URL" : ""}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("secret-key", "correct")

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sh.HandleRegister)

	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body.
	expected := `{"result":"success"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandleDeregisterFail(t *testing.T) {
	setupServiceHandler()
	var sh ServiceHandler
	sh.Registration = RegisterMock{Work: ReturnError}

	// Create empty request for handler.
	req, err := http.NewRequest("DELETE", "/deregister/test", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("secret-key", "correct")

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sh.HandleDeregister)

	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body.
	expected := `{"result":"failure"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandleDeregisterJSONFail(t *testing.T) {
	setupServiceHandler()
	var sh ServiceHandler
	sh.Registration = RegisterMock{Work: ReturnError}
	jsonMarshal = func(v interface{}) ([]byte, error) {
		return nil, errors.New("test")
	}

	// Create empty request for handler.
	req, err := http.NewRequest("DELETE", "/deregister/test", strings.NewReader(`{"URL" : ""}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("secret-key", "correct")

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sh.HandleDeregister)

	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	// Check the response body.
	expected := "test\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandleDeregisterKeyFail(t *testing.T) {
	setupServiceHandler()
	var sh ServiceHandler
	sh.Registration = RegisterMock{Work: ReturnError}

	// Create empty request for handler.
	req, err := http.NewRequest("DELETE", "/deregister/test", strings.NewReader(`{"URL" : ""}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("secret-key", "wrong")

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sh.HandleDeregister)

	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	// Check the response body.
	expected := `{"result":"failure","reason":"Incorrect Key"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandleDeregisterSuccess(t *testing.T) {
	setupServiceHandler()
	var sh ServiceHandler
	sh.Registration = RegisterMock{Work: ReturnNoError}

	// Create empty request for handler.
	req, err := http.NewRequest("DELETE", "/deregister/test", strings.NewReader(`{"URL" : ""}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("secret-key", "correct")

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sh.HandleDeregister)

	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body.
	expected := `{"result":"success"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
