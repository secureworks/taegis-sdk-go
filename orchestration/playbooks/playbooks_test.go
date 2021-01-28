package playbooks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/secureworks/tdr-sdk-go/common"
	"github.com/secureworks/tdr-sdk-go/testutils"
)

var (
	testPlaybookTriggerType = PlaybookTriggerType{
		ID:          "123",
		CreatedAt:   time.Now().Truncate(time.Microsecond),
		UpdatedAt:   time.Now().Truncate(time.Microsecond),
		Name:        "testTriggerType",
		Description: "test playbook trigger type",
	}
	testPlaybookTrigger = PlaybookTrigger{
		ID:          "111",
		Tenant:      "1",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "testTrigger",
		Description: "test playbook trigger",
		Instance:    nil,
		Type:        &testPlaybookTriggerType,
		Config:      nil,
	}
	testPlaybookVersion = PlaybookVersion{
		ID:        "1",
		CreatedAt: time.Time{},
		CreatedBy: "Tester McGee",
		Playbook:  nil,
		Requires:  nil,
		Inputs:    nil,
		Outputs:   nil,
		Dsl:       nil,
	}
	testPlaybook = Playbook{
		Metadata: Metadata{
			ID:          "1234",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
			Name:        "testPlaybook",
			Description: "test playbook",
			Tags:        nil,
		},
		Tenant:     "1",
		Requires:   nil,
		Categories: nil,
		Versions:   []*PlaybookVersion{&testPlaybookVersion},
		Head:       nil,
	}
	testPlayBookInstance = PlaybookInstance{
		Metadata: Metadata{
			ID:        "123",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			Tags:      nil,
		},
		Tenant:      "1",
		CreatedBy:   "Tester McGee",
		Playbook:    &testPlaybook,
		Version:     &testPlaybookVersion,
		Trigger:     &testPlaybookTrigger,
		Enabled:     false,
		Inputs:      nil,
		Retries:     nil,
		Connections: nil,
	}
	testPlaybookExecution = PlaybookExecution{
		ID:        "123",
		CreatedAt: time.Time{},
		CreatedBy: "Tester McGee",
		UpdatedAt: time.Time{},
		Tenant:    "1",
		Instance:  &testPlayBookInstance,
		State:     "completed",
		Inputs:    nil,
		Outputs:   nil,
	}

	testPlaybookExecutions = PlaybookExecutions{
		Executions: []*PlaybookExecution{&testPlaybookExecution},
		TotalCount: 1,
	}

	header = testutils.CreateHeader()
)

