package main

import (
	"fmt"
	"os"

	"github.com/secureworks/taegis-sdk-go/graphql"

	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/orchestration/connectors"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofrs/uuid"
)

const (
	TENANT_ID = "123456789"
)

func errTest(err error) {
	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}
}
func main() {

	connectorsSvc := connectors.New("test.com")

	testObjMetaInput := common.ObjectMetaInput{
		Name:        "test",
		Description: nil,
		Tags:        nil,
	}
	//Define a ConnectionMethod
	connectionMethodInput := connectors.ConnectionMethodInput{
		ObjectMetaInput: testObjMetaInput,
		URL:             "test.com",
	}

	definedConnectionMethod, err := connectorsSvc.DefineConnectionMethod(&connectionMethodInput)
	errTest(err)
	//connectionMethod represents a ConnectionMethod Output
	//use dot notations to access specific ConnectionMethod fields
	fmt.Printf("CONNECTION METHOD NAME: %v\n", definedConnectionMethod.Name)
	fmt.Printf("CONNECTION METHOD DESCRIPTION: %v\n", definedConnectionMethod.Description)

	connectionMethod, err := connectorsSvc.GetConnectionMethod("test")
	errTest(err)
	//connectionMethod represents a ConnectionMethod Output
	//use dot notations to access specific ConnectionMethod fields
	fmt.Printf("CONNECTION METHOD NAME: %v\n", connectionMethod.Name)
	fmt.Printf("CONNECTION METHOD DESCRIPTION: %v\n", connectionMethod.Description)

	//Remove a ConnectionMethod
	removedConnectionMethod, err := connectorsSvc.RemoveConnectionMethod(connectionMethod.ID)
	//connectionMethod represents a ConnectionMethod Output
	//use dot notations to access specific ConnectionMethod fields
	fmt.Printf("CONNECTION METHOD NAME: %v\n", removedConnectionMethod.Name)
	fmt.Printf("CONNECTION METHOD DESCRIPTION: %v\n", removedConnectionMethod.Description)

	connectionInput := connectors.GetConnectionsInput{
		ConnectionIDs:         []string{"1234"},
		ConnectorIDs:          nil,
		ConnectorInterfaceIDs: nil,
	}
	connections, err := connectorsSvc.GetConnections(&connectionInput)

	errTest(err)

	for _, connection := range connections {
		//connection represents a Connection Output
		//use dot notations to access specific Connection fields
		fmt.Printf("CONNECTION NAME: %v\n", connection.Name)
		fmt.Printf("CONNECTION DESCRIPTION: %v\n", connection.Description)
	}

	connectorInput := connectors.GetConnectorsInput{
		ConnectionMethodIDs:   []string{"test"},
		ConnectorInterfaceIDs: nil,
		ConnectorCategoryIDs:  nil,
		Tags:                  nil,
	}

	//connectors represents an array of Connector

	connectorsArr, err := connectorsSvc.GetConnectors(&connectorInput)

	errTest(err)

	for _, connector := range connectorsArr {
		//connector represents a Connector Output
		//use dot notations to access specific Connector fields
		fmt.Printf("CONNECTOR NAME: %v\n", connector.Name)
		fmt.Printf("CONNECTOR DESCRIPTION: %v\n", connector.Description)
	}

	context, err := connectorsSvc.GetContext(&connectionInput, graphql.RequestWithTenant(TENANT_ID))

	errTest(err)

	//context represents a Context Output
	//use dot notation to access specific Context fields
	fmt.Printf("CONTEXT TENANT ID: %v\n", context.TenantID)
	fmt.Printf("CONTEXT USER ID: %v\n", context.UserID)

	for _, connection := range context.Connections {
		//connection represents a Connection Output
		//use dot notations to access specific Connection fields
		fmt.Printf("CONTEXT CONNECTION NAME: %v\n", connection.Name)
		fmt.Printf("CONTEXT CONNECTION DESCRIPTION: %v\n", connection.Description)
	}

	//Create a ConnectorInterface
	connectorInterfaceInput := connectors.ConnectorInterfaceInput{
		ObjectMetaInput: testObjMetaInput,
		Categories:      nil,
		Actions:         nil,
		AllTenants:      nil,
	}
	connectorInterface, err := connectorsSvc.CreateConnectorInterface(&connectorInterfaceInput)
	errTest(err)
	//connectorInterface represents a ConnectorInterface Output
	//use dot notations to access specific ConnectorInterface fields
	fmt.Printf("CONNECTOR INTERFACE NAME: %v\n", connectorInterface.Name)
	fmt.Printf("CONNECTOR INTERFACE DESCRIPTION: %v\n", connectorInterface.Description)

	//Get ConnectorInterface
	connectorInterface, err = connectorsSvc.GetConnectorInterface(connectorInterface.ID)
	errTest(err)
	//connectorInterface represents a ConnectorInterface Output
	//use dot notations to access specific ConnectorInterface fields
	fmt.Printf("CONNECTOR INTERFACE NAME: %v\n", connectorInterface.Name)
	fmt.Printf("CONNECTOR INTERFACE DESCRIPTION: %v\n", connectorInterface.Description)

	//Update ConnectorInterface
	connectorInterface, err = connectorsSvc.UpdateConnectorInterface(connectorInterface.ID, &connectorInterfaceInput)
	errTest(err)
	//connectorInterface represents a ConnectorInterface Output
	//use dot notations to access specific ConnectorInterface fields
	fmt.Printf("CONNECTOR INTERFACE NAME: %v\n", connectorInterface.Name)
	fmt.Printf("CONNECTOR INTERFACE DESCRIPTION: %v\n", connectorInterface.Description)

	//Delete ConnectorInterface
	connectorInterface, err = connectorsSvc.DeleteConnectorInterface(connectorInterface.ID)
	errTest(err)
	//connectorInterface represents a ConnectorInterface Output
	//use dot notations to access specific ConnectorInterface fields
	fmt.Printf("DELETED CONNECTOR INTERFACE NAME: %v\n", connectorInterface.Name)
	fmt.Printf("DELETED CONNECTOR INTERFACE DESCRIPTION: %v\n", connectorInterface.Description)

	//Create Connector
	connectorIn := connectors.ConnectorInput{
		ObjectMetaInput: testObjMetaInput,
		Implements:      nil,
		Parameters:      nil,
		AuthTypes:       nil,
		AllTenants:      nil,
		Actions:         nil,
		Documentation:   nil,
		Categories:      nil,
	}
	connector, err := connectorsSvc.CreateConnector("123", &connectorIn)
	errTest(err)
	//connector represents a Connector Output
	//use dot notations to access specific Connector fields
	fmt.Printf("CONNECTOR NAME: %v\n", connector.Name)
	fmt.Printf("CONNECTOR DESCRIPTION: %v\n", connector.Description)

	//Get Connector
	connector, err = connectorsSvc.GetConnector(connector.ID)
	errTest(err)
	//connector represents a Connector Output
	//use dot notations to access specific Connector fields
	fmt.Printf("CONNECTOR NAME: %v\n", connector.Name)
	fmt.Printf("CONNECTOR DESCRIPTION: %v\n", connector.Description)
	id, _ := uuid.NewV4()

	connectorUpdateInput := connectors.ConnectorUpdateInput{
		Name:          "new name",
		Description:   "updated connector",
		Implements:    nil,
		AuthTypes:     []connectors.AuthType{connectors.AuthTypeBasic},
		Documentation: "",
		Actions:       nil,
		Categories:    common.IDs{"123"},
		Tags:          nil,
	}
	//Update Connector
	connector, err = connectorsSvc.UpdateConnector(connector.ID, &connectorUpdateInput)
	errTest(err)
	//connector represents a Connector Output
	//use dot notations to access specific Connector fields
	fmt.Printf("UPDATED CONNECTOR NAME: %v\n", connector.Name)
	fmt.Printf("UPDATED CONNECTOR DESCRIPTION: %v\n", connector.Description)

	//Delete Connector
	connector, err = connectorsSvc.DeleteConnector(connector.ID)
	errTest(err)
	//connector represents a Connector Output
	//use dot notations to access specific Connector fields
	fmt.Printf("DELETED CONNECTOR NAME: %v\n", connector.Name)
	fmt.Printf("DELETED CONNECTOR DESCRIPTION: %v\n", connector.Description)

	//Get ConnectorCategory

	connectorCategory, err := connectorsSvc.GetConnectorCategory(id.String())
	errTest(err)
	//connectorCategory represents a ConnectorCategory Output
	//use dot notations to access specific ConnectorCategory fields
	fmt.Printf("CONNECTOR CATEGORY NAME: %v\n", connectorCategory.Name)
	fmt.Printf("CONNECTOR CATEGORY DESCRIPTION: %v\n", connectorCategory.Description)

	conInput := connectors.ConnectionInput{
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Config:      nil,
		Credentials: nil,
		AuthType:    connectors.AuthTypeBasic,
		AuthURL:     nil,
		Actions:     nil,
	}
	//Create Connection
	connection, err := connectorsSvc.CreateConnection("123", &conInput)
	errTest(err)
	//connection represents a Connection Output
	//use dot notations to access specific Connection fields
	fmt.Printf("CONNECTION NAME: %v\n", connection.Name)
	fmt.Printf("CONNECTION DESCRIPTION: %v\n", connection.Description)

	//Get Connection
	connection, err = connectorsSvc.GetConnection(connection.ID)
	errTest(err)
	//connection represents a Connection Output
	//use dot notations to access specific Connection fields
	fmt.Printf("CONNECTION NAME: %v\n", connection.Name)
	fmt.Printf("CONNECTION DESCRIPTION: %v\n", connection.Description)

	//Update Connection
	connection, err = connectorsSvc.UpdateConnection(connection.ID, &conInput)
	errTest(err)
	//connection represents a Connection Output
	//use dot notations to access specific Connection fields
	fmt.Printf("CONNECTION NAME: %v\n", connection.Name)
	fmt.Printf("CONNECTION DESCRIPTION: %v\n", connection.Description)

	//Execute Connection Action
	action, err := connectorsSvc.ExecuteConnectionAction(connection.ID, "test", nil)
	errTest(err)
	//action represents an interface
	//Cast this interface to your expected type
	fmt.Printf("ACTION INTERFACE: %v\n", action)

	//Delete Connection
	connection, err = connectorsSvc.DeleteConnection(connection.ID)
	//connection represents a Connection Output
	//use dot notations to access specific Connection fields
	fmt.Printf("DELETED CONNECTION NAME: %v\n", connection.Name)
	fmt.Printf("DELETED CONNECTION DESCRIPTION: %v\n", connection.Description)

}
