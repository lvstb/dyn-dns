package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWanIP(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("192.168.1.1"))
	}))
	defer server.Close()

	// Test
	ip, err := getWanIP(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if ip != "192.168.1.1" {
		t.Errorf("Expected IP 192.168.1.1, got %s", ip)
	}
}