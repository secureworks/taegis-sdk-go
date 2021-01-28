package playbooks

import (
	"time"

	"github.com/secureworks/tdr-sdk-go/common"
	"github.com/secureworks/tdr-sdk-go/orchestration/connectors"
)

type Metadata struct {
	ID          string      `json:"id" gorm:"type:uuid;primary_key"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Tags        common.Tags `json:"tags"`
	Sequence    *int64      `json:"sequence,omitempty"`
}

// PlaybookTriggerType defines an available triggering mechanism
type PlaybookTriggerType struct {
	ID          string        `json:"id"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Parameters  common.Object `json:"parameters"`
}

// PlaybookTrigger defines a set of attributes common to different trigger types
type PlaybookTrigger struct {
	ID          string               `json:"id"`
	Tenant      string               `json:"tenant"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Instance    *PlaybookInstance    `json:"instance"`
	Type        *PlaybookTriggerType `json:"type"`
	Config      common.Object        `json:"config"`
}

// Playbook is an entry in catalog of available playbooks"
type Playbook struct {
	Metadata   `json:",inline" yaml:",inline"`
	Tenant     string                           `json:"tenant"`
	Requires   []*connectors.ConnectorInterface `json:"requires"`   // connector interfaces used by playbook dsl (Contains only the IDs)
	Categories []*connectors.ConnectorCategory  `json:"categories"` // connector categories supported by those types (Contains only the IDs)
	Versions   []*PlaybookVersion               `json:"versions"`   // versions in sorted order most recent first
	Head       *PlaybookVersion                 `json:"head"`       // defaults to most recent version but can be rolled back to a previous version
	Title      *string                          `json:"title"`
}

// PlaybookVersion maintains a change record of the playbook definition. Multiple versions of a playbook could be in use concurrently"
type PlaybookVersion struct {
	ID        string                           `json:"id"`
	CreatedAt time.Time                        `json:"createdAt"`
	CreatedBy string                           `json:"createdBy"`
	Playbook  *Playbook                        `json:"playbook"`
	Requires  []*connectors.ConnectorInterface `json:"retries"`
	Inputs    common.Object                    `json:"inputs"`
	Outputs   common.Object                    `json:"outputs"`
	Dsl       common.Object                    `json:"dsl"`
}

// PlaybookInstance defines the configuration of a playbook in a user account"
type PlaybookInstance struct {
	Metadata    `json:",inline" yaml:",inline"`
	Tenant      string                   `json:"tenant"`
	CreatedBy   string                   `json:"createdBy"`
	Playbook    *Playbook                `json:"playbook"`
	Version     *PlaybookVersion         `json:"version"`
	Trigger     *PlaybookTrigger         `json:"trigger"`
	Enabled     bool                     `json:"enabled"`
	Inputs      common.Object            `json:"inputs"`
	Retries     *PlaybookRetries         `json:"retries"`
	Connections []*connectors.Connection `json:"connections"`
}

// PlaybookExecution represents the state of a current playbook execution
type PlaybookExecution struct {
	ID        string            `json:"id"`
	CreatedAt time.Time         `json:"createdAt"`
	CreatedBy string            `json:"createdBy"`
	UpdatedAt time.Time         `json:"updatedAt"`
	Tenant    string            `json:"tenant"`
	Instance  *PlaybookInstance `json:"instance"`
	State     string            `json:"state"`
	Inputs    common.Object     `json:"inputs"`
	Outputs   common.Object     `json:"outputs"`
}

type PlaybookExecutions struct {
	Executions []*PlaybookExecution
	TotalCount int
}

type PlaybookRetries struct {
	InitialInterval    int     `json:"initialInterval"`
	MaximumInterval    int     `json:"maximumInterval"`
	BackoffCoefficient float32 `json:"backOffCoefficient"`
	MaximumRetries     int     `json:"maximumRetries"`
	MaximumDuration    int     `json:"maximumDuration"`
}
