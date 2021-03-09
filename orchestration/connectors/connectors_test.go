package connectors

import (
	"testing"
	"time"

	"github.com/secureworks/taegis-sdk-go/client"

	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	tenantId            = "1"
	testConnectorAction = ConnectorAction{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "testAction",
		Description: "testing actions",
		Tags:        nil,
		Interface:   nil,
		Inputs:      nil,
		Outputs:     nil,
	}
	testConnectorInterface = ConnectorInterface{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Categories:  nil,
		Actions:     []*ConnectorAction{&testConnectorAction},
	}

	testConnectorActionDefinition = ConnectorActionDefinition{
		ID:        "1234",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Name:      "test",
		Action:    &testConnectorAction,
	}

	TestConnector = Connector{
		ID:          "1",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "TestConnector",
		Description: "Testing",
		Tags:        nil,
		Docs:        nil,
		Method:      nil,
		Implements:  []*ConnectorInterface{&testConnectorInterface},
		Actions:     []*ConnectorActionDefinition{&testConnectorActionDefinition},
		Parameters:  nil,
		AuthTypes:   []AuthType{AuthTypeAPIKey},
	}

	testConnection = Connection{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Connector:   &TestConnector,
		Actions:     nil,
		AuthType:    AuthTypeClientCerts,
		AuthURL:     nil,
		Config:      nil,
		Credentials: nil,
	}

	testConnectionMethod = ConnectionMethod{
		ID:          "123",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Parameters:  nil,
		GraphQLURL:  "",
	}

	testConnectorCategory = ConnectorCategory{
		ID:          "123",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
	}

	header = testutils.CreateHeader()
)

func TestConnectorSvc_GetConnectionMethod(t *testing.T) {
	testConnectionMethod := &ConnectionMethod{
		ID:          "123",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Parameters:  nil,
	}

	r := struct {
		Out *ConnectionMethod `json:"connectionMethod"`
	}{
		Out: testConnectionMethod,
	}
	m := testutils.NewMockGQLOutput(t, header, r)
	defer m.Close()

	c := New(m.URL)

	connectionMethod, err := c.GetConnectionMethod("test")

	assert.Nil(t, err)
	assert.Equal(t, "test", connectionMethod.Name)
}

func TestConnectorSvc_GetConnectionMethodErr(t *testing.T) {

	m := testutils.NewMockGQLError(t, header)
	defer m.Close()

	c := New(m.URL)

	connectionMethod, err := c.GetConnectionMethod("blah")

	assert.NotNil(t, err)
	assert.Nil(t, connectionMethod)
}

func TestConnectorSvc_GetConnections(t *testing.T) {
	r := struct {
		Out []*Connection `json:"connections"`
	}{
		Out: []*Connection{&testConnection},
	}

	m := testutils.NewMockGQLOutput(t, header, r)
	defer m.Close()

	c := New(m.URL)
	input := GetConnectionsInput{
		ConnectionIDs:         []string{"1234"},
		ConnectorIDs:          []string{"1234"},
		ConnectorInterfaceIDs: []string{"1234"},
	}
	connections, err := c.GetConnections(&input)
	assert.Nil(t, err)
	assert.Len(t, connections, 1)
	assert.Equal(t, "test", connections[0].Name)
}

func TestConnectorSvc_GetConnectionsErr(t *testing.T) {
	m := testutils.NewMockGQLError(t, header)
	defer m.Close()

	c := New(m.URL)

	input := GetConnectionsInput{
		ConnectionIDs:         []string{"test"},
		ConnectorIDs:          []string{"123"},
		ConnectorInterfaceIDs: []string{"blah"},
	}
	connections, err := c.GetConnections(&input)

	assert.NotNil(t, err)
	assert.Nil(t, connections)
}

