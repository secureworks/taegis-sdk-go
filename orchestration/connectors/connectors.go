package connectors

import (
	"context"

	"github.com/secureworks/taegis-sdk-go/client"
	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/graphql"
)

type Service interface {
	GetConnectionMethod(name string, opts ...graphql.RequestOption) (*ConnectionMethod, error)
	GetConnectors(input *GetConnectorsInput, opts ...graphql.RequestOption) ([]*Connector, error)
	GetConnections(input *GetConnectionsInput, opts ...graphql.RequestOption) ([]*Connection, error)
	GetContext(input *GetConnectionsInput, opts ...graphql.RequestOption) (*Context, error)
	DefineConnectionMethod(input *ConnectionMethodInput, opts ...graphql.RequestOption) (*ConnectionMethod, error)
	RemoveConnectionMethod(connectionID string, opts ...graphql.RequestOption) (*ConnectionMethod, error)
	CreateConnectorInterface(input *ConnectorInterfaceInput, opts ...graphql.RequestOption) (*ConnectorInterface, error)
	UpdateConnectorInterface(connectorID string, input *ConnectorInterfaceInput, opts ...graphql.RequestOption) (*ConnectorInterface, error)
	DeleteConnectorInterface(connectorID string, opts ...graphql.RequestOption) (*ConnectorInterface, error)
	CreateConnector(connectionMethodID string, input *ConnectorInput, opts ...graphql.RequestOption) (*Connector, error)
	UpdateConnector(connectorID string, input *ConnectorUpdateInput, opts ...graphql.RequestOption) (*Connector, error)
	DeleteConnector(connectorID string, opts ...graphql.RequestOption) (*Connector, error)
	CreateConnection(connectorID string, input *ConnectionInput, opts ...graphql.RequestOption) (*Connection, error)
	UpdateConnection(connectionID string, input *ConnectionInput, opts ...graphql.RequestOption) (*Connection, error)
	DeleteConnection(connectionID string, opts ...graphql.RequestOption) (*Connection, error)
	ExecuteConnectionAction(connectionID string, actionName string, inputs interface{}, opts ...graphql.RequestOption) (*interface{}, error)
	GetConnectorCategory(connectorCategoryID string, opts ...graphql.RequestOption) (*ConnectorCategory, error)
	GetConnectorInterface(connectorInterfaceID string, opts ...graphql.RequestOption) (*ConnectorInterface, error)
	GetConnector(connectorID string, opts ...graphql.RequestOption) (*Connector, error)
	GetConnection(connectionID string, opts ...graphql.RequestOption) (*Connection, error)
	ConnectorCreated(ctx context.Context, connectorMethods common.IDs, allTenants bool, options ...graphql.SubscriptionOption) (Subscription, error)
	ConnectorDeleted(ctx context.Context, connectorMethods common.IDs, allTenants bool, options ...graphql.SubscriptionOption) (Subscription, error)
	ConnectorUpdated(ctx context.Context, connectorMethods common.IDs, allTenants bool, options ...graphql.SubscriptionOption) (Subscription, error)
}

var _ Service = (*connectorSvc)(nil)

const (
	idVarName    = "id"
	inputVarName = "input"
)

type connectorSvc struct {
	client *client.Client
	url    string
}

func New(url string, opts ...client.Option) *connectorSvc {
	client := client.NewClient(opts...)
	return &connectorSvc{
		client: client,
		url:    url,
	}
}

func (o *connectorSvc) GetConnectionMethod(name string, opts ...graphql.RequestOption) (*ConnectionMethod, error) {
	const nameVar = "name"
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(connectionMethodQuery, nameVar), opts...)
	req.Var(nameVar, name)

	var data struct {
		ConnectionMethod ConnectionMethod `json:"connectionMethod"`
	}
	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return &data.ConnectionMethod, nil
}

