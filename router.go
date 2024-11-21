package jsonrpc

import (
	"context"
	"encoding/json"
	"net/http"
)

// Router is a JSON-RPC 2.0 router.
type Router struct {
	methodHandlers map[string]Handler
	defaultHandler *Handler
	middlewares    []Middleware
	agent          string
	vars           func(*http.Request) map[string]string
}

// Handler is a JSON-RPC 2.0 handler.
type Handler func(context.Context, *Request) *Response

// Middleware is a JSON-RPC 2.0 middleware.
type Middleware func(Handler) Handler

// NewRouter creates a new JSON-RPC 2.0 router with the given RouterOption.
func NewRouter(opts ...RouterOption) *Router {
	r := &Router{
		methodHandlers: make(map[string]Handler),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// HandleFunc registers a JSON-RPC 2.0 handler for the given method.
func (r *Router) HandleFunc(method string, handler Handler) {
	r.methodHandlers[method] = handler
}

// DefaultHandler sets the default JSON-RPC 2.0 handler.
func (r *Router) DefaultHandler(handler Handler) {
	r.defaultHandler = &handler
}

// Use registers a JSON-RPC 2.0 Middleware.
func (r *Router) Use(middleware Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

// ServeHTTP implements the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var request Request
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		r.writeError(w, -32700, "Parse error", err.Error())
		return
	}

	if r.vars != nil {
		request.Vars = r.vars(req)
	}

	handler, ok := r.methodHandlers[request.Method]
	if !ok {
		if r.defaultHandler == nil {
			r.writeError(w, -32601, "Method not found", nil)
			return
		}
		handler = *r.defaultHandler
	}

	for _, middleware := range r.middlewares {
		handler = middleware(handler)
	}

	resp := handler(ctx, &request)
	r.writeResponse(w, resp)
}

func (r *Router) writeResponse(w http.ResponseWriter, resp *Response) {
	w.Header().Set("Content-Type", "application/json")
	if r.agent != "" {
		w.Header().Set("User-Agent", r.agent)
	}
	w.WriteHeader(resp.StatusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

func (r *Router) writeError(w http.ResponseWriter, code int, message string, data interface{}) {
	resp := &Response{
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
	if data != nil {
		resp.Error.Data = data
	}
	w.Header().Set("Content-Type", "application/json")
	if r.agent != "" {
		w.Header().Set("User-Agent", r.agent)
	}
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(resp)
}
