package graphql_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/secureworks/taegis-sdk-go/client"
	"github.com/secureworks/taegis-sdk-go/graphql"
	"github.com/secureworks/taegis-sdk-go/testutils"
	"github.com/stretchr/testify/assert"
)

func newHeader() http.Header {
	result := http.Header{}
	result.Add("my_header", "val")
	return result
}

func TestExecuteQuery_SuccessBadResponseBody(t *testing.T) {
	server := testutils.NewMockGQLServer(t, newHeader(), http.StatusOK, []byte("status=success"))
	defer server.Close()

	result := struct {
		Value *string `mapstructure:"value"`
	}{}

	err := graphql.ExecuteQuery(client.NewClient(), server.URL, graphql.NewRequest("testQuery", graphql.RequestWithHeader(newHeader())), &result)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "error decoding response"))
}

func TestExecuteQuery_Success(t *testing.T) {
	type testStruct struct {
		Value *string `json:"value" mapstructure:"value"`
	}
	val := "val"
	expected := testStruct{Value: &val}

	server := testutils.NewMockGQLOutput(t, newHeader(), expected)
	defer server.Close()

	result := testStruct{}

	err := graphql.ExecuteQuery(client.NewClient(), server.URL, graphql.NewRequest("testQuery", graphql.RequestWithHeader(newHeader())), &result)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestExecuteQueryIgnoreOutput_Success(t *testing.T) {
	type testStruct struct {
		Value *string `json:"value" mapstructure:"value"`
	}
	val := "blabla"
	expected := testStruct{Value: &val}

	server := testutils.NewMockGQLOutput(t, newHeader(), expected)
	defer server.Close()

	err := graphql.ExecuteQuery(client.NewClient(), server.URL, graphql.NewRequest("testQuery", graphql.RequestWithHeader(newHeader())), nil)
	assert.Nil(t, err)
}

func TestExecuteQuery_GraphqlError(t *testing.T) {
	server := testutils.NewMockGQLError(t, newHeader())
	defer server.Close()

	err := graphql.ExecuteQuery(client.NewClient(), server.URL, graphql.NewRequest("testQuery", graphql.RequestWithHeader(newHeader())), nil)
	assert.NotNil(t, err)
	assert.Equal(t, "1 error occurred:\n\t* message: expected at least one definition, found }\n\n", err.Error())
}

func TestExecuteQuery_ServerErrors(t *testing.T) {
	errMsg := "unknown field 'foobar'"
	rsp := graphql.Response{
		Error: []graphql.Error{
			{
				Message: errMsg,
			},
		},
	}
	b, err := json.Marshal(rsp)
	assert.Nil(t, err)
	server := testutils.NewMockGQLServer(t, newHeader(), http.StatusInternalServerError, b)

	result := struct {
		Value *string `mapstructure:"value"`
	}{}

	err = graphql.ExecuteQuery(client.NewClient(), server.URL, graphql.NewRequest("testQuery", graphql.RequestWithHeader(newHeader())), &result)
	assert.NotNil(t, err)
	errStr := err.Error()
	// The error should contain two messages
	assert.Contains(t, errStr, "server responded with an error: 500")
	assert.Contains(t, errStr, errMsg)

	server.Close()
	err = graphql.ExecuteQuery(client.NewClient(), server.URL, graphql.NewRequest("testQuery", graphql.RequestWithHeader(newHeader())), &result)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "server connection error")
}

func TestExecuteQueryWithTenant(t *testing.T) {
	type testStruct struct {
		Value *string `json:"value" mapstructure:"value"`
	}
	val := "val"
	expected := testStruct{Value: &val}

	server := testutils.NewMockGQLOutput(t, newHeader(), expected)
	defer server.Close()

	result := testStruct{}
	const expectedTenant = "tenant"
	tenant, err := graphql.ExecuteQueryWithTenant(client.NewClient(), server.URL, graphql.NewRequest("testQuery", graphql.RequestWithHeader(newHeader()), graphql.RequestWithTenant(expectedTenant)), &result)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, expectedTenant, tenant)

	tenant, err = graphql.ExecuteQueryWithTenant(client.NewClient(client.WithTenant(expectedTenant)), server.URL, graphql.NewRequest("testQuery", graphql.RequestWithHeader(newHeader())), &result)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, expectedTenant, tenant)

	//should fail when no tenant header specified
	_, err = graphql.ExecuteQueryWithTenant(client.NewClient(), server.URL, graphql.NewRequest("testQuery", graphql.RequestWithHeader(newHeader())), &result)
	assert.Error(t, err)
}