func (o *connectorSvc) GetConnectors(input *GetConnectorsInput, opts ...graphql.RequestOption) ([]*Connector, error) {
	const connectorMethodIDsVar = "connectorMethodIDs"
	const connectorInterfaceIDsVar = "connectorInterfaceIDs"
	const connectorCategoryIDsVar = "connectorCategoryIDs"
	const tagsVar = "tags"
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(connectorsQuery, connectorMethodIDsVar, connectorInterfaceIDsVar, connectorCategoryIDsVar, tagsVar), opts...)
	req.Var(connectorMethodIDsVar, input.ConnectionMethodIDs)
	req.Var(connectorInterfaceIDsVar, input.ConnectorInterfaceIDs)
	req.Var(connectorCategoryIDsVar, input.ConnectorCategoryIDs)
	req.Var(tagsVar, input.Tags)

	var data struct {
		Connectors []*Connector `json:"connectors"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}
	return data.Connectors, nil
}

func (o *connectorSvc) GetConnections(input *GetConnectionsInput, opts ...graphql.RequestOption) ([]*Connection, error) {
	const connectionIDsVar = "connectionIDs"
	const connectorIDsVar = "connectorIDs"
	const connectorInterfaceIDsVar = "connectorInterfaceIDs"

	req := graphql.NewRequest(graphql.AddVarNamesToQuery(connectionsQuery, connectionIDsVar, connectorIDsVar, connectorInterfaceIDsVar), opts...)
	req.Var(connectionIDsVar, input.ConnectionIDs)
	req.Var(connectorIDsVar, input.ConnectorIDs)
	req.Var(connectorInterfaceIDsVar, input.ConnectorInterfaceIDs)

	var data struct {
		Connections []*Connection `json:"connections"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}
	return data.Connections, nil
}

func (o *connectorSvc) GetContext(input *GetConnectionsInput, opts ...graphql.RequestOption) (*Context, error) {
	const connectionIDsVar = "connectionIDs"
	const connectorIDsVar = "connectorIDs"
	const connectorInterfaceIDsVar = "connectorInterfaceIDs"

	req := graphql.NewRequest(graphql.AddVarNamesToQuery(getContextQuery, connectionIDsVar, connectorIDsVar, connectorInterfaceIDsVar), opts...)
	req.Var(connectionIDsVar, input.ConnectionIDs)
	req.Var(connectorIDsVar, input.ConnectorIDs)
	req.Var(connectorInterfaceIDsVar, input.ConnectorInterfaceIDs)

	var data struct {
		Connections []*Connection `json:"connections"`
	}
	tenantID, err := graphql.ExecuteQueryWithTenant(o.client, o.url, req, &data)
	if err != nil {
		return nil, err
	}

	return &Context{TenantID: tenantID, Connections: data.Connections}, nil
}

func (o *connectorSvc) DefineConnectionMethod(input *ConnectionMethodInput, opts ...graphql.RequestOption) (*ConnectionMethod, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(defineConnectionMethodMutation, inputVarName), opts...)
	req.Var(inputVarName, input)

	var data struct {
		ConnectionMethod *ConnectionMethod `json:"defineConnectionMethod"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.ConnectionMethod, nil
}

func (o *connectorSvc) RemoveConnectionMethod(connectionID string, opts ...graphql.RequestOption) (*ConnectionMethod, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(removeConnectionMethodMutation, idVarName), opts...)
	req.Var(idVarName, connectionID)

	var data struct {
		ConnectionMethod *ConnectionMethod `json:"removeConnectionMethod"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.ConnectionMethod, nil
}

func (o *connectorSvc) CreateConnectorInterface(input *ConnectorInterfaceInput, opts ...graphql.RequestOption) (*ConnectorInterface, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(createConnectorInterfaceMutation, inputVarName), opts...)
	req.Var(inputVarName, input)

	var data struct {
		ConnectorInterface *ConnectorInterface `json:"createConnectorInterface"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.ConnectorInterface, nil
}

func (o *connectorSvc) UpdateConnectorInterface(connectorInterfaceID string, input *ConnectorInterfaceInput, opts ...graphql.RequestOption) (*ConnectorInterface, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(updateConnectorInterfaceMutation, idVarName, inputVarName), opts...)
	req.Var(idVarName, connectorInterfaceID)
	req.Var(inputVarName, input)

	var data struct {
		ConnectorInterface *ConnectorInterface `json:"updateConnectorInterface"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.ConnectorInterface, nil
}

func (o *connectorSvc) DeleteConnectorInterface(connectorInterfaceID string, opts ...graphql.RequestOption) (*ConnectorInterface, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(deleteConnectorInterfaceMutation, idVarName), opts...)
	req.Var(idVarName, connectorInterfaceID)

	var data struct {
		ConnectorInterface *ConnectorInterface `json:"deleteConnectorInterface"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.ConnectorInterface, nil
}

func (o *connectorSvc) CreateConnector(connectionMethodID string, input *ConnectorInput, opts ...graphql.RequestOption) (*Connector, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(createConnectorMutation, idVarName, inputVarName), opts...)
	req.Var(idVarName, connectionMethodID)
	req.Var(inputVarName, input)

	var data struct {
		Connector *Connector `json:"createConnector"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.Connector, nil
}

