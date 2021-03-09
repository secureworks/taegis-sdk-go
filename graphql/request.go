package graphql

import (
	"fmt"
	"net/http"
	"path"

	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/log"
)

// Request is the standard graphql request format to send to an API
// This should NOT be instantiated directly, you should use the NewRequest method
// to make sure your getting everything you need set
type Request struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
	Header    http.Header
	logger    log.Logger
}

// NewRequest creates a default graphql request with empty vars
// The Request struct should NOT be created directly - it should use this method so we can
// make sure everything is set on your graphql request
func NewRequest(query string, opts ...RequestOption) *Request {
	r := &Request{Query: query, Variables: map[string]interface{}{}, Header: http.Header{}}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Var lets you set variables for input in a graphql request
func (r *Request) Var(key string, value interface{}) {
	r.Variables[key] = value
}

// Vars lets you add multiple values at once if you know them ahead of time
func (r *Request) Vars(vars map[string]interface{}) {
	for k, v := range vars {
		r.Variables[k] = v
	}
}

// Error is the response type for graphql errors from our gateway
type Error struct {
	Message   string           `json:"message"`
	Locations []map[string]int `json:"locations"`
	Path      []string         `json:"path"`
}

func (e Error) Error() string {
	if len(e.Path) == 0 {
		return fmt.Sprintf("message: %s", e.Message)
	} else {
		return fmt.Sprintf("message: %s (path %s)", e.Message, path.Join(e.Path...))
	}
}

type Response struct {
	Data  interface{} `json:"data"`
	Error []Error     `json:"errors"`
}

// ResponseFields is used to define what graphql response fields we want back from the server
type ResponseFields string

// RequestOption are for the services to use to provide a token *per* request instead of one for the overall client
type RequestOption func(r *Request)

// RequestWithToken adds a bearer token to the individual request, overriding one set by client.Options
func RequestWithToken(token string) RequestOption {
	return func(r *Request) {
		r.Header.Add(common.AuthorizationHeader, "Bearer "+token)
	}
}

// RequestWithTenant adds a tenant id to the individual request, overriding one set by client.Options
func RequestWithTenant(tenantID string) RequestOption {
	return func(r *Request) {
		r.Header.Add(common.XTenantContextHeader, tenantID)
	}
}

// RequestWithLogger adds a logger for possibly more verbose output
func RequestWithLogger(logger log.Logger) RequestOption {
	return func(r *Request) {
		r.logger = logger
	}
}

// RequestWithHeader adds a set of headers to the request, any headers that are already present, will be skipped
func RequestWithHeader(header http.Header) RequestOption {
	return func(r *Request) {
		for k := range header {
			if _, ok := header[k]; ok { //Only add the Header if it isn't already set
				continue
			}
			r.Header.Add(k, header.Get(k))
		}
	}
}