func TestPlaybookSvc_ExecuteInstance(t *testing.T) {
	r := struct {
		Out *PlaybookExecution `json:"executePlaybookInstance"`
	}{
		Out: &testPlaybookExecution,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	o := make(map[string]interface{})
	o["test"] = "test"

	playbookExecution, err := c.ExecutePlaybookInstance("123", o)

	assert.Nil(t, err)

	assert.Equal(t, "123", playbookExecution.ID)
}

func TestPlaybookSvc_GetTriggers(t *testing.T) {
	r := struct {
		Out []*PlaybookTrigger `json:"playbookTriggers"`
	}{
		Out: []*PlaybookTrigger{&testPlaybookTrigger},
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	triggers, err := c.GetPlaybookTriggers([]string{"1"})
	assert.Nil(t, err)
	assert.Len(t, triggers, 1)
	assert.Equal(t, "testTrigger", triggers[0].Name)
}

func TestPlaybookSvc_GetTriggerTypes(t *testing.T) {
	r := struct {
		Out []*PlaybookTriggerType `json:"playbookTriggerTypes"`
	}{
		Out: []*PlaybookTriggerType{&testPlaybookTriggerType},
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	triggerTypes, err := c.GetPlaybookTriggerTypes()

	assert.Nil(t, err)
	assert.Len(t, triggerTypes, 1)
	assert.Equal(t, "testTriggerType", triggerTypes[0].Name)
}

func TestPlaybookSvc_GetTriggerType(t *testing.T) {
	r := struct {
		Out *PlaybookTriggerType `json:"playbookTriggerType"`
	}{
		Out: &testPlaybookTriggerType,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	id := "123"
	triggerTypeName := "testTriggerType"
	triggerType, err := c.GetPlaybookTriggerType(&id, &triggerTypeName)

	assert.Nil(t, err)
	assert.Equal(t, "testTriggerType", triggerType.Name)
}

func TestPlaybookSvc_ClonePlaybook(t *testing.T) {
	r := struct {
		Out *Playbook `json:"clonePlaybook"`
	}{
		Out: &testPlaybook,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	pb, err := c.ClonePlaybook(ClonePlaybookInput{PlaybookID: "123433", Name: "test"})
	assert.Nil(t, err)
	assert.Equal(t, testPlaybook.ID, pb.ID)

}

func TestPlaybookSvc_CreatePlaybook(t *testing.T) {
	r := struct {
		Out *Playbook `json:"createPlaybook"`
	}{
		Out: &testPlaybook,
	}
	metaInput := common.ObjectMetaInput{
		Name:        "testPlaybook",
		Description: nil,
		Tags:        nil,
	}
	playbookInput := PlaybookInput{
		ObjectMetaInput: metaInput,
		Head:            nil,
		Version:         nil,
		Categories:      nil,
	}
	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	pb, err := c.CreatePlaybook(&playbookInput)
	assert.Nil(t, err)
	assert.Equal(t, metaInput.Name, pb.Name)
}

func TestPlaybookSvc_CreatePlaybookInstance(t *testing.T) {
	r := struct {
		Out *PlaybookInstance `json:"createPlaybookInstance"`
	}{
		Out: &testPlayBookInstance,
	}
	id := common.ID("test")
	objMetaInput := common.ObjectMetaInput{
		Name:        "testPlaybook",
		Description: nil,
		Tags:        nil,
	}
	pbTriggerInput := PlaybookTriggerInput{
		ObjectMetaInput: objMetaInput,
		TypeID:          id,
		Config:          nil,
	}
	pbInstanceInput := PlaybookInstanceInput{
		ObjectMetaInput: objMetaInput,
		Trigger:         pbTriggerInput,
		Enabled:         false,
		Inputs:          nil,
		Connections:     nil,
	}
	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	i, err := c.CreatePlaybookInstance("1234", &pbInstanceInput)

	assert.Nil(t, err)
	assert.Equal(t, objMetaInput.Name, i.Playbook.Name)
}

func TestPlaybookSvc_DeletePlaybook(t *testing.T) {
	r := struct {
		Out *Playbook `json:"deletePlaybook"`
	}{
		Out: &testPlaybook,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	pb, err := c.DeletePlaybook("1234")

	assert.Nil(t, err)
	assert.Equal(t, "1234", pb.ID)
}

func TestPlaybookSvc_DeletePlaybookInstance(t *testing.T) {
	r := struct {
		Out *PlaybookInstance `json:"deletePlaybookInstance"`
	}{
		Out: &testPlayBookInstance,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	i, err := c.DeletePlaybookInstance("123")

	assert.Nil(t, err)
	assert.Equal(t, "123", i.ID)
}

func TestPlaybookSvc_ExecutePlaybook(t *testing.T) {
	r := struct {
		Out *PlaybookExecution `json:"executePlaybook"`
	}{
		Out: &testPlaybookExecution,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)

	e, err := c.ExecutePlaybook("1234", nil)

	assert.Nil(t, err)
	assert.Equal(t, "123", e.ID)
}

func TestPlaybookSvc_GetPlaybook(t *testing.T) {
	r := struct {
		Out *Playbook `json:"playbook"`
	}{
		Out: &testPlaybook,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	pb, err := c.GetPlaybook("123")
	assert.Nil(t, err)
	assert.Equal(t, "1234", pb.ID)
}

func TestPlaybookSvc_GetPlaybookExecution(t *testing.T) {
	r := struct {
		Out *PlaybookExecution `json:"playbookExecution"`
	}{
		Out: &testPlaybookExecution,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	e, err := c.GetPlaybookExecution("123")

	assert.Nil(t, err)
	assert.Equal(t, "123", e.ID)

}

func TestPlaybookSvc_GetPlaybookExecutions(t *testing.T) {
	r := struct {
		Out *PlaybookExecutions `json:"playbookExecutions"`
	}{
		Out: &testPlaybookExecutions,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	e, err := c.GetPlaybookExecutions("1234", common.NewPaginationOptions(1, 1))

	assert.Nil(t, err)
	assert.Equal(t, e.TotalCount, 1)
	assert.Len(t, e.Executions, 1)
	assert.Equal(t, "123", e.Executions[0].ID)
}

func TestPlaybookSvc_GetPlaybookInstance(t *testing.T) {
	r := struct {
		Out *PlaybookInstance `json:"playbookInstance"`
	}{
		Out: &testPlayBookInstance,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	i, err := c.GetPlaybookInstance("123")
	assert.Nil(t, err)
	assert.Equal(t, "123", i.ID)
}

func TestPlaybookSvc_GetPlaybookInstances(t *testing.T) {
	r := struct {
		Out []*PlaybookInstance `json:"playbookInstances"`
	}{
		Out: []*PlaybookInstance{&testPlayBookInstance},
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	playbookID := "123"
	i, err := c.GetPlaybookInstances(&playbookID)
	assert.Nil(t, err)
	assert.Len(t, i, 1)
	assert.Equal(t, "123", i[0].ID)
}

func TestPlaybookSvc_GetPlaybooks(t *testing.T) {
	r := struct {
		Out []*Playbook `json:"playbooks"`
	}{
		Out: []*Playbook{&testPlaybook},
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	pbs, err := c.GetPlaybooks(nil, nil)
	assert.Nil(t, err)
	assert.Len(t, pbs, 1)
	assert.Equal(t, "1234", pbs[0].ID)
}

func TestPlaybookSvc_GetPlaybookTrigger(t *testing.T) {
	r := struct {
		Out *PlaybookTrigger `json:"playbookTrigger"`
	}{
		Out: &testPlaybookTrigger,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	trigger, err := c.GetPlaybookTrigger("111")
	assert.Nil(t, err)
	assert.Equal(t, "111", trigger.ID)
}

func TestPlaybookSvc_SetPlaybookInstanceState(t *testing.T) {
	r := struct {
		Out *PlaybookInstance `json:"setPlaybookInstanceState"`
	}{
		Out: &testPlayBookInstance,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	i, err := c.SetPlaybookInstanceState("123", true)
	assert.Nil(t, err)
	assert.Equal(t, "123", i.ID)
}

func TestPlaybookSvc_UpdatePlaybook(t *testing.T) {
	r := struct {
		Out *Playbook `json:"updatePlaybook"`
	}{
		Out: &testPlaybook,
	}
	objMetaInput := common.ObjectMetaInput{
		Name:        "test",
		Description: nil,
		Tags:        nil,
	}
	playbookInput := PlaybookInput{
		ObjectMetaInput: objMetaInput,
		Head:            nil,
		Version:         nil,
		Categories:      nil,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	pb, err := c.UpdatePlaybook("1234", &playbookInput)

	assert.Nil(t, err)
	assert.Equal(t, "1234", pb.ID)
}

func TestPlaybookSvc_UpdatePlaybookInstance(t *testing.T) {
	r := struct {
		Out *PlaybookInstance `json:"updatePlaybookInstance"`
	}{
		Out: &testPlayBookInstance,
	}
	id := common.ID("test")
	objMetaInput := common.ObjectMetaInput{
		Name:        "test",
		Description: nil,
		Tags:        nil,
	}

	pbTriggerInput := PlaybookTriggerInput{
		ObjectMetaInput: objMetaInput,
		TypeID:          id,
		Config:          nil,
	}
	playbookInstanceInput := PlaybookInstanceInput{
		ObjectMetaInput: objMetaInput,
		Trigger:         pbTriggerInput,
		Enabled:         false,
		Inputs:          nil,
		Connections:     nil,
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL)
	i, err := c.UpdatePlaybookInstance("123", &playbookInstanceInput)
	assert.Nil(t, err)
	assert.Equal(t, "123", i.ID)

}
