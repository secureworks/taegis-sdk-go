package connectors

import (
	"strings"
)

// Context is the execution context for the current playbook
type Context struct {
	TenantID    string
	UserID      string
	Connections []*Connection
}

const GlobalNamespace = "global"

func (ctx *Context) LookupAction(actionLookup *ActionLookup) (*Connection, *ConnectorAction) {
	globalConnectionIndex := -1
	if actionLookup.ConnectorInterfaceActionName != "" {
		//If connectorName is set, look for a match on connector name, then look up action
		if actionLookup.ConnectorName != "" {
			for i, c := range ctx.Connections {
				if actionLookup.Namespace != "" {
					namespace := GlobalNamespace
					if c.Connector.TenantID != nil {
						namespace = *c.Connector.TenantID
					}
					if namespace != actionLookup.Namespace {
						continue
					}
				}

				if strings.EqualFold(c.Connector.Name, actionLookup.ConnectorName) {
					if c.Connector.TenantID != nil {
						return c, c.Connector.LookupAction(actionLookup.ConnectorInterfaceActionName)
					}
					globalConnectionIndex = i
				}
			}
			if globalConnectionIndex >= 0 {
				c := ctx.Connections[globalConnectionIndex]
				return c, c.Connector.LookupAction(actionLookup.ConnectorInterfaceActionName)
			}
		}

		// Only loop through connectorInterfaces if necessary
		if actionLookup.ImplementedConnectorInterfaceName != "" {
			for j, c := range ctx.Connections {
				for _, i := range c.Connector.Implements {
					if actionLookup.Namespace != "" {
						namespace := GlobalNamespace
						if c.Connector.TenantID != nil {
							namespace = *c.Connector.TenantID
						}
						if namespace != actionLookup.Namespace {
							continue
						}
					}

					if strings.EqualFold(i.Name, actionLookup.ImplementedConnectorInterfaceName) {
						if c.Connector.TenantID != nil {
							return c, i.LookupAction(actionLookup.ConnectorInterfaceActionName)
						}
						globalConnectionIndex = j
					}
				}
			}
			if globalConnectionIndex >= 0 {
				c := ctx.Connections[globalConnectionIndex]
				return c, c.Connector.LookupAction(actionLookup.ConnectorInterfaceActionName)
			}
		}
	}
	return nil, nil
}