func (o *connectorSvc) UpdateConnector(connectorID string, input *ConnectorUpdateInput, opts ...graphql.RequestOption) (*Connector, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(updateConnectorMutation, idVarName, inputVarName), opts...)
	req.Var(idVarName, connectorID)
	req.Var(inputVarName, input)

	var data struct {
		Connector *Connector `json:"updateConnector"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.Connector, nil
}

func (o *connectorSvc) DeleteConnector(connectorID string, opts ...graphql.RequestOption) (*Connector, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(deleteConnectorMutation, idVarName), opts...)
	req.Var(idVarName, connectorID)

	var data struct {
		Connector *Connector `json:"deleteConnector"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.Connector, nil
}

func (o *connectorSvc) CreateConnection(connectorID string, input *ConnectionInput, opts ...graphql.RequestOption) (*Connection, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(createConnectionMutation, idVarName, inputVarName), opts...)
	req.Var(idVarName, connectorID)
	req.Var(inputVarName, input)

	var data struct {
		Connection *Connection `json:"createConnection"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.Connection, nil
}

func (o *connectorSvc) UpdateConnection(connectionID string, input *ConnectionInput, opts ...graphql.RequestOption) (*Connection, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(updateConnectionMutation, idVarName, inputVarName), opts...)
	req.Var(idVarName, connectionID)
	req.Var(inputVarName, input)

	var data struct {
		Connection *Connection `json:"updateConnection"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.Connection, nil
}

func (o *connectorSvc) DeleteConnection(connectionID string, opts ...graphql.RequestOption) (*Connection, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(deleteConnectionMutation, idVarName), opts...)
	req.Var(idVarName, connectionID)

	var data struct {
		Connection *Connection `json:"deleteConnection"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.Connection, nil
}

func (o *connectorSvc) ValidateConnection(connectionID string, opts ...graphql.RequestOption) (*Connection, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(validateConnectionMutation, idVarName), opts...)
	req.Var(idVarName, connectionID)

	var data struct {
		Connection *Connection `json:"validateConnection"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.Connection, nil
}

func (o *connectorSvc) ValidateConnectionInput(connectorID string, input *ConnectionInput, opts ...graphql.RequestOption) (*Connector, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(validateConnectionInputMutation, idVarName, inputVarName), opts...)
	req.Var(idVarName, connectorID)
	req.Var(inputVarName, input)

	var data struct {
		Connector *Connector `json:"validateConnectionInput"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.Connector, nil
}

func (o *connectorSvc) ExecuteConnectionAction(connectionID string, actionName string, inputs interface{}, opts ...graphql.RequestOption) (*interface{}, error) {
	const actionNameVar = "action_name"

	req := graphql.NewRequest(graphql.AddVarNamesToQuery(executeConnectionActionMutation, idVarName, actionNameVar, inputVarName), opts...)
	req.Var(idVarName, connectionID)
	req.Var(actionNameVar, actionName)
	req.Var(inputVarName, inputs)

	var data struct {
		Interface *interface{} `json:"executeConnectionAction"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.Interface, nil
}

func (o *connectorSvc) GetConnectorCategory(connectorCategoryID string, opts ...graphql.RequestOption) (*ConnectorCategory, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(connectorCategoryQuery, idVarName), opts...)
	req.Var(idVarName, connectorCategoryID)

	var data struct {
		ConnectorCategory *ConnectorCategory `json:"connectorCategory"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.ConnectorCategory, nil
}

func (o *connectorSvc) GetConnectorInterface(connectorInterfaceID string, opts ...graphql.RequestOption) (*ConnectorInterface, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(connectorInterfaceQuery, idVarName), opts...)
	req.Var(idVarName, connectorInterfaceID)

	var data struct {
		ConnectorInterface *ConnectorInterface `json:"connectorInterface"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.ConnectorInterface, nil
}

func (o *connectorSvc) GetConnector(connectorID string, opts ...graphql.RequestOption) (*Connector, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(connectorQuery, idVarName), opts...)
	req.Var(idVarName, connectorID)

	var data struct {
		Connector *Connector `json:"connector"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.Connector, nil
}

func (o *connectorSvc) GetConnection(connectionID string, opts ...graphql.RequestOption) (*Connection, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(connectionQuery, idVarName), opts...)
	req.Var(idVarName, connectionID)

	var data struct {
		Connection *Connection `json:"connection"`
	}

	if err := graphql.ExecuteQuery(o.client, o.url, req, &data); err != nil {
		return nil, err
	}

	return data.Connection, nil
}
