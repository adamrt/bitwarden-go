package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIconHandlerSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "/icons/google.com/icon.png", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(iconHandler)
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMovedPermanently)
	}
}

func TestIconHandlerNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/icons/google.com/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(iconHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}
