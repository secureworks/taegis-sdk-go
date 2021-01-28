package connectors

import "github.com/secureworks/tdr-sdk-go/common"

type GetConnectionsInput struct {
	ConnectionIDs         []string
	ConnectorIDs          []string
	ConnectorInterfaceIDs []string
}

type ConnectionMethodInput struct {
	common.ObjectMetaInput
	URL        string
	parameters common.Object
}

// ConnectorActionInput defines the mutable fields of a connector action
type ConnectorActionInput struct {
	common.ObjectMetaInput
	Inputs  common.Object `json:"inputs"`
	Outputs common.Object `json:"outputs"`
}

type ConnectorInterfaceInput struct {
	common.ObjectMetaInput
	Categories common.IDs              `json:"categories"`
	Actions    []*ConnectorActionInput `json:"actions"`
	AllTenants *bool                   `json:"all_tenants"`
}

// ConnectorActionInput defines the mutable fields of a connector action
type ConnectorActionDefinitionInput struct {
	Action common.ID     `json:"action"`
	Config common.Object `json:"config"`
}

// ConnectorInput defines the mutable fields of a connector
type ConnectorInput struct {
	common.ObjectMetaInput
	Implements    common.IDs                        `json:"implements"`
	Parameters    common.Object                     `json:"parameters"`
	AuthTypes     []AuthType                        `json:"authTypes"`
	AllTenants    *bool                             `json:"allTenants"`
	Actions       []*ConnectorActionDefinitionInput `json:"actions"`
	Documentation *string                           `json:"documentation"`
	Categories    common.IDs                        `json:"categories"`
	Title         *string                           `json:"title"`
}

type ConnectorUpdateInput struct {
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Implements    common.IDs `json:"implements"`
	parameters    common.Object
	AuthTypes     []AuthType                       `json:"authTypes"`
	Documentation string                           `json:"documentation"`
	Actions       []ConnectorActionDefinitionInput `json:"actions"`
	Categories    common.IDs                       `json:"categories"`
	Tags          common.Tags                      `json:"tags"`
	Title         *string                          `json:"title"`
}
