package playbooks

import (
	"context"

	"github.com/secureworks/taegis-sdk-go/client"
	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/graphql"
)

type Service interface {
	GetPlaybook(playbookID string, opts ...graphql.RequestOption) (*Playbook, error)
	GetPlaybooks(categoryID *string, tags *common.Tags, opts ...graphql.RequestOption) ([]*Playbook, error)
	GetPlaybookInstance(playbookInstanceID string, opts ...graphql.RequestOption) (*PlaybookInstance, error)
	GetPlaybookInstances(playbookID *string, opts ...graphql.RequestOption) ([]*PlaybookInstance, error)
	GetPlaybookExecution(playbookExecutionID string, opts ...graphql.RequestOption) (*PlaybookExecution, error)
	GetPlaybookExecutions(playbookInstanceID string, pagination common.Pagination, opts ...graphql.RequestOption) (*PlaybookExecutions, error)
	GetPlaybookTriggerType(id *string, name *string, opts ...graphql.RequestOption) (*PlaybookTriggerType, error)
	GetPlaybookTriggerTypes(opts ...graphql.RequestOption) ([]*PlaybookTriggerType, error)
	GetPlaybookTrigger(playbookTriggerID string, opts ...graphql.RequestOption) (*PlaybookTrigger, error)
	GetPlaybookTriggers(triggerTypeIDs []string, opts ...graphql.RequestOption) ([]*PlaybookTrigger, error)
	CreatePlaybookInstance(playbookID string, input *PlaybookInstanceInput, opts ...graphql.RequestOption) (*PlaybookInstance, error)
	UpdatePlaybookInstance(playbookInstanceID string, input *PlaybookInstanceInput, opts ...graphql.RequestOption) (*PlaybookInstance, error)
	DeletePlaybookInstance(playbookInstanceID string, opts ...graphql.RequestOption) (*PlaybookInstance, error)
	ExecutePlaybookInstance(playbookInstanceID string, params common.Object, opts ...graphql.RequestOption) (*PlaybookExecution, error)
	CreatePlaybook(input *PlaybookInput, opts ...graphql.RequestOption) (*Playbook, error)
	ClonePlaybook(input ClonePlaybookInput, opts ...graphql.RequestOption) (*Playbook, error)
	UpdatePlaybook(playbookID string, input *PlaybookInput, opts ...graphql.RequestOption) (*Playbook, error)
	DeletePlaybook(playbookID string, opts ...graphql.RequestOption) (*Playbook, error)
	ExecutePlaybook(playbookID string, parameters common.Object, opts ...graphql.RequestOption) (*PlaybookExecution, error)
	SetPlaybookInstanceState(playbookInstanceID string, enabled bool, opts ...graphql.RequestOption) (*PlaybookInstance, error)
	PlaybookInstanceCreated(ctx context.Context, playbooks common.IDs, options ...graphql.SubscriptionOption) (Subscription, error)
	PlaybookInstanceUpdated(ctx context.Context, playbooks common.IDs, options ...graphql.SubscriptionOption) (Subscription, error)
	PlaybookInstanceDeleted(ctx context.Context, playbooks common.IDs, options ...graphql.SubscriptionOption) (Subscription, error)
}

var _ Service = (*playbookSvc)(nil)

type playbookSvc struct {
	client *client.Client
	url    string
}

func New(url string, opts ...client.Option) *playbookSvc {
	client := client.NewClient(opts...)
	return &playbookSvc{
		client: client,
		url:    url,
	}
}

