package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/route53"
)

type mockR53set interface {
	ChangeResourceRecordSets(context.Context, *route53.ChangeResourceRecordSetsInput, ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error)
}

type mockRoute53Client struct {
	ChangeResourceRecordSetsFunc func(context.Context, *route53.ChangeResourceRecordSetsInput, ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error)
}

func (m mockRoute53Client) ChangeResourceRecordSets(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error) {
	return m.ChangeResourceRecordSetsFunc(ctx, params, optFns...)
}

func TestSetDNSRecord(t *testing.T) {
	ctx := context.Background()
	mockClient := &mockRoute53Client{
		ChangeResourceRecordSetsFunc: func(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error) {
			return &route53.ChangeResourceRecordSetsOutput{}, nil
		},
	}

	// Test success case
	if err := setDNSRecord(ctx, "test.example.com", "192.168.1.1", "Z1234567890", mockClient); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test error case
	mockClient.ChangeResourceRecordSetsFunc = func(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error) {
		return nil, fmt.Errorf("AWS error")
	}
	if err := setDNSRecord(ctx, "test.example.com", "192.168.1.1", "Z1234567890", mockClient); err == nil {
		t.Error("Expected an error, got nil")
	}
}

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