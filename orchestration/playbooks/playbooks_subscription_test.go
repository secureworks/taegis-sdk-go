package playbooks

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/secureworks/tdr-sdk-go/client"

	"github.com/secureworks/tdr-sdk-go/graphql"
	"github.com/stretchr/testify/require"
)

// ran into some issues when comparing using require.Equal in ci, so for the sake of good use of my time and sanity
// im just comparing the resulting jsons
func TestCreatedBasicFlow(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(playbookInstanceCreateSubQuery, playbookIDs),
		map[string]interface{}{playbookIDs: nil},
		&playbookInstanceCreatedEvent{PlaybookInstanceEvent: &testPlayBookInstance},
	)

	defer s.Close()
	svc := New(s.URL, client.WithHTTPTimeout(5*time.Second))

	sub, err := svc.PlaybookInstanceCreated(context.TODO(), nil)
	require.NoError(t, err)
	defer sub.Close()

	c, err := sub.Next(context.TODO())
	require.NoError(t, err)
	expectedData, err := json.Marshal(testPlayBookInstance)
	require.NoError(t, err)
	actualData, err := json.Marshal(c)
	require.JSONEq(t, string(expectedData), string(actualData))
}

func TestCreatedReturnsError(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(playbookInstanceCreateSubQuery, playbookIDs),
		map[string]interface{}{playbookIDs: nil},
		errors.New("failed"),
	)

	defer s.Close()
	svc := New(s.URL, client.WithHTTPTimeout(5*time.Second))

	sub, err := svc.PlaybookInstanceCreated(context.TODO(), nil)
	require.NoError(t, err)
	defer sub.Close()

	_, err = sub.Next(context.TODO())
	require.Error(t, err)
}

func TestUpdatedBasicFlow(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(playbookInstanceUpdateSubQuery, playbookIDs),
		map[string]interface{}{playbookIDs: nil},
		&playbookInstanceUpdatedEvent{PlaybookInstanceEvent: &testPlayBookInstance},
	)

	defer s.Close()
	svc := New(s.URL, client.WithHTTPTimeout(5*time.Second))

	sub, err := svc.PlaybookInstanceUpdated(context.TODO(), nil)
	require.NoError(t, err)
	defer sub.Close()

	c, err := sub.Next(context.TODO())
	require.NoError(t, err)
	expectedData, err := json.Marshal(testPlayBookInstance)
	require.NoError(t, err)
	actualData, err := json.Marshal(c)
	require.JSONEq(t, string(expectedData), string(actualData))
}

func TestDeletedBasicFlow(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(playbookInstanceDeleteSubQuery, playbookIDs),
		map[string]interface{}{playbookIDs: nil},
		&playbookInstanceDeletedEvent{PlaybookInstanceEvent: &testPlayBookInstance},
	)

	defer s.Close()
	svc := New(s.URL, client.WithHTTPTimeout(5*time.Second))

	sub, err := svc.PlaybookInstanceDeleted(context.TODO(), nil)
	require.NoError(t, err)
	defer sub.Close()

	c, err := sub.Next(context.TODO())
	require.NoError(t, err)
	expectedData, err := json.Marshal(testPlayBookInstance)
	require.NoError(t, err)
	actualData, err := json.Marshal(c)
	require.JSONEq(t, string(expectedData), string(actualData))
}
