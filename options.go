package moacrpc

import (
	"io"
	"net/http"
)

type httpClient interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

type logger interface {
	Println(v ...interface{})
}

// WithHttpClient set custom http client
func WithHttpClient(client httpClient) func(rpc *MoacRPC) {
	return func(rpc *MoacRPC) {
		rpc.client = client
	}
}

// WithLogger set custom logger
func WithLogger(l logger) func(rpc *MoacRPC) {
	return func(rpc *MoacRPC) {
		rpc.log = l
	}
}

// WithDebug set debug flag
func WithDebug(enabled bool) func(rpc *MoacRPC) {
	return func(rpc *MoacRPC) {
		rpc.Debug = enabled
	}
}
