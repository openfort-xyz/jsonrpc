package jsonrpc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_UsesDefaultHttpClient(t *testing.T) {
	client := NewClient()
	if client.http != http.DefaultClient {
		t.Errorf("expected default http client, got %v", client.http)
	}
}

func TestNewClient_AppliesOptions(t *testing.T) {
	client := NewClient(func(c *Client) {
		c.agent = "test-agent"
	})
	if client.agent != "test-agent" {
		t.Errorf("expected agent to be 'test-agent', got %v", client.agent)
	}
}

func TestDo_ValidRequest_ReturnsResponse(t *testing.T) {
	client := NewClient()
	req := NewRequestWithContext(context.Background(), 1, "test.method", "http://example.com", json.RawMessage(`{}`))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success"}`))
	}))
	defer server.Close()

	req.url = server.URL
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Result != "success" {
		t.Errorf("expected success, got %v", resp.Result)
	}
}

func TestDo_InvalidJSON_ReturnsError(t *testing.T) {
	client := NewClient()
	req := NewRequestWithContext(context.Background(), 1, "test.method", "http://example.com", json.RawMessage(`{`))

	_, err := client.Do(req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDo_Non200StatusCode_ReturnsError(t *testing.T) {
	client := NewClient()
	req := NewRequestWithContext(context.Background(), 1, "test.method", "http://example.com", json.RawMessage(`{}`))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`internal server error`))
	}))
	defer server.Close()

	req.url = server.URL
	_, err := client.Do(req)
	if err == nil || (err.Error() != "internal server error") {
		t.Fatalf("expected internal server error, got %v", err)
	}
}

func TestDo_ReadResponseBodyError_ReturnsError(t *testing.T) {
	client := NewClient()
	req := NewRequestWithContext(context.Background(), 1, "test.method", "http://example.com", json.RawMessage(`{}`))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success"}`))
	}))
	defer server.Close()

	req.url = server.URL
	server.Close()

	_, err := client.Do(req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
