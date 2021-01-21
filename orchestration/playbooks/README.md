# Playbook Client

## Types
### PlaybookTriggerType:
    - ID          string
    - CreatedAt   time.Time
    - UpdatedAt   time.Time
    - Name        string
    - Description string
    - Parameters  Object
### PlaybookTrigger:
    - ID          string
	- Tenant      string
	- CreatedAt   time.Time
	- UpdatedAt   time.Time
	- Name        string
	- Description string
	- Instance    *PlaybookInstance
	- Type        *PlaybookTriggerType
	- Config      Object
### Playbook:
    - ID          string
	- Tenant      string
	- CreatedAt   time.Time
	- UpdatedAt   time.Time
	- Name        string
	- Description string
	- Tags        common.Tags
	- Requires    []*connectors.ConnectorInterface // connector interfaces used by playbook dsl
	- Categories  []*connectors.ConnectorCategory  // connector categories supported by those types
	- Versions    []*PlaybookVersion               // versions in sorted order most recent first
	- Head        *PlaybookVersion                 // defaults to most recent version but can be rolled back to a previous version
### PlaybookVersion:
    - id        string
	- CreatedAt time.Time
	- CreatedBy string
	- playbook  *Playbook
	- requires  []*connectors.ConnectorInterface
	- inputs    Object
	- outputs   Object
	- dsl       Object
### PlaybookInstance:
    - ID          string
	- Tenant      string
	- CreatedAt   time.Time
	- CreatedBy   string
	- UpdatedAt   time.Time
	- Playbook    *Playbook
	- Version     *PlaybookVersion
	- Tags        common.Tags
	- Trigger     *PlaybookTrigger
	- Eenabled    bool
	- Inputs      Object
	- Retries     *PlaybookRetries
	- Connections []*connectors.Connection
### PlaybookExecution:
    - ID        string
	- CreatedAt time.Time
	- CreatedBy string
	- UpdatedAt time.Time
	- Tenant    string
	- Instance  *PlaybookInstance
	- State     string
	- Inputs    Object
	- Outputs   Object
### PlaybookExecutions:
	Executions []*PlaybookExecution
	TotalCount int
### PlaybookRetries:
    - InitialInterval    int
	- MaximumInterval    int
	- BackoffCoefficient float32
	- MaximumRetries     int
	- MaximumDuration    int
	
## Functions
- New(url string, timeout time.Duration) Client
- GetPlaybook(tenantID string, playbookID string) (*Playbook, error)
- GetPlaybooks(tenantID string, categoryID *string, tags *common.Tags) ([]*Playbook, error)
- GetPlaybookInstance(tenantID string, playbookInstanceID string) (*PlaybookInstance, error)
- GetPlaybookInstances(tenantID string, playbookID string) ([]*PlaybookInstance, error)
- GetPlaybookExecution(tenantID string, playbookExecutionID string) (*PlaybookExecution, error)
- GetPlaybookExecutions(tenantID string, playbookInstanceID string, pagination common.Pagination) (*PlaybookExecutions, error)
- GetPlaybookTriggerType(id *string, name *string) (*PlaybookTriggerType, error)
- GetPlaybookTriggerTypes() ([]*PlaybookTriggerType, error)
- GetPlaybookTrigger(tenantID string, playbookTriggerID string) (*PlaybookTrigger, error)
- GetPlaybookTriggers(triggerTypeIDs []string) ([]*PlaybookTrigger, error)
- CreatePlaybookInstance(tenantID string, playbookID string, input *PlaybookInstanceInput) (*PlaybookInstance, error)
- UpdatePlaybookInstance(tenantID string, playbookInstanceID string, input *PlaybookInstanceInput) (*PlaybookInstance, error)
- DeletePlaybookInstance(tenantID string, playbookInstanceID string) (*PlaybookInstance, error)
- ExecutePlaybookInstance(tenantID string, playbookInstanceID string, params common.Object) (*PlaybookExecution, error)
- CreatePlaybook(tenantID string, input *PlaybookInput) (*Playbook, error)
- ClonePlaybook(tenantID string, playbookID string) (*Playbook, error)
- UpdatePlaybook(tenantID string, playbookID string, input *PlaybookInput) (*Playbook, error)
- DeletePlaybook(tenantID string, playbookID string) (*Playbook, error)
- ExecutePlaybook(tenantID string, playbookID string, parameters common.Object) (*PlaybookExecution, error)
- SetPlaybookInstanceState(tenantID string, playbookInstanceID string, enabled bool) (*PlaybookInstance, error)
- PlaybookInstanceCreated(ctx context.Context, playbooks common.IDs) (Subscription, error)
- PlaybookInstanceUpdated(ctx context.Context, playbooks common.IDs) (Subscription, error)
- PlaybookInstanceDeleted(ctx context.Context, playbooks common.IDs) (Subscription, error)
