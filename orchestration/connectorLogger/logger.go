package connectorLogger

import (
	"github.com/secureworks/taegis-sdk-go/client"
	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/graphql"
)

type Service interface {
	GetAllConnectorLogs(args ConnectorLogQueryInput, pagination common.Pagination, opts ...graphql.RequestOption) (*ConnectorLogEntries, error)
}

var _ Service = (*connectorLoggerSvc)(nil)

type connectorLoggerSvc struct {
	client *client.Client
	url    string
}

func New(url string, opts ...client.Option) *connectorLoggerSvc {
	client := client.NewClient(opts...)
	return &connectorLoggerSvc{
		client: client,
		url:    url,
	}
}

func (loggerService *connectorLoggerSvc) GetAllConnectorLogs(args ConnectorLogQueryInput, pagination common.Pagination, opts ...graphql.RequestOption) (*ConnectorLogEntries, error) {
	req := graphql.NewRequest(getAllConnectorLogsQuery, opts...)
	req.Var("args", args)
	req.Var("pagination", pagination)

	var data struct {
		ConnectorLogs *ConnectorLogEntries `json:"getAllConnectorLogs"`
	}
	if err := graphql.ExecuteQuery(loggerService.client, loggerService.url, req, &data); err != nil {
		return nil, err
	}

	return data.ConnectorLogs, nil
}
