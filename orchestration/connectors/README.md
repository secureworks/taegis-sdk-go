# Connectors SDK

## Types
  ### Context:
    - TenantID     string
    - UserID       string
    - Connections []*Connection
  ### Connection:
    - ID          string             
    - CreatedAt   time.Time          
    - UpdatedAt   time.Time          
    - Name        string           
    - Description string             
    - Tags        common.Tags           
    - Connector   *Connector         
    - Actions     []*ConnectorAction 
    - AuthType    string             
    - AuthURL     *string            
    - Config      Object             
    - Credentials Object    
  ### ConnectionInput
    - Name        string
    - Description string  
    - Tags        common.Tags 
    - Config      Object   
    - Credentials Object   
    - AuthType    string   
    - AuthURL     *string  
    - Actions     []string    
    

## Functions
- New(url string) Client
- GetConnectionMethod(name string) (*ConnectionMethod, error) 
- GetConnectors(connectorMethodID string) ([]*Connector, error)
- GetConnections(tenantID string, connectionIDs []string, connectorIDs []string, connectorInterfaceIDs []string) ([]*Connection, error)
- GetContext(tenantID string, connectionIDs []string, connectorIDs []string, connectorInterfaceIDs []string) (*Context, error)



