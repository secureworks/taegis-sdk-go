package connectorLogger

import (
	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/graphql"
	"github.com/secureworks/taegis-sdk-go/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	tenantId         = "1"
	testConnectorLog = ConnectorLogEntry{
		ID:        "ID",
		TenantID:  "TenantID",
		Connector: "Connector",
		User:      "User",
		Level:     "Level",
		Message:   common.Object{"message": "test"},
		RawError:  "RawError",
		CreatedAt: time.Time{},
		WrittenAt: time.Time{},
	}
	header = testutils.CreateHeader()
)

func TestConnectorLoggerSvc_GetAllConnectorLogs(t *testing.T) {
	r := struct {
		Out *ConnectorLogEntries `json:"getAllConnectorLogs"`
	}{
		Out: &ConnectorLogEntries{Entries: []*ConnectorLogEntry{&testConnectorLog}},
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	log, err := c.GetAllConnectorLogs(ConnectorLogQueryInput{}, common.NewPaginationOptions(1, 1), graphql.RequestWithTenant(tenantId))
	assert.Nil(t, err)
	assert.Len(t, log.Entries, 1)
	assert.Equal(t, "Connector", log.Entries[0].Connector)
}
