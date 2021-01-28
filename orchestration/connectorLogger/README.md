# Playbook Client

## Types
### ConnectorLogEntry:
    - ID        string       
    - TenantID  string       
    - Connector string       
    - User      string       
    - Level     string       
    - Message   Object
    - RawError  string       
    - CreatedAt time.Time    
    - WrittenAt time.Time    

### ConnectorLogEntries:
	- Nodes      []*ConnectorLogEntry
	- TotalCount int

### Pagination:
    - Page    *int 
    - PerPage *int 


## Functions
- New(url string, timeout time.Duration) Client
- GetAllConnectorLogs(tenantID string, args GetConnectorLoggerInput, pagination common.Pagination) (*ConnectorLogEntries, error)
