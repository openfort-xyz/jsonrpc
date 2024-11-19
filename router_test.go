package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP_ValidRequest_CallsCorrectHandler(t *testing.T) {
	router := NewRouter()
	router.HandleFunc("test.method", func(ctx context.Context, req *Request) *Response {
		return &Response{Result: "success"}
	})

	reqBody := `{"jsonrpc":"2.0","method":"test.method","id":1}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	var response Response
	_ = json.NewDecoder(resp.Body).Decode(&response)

	if response.Result != "success" {
		t.Errorf("expected success, got %v", response.Result)
	}
}

func TestServeHTTP_InvalidJSON_ReturnsParseError(t *testing.T) {
	router := NewRouter()

	reqBody := `{"jsonrpc":"2.0","method":"test.method","id":1`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %v", resp.StatusCode)
	}

	var response Response
	_ = json.NewDecoder(resp.Body).Decode(&response)

	if response.Error == nil || response.Error.Code != -32700 {
		t.Errorf("expected parse error, got %v", response.Error)
	}
}

func TestServeHTTP_MethodNotFound_ReturnsMethodNotFoundError(t *testing.T) {
	router := NewRouter()

	reqBody := `{"jsonrpc":"2.0","method":"unknown.method","id":1}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %v", resp.StatusCode)
	}

	var response Response
	_ = json.NewDecoder(resp.Body).Decode(&response)

	if response.Error == nil || response.Error.Code != -32601 {
		t.Errorf("expected method not found error, got %v", response.Error)
	}
}

func TestServeHTTP_UsesMiddleware(t *testing.T) {
	router := NewRouter()
	router.HandleFunc("test.method", func(ctx context.Context, req *Request) *Response {
		return &Response{Result: "success"}
	})
	router.Use(func(next Handler) Handler {
		return func(ctx context.Context, req *Request) *Response {
			resp := next(ctx, req)
			resp.Result = "modified by middleware"
			return resp
		}
	})

	reqBody := `{"jsonrpc":"2.0","method":"test.method","id":1}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	var response Response
	_ = json.NewDecoder(resp.Body).Decode(&response)

	if response.Result != "modified by middleware" {
		t.Errorf("expected modified by middleware, got %v", response.Result)
	}
}
