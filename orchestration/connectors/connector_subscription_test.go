package connectors

import (
	"context"
	"errors"
	"testing"

	"github.com/secureworks/taegis-sdk-go/graphql"
	"github.com/stretchr/testify/require"
)

func TestCreatedBasicFlow(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(connectorCreatedSubQuery, methodsVarName, allTenantsVarName),
		map[string]interface{}{methodsVarName: nil, allTenantsVarName: false},
		&connectorCreatedEvent{ConnectorEvent: &TestConnector},
	)

	defer s.Close()
	svc := New(s.URL)

	sub, err := svc.ConnectorCreated(context.TODO(), nil, false)
	require.NoError(t, err)
	defer sub.Close()

	c, err := sub.Next(context.TODO())
	require.NoError(t, err)
	require.Equal(t, &TestConnector, c)
}

func TestCreatedReturnsError(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(connectorCreatedSubQuery, methodsVarName, allTenantsVarName),
		map[string]interface{}{methodsVarName: nil, allTenantsVarName: false},
		errors.New("failed"),
	)

	defer s.Close()
	svc := New(s.URL)

	sub, err := svc.ConnectorCreated(context.TODO(), nil, false)
	require.NoError(t, err)
	defer sub.Close()

	_, err = sub.Next(context.TODO())
	require.Error(t, err)
}

func TestUpdatedBasicFlow(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(connectorUpdatedSubQuery, methodsVarName, allTenantsVarName),
		map[string]interface{}{methodsVarName: nil, allTenantsVarName: false},
		&connectorUpdatedEvent{ConnectorEvent: &TestConnector},
	)

	defer s.Close()
	svc := New(s.URL)

	sub, err := svc.ConnectorUpdated(context.TODO(), nil, false)
	require.NoError(t, err)
	defer sub.Close()

	c, err := sub.Next(context.TODO())
	require.NoError(t, err)
	require.Equal(t, &TestConnector, c)
}

func TestDeletedBasicFlow(t *testing.T) {
	s := graphql.NewMockSubServer(t,
		graphql.AddVarNamesToQuery(connectorDeletedSubQuery, methodsVarName, allTenantsVarName),
		map[string]interface{}{methodsVarName: nil, allTenantsVarName: false},
		&connectorDeletedEvent{ConnectorEvent: &TestConnector},
	)

	defer s.Close()
	svc := New(s.URL)

	sub, err := svc.ConnectorDeleted(context.TODO(), nil, false)
	require.NoError(t, err)
	defer sub.Close()

	c, err := sub.Next(context.TODO())
	require.NoError(t, err)
	require.Equal(t, &TestConnector, c)
}