func (playbookService *playbookSvc) ExecutePlaybookInstance(id string, params common.Object, opts ...graphql.RequestOption) (*PlaybookExecution, error) {
	req := graphql.NewRequest(executePlaybookInstanceMutation, opts...)
	req.Var("playbookInstanceId", id)
	req.Var("parameters", params)

	var data struct {
		PlaybookExecution *PlaybookExecution `json:"executePlaybookInstance"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}
	return data.PlaybookExecution, nil
}

func (playbookService *playbookSvc) CreatePlaybook(input *PlaybookInput, opts ...graphql.RequestOption) (*Playbook, error) {
	req := graphql.NewRequest(createPlaybookMutation, opts...)
	req.Var("playbook", input)

	var data struct {
		Playbook *Playbook `json:"createPlaybook"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.Playbook, nil
}

func (playbookService *playbookSvc) ClonePlaybook(input ClonePlaybookInput, opts ...graphql.RequestOption) (*Playbook, error) {
	req := graphql.NewRequest(clonePlaybookMutation, opts...)
	req.Var("input", input)

	var data struct {
		Playbook *Playbook `json:"clonePlaybook"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.Playbook, nil
}

func (playbookService *playbookSvc) UpdatePlaybook(playbookID string, input *PlaybookInput, opts ...graphql.RequestOption) (*Playbook, error) {
	req := graphql.NewRequest(updatePlaybookMutation, opts...)
	req.Var("playbookId", playbookID)
	req.Var("playbook", input)

	var data struct {
		Playbook *Playbook `json:"updatePlaybook"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.Playbook, nil
}

func (playbookService *playbookSvc) DeletePlaybook(playbookID string, opts ...graphql.RequestOption) (*Playbook, error) {
	req := graphql.NewRequest(deletePlaybookMutation, opts...)
	req.Var("playbookId", playbookID)

	var data struct {
		Playbook *Playbook `json:"deletePlaybook"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.Playbook, nil
}

func (playbookService *playbookSvc) ExecutePlaybook(playbookID string, parameters common.Object, opts ...graphql.RequestOption) (*PlaybookExecution, error) {
	req := graphql.NewRequest(executePlaybookMutation, opts...)
	req.Var("playbookId", playbookID)
	req.Var("parameters", parameters)

	var data struct {
		PlaybookExecution *PlaybookExecution `json:"executePlaybook"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookExecution, nil
}

func (playbookService *playbookSvc) CreatePlaybookInstance(playbookID string, input *PlaybookInstanceInput, opts ...graphql.RequestOption) (*PlaybookInstance, error) {
	req := graphql.NewRequest(createPlaybookInstanceMutation, opts...)
	req.Var("playbookId", playbookID)
	req.Var("instance", input)

	var data struct {
		PlaybookInstance *PlaybookInstance `json:"createPlaybookInstance"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookInstance, nil
}

func (playbookService *playbookSvc) UpdatePlaybookInstance(playbookInstanceID string, input *PlaybookInstanceInput, opts ...graphql.RequestOption) (*PlaybookInstance, error) {
	req := graphql.NewRequest(updatePlaybookInstanceMutation, opts...)
	req.Var("playbookInstanceId", playbookInstanceID)
	req.Var("instance", input)

	var data struct {
		PlaybookInstance *PlaybookInstance `json:"updatePlaybookInstance"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookInstance, nil
}

func (playbookService *playbookSvc) DeletePlaybookInstance(playbookInstanceID string, opts ...graphql.RequestOption) (*PlaybookInstance, error) {
	req := graphql.NewRequest(deletePlaybookInstanceMutation, opts...)
	req.Var("playbookInstanceId", playbookInstanceID)

	var data struct {
		PlaybookInstance *PlaybookInstance `json:"deletePlaybookInstance"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookInstance, nil
}

func (playbookService *playbookSvc) SetPlaybookInstanceState(playbookInstanceID string, enabled bool, opts ...graphql.RequestOption) (*PlaybookInstance, error) {
	req := graphql.NewRequest(setPlaybookInstanceStateMutation, opts...)
	req.Var("playbookInstanceId", playbookInstanceID)
	req.Var("enabled", enabled)

	var data struct {
		PlaybookInstance *PlaybookInstance `json:"setPlaybookInstanceState"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookInstance, nil
}

func (playbookService *playbookSvc) GetPlaybook(playbookID string, opts ...graphql.RequestOption) (*Playbook, error) {
	req := graphql.NewRequest(getPlaybookQuery, opts...)
	req.Var("playbookId", playbookID)

	var data struct {
		Playbook *Playbook `json:"playbook"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.Playbook, nil
}

func (playbookService *playbookSvc) GetPlaybooks(categoryID *string, tags *common.Tags, opts ...graphql.RequestOption) ([]*Playbook, error) {
	req := graphql.NewRequest(getPlaybooksQuery, opts...)
	req.Var("categoryId", categoryID)
	req.Var("tags", tags)

	var data struct {
		Playbooks []*Playbook `json:"playbooks"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.Playbooks, nil
}

func (playbookService *playbookSvc) GetPlaybookInstance(playbookInstanceID string, opts ...graphql.RequestOption) (*PlaybookInstance, error) {
	req := graphql.NewRequest(getPlaybookInstanceQuery, opts...)
	req.Var("playbookInstanceId", playbookInstanceID)

	var data struct {
		PlaybookInstance *PlaybookInstance `json:"playbookInstance"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookInstance, nil
}

func (playbookService *playbookSvc) GetPlaybookInstances(playbookID *string, opts ...graphql.RequestOption) ([]*PlaybookInstance, error) {
	req := graphql.NewRequest(getPlaybookInstancesQuery, opts...)
	if playbookID != nil {
		req.Var("playbookId", playbookID)
	}

	var data struct {
		PlaybookInstances []*PlaybookInstance `json:"playbookInstances"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookInstances, nil
}

func (playbookService *playbookSvc) GetPlaybookExecution(playbookExecutionID string, opts ...graphql.RequestOption) (*PlaybookExecution, error) {
	req := graphql.NewRequest(getPlaybookExecutionQuery, opts...)
	req.Var("playbookExecutionId", playbookExecutionID)

	var data struct {
		PlaybookExecution *PlaybookExecution `json:"playbookExecution"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookExecution, nil
}

func (playbookService *playbookSvc) GetPlaybookExecutions(playbookInstanceID string, pagination common.Pagination, opts ...graphql.RequestOption) (*PlaybookExecutions, error) {
	req := graphql.NewRequest(getPlaybookExecutionsQuery, opts...)
	req.Var("playbookInstanceId", playbookInstanceID)
	req.Var("pagination", pagination)

	var data struct {
		PlaybookExecutions *PlaybookExecutions `json:"playbookExecutions"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookExecutions, nil
}

func (playbookService *playbookSvc) GetPlaybookTrigger(playbookTriggerID string, opts ...graphql.RequestOption) (*PlaybookTrigger, error) {
	req := graphql.NewRequest(getPlaybookTriggerQuery, opts...)
	req.Var("playbookTriggerId", playbookTriggerID)

	var data struct {
		PlaybookTrigger *PlaybookTrigger `json:"playbookTrigger"`
	}
	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookTrigger, nil
}

func (playbookService *playbookSvc) GetPlaybookTriggers(ids []string, opts ...graphql.RequestOption) ([]*PlaybookTrigger, error) {
	req := graphql.NewRequest(getPlaybookTriggersQuery, opts...)
	req.Var("playbookTriggerTypeIds", ids)
	var data struct {
		PlaybookTriggers []*PlaybookTrigger `json:"playbookTriggers"`
	}

	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookTriggers, nil
}

func (playbookService *playbookSvc) GetPlaybookTriggerType(id *string, name *string, opts ...graphql.RequestOption) (*PlaybookTriggerType, error) {
	req := graphql.NewRequest(getTriggerTypeQuery, opts...)
	req.Var("playbookTriggerTypeId", id)
	req.Var("playbookTriggerTypeName", name)

	var data struct {
		PlaybookTriggerType *PlaybookTriggerType `json:"playbookTriggerType"`
	}

	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookTriggerType, nil
}

func (playbookService *playbookSvc) GetPlaybookTriggerTypes(opts ...graphql.RequestOption) ([]*PlaybookTriggerType, error) {
	req := graphql.NewRequest(getTriggerTypesQuery, opts...)
	var data struct {
		PlaybookTriggerTypes []*PlaybookTriggerType `json:"playbookTriggerTypes"`
	}

	if err := graphql.ExecuteQuery(playbookService.client, playbookService.url, req, &data); err != nil {
		return nil, err
	}

	return data.PlaybookTriggerTypes, nil
}
