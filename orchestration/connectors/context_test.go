package connectors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContext_LookupAction_No_Namespace(t *testing.T) {
	tenant := "tenant"
	testAction := ConnectorAction{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "testAction",
		Description: "test",
		Tags:        nil,
		Interface:   nil,
		Inputs:      nil,
		Outputs:     nil,
	}

	testConnectorInterface := ConnectorInterface{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Categories:  nil,
		Actions:     []*ConnectorAction{&testAction},
	}

	testConnector := Connector{
		ID:          "1",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "TestConnector2",
		Description: "Testing",
		Tags:        nil,
		Docs:        nil,
		Method:      nil,
		Implements:  []*ConnectorInterface{&testConnectorInterface},
		Actions:     nil,
		Parameters:  nil,
		AuthTypes:   nil,
		TenantID:    &tenant,
	}

	testConnection := Connection{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Connector:   &testConnector,
		Actions:     nil,
		AuthType:    AuthTypeNone,
		AuthURL:     nil,
		Config:      nil,
		Credentials: nil,
	}

	context := Context{
		TenantID:    tenant,
		UserID:      "123",
		Connections: []*Connection{&testConnection},
	}
	actionLookup := ActionLookup{
		ConnectorName:                "testConnector2",
		ConnectorInterfaceActionName: "testAction",
	}

	// Test lookup on connection name and action name
	connection, connectorAction := context.LookupAction(&actionLookup)

	assert.NotNil(t, connection)
	assert.NotNil(t, connectorAction)

	assert.Equal(t, "test", connection.Name)
	assert.Equal(t, "testAction", connectorAction.Name)

	actionLookup.ConnectorInterfaceActionName = ""
	// Test that a lookup with no action name returns nil
	connection, connectorAction = context.LookupAction(&actionLookup)

	assert.Nil(t, connection)
	assert.Nil(t, connectorAction)

	actionLookup.ConnectorInterfaceActionName = "testAction"

	// Test match on connection name and action
	connection, connectorAction = context.LookupAction(&actionLookup)

	assert.NotNil(t, connection)
	assert.NotNil(t, connectorAction)

	assert.Equal(t, "test", connection.Name)
	assert.Equal(t, "testAction", connectorAction.Name)
	actionLookup.ConnectorName = ""
	actionLookup.ConnectorInterfaceActionName = "blah"
	actionLookup.ImplementedConnectorInterfaceName = "testConnector2"
	// Test that matching connectorInterface name matches, but action name does not
	connection, connectorAction = context.LookupAction(&actionLookup)
	assert.Nil(t, connection)
	assert.Nil(t, connectorAction)
}

func TestContext_LookupAction_No_Namespace_Prioritizes_Tenant(t *testing.T) {
	tenant := "tenant"
	testAction := &ConnectorAction{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "testAction",
		Description: "test",
		Tags:        nil,
		Interface:   nil,
		Inputs:      nil,
		Outputs:     nil,
	}

	testGlobalConnectorInterface := &ConnectorInterface{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Categories:  nil,
		Actions:     []*ConnectorAction{testAction},
	}

	testGlobalConnector := &Connector{
		ID:          "1",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "TestConnector2",
		Description: "Testing",
		Tags:        nil,
		Docs:        nil,
		Method:      nil,
		Implements:  []*ConnectorInterface{testGlobalConnectorInterface},
		Actions:     nil,
		Parameters:  nil,
		AuthTypes:   nil,
	}

	testGlobalConnection := &Connection{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Connector:   testGlobalConnector,
		Actions:     nil,
		AuthType:    AuthTypeNone,
		AuthURL:     nil,
		Config:      nil,
		Credentials: nil,
	}

	testTenantConnectorInterface := &ConnectorInterface{
		ID:          "12345",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Categories:  nil,
		Actions:     []*ConnectorAction{testAction},
		TenantID:    &tenant,
	}

	testTenantConnector := &Connector{
		ID:          "1",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "TestConnector2",
		Description: "Testing",
		Tags:        nil,
		Docs:        nil,
		Method:      nil,
		Implements:  []*ConnectorInterface{testTenantConnectorInterface},
		Actions:     nil,
		Parameters:  nil,
		AuthTypes:   nil,
		TenantID:    &tenant,
	}

	testTenantConnection := &Connection{
		ID:          "12345",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Connector:   testTenantConnector,
		Actions:     nil,
		AuthType:    AuthTypeNone,
		AuthURL:     nil,
		Config:      nil,
		Credentials: nil,
	}

	context := Context{
		TenantID:    tenant,
		UserID:      "123",
		Connections: []*Connection{testGlobalConnection, testTenantConnection},
	}
	actionLookup := ActionLookup{
		ConnectorName:                "testConnector2",
		ConnectorInterfaceActionName: "testAction",
	}

	// Test lookup on connection name and action name
	connection, connectorAction := context.LookupAction(&actionLookup)

	assert.NotNil(t, connection)
	assert.NotNil(t, connectorAction)

	assert.Equal(t, "test", connection.Name)
	assert.Equal(t, "testAction", connectorAction.Name)
	assert.Equal(t, testTenantConnection.ID, connection.ID)

	actionLookup.ConnectorInterfaceActionName = ""
	// Test that a lookup with no action name returns nil
	connection, connectorAction = context.LookupAction(&actionLookup)

	assert.Nil(t, connection)
	assert.Nil(t, connectorAction)

	actionLookup.ConnectorInterfaceActionName = "testAction"

	// Test match on connection name and action
	connection, connectorAction = context.LookupAction(&actionLookup)

	assert.NotNil(t, connection)
	assert.NotNil(t, connectorAction)

	assert.Equal(t, "test", connection.Name)
	assert.Equal(t, "testAction", connectorAction.Name)
	assert.Equal(t, testTenantConnection.ID, connection.ID)

	actionLookup.ConnectorName = ""
	actionLookup.ConnectorInterfaceActionName = "blah"
	actionLookup.ImplementedConnectorInterfaceName = "testConnector2"
	// Test that matching connectorInterface name matches, but action name does not
	connection, connectorAction = context.LookupAction(&actionLookup)
	assert.Nil(t, connection)
	assert.Nil(t, connectorAction)
}

