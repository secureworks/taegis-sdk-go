package mocks

import (
	"context"

	"github.com/secureworks/tdr-sdk-go/graphql"

	"github.com/secureworks/tdr-sdk-go/common"
	"github.com/secureworks/tdr-sdk-go/orchestration/playbooks"
)

var _ playbooks.Service = (*Service)(nil)

type Service struct {
	GetPlaybookError               error
	GetPlaybooksError              error
	GetPlaybookInstanceError       error
	GetPlaybookInstancesError      error
	GetPlaybookExecutionError      error
	GetPlaybookExecutionsError     error
	GetTriggerTypeError            error
	GetTriggerTypesError           error
	GetPlaybookTriggerError        error
	GetTriggersError               error
	ExecutePlaybookInstanceError   error
	CreatePlaybookInstanceError    error
	UpdatePlaybookInstanceError    error
	DeletePlaybookInstanceError    error
	CreatePlaybookError            error
	UpdatePlaybookError            error
	DeletePlaybookError            error
	ClonePlaybookError             error
	ExecutePlaybookError           error
	SetPlaybookInstanceStateError  error
	PlaybookInstanceCreatedError   error
	PlaybookInstanceUpdatedError   error
	PlaybookInstanceDeletedError   error
	GetPlaybookResult              *playbooks.Playbook
	GetPlaybooksResult             []*playbooks.Playbook
	GetPlaybookInstanceResult      *playbooks.PlaybookInstance
	GetPlaybookInstancesResult     []*playbooks.PlaybookInstance
	GetPlaybookExecutionResult     *playbooks.PlaybookExecution
	GetPlaybookExecutionsResult    *playbooks.PlaybookExecutions
	GetTriggerTypeResult           *playbooks.PlaybookTriggerType
	GetTriggerTypesResult          []*playbooks.PlaybookTriggerType
	GetTriggerResult               *playbooks.PlaybookTrigger
	GetTriggersResult              []*playbooks.PlaybookTrigger
	ExecutePlaybookInstanceResult  *playbooks.PlaybookExecution
	CreatePlaybookInstanceResult   *playbooks.PlaybookInstance
	UpdatePlaybookInstanceResult   *playbooks.PlaybookInstance
	DeletePlaybookInstanceResult   *playbooks.PlaybookInstance
	CreatePlaybookResult           *playbooks.Playbook
	UpdatePlabookResult            *playbooks.Playbook
	DeletePlaybookResult           *playbooks.Playbook
	ClonePlaybookResult            *playbooks.Playbook
	ExecutePlaybookResult          *playbooks.PlaybookExecution
	SetPlaybookInstanceStateResult *playbooks.PlaybookInstance
	PlaybookInstanceCreatedResult  chan interface{}
	PlaybookInstanceUpdatedResult  chan interface{}
	PlaybookInstanceDeletedResult  chan interface{}
}

func (m *Service) GetPlaybook(playbookInstanceID string, _ ...graphql.RequestOption) (*playbooks.Playbook, error) {
	return m.GetPlaybookResult, m.GetPlaybookError
}

func (m *Service) GetPlaybooks(categoryID *string, tags *common.Tags, _ ...graphql.RequestOption) ([]*playbooks.Playbook, error) {
	return m.GetPlaybooksResult, m.GetPlaybooksError
}

func (m *Service) GetPlaybookInstance(_ string, _ ...graphql.RequestOption) (*playbooks.PlaybookInstance, error) {
	return m.GetPlaybookInstanceResult, m.GetPlaybookInstanceError
}

func (m *Service) GetPlaybookInstances(_ *string, _ ...graphql.RequestOption) ([]*playbooks.PlaybookInstance, error) {
	return m.GetPlaybookInstancesResult, m.GetPlaybookInstancesError
}

func (m *Service) GetPlaybookExecution(_ string, _ ...graphql.RequestOption) (*playbooks.PlaybookExecution, error) {
	return m.GetPlaybookExecutionResult, m.GetPlaybookExecutionError
}

func (m *Service) GetPlaybookExecutions(_ string, _ common.Pagination, _ ...graphql.RequestOption) (*playbooks.PlaybookExecutions, error) {
	return m.GetPlaybookExecutionsResult, m.GetPlaybookExecutionsError
}

