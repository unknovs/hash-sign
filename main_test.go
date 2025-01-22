package main

import (
	"net/http"
	"testing"
	"time"
)

func TestServerListening(t *testing.T) {
	println("Starting main.go tests!!!")
	go main()
	time.Sleep(1 * time.Second)
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Errorf("Failed to send GET request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}