func TestContext_LookupAction_GlobalTenant(t *testing.T) {
	testAction := ConnectorAction{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "testAction",
		Description: "test",
		Tags:        nil,
		Interface:   nil,
		Inputs:      nil,
		Outputs:     nil,
	}

	testConnectorInterface := ConnectorInterface{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Categories:  nil,
		Actions:     []*ConnectorAction{&testAction},
	}

	testConnector := Connector{
		ID:          "1",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "TestConnector2",
		Description: "Testing",
		Tags:        nil,
		Docs:        nil,
		Method:      nil,
		Implements:  []*ConnectorInterface{&testConnectorInterface},
		Actions:     nil,
		Parameters:  nil,
		AuthTypes:   nil,
	}

	testConnection := Connection{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Connector:   &testConnector,
		Actions:     nil,
		AuthType:    AuthTypeNone,
		AuthURL:     nil,
		Config:      nil,
		Credentials: nil,
	}

	context := Context{
		TenantID:    "1234",
		UserID:      "123",
		Connections: []*Connection{&testConnection},
	}
	actionLookup := ActionLookup{
		ConnectorName:                "testConnector2",
		ConnectorInterfaceActionName: "testAction",
		Namespace:                    GlobalNamespace,
	}

	// Test lookup on connection name and action name
	connection, connectorAction := context.LookupAction(&actionLookup)

	assert.NotNil(t, connection)
	assert.NotNil(t, connectorAction)

	assert.Equal(t, "test", connection.Name)
	assert.Equal(t, "testAction", connectorAction.Name)

	actionLookup.ConnectorInterfaceActionName = ""
	// Test that a lookup with no action name returns nil
	connection, connectorAction = context.LookupAction(&actionLookup)

	assert.Nil(t, connection)
	assert.Nil(t, connectorAction)

	actionLookup.ConnectorInterfaceActionName = "testAction"

	// Test match on connection name and action
	connection, connectorAction = context.LookupAction(&actionLookup)

	assert.NotNil(t, connection)
	assert.NotNil(t, connectorAction)

	assert.Equal(t, "test", connection.Name)
	assert.Equal(t, "testAction", connectorAction.Name)
	actionLookup.ConnectorName = ""
	actionLookup.ConnectorInterfaceActionName = "blah"
	actionLookup.ImplementedConnectorInterfaceName = "testConnector2"
	// Test that matching connectorInterface name matches, but action name does not
	connection, connectorAction = context.LookupAction(&actionLookup)
	assert.Nil(t, connection)
	assert.Nil(t, connectorAction)
}

func TestContext_LookupAction_Tenant_Namespace(t *testing.T) {
	tenant := "123456789"
	testAction := ConnectorAction{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "testAction",
		Description: "test",
		Tags:        nil,
		Interface:   nil,
		Inputs:      nil,
		Outputs:     nil,
	}

	testConnectorInterface := ConnectorInterface{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		TenantID:    &tenant,
		Tags:        nil,
		Categories:  nil,
		Actions:     []*ConnectorAction{&testAction},
	}

	testConnector := Connector{
		ID:          "1",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "TestConnector2",
		Description: "Testing",
		Tags:        nil,
		Docs:        nil,
		Method:      nil,
		Implements:  []*ConnectorInterface{&testConnectorInterface},
		Actions:     nil,
		Parameters:  nil,
		AuthTypes:   nil,
		TenantID:    &tenant,
	}

	testConnection := Connection{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Connector:   &testConnector,
		Actions:     nil,
		AuthType:    AuthTypeNone,
		AuthURL:     nil,
		Config:      nil,
		Credentials: nil,
	}

	context := Context{
		TenantID:    tenant,
		UserID:      "123",
		Connections: []*Connection{&testConnection},
	}
	actionLookup := ActionLookup{
		ConnectorName:                "TestConnector2",
		ConnectorInterfaceActionName: "testAction",
		Namespace:                    GlobalNamespace,
	}

	// Test lookup on connection name and action name
	connection, connectorAction := context.LookupAction(&actionLookup)

	assert.Nil(t, connection)
	assert.Nil(t, connectorAction)

	actionLookup.Namespace = tenant

	// Test lookup on connection name and action name
	connection, connectorAction = context.LookupAction(&actionLookup)

	assert.NotNil(t, connection)
	assert.NotNil(t, connectorAction)
	assert.Equal(t, "test", connection.Name)
	assert.Equal(t, "testAction", connectorAction.Name)

	actionLookup.ConnectorName = ""
	actionLookup.ImplementedConnectorInterfaceName = "test"

	// Test that matching connectorInterface name matches, but action name does not
	connection, connectorAction = context.LookupAction(&actionLookup)
	assert.NotNil(t, connection)
	assert.NotNil(t, connectorAction)
	assert.Equal(t, "test", connection.Name)
	assert.Equal(t, "testAction", connectorAction.Name)

	actionLookup.Namespace = GlobalNamespace

	// Test lookup on connection name and action name
	connection, connectorAction = context.LookupAction(&actionLookup)
	assert.Nil(t, connection)
	assert.Nil(t, connectorAction)
}