func (m *Service) GetPlaybookTriggerType(_, _ *string, _ ...graphql.RequestOption) (*playbooks.PlaybookTriggerType, error) {
	return m.GetTriggerTypeResult, m.GetTriggerTypeError
}

func (m *Service) GetPlaybookTriggerTypes(_ ...graphql.RequestOption) ([]*playbooks.PlaybookTriggerType, error) {
	return m.GetTriggerTypesResult, m.GetTriggerTypesError
}

func (m *Service) GetPlaybookTrigger(_ string, _ ...graphql.RequestOption) (*playbooks.PlaybookTrigger, error) {
	return m.GetTriggerResult, m.GetPlaybookTriggerError
}

func (m *Service) GetPlaybookTriggers(_ []string, _ ...graphql.RequestOption) ([]*playbooks.PlaybookTrigger, error) {
	return m.GetTriggersResult, m.GetTriggersError
}

func (m *Service) ExecutePlaybookInstance(_ string, _ common.Object, _ ...graphql.RequestOption) (*playbooks.PlaybookExecution, error) {
	return m.ExecutePlaybookInstanceResult, m.ExecutePlaybookInstanceError
}

func (m *Service) CreatePlaybookInstance(_ string, _ *playbooks.PlaybookInstanceInput, _ ...graphql.RequestOption) (*playbooks.PlaybookInstance, error) {
	return m.CreatePlaybookInstanceResult, m.CreatePlaybookInstanceError
}

func (m *Service) UpdatePlaybookInstance(_ string, _ *playbooks.PlaybookInstanceInput, _ ...graphql.RequestOption) (*playbooks.PlaybookInstance, error) {
	return m.UpdatePlaybookInstanceResult, m.UpdatePlaybookInstanceError
}

func (m *Service) DeletePlaybookInstance(_ string, _ ...graphql.RequestOption) (*playbooks.PlaybookInstance, error) {
	return m.DeletePlaybookInstanceResult, m.DeletePlaybookInstanceError
}

func (m *Service) CreatePlaybook(_ *playbooks.PlaybookInput, _ ...graphql.RequestOption) (*playbooks.Playbook, error) {
	return m.CreatePlaybookResult, m.CreatePlaybookError
}

func (m *Service) ClonePlaybook(_ playbooks.ClonePlaybookInput, _ ...graphql.RequestOption) (*playbooks.Playbook, error) {
	return m.ClonePlaybookResult, m.ClonePlaybookError
}

func (m *Service) UpdatePlaybook(_ string, _ *playbooks.PlaybookInput, _ ...graphql.RequestOption) (*playbooks.Playbook, error) {
	return m.UpdatePlabookResult, m.UpdatePlaybookError
}

func (m *Service) DeletePlaybook(_ string, _ ...graphql.RequestOption) (*playbooks.Playbook, error) {
	return m.DeletePlaybookResult, m.DeletePlaybookError
}

func (m *Service) ExecutePlaybook(_ string, _ common.Object, _ ...graphql.RequestOption) (*playbooks.PlaybookExecution, error) {
	return m.ExecutePlaybookResult, m.ExecutePlaybookError
}

func (m *Service) SetPlaybookInstanceState(_ string, _ bool, _ ...graphql.RequestOption) (*playbooks.PlaybookInstance, error) {
	return m.SetPlaybookInstanceStateResult, m.SetPlaybookInstanceStateError
}

func (m *Service) PlaybookInstanceCreated(ctx context.Context, _ common.IDs, _ ...graphql.SubscriptionOption) (playbooks.Subscription, error) {
	return &subscription{instances: m.PlaybookInstanceCreatedResult}, m.PlaybookInstanceCreatedError
}

func (m *Service) PlaybookInstanceDeleted(ctx context.Context, _ common.IDs, _ ...graphql.SubscriptionOption) (playbooks.Subscription, error) {
	return &subscription{instances: m.PlaybookInstanceDeletedResult}, m.PlaybookInstanceDeletedError
}

func (m *Service) PlaybookInstanceUpdated(ctx context.Context, _ common.IDs, _ ...graphql.SubscriptionOption) (playbooks.Subscription, error) {
	return &subscription{instances: m.PlaybookInstanceUpdatedResult}, m.PlaybookInstanceUpdatedError
}
