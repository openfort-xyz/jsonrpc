package jsonrpc

import "net/http"

// ClientOption is a function that configures a Client.
type ClientOption func(c *Client)

// WithHTTPClient sets the HTTP client for the Client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.http = httpClient
	}
}

// WithClientAgent sets the User-Agent for the Client.
func WithClientAgent(agent string) ClientOption {
	return func(c *Client) {
		c.agent = agent
	}
}

// RouterOption is a function that configures a Router.
type RouterOption func(r *Router)

// WithRouterAgent sets the User-Agent for the Router.
func WithRouterAgent(agent string) RouterOption {
	return func(r *Router) {
		r.agent = agent
	}
}
