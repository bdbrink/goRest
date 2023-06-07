package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	// Create a request to pass to the handler function
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create an HTTP handler function based on the helloHandler function
	handler := http.HandlerFunc(helloHandler)

	// Serve the HTTP request and record the response
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v, want %v",
			rr.Body.String(), expected)
	}
}
