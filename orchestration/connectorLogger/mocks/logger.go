package mocks

import (
	"github.com/secureworks/tdr-sdk-go/common"
	"github.com/secureworks/tdr-sdk-go/graphql"
	"github.com/secureworks/tdr-sdk-go/orchestration/connectorLogger"
)

var _ connectorLogger.Service = (*Service)(nil)

type Service struct {
	GetAllConnectorLogsError  error
	GetAllConnectorLogsResult *connectorLogger.ConnectorLogEntries
}

func (m *Service) GetAllConnectorLogs(_ connectorLogger.ConnectorLogQueryInput, _ common.Pagination, _ ...graphql.RequestOption) (*connectorLogger.ConnectorLogEntries, error) {
	return m.GetAllConnectorLogsResult, m.GetAllConnectorLogsError
}
