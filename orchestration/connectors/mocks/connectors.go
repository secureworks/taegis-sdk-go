package mocks

import (
	"context"

	"github.com/secureworks/taegis-sdk-go/graphql"

	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/orchestration/connectors"
)

var _ connectors.Service = (*Service)(nil)

type Service struct {
	GetConnectionMethodError       error
	GetConnectorsError             error
	GetConnectionsError            error
	GetContextError                error
	DefineConnectionMethodError    error
	RemoveConnectionMethodError    error
	CreateConnectorInterfaceError  error
	UpdateConnectorInterfaceError  error
	DeleteConnectorInterfaceError  error
	CreateConnectorError           error
	UpdateConnectorError           error
	DeleteConnectorError           error
	CreateConnectionError          error
	UpdateConnectionError          error
	DeleteConnectionError          error
	ExecuteConnectionActionError   error
	GetConnectorCategoryError      error
	GetConnectorInterfaceError     error
	GetConnectorError              error
	GetConnectionError             error
	ConnectorCreatedError          error
	ConnectorUpdatedError          error
	ConnectorDeletedError          error
	GetConnectionMethodResult      *connectors.ConnectionMethod
	GetConnectorsResult            []*connectors.Connector
	GetConnectionsResult           []*connectors.Connection
	GetContextResult               *connectors.Context
	DefineConnectionMethodResult   *connectors.ConnectionMethod
	RemoveConnectionMethodResult   *connectors.ConnectionMethod
	CreateConnectorInterfaceResult *connectors.ConnectorInterface
	UpdateConnectorInterfaceResult *connectors.ConnectorInterface
	DeleteConnectorInterfaceResult *connectors.ConnectorInterface
	CreateConnectorResult          *connectors.Connector
	UpdateConnectorResult          *connectors.Connector
	DeleteConnectorResult          *connectors.Connector
	CreateConnectionResult         *connectors.Connection
	UpdateConnectionResult         *connectors.Connection
	DeleteConnectionResult         *connectors.Connection
	ExecuteConnectionActionResult  *interface{}
	GetConnectorCategoryResult     *connectors.ConnectorCategory
	GetConnectorInterfaceResult    *connectors.ConnectorInterface
	GetConnectorResult             *connectors.Connector
	GetConnectionResult            *connectors.Connection
	ConnectorCreatedResult         chan *connectors.Connector
	ConnectorUpdatedResult         chan *connectors.Connector
	ConnectorDeletedResult         chan *connectors.Connector
}

func (m *Service) GetConnectionMethod(_ string, _ ...graphql.RequestOption) (*connectors.ConnectionMethod, error) {
	return m.GetConnectionMethodResult, m.GetConnectionMethodError
}

func (m *Service) GetConnectors(_ *connectors.GetConnectorsInput, _ ...graphql.RequestOption) ([]*connectors.Connector, error) {
	return m.GetConnectorsResult, m.GetConnectorsError
}

func (m *Service) GetConnections(_ *connectors.GetConnectionsInput, _ ...graphql.RequestOption) ([]*connectors.Connection, error) {
	return m.GetConnectionsResult, m.GetConnectionsError
}

func (m *Service) GetContext(_ *connectors.GetConnectionsInput, _ ...graphql.RequestOption) (*connectors.Context, error) {
	return m.GetContextResult, m.GetContextError
}

func (m *Service) DefineConnectionMethod(_ *connectors.ConnectionMethodInput, _ ...graphql.RequestOption) (*connectors.ConnectionMethod, error) {
	return m.DefineConnectionMethodResult, m.DefineConnectionMethodError
}

func (m *Service) RemoveConnectionMethod(_ string, _ ...graphql.RequestOption) (*connectors.ConnectionMethod, error) {
	return m.RemoveConnectionMethodResult, m.RemoveConnectionMethodError
}

