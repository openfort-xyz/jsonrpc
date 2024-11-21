package jsonrpc

import (
	"context"
	"encoding/json"
)

const jsonRPCVersion = "2.0"

// Request is a JSON-RPC request.
type Request struct {
	ctx     context.Context
	url     string
	Vars    map[string]string `json:"-"`
	JSONRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  json.RawMessage   `json:"params"`
	ID      interface{}       `json:"id"`
}

// Response is a JSON-RPC response.
type Response struct {
	Raw        []byte      `json:"-"`
	StatusCode int         `json:"-"`
	JSONRPC    string      `json:"jsonrpc"`
	Result     interface{} `json:"result,omitempty"`
	Error      *Error      `json:"error,omitempty"`
	ID         interface{} `json:"id"`
}

// Error is a JSON-RPC error.
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
