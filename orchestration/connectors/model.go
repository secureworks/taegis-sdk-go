package connectors

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/secureworks/tdr-sdk-go/common"
)

// Enumeration of supported auth types
type AuthType string

const (
	AuthTypeNone                    AuthType = "None"
	AuthTypePlatform                AuthType = "Platform"
	AuthTypeRaw                     AuthType = "Raw"
	AuthTypeBasic                   AuthType = "Basic"
	AuthTypeAPIKey                  AuthType = "APIKey"
	AuthTypeClientCerts             AuthType = "ClientCerts"
	AuthTypeOAuthClientCreds        AuthType = "OAuthClientCreds"
	AuthTypeOAuthOwnerPasswordCreds AuthType = "OAuthPassword"
	AuthTypeOAuthAuthCodeCreds      AuthType = "OAuthAuthCode"
)

var AllAuthType = []AuthType{
	AuthTypeNone,
	AuthTypePlatform,
	AuthTypeRaw,
	AuthTypeBasic,
	AuthTypeAPIKey,
	AuthTypeClientCerts,
	AuthTypeOAuthClientCreds,
	AuthTypeOAuthOwnerPasswordCreds,
	AuthTypeOAuthAuthCodeCreds,
}

var OAuthType = []AuthType{
	AuthTypeOAuthClientCreds,
	AuthTypeOAuthOwnerPasswordCreds,
	AuthTypeOAuthAuthCodeCreds,
	AuthTypePlatform,
}

// ConnectorCategory is a grouping/categorization of available connectors (e.g. IP reputation services, DNS lookup, etc)
type ConnectorCategory struct {
	ID          string      `json:"id"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Tags        common.Tags `json:"tags,omitempty"`
}

// ConnectorAction defines a method or activity that can be called on a connector and its corresponding input and output
type ConnectorAction struct {
	ID          string              `json:"id"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Tags        common.Tags         `json:"tags,omitempty"`
	Interface   *ConnectorInterface `json:"interface"`
	Inputs      common.Object       `json:"inputs,omitempty"`
	Outputs     common.Object       `json:"outputs,omitempty"`
}

// ConnectorInterface defines an abstract type (set of actions) that could be implemented by multiple connectors
type ConnectorInterface struct {
	ID          string               `json:"id"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Tags        common.Tags          `json:"tags,omitempty"`
	Categories  []*ConnectorCategory `json:"categories,omitempty"`
	Actions     []*ConnectorAction   `json:"actions,omitempty"`
	TenantID    *string              `json:"tenantId,omitempty"`
}

// ConnectionMethod references a service that implements connectors of a specific connection method (e.g. http, grpc, graphql)
type ConnectionMethod struct {
	ID          string        `json:"id"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Tags        common.Tags   `json:"tags,omitempty"`
	Parameters  common.Object `json:"parameters,omitempty"`
	GraphQLURL  string        `json:"graphqlUrl"`
	Connectors  []*Connector  `json:"connectors"`
}

// Connector is an entry in catalog of available connectors (e.g. service now connector based on generic http connector)
type Connector struct {
	ID          string                       `json:"id"`
	CreatedAt   time.Time                    `json:"createdAt"`
	UpdatedAt   time.Time                    `json:"updatedAt"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Tags        common.Tags                  `json:"tags,omitempty"`
	Docs        *string                      `json:"documentation,omitempty"`
	Method      *ConnectionMethod            `json:"method,omitempty"`
	Implements  []*ConnectorInterface        `json:"implements,omitempty"`
	Actions     []*ConnectorActionDefinition `json:"actions,omitempty"`
	Connections []*Connection                `json:"connections,omitempty"`
	Parameters  common.Object                `json:"parameters,omitempty"`
	AuthTypes   []AuthType                   `json:"authTypes"`
	Sequence    *int64                       `json:"sequence,omitempty"`
	Title       *string                      `json:"title"`
	TenantID    *string                      `json:"tenant,omitempty"`
}

// ConnectorActionDefinition defines the configuration of a connector action implementation"
type ConnectorActionDefinition struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
	Action    *ConnectorAction `json:"action"`
	Config    common.Object    `json:"config"`
}

// Connection is a per-tenant configuration of a connector
type Connection struct {
	ID          string             `json:"id"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Tags        common.Tags        `json:"tags,omitempty"`
	Connector   *Connector         `json:"connector,omitempty"`
	Actions     []*ConnectorAction `json:"actions,omitempty"`
	AuthType    AuthType           `json:"authType"`
	AuthURL     *string            `json:"authUrl,omitempty"`
	Config      common.Object      `json:"config,omitempty"`
	Credentials common.Object      `json:"credentials,omitempty"`
	Sequence    *int64             `json:"sequence,omitempty"`
}

// ConnectionInput defines the mutable fields of a connection
type ConnectionInput struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Tags        common.Tags   `json:"tags"`
	Config      common.Object `json:"config"`
	Credentials common.Object `json:"credentials"`
	AuthType    AuthType      `json:"authType"`
	AuthURL     *string       `json:"authUrl"`
	Actions     []string      `json:"actions"`
}

type ActionLookup struct {
	ConnectorName                     string
	ImplementedConnectorInterfaceName string
	ConnectorInterfaceActionName      string
	Namespace                         string
}

type GetConnectorsInput struct {
	ConnectionMethodIDs   []string
	ConnectorInterfaceIDs []string
	ConnectorCategoryIDs  []string
	Tags                  []string
}

func (c *Connector) LookupAction(name string) *ConnectorAction {
	for _, t := range c.Implements {
		for _, a := range t.Actions {
			if strings.EqualFold(a.Name, name) {
				return a
			}
		}
	}
	return nil
}

func (c *ConnectorInterface) LookupAction(name string) *ConnectorAction {
	for _, a := range c.Actions {
		if strings.EqualFold(a.Name, name) {
			return a
		}
	}
	return nil
}

func (a AuthType) IsValid() bool {
	for _, at := range AllAuthType {
		if a == at {
			return true
		}
	}
	return false
}

func (a AuthType) String() string {
	return string(a)
}

func (a *AuthType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*a = AuthType(str)
	if !a.IsValid() {
		return fmt.Errorf("%s is not a valid AuthType", str)
	}
	return nil
}

func (a AuthType) MarshalGQL(w io.Writer) {
	_, _ = fmt.Fprint(w, strconv.Quote(a.String()))
}

func (a *AuthType) UnmarshalJSON(v []byte) error {
	*a = AuthType(strings.Trim(string(v), "'\""))
	if !a.IsValid() {
		return fmt.Errorf("%s is not a valid AuthType", a.String())
	}
	return nil
}

func (a AuthType) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(a.String())), nil
}

func (a AuthType) IsOAuthType() bool {
	for _, auth := range OAuthType {
		if auth == a {
			return true
		}
	}
	return false
}