func TestConnectorSvc_GetConnectors(t *testing.T) {
	r := struct {
		Out []*Connector `json:"connectors"`
	}{
		Out: []*Connector{&TestConnector},
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	input := &GetConnectorsInput{
		ConnectionMethodIDs:   []string{"1"},
		ConnectorInterfaceIDs: []string{"1234"},
	}
	connectors, err := c.GetConnectors(input)
	assert.Nil(t, err)
	assert.Len(t, connectors, 1)
	assert.Equal(t, "TestConnector", connectors[0].Name)
}

func TestConnectorSvc_GetConnectorsErr(t *testing.T) {
	m := testutils.NewMockGQLError(t, header)
	defer m.Close()

	c := New(m.URL)
	input := &GetConnectorsInput{
		ConnectionMethodIDs: []string{"123"},
	}
	connectors, err := c.GetConnectors(input)

	assert.NotNil(t, err)
	assert.Nil(t, connectors)
}

func TestConnectorSvc_GetContext(t *testing.T) {
	testContext := Context{
		TenantID:    "1",
		UserID:      "1",
		Connections: []*Connection{&testConnection},
	}

	r := struct {
		Out *Context `json:"context"`
	}{
		Out: &testContext,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL, client.WithTenant(testContext.TenantID))
	input := GetConnectionsInput{
		ConnectionIDs:         []string{"1234"},
		ConnectorIDs:          []string{"1"},
		ConnectorInterfaceIDs: []string{"1234"},
	}
	context, err := c.GetContext(&input)

	assert.Nil(t, err)
	assert.Equal(t, "1", context.TenantID)
}

func TestConnectorSvc_GetContextErr(t *testing.T) {
	m := testutils.NewMockGQLError(t, header)
	defer m.Close()

	c := New(m.URL, client.WithTenant("test"))
	input := GetConnectionsInput{
		ConnectionIDs:         []string{"test"},
		ConnectorIDs:          []string{"123"},
		ConnectorInterfaceIDs: []string{"blah"},
	}
	context, err := c.GetContext(&input)

	assert.NotNil(t, err)
	assert.Nil(t, context)
}

func TestConnectorSvc_CreateConnection(t *testing.T) {

	r := struct {
		Out *Connection `json:"createConnection"`
	}{
		Out: &testConnection,
	}

	input := ConnectionInput{
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Config:      nil,
		Credentials: nil,
		AuthType:    AuthTypeNone,
		AuthURL:     nil,
		Actions:     nil,
	}
	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	connection, err := c.CreateConnection("123", &input)
	assert.Nil(t, err)
	assert.Equal(t, "1234", connection.ID)
}

func TestConnectorSvc_CreateConnector(t *testing.T) {
	r := struct {
		Out *Connector `json:"createConnector"`
	}{
		Out: &TestConnector,
	}
	objMetaInput := common.ObjectMetaInput{
		Name:        "test",
		Description: nil,
		Tags:        nil,
	}
	input := ConnectorInput{
		ObjectMetaInput: objMetaInput,
		Implements:      nil,
		Parameters:      nil,
		AuthTypes:       []AuthType{AuthTypeNone},
		AllTenants:      nil,
		Actions:         nil,
		Documentation:   nil,
		Categories:      nil,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	connector, err := c.CreateConnector("123", &input)

	assert.Nil(t, err)
	assert.Equal(t, "1", connector.ID)
}

func TestConnectorSvc_CreateConnectorInterface(t *testing.T) {
	r := struct {
		Out *ConnectorInterface `json:"createConnectorInterface"`
	}{
		Out: &testConnectorInterface,
	}
	objMetaInput := common.ObjectMetaInput{
		Name:        "test",
		Description: nil,
		Tags:        nil,
	}
	input := ConnectorInterfaceInput{
		ObjectMetaInput: objMetaInput,
		Categories:      nil,
		Actions:         nil,
		AllTenants:      nil,
	}
	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	i, err := c.CreateConnectorInterface(&input)
	assert.Nil(t, err)
	assert.Equal(t, "1234", i.ID)
}

func TestConnectorSvc_DefineConnectionMethod(t *testing.T) {
	r := struct {
		Out *ConnectionMethod `json:"defineConnectionMethod"`
	}{
		Out: &testConnectionMethod,
	}
	objMetaInput := common.ObjectMetaInput{
		Name:        "test",
		Description: nil,
		Tags:        nil,
	}
	input := ConnectionMethodInput{
		ObjectMetaInput: objMetaInput,
		URL:             "",
		parameters:      nil,
	}
	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	i, err := c.DefineConnectionMethod(&input)

	assert.Nil(t, err)

	assert.Equal(t, "123", i.ID)

}

func TestConnectorSvc_DeleteConnection(t *testing.T) {
	r := struct {
		Out *Connection `json:"deleteConnection"`
	}{
		Out: &testConnection,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	con, err := c.DeleteConnection("1234")

	assert.Nil(t, err)
	assert.Equal(t, "1234", con.ID)
}

func TestConnectorSvc_DeleteConnector(t *testing.T) {
	r := struct {
		Out *Connector `json:"deleteConnector"`
	}{
		Out: &TestConnector,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	connector, err := c.DeleteConnector("1")

	assert.Nil(t, err)
	assert.Equal(t, "1", connector.ID)
}

func TestConnectorSvc_UpdateConnectorInterface(t *testing.T) {
	r := struct {
		Out *ConnectorInterface `json:"updateConnectorInterface"`
	}{
		Out: &testConnectorInterface,
	}
	objMetaInput := common.ObjectMetaInput{
		Name:        "test",
		Description: nil,
		Tags:        nil,
	}
	input := ConnectorInterfaceInput{
		ObjectMetaInput: objMetaInput,
		Categories:      nil,
		Actions:         nil,
		AllTenants:      nil,
	}
	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	i, err := c.UpdateConnectorInterface("1234", &input)
	assert.Nil(t, err)

	assert.Equal(t, "1234", i.ID)

}

func TestConnectorSvc_DeleteConnectorInterface(t *testing.T) {
	r := struct {
		Out *ConnectorInterface `json:"deleteConnectorInterface"`
	}{
		Out: &testConnectorInterface,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	i, err := c.DeleteConnectorInterface("1234")
	assert.Nil(t, err)
	assert.Equal(t, "1234", i.ID)

}

func TestConnectorSvc_ExecuteConnectionAction(t *testing.T) {

	r := struct {
		Out *ConnectorAction `json:"executeConnectionAction"`
	}{
		Out: &testConnectorAction,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	a, err := c.ExecuteConnectionAction("123", "test", nil)
	assert.Nil(t, err)

	assert.NotNil(t, a)
}

func TestConnectorSvc_GetConnection(t *testing.T) {
	r := struct {
		Out *Connection `json:"connection"`
	}{
		Out: &testConnection,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	con, err := c.GetConnection("1234")
	assert.Nil(t, err)
	assert.Equal(t, "1234", con.ID)
}

func TestConnectorSvc_GetConnector(t *testing.T) {
	r := struct {
		Out *Connector `json:"connector"`
	}{
		Out: &TestConnector,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	con, err := c.GetConnector("1234")
	assert.Nil(t, err)
	assert.Equal(t, "1", con.ID)
}

func TestConnectorSvc_GetConnectorCategory(t *testing.T) {
	r := struct {
		Out *ConnectorCategory `json:"connectorCategory"`
	}{
		Out: &testConnectorCategory,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	con, err := c.GetConnectorCategory("123")
	assert.Nil(t, err)
	assert.Equal(t, "123", con.ID)
}

func TestConnectorSvc_GetConnectorInterface(t *testing.T) {
	r := struct {
		Out *ConnectorInterface `json:"connectorInterface"`
	}{
		Out: &testConnectorInterface,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	con, err := c.GetConnectorInterface("1234")
	assert.Nil(t, err)
	assert.Equal(t, "1234", con.ID)
}

func TestConnectorSvc_RemoveConnectionMethod(t *testing.T) {
	r := struct {
		Out *ConnectionMethod `json:"removeConnectionMethod"`
	}{
		Out: &testConnectionMethod,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	con, err := c.RemoveConnectionMethod("123")
	assert.Nil(t, err)
	assert.Equal(t, "123", con.ID)
}

func TestConnectorSvc_UpdateConnection(t *testing.T) {
	r := struct {
		Out *Connection `json:"updateConnection"`
	}{
		Out: &testConnection,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()
	input := ConnectionInput{
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Config:      nil,
		Credentials: nil,
		AuthType:    AuthTypeNone,
		AuthURL:     nil,
		Actions:     nil,
	}

	c := New(s.URL)
	con, err := c.UpdateConnection("1234", &input)
	assert.Nil(t, err)

	assert.Equal(t, "1234", con.ID)
}

func TestConnectorSvc_UpdateConnector(t *testing.T) {
	r := struct {
		Out *Connector `json:"updateConnector"`
	}{
		Out: &TestConnector,
	}
	input := ConnectorUpdateInput{
		Name:          "test",
		Description:   "test",
		Implements:    nil,
		parameters:    nil,
		AuthTypes:     nil,
		Documentation: "test",
		Actions:       nil,
		Categories:    nil,
		Tags:          nil,
	}
	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	con, err := c.UpdateConnector("1", &input)
	assert.Nil(t, err)
	assert.Equal(t, "1", con.ID)
}

func TestConnectorSvc_ValidateConnection(t *testing.T) {
	r := struct {
		Out *Connection `json:"validateConnection"`
	}{
		Out: &testConnection,
	}
	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	con, err := c.ValidateConnection("1234")
	assert.Nil(t, err)
	assert.Equal(t, "1234", con.ID)
}

func TestConnectorSvc_ValidateConnectionInput(t *testing.T) {
	r := struct {
		Out *Connector `json:"validateConnectionInput"`
	}{
		Out: &TestConnector,
	}

	input := ConnectionInput{
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Config:      nil,
		Credentials: nil,
		AuthType:    AuthTypeNone,
		AuthURL:     nil,
		Actions:     nil,
	}
	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	con, err := c.ValidateConnectionInput("1", &input)
	assert.Nil(t, err)
	assert.Equal(t, "1", con.ID)
}

func TestNew(t *testing.T) {
	c := New("testurl")
	assert.NotNil(t, c)
}
