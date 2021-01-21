package testutils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/secureworks/tdr-sdk-go/graphql"
	"github.com/stretchr/testify/assert"
)

func NewMockServer(t *testing.T, expectedHeader http.Header, expectedMethod string, expectJSON bool, responseStatus int, response []byte) *httptest.Server {
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedMethod, r.Method)
		if expectJSON {
			assert.Contains(t, r.Header.Get("Content-Type"), "application/json")
		} else {
			for k := range expectedHeader {
				assert.Contains(t, r.Header.Get(k), expectedHeader.Get(k))
			}
		}

		w.WriteHeader(responseStatus)
		_, _ = w.Write(response)
	}

	return httptest.NewServer(http.HandlerFunc(fakeHandler))

}

func NewMockGQLServer(t *testing.T, expectedHeader http.Header, responseStatus int, response []byte) *httptest.Server {
	return NewMockServer(t, expectedHeader, http.MethodPost, true, responseStatus, response)
}

func NewMockGQLOutput(t *testing.T, expectedHeader http.Header, output interface{}) *httptest.Server {
	data, _ := json.Marshal(graphql.Response{Data: output})
	return NewMockGQLServer(t, expectedHeader, http.StatusOK, data)
}

func NewMockGQLError(t *testing.T, expectedHeader http.Header) *httptest.Server {
	data, _ := json.Marshal(graphql.Response{Error: []graphql.Error{{
		Message: "expected at least one definition, found }",
		Locations: []map[string]int{
			{
				"line":   4,
				"column": 36,
			},
		},
	}}})
	return NewMockGQLServer(t, expectedHeader, http.StatusOK, data)
}

func CreateHeader() http.Header {
	h := http.Header{}
	h.Add("Content-Type", "application/json")
	return h
}

// ToGenericMap converts a struct to the generic `map[string]interface{}` used as
// the default for json.Unmarshal. Its main purpose is to convert GraphQL input
// values when writing unit tests for GraphQL API clients. It cannot be Object
// above because the type does not match for equality checks.
func ToGenericMap(i interface{}) map[string]interface{} {
	b, _ := json.Marshal(i)

	o := map[string]interface{}{}
	_ = json.Unmarshal(b, &o)

	return o
}
