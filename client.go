package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Client is a JSON-RPC client.
type Client struct {
	http  *http.Client
	agent string
}

// NewClient creates a new JSON-RPC client with the given ClientOption.
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		http: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// NewRequestWithContext creates a new JSON-RPC request with the given context, id, method, url, and params.
func NewRequestWithContext(ctx context.Context, id interface{}, method, url string, params json.RawMessage) *Request {
	return &Request{
		ctx:     ctx,
		url:     url,
		ID:      id,
		JSONRPC: jsonRPCVersion,
		Method:  method,
		Params:  params,
	}
}

// Do sends the given JSON-RPC Request and returns the JSON-RPC Response.
func (c *Client) Do(req *Request) (*Response, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(req.ctx, http.MethodPost, req.url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.agent != "" {
		httpReq.Header.Set("User-Agent", c.agent)
	}

	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		responseText, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(responseText))
	}

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp Response
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}

	resp.Raw = respBody
	resp.StatusCode = httpResp.StatusCode

	return &resp, nil
}
