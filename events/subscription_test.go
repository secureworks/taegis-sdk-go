package events

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"

	"github.com/secureworks/taegis-sdk-go/client"
	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/graphql"
)

var (
	now                      = time.Now()
	testRow                  = common.Object{"foo": "bar"}
	testEventQueryResultsOne = &EventQueryResults{
		Query: &EventQuery{},
		Result: &EventQueryResult{
			ID:        "ID One",
			Type:      "Type",
			Submitted: now,
			Completed: &now,
			Expires:   &now,
			Status:    "SUCCEEDED",
			Facets:    common.Object{"Facets": "test"},
			Rows:      []common.Object{testRow},
		},
	}

	testNextPage             = "next:page"
	testEventQueryResultsTwo = &EventQueryResults{
		Query: &EventQuery{},
		Result: &EventQueryResult{
			ID:        "ID Two",
			Type:      "Type",
			Submitted: now,
			Completed: &now,
			Expires:   &now,
			Status:    "SUCCEEDED",
			Facets:    common.Object{"Facets": "test"},
			Rows:      []common.Object{testRow},
		},
		Next: &testNextPage,
	}
	testEventQueryResultsThree = &EventQueryResults{ //Running Status
		Query: &EventQuery{
			Status: "RUNNING",
		},
	}
	testClosureEventQueryResult = &EventQueryResults{
		Query: &EventQuery{
			Status: "SUCCEEDED",
		},
	}
)

func TestEventQuery_Next(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(eventQuerySubscription, Query, Metadata, Options),
		map[string]interface{}{Query: "", Metadata: nil, Options: nil},
		&eventQueryResult{EventQueryResults: testEventQueryResultsOne},
	)

	defer s.Close()
	svc := New(s.URL, client.WithHTTPTimeout(5*time.Second))

	sub, err := svc.EventQuery(context.Background(), "", nil, nil)
	require.NoError(t, err)
	defer sub.Close()

	c, err := sub.Next(context.Background())
	require.NoError(t, err)
	expectedData, err := json.Marshal(testEventQueryResultsOne)
	require.NoError(t, err)
	actualData, err := json.Marshal(c)
	require.JSONEq(t, string(expectedData), string(actualData))
}

func TestEventQuery_NextError(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(eventQuerySubscription, Query, Metadata, Options),
		map[string]interface{}{Query: "", Metadata: nil, Options: nil},
		errors.New("failed"),
	)

	defer s.Close()
	svc := New(s.URL, client.WithHTTPTimeout(5*time.Second))

	sub, err := svc.EventQuery(context.Background(), "", nil, nil)
	require.NoError(t, err)
	defer sub.Close()

	_, err = sub.Next(context.Background())
	require.Error(t, err)
}

func TestEventPage_Next(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(eventPageSubscription, PageID),
		map[string]interface{}{PageID: "ID"},
		&eventPageResult{EventQueryResults: testEventQueryResultsOne},
	)

	defer s.Close()
	svc := New(s.URL, client.WithHTTPTimeout(5*time.Second))

	sub, err := svc.EventPage(context.Background(), "ID")
	require.NoError(t, err)
	defer sub.Close()

	c, err := sub.Next(context.Background())
	require.NoError(t, err)
	expectedData, err := json.Marshal(testEventQueryResultsOne)
	require.NoError(t, err)
	actualData, err := json.Marshal(c)
	require.JSONEq(t, string(expectedData), string(actualData))
}

func TestEventPage_NextError(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(eventPageSubscription, PageID),
		map[string]interface{}{PageID: ""},
		errors.New("failed"),
	)

	defer s.Close()
	svc := New(s.URL, client.WithHTTPTimeout(5*time.Second))

	sub, err := svc.EventPage(context.Background(), "")
	require.NoError(t, err)
	defer sub.Close()

	_, err = sub.Next(context.Background())
	require.Error(t, err)
}

