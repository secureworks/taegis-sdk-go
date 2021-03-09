package connectorLogger

import (
	"time"

	"github.com/secureworks/taegis-sdk-go/common"
)

// ConnectorLogEntry defines a single log entry
type ConnectorLogEntry struct {
	ID        string        `json:"id"`
	TenantID  string        `json:"tenant_id"`
	Connector string        `json:"connector"`
	User      string        `json:"user"`
	Level     string        `json:"level"`
	Message   common.Object `json:"message"`
	RawError  string        `json:"description"`
	CreatedAt time.Time     `json:"created_at"`
	WrittenAt time.Time     `json:"written_at"`
}

// ConnectorLogEntries defines a list of logs along with other metadata
type ConnectorLogEntries struct {
	Entries    []*ConnectorLogEntry `json:"entries"`
	TotalCount int                  `json:"totalCount"`
}
