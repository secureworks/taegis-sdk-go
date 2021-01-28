package testutils

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/secureworks/tdr-sdk-go/graphql"
)

// MockGraphQLHandler provides a handler to mock a GraphQL server in
// combination with httptest.Server. It allows the response, headers, expected
// input variables and errors to be set up and changed between each test. It is
// designed to be run in a single test with sub-tests.
type MockGraphQLHandler struct {
	t *testing.T

	ResponseStatus    int
	Response          interface{}
	Errors            []graphql.Error
	ExpectedHeaders   http.Header
	ExpectedVariables map[string]interface{}
}

// NewMockGraphQLHandler sets up a MockGraphQLHandler to use the given testing.T
// for assertions with a default ResponseStatus of http.StatusOK.
func NewMockGraphQLHandler(t *testing.T) *MockGraphQLHandler {
	return &MockGraphQLHandler{
		t:              t,
		ResponseStatus: http.StatusOK,
	}
}

// ServeHTTP is the standard net/http handler which validates the incoming
// GraphQL request and returns a mock response or error.
func (g *MockGraphQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check headers
	assert.Contains(g.t, r.Header.Get("Content-Type"), "application/json")
	for k := range g.ExpectedHeaders {
		assert.Contains(g.t, r.Header.Get(k), g.ExpectedHeaders.Get(k), "header value %s", k)
	}

	// Read the request
	dec := json.NewDecoder(r.Body)
	gr := graphql.Request{}
	err := dec.Decode(&gr)
	assert.Nil(g.t, err)

	// Make sure the variables are what is expected if they are defined
	if g.ExpectedVariables != nil {
		assert.Equal(g.t, g.ExpectedVariables, gr.Variables)
	}

	var res graphql.Response
	if g.Errors != nil {
		res.Error = g.Errors
	} else {
		res.Data = g.Response
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(g.ResponseStatus)
	err = enc.Encode(res)
	assert.Nil(g.t, err)
}

// WithErrors will set up the provided errors to only be present during the
// provided callback, making it easy to set up errors for a single test.
func (g *MockGraphQLHandler) WithErrors(errors []graphql.Error, cb func()) {
	g.Errors = errors
	cb()
	g.Errors = nil
}