func TestEventPage_GetAllEventResults_EventQuery(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(eventQuerySubscription, Query, Metadata, Options),
		map[string]interface{}{Query: "query", Metadata: nil, Options: nil},
		&eventQueryResult{EventQueryResults: testEventQueryResultsOne},
		&eventQueryResult{EventQueryResults: testEventQueryResultsTwo},
		&eventQueryResult{EventQueryResults: testEventQueryResultsThree},  //Result with RUNNING status
		&eventQueryResult{EventQueryResults: testClosureEventQueryResult}, //testClosureEventQueryResult is the signal that all events have been returned. (what the sdk will return in a live situation)
	)

	defer s.Close()
	svc := New(s.URL, client.WithHTTPTimeout(5*time.Second))

	sub, err := svc.EventQuery(context.Background(), "query", nil, nil)
	require.NoError(t, err)
	defer sub.Close()

	c, next, err := sub.GetAllEventResults(context.Background())
	require.NoError(t, err)
	expectedData, err := json.Marshal(Results{testEventQueryResultsOne, testEventQueryResultsTwo})
	require.NoError(t, err)
	actualData, err := json.Marshal(c)
	require.JSONEq(t, string(expectedData), string(actualData))
	require.Equal(t, &testNextPage, next)
}

func TestEventPage_GetAllEventResults_EventPage(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(eventPageSubscription, PageID), map[string]interface{}{PageID: "page"},
		&eventPageResult{EventQueryResults: testEventQueryResultsOne},
		&eventPageResult{EventQueryResults: testEventQueryResultsTwo},
		&eventPageResult{EventQueryResults: testEventQueryResultsThree},  //Result with RUNNING status
		&eventPageResult{EventQueryResults: testClosureEventQueryResult}, //testClosureEventQueryResult is the signal that all events have been returned. (what the sdk will return in a live situation)
	)

	defer s.Close()
	svc := New(s.URL, client.WithHTTPTimeout(5*time.Second))

	sub, err := svc.EventPage(context.Background(), "page")
	require.NoError(t, err)
	defer sub.Close()

	c, next, err := sub.GetAllEventResults(context.Background())
	require.NoError(t, err)
	expectedData, err := json.Marshal(Results{testEventQueryResultsOne, testEventQueryResultsTwo})
	require.NoError(t, err)
	actualData, err := json.Marshal(c)
	require.JSONEq(t, string(expectedData), string(actualData))
	require.Equal(t, &testNextPage, next)
}

func TestEventPage_GetAllEventResults_Error(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(eventQuerySubscription, Query, Metadata, Options),
		map[string]interface{}{Query: "query", Metadata: nil, Options: nil},
		&eventQueryResult{EventQueryResults: testEventQueryResultsOne},
		&multierror.Error{Errors: []error{errors.New("test error")}}, //Error in the middle of results, client should still get testEventQueryResultsOne
		&eventQueryResult{EventQueryResults: testEventQueryResultsTwo},
	)

	defer s.Close()
	svc := New(s.URL, client.WithHTTPTimeout(5*time.Second))

	sub, err := svc.EventQuery(context.Background(), "query", nil, nil)
	require.NoError(t, err)
	defer sub.Close()

	c, next, err := sub.GetAllEventResults(context.Background())
	require.Error(t, err)
	expectedData, err := json.Marshal(Results{testEventQueryResultsOne})
	require.NoError(t, err)
	actualData, err := json.Marshal(c)
	require.JSONEq(t, string(expectedData), string(actualData))
	require.Nil(t, next)
}

func TestEventPage_GetAllEventResults_NilResultReturned(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(eventQuerySubscription, Query, Metadata, Options),
		map[string]interface{}{Query: "query", Metadata: nil, Options: nil},
		&eventQueryResult{EventQueryResults: testEventQueryResultsOne},
		nil, //Nil causes Error in the middle of results, client should still get testEventQueryResultsOne
		&eventQueryResult{EventQueryResults: testEventQueryResultsTwo},
	)

	defer s.Close()
	svc := New(s.URL, client.WithHTTPTimeout(5*time.Second))

	sub, err := svc.EventQuery(context.Background(), "query", nil, nil)
	require.NoError(t, err)
	defer sub.Close()

	c, next, err := sub.GetAllEventResults(context.Background())
	require.Error(t, err)
	expectedData, err := json.Marshal(Results{testEventQueryResultsOne})
	require.NoError(t, err)
	actualData, err := json.Marshal(c)
	require.JSONEq(t, string(expectedData), string(actualData))
	require.Nil(t, next)
}