func (m *Service) CreateConnectorInterface(_ *connectors.ConnectorInterfaceInput, _ ...graphql.RequestOption) (*connectors.ConnectorInterface, error) {
	return m.CreateConnectorInterfaceResult, m.CreateConnectorInterfaceError
}

func (m *Service) UpdateConnectorInterface(_ string, _ *connectors.ConnectorInterfaceInput, _ ...graphql.RequestOption) (*connectors.ConnectorInterface, error) {
	return m.UpdateConnectorInterfaceResult, m.UpdateConnectorInterfaceError
}

func (m *Service) DeleteConnectorInterface(_ string, _ ...graphql.RequestOption) (*connectors.ConnectorInterface, error) {
	return m.DeleteConnectorInterfaceResult, m.DeleteConnectorInterfaceError
}

func (m *Service) CreateConnector(_ string, _ *connectors.ConnectorInput, _ ...graphql.RequestOption) (*connectors.Connector, error) {
	return m.CreateConnectorResult, m.CreateConnectorError
}

func (m *Service) UpdateConnector(_ string, _ *connectors.ConnectorUpdateInput, _ ...graphql.RequestOption) (*connectors.Connector, error) {
	return m.UpdateConnectorResult, m.UpdateConnectorError
}

func (m *Service) DeleteConnector(_ string, _ ...graphql.RequestOption) (*connectors.Connector, error) {
	return m.DeleteConnectorResult, m.DeleteConnectorError
}

func (m *Service) CreateConnection(_ string, _ *connectors.ConnectionInput, _ ...graphql.RequestOption) (*connectors.Connection, error) {
	return m.CreateConnectionResult, m.CreateConnectionError
}

func (m *Service) UpdateConnection(_ string, _ *connectors.ConnectionInput, _ ...graphql.RequestOption) (*connectors.Connection, error) {
	return m.UpdateConnectionResult, m.UpdateConnectionError
}

func (m *Service) DeleteConnection(_ string, _ ...graphql.RequestOption) (*connectors.Connection, error) {
	return m.DeleteConnectionResult, m.DeleteConnectionError
}

func (m *Service) ExecuteConnectionAction(_ string, _ string, _ interface{}, _ ...graphql.RequestOption) (*interface{}, error) {
	return m.ExecuteConnectionActionResult, m.ExecuteConnectionActionError
}

func (m *Service) GetConnectorCategory(_ string, _ ...graphql.RequestOption) (*connectors.ConnectorCategory, error) {
	return m.GetConnectorCategoryResult, m.GetConnectorCategoryError
}

func (m *Service) GetConnectorInterface(_ string, _ ...graphql.RequestOption) (*connectors.ConnectorInterface, error) {
	return m.GetConnectorInterfaceResult, m.GetConnectorInterfaceError
}

func (m *Service) GetConnector(_ string, _ ...graphql.RequestOption) (*connectors.Connector, error) {
	return m.GetConnectorResult, m.GetConnectorError
}

func (m *Service) GetConnection(_ string, _ ...graphql.RequestOption) (*connectors.Connection, error) {
	return m.GetConnectionResult, m.GetConnectionError
}

func (m *Service) ConnectorCreated(ctx context.Context, connectorMethods common.IDs, allTenants bool, _ ...graphql.SubscriptionOption) (connectors.Subscription, error) {
	return &subscription{err: m.ConnectorCreatedError, connectors: m.ConnectorCreatedResult}, nil
}

func (m *Service) ConnectorDeleted(ctx context.Context, connectorMethods common.IDs, allTenants bool, _ ...graphql.SubscriptionOption) (connectors.Subscription, error) {
	return &subscription{err: m.ConnectorDeletedError, connectors: m.ConnectorDeletedResult}, nil
}

func (m *Service) ConnectorUpdated(ctx context.Context, connectorMethods common.IDs, allTenants bool, _ ...graphql.SubscriptionOption) (connectors.Subscription, error) {
	return &subscription{err: m.ConnectorUpdatedError, connectors: m.ConnectorUpdatedResult}, nil
}
