package main

import (
	"fmt"
	"os"
	"time"

	"github.com/secureworks/taegis-sdk-go/client"

	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/orchestration/playbooks"
	"github.com/davecgh/go-spew/spew"
)

func testErr(err error) {
	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}
}
func main() {
	playbookSvc := playbooks.New("test.com", client.WithHTTPTimeout(5*time.Second))

	id := "123"
	name := "test"

	playbookTriggerType, err := playbookSvc.GetPlaybookTriggerType(&id, &name)

	testErr(err)
	//playbookTriggerType represents a PlaybookTriggerType Output
	//use dot notations to access specific PlaybookTriggerType fields
	fmt.Printf("PLAYBOOK TRIGGER TYPE NAME: %v\n", playbookTriggerType.Name)
	fmt.Printf("PLAYBOOK TRIGGER TYPE DESCRIPTION: %v\n", playbookTriggerType.Description)

	//playbookTriggers represents an array of PlaybookTriggers
	playbookTriggers, err := playbookSvc.GetPlaybookTriggers([]string{"123"})

	testErr(err)

	for _, playbookTrigger := range playbookTriggers {
		//playbookTrigger represents a PlaybookTrigger Output
		//use dot notations to access specific PlaybookTrigger fields
		fmt.Printf("PLAYBOOK TRIGGER NAME: %v\n", playbookTrigger.Name)
		fmt.Printf("PLAYBOOK TRIGGER DESCRIPTION: %v\n", playbookTrigger.Description)
	}

	//triggerTypes represents an array of PlaybookTriggerTypes
	triggerTypes, err := playbookSvc.GetPlaybookTriggerTypes()

	testErr(err)

	for _, triggerType := range triggerTypes {
		//triggerType represents a PlaybookTriggerType Output
		//use dot notations to access specific PlaybookTriggerType fields
		fmt.Printf("PLAYBOOK TRIGGER TYPE NAME: %v\n", triggerType.Name)
		fmt.Printf("PLAYBOOK TRIGGER TYPE DESCRIPTION: %v\n", triggerType.Description)
	}

	playbookExecution, err := playbookSvc.ExecutePlaybookInstance("test", nil)

	testErr(err)

	//playbookExecution represents a PlaybookExecution Output
	//use dot notations to access specific PlaybookExecution fields
	fmt.Printf("PLAYBOOK EXECUTION ID: %v\n", playbookExecution.ID)
	fmt.Printf("PLAYBOOK EXECUTION STATE: %v\n", playbookExecution.State)
	testObjMetaInput := common.ObjectMetaInput{
		Name:        "test",
		Description: nil,
		Tags:        nil,
	}
	testPlaybook := playbooks.PlaybookInput{
		ObjectMetaInput: testObjMetaInput,
		Head:            nil,
		Version:         nil,
		Categories:      nil,
	}
	//Create Playbook
	playbook, err := playbookSvc.CreatePlaybook(&testPlaybook)
	testErr(err)

	//playbook represents a Playbook Output
	//use dot notations to access specific Playbook fields
	fmt.Printf("PLAYBOOK ID: %v\n", playbook.ID)
	fmt.Printf("PLAYBOOK NAME: %v\n", playbook.Name)

	//Get Playbook
	playbook, err = playbookSvc.GetPlaybook("1234")
	testErr(err)

	//playbook represents a Playbook Output
	//use dot notations to access specific Playbook fields
	fmt.Printf("PLAYBOOK ID: %v\n", playbook.ID)
	fmt.Printf("PLAYBOOK NAME: %v\n", playbook.Name)

	//Update Playbook
	playbook, err = playbookSvc.UpdatePlaybook(playbook.ID, &testPlaybook)
	testErr(err)

	//playbook represents a Playbook Output
	//use dot notations to access specific Playbook fields
	fmt.Printf("PLAYBOOK ID: %v\n", playbook.ID)
	fmt.Printf("PLAYBOOK NAME: %v\n", playbook.Name)

	//Execute Playbook
	execution, err := playbookSvc.ExecutePlaybook(playbook.ID, nil)
	testErr(err)

	//execution represents a PlaybookExecution Output
	//use dot notations to access specific PlaybookExecution fields
	fmt.Printf("PLAYBOOK EXECUTION ID: %v\n", execution.ID)
	fmt.Printf("PLAYBOOK EXECUTION STATE: %v\n", execution.State)

	//Clone Playbook
	playbook, err = playbookSvc.ClonePlaybook(playbooks.ClonePlaybookInput{PlaybookID: common.ID(playbook.ID), Name: "test"})
	testErr(err)

	//playbook represents a Playbook Output
	//use dot notations to access specific Playbook fields
	fmt.Printf("CLONED PLAYBOOK ID: %v\n", playbook.ID)
	fmt.Printf("CLONED PLAYBOOK NAME: %v\n", playbook.Name)

	triggerID := common.ID("triggerID")
	playbookTriggerInput := playbooks.PlaybookTriggerInput{
		ObjectMetaInput: testObjMetaInput,
		TypeID:          triggerID,
		Config:          nil,
	}
	playbookInstanceInput := playbooks.PlaybookInstanceInput{
		ObjectMetaInput: testObjMetaInput,
		Trigger:         playbookTriggerInput,
		Enabled:         false,
		Inputs:          nil,
		Connections:     nil,
	}
	//Create Playbook Instance
	instance, err := playbookSvc.CreatePlaybookInstance(playbook.ID, &playbookInstanceInput)
	testErr(err)
	//instance represents a PlaybookInstance Output
	//use dot notations to access specific PlaybookInstance fields
	fmt.Printf("PLAYBOOK INSTANCE ID: %v\n", instance.ID)

	//Get Playbook Instance
	instance, err = playbookSvc.GetPlaybookInstance(instance.ID)
	testErr(err)
	//instance represents a PlaybookInstance Output
	//use dot notations to access specific PlaybookInstance fields
	fmt.Printf("PLAYBOOK INSTANCE ID: %v\n", instance.ID)

	//Update Playbook Instance
	instance, err = playbookSvc.UpdatePlaybookInstance(instance.ID, &playbookInstanceInput)
	testErr(err)
	//instance represents a PlaybookInstance Output
	//use dot notations to access specific PlaybookInstance fields
	fmt.Printf("PLAYBOOK INSTANCE ID: %v\n", instance.ID)

	//Set PlaybookInstance State
	instance, err = playbookSvc.SetPlaybookInstanceState(instance.ID, true)
	testErr(err)
	//instance represents a PlaybookInstance Output
	//use dot notations to access specific PlaybookInstance fields
	fmt.Printf("PLAYBOOK INSTANCE ID: %v\n", instance.ID)

	//Get Playbooks
	playbooks, err := playbookSvc.GetPlaybooks(nil, nil)
	testErr(err)
	for _, pb := range playbooks {
		//pb represents a Playbook Output
		//use dot notations to access specific Playbook fields
		fmt.Printf("CLONED PLAYBOOK ID: %v\n", pb.ID)
		fmt.Printf("CLONED PLAYBOOK NAME: %v\n", pb.Name)
	}

	//Get PlaybookInstances
	instances, err := playbookSvc.GetPlaybookInstances(&playbook.ID)
	testErr(err)

	for _, i := range instances {
		//i represents a PlaybookInstance Output
		//use dot notations to access specific PlaybookInstance fields
		fmt.Printf("PLAYBOOK INSTANCE ID: %v\n", i.ID)
	}

	//Get Playbook Execution
	execution, err = playbookSvc.GetPlaybookExecution(execution.ID)
	testErr(err)

	//execution represents a PlaybookExecution Output
	//use dot notations to access specific PlaybookExecution fields
	fmt.Printf("PLAYBOOK EXECUTION ID: %v\n", execution.ID)
	fmt.Printf("PLAYBOOK EXECUTION STATE: %v\n", execution.State)

	//Get Playbook Trigger
	trigger, err := playbookSvc.GetPlaybookTrigger("1234")
	testErr(err)

	//trigger represents a PlaybookTrigger Output
	//use dot notations to access specific PlaybookTrigger fields
	fmt.Printf("PLAYBOOK TRIGGER ID: %v\n", trigger.ID)
	fmt.Printf("PLAYBOOK TRIGGER NAME: %v\n", trigger.Name)

	//Get Playbook Executions
	executions, err := playbookSvc.GetPlaybookExecutions(playbook.ID, common.NewPaginationOptions(1, 1))
	testErr(err)
	fmt.Printf("TOTAL COUNT: %v\n", executions.TotalCount)
	for _, ex := range executions.Executions {
		//ex represents a PlaybookExecution Output
		//use dot notations to access specific PlaybookExecution fields
		fmt.Printf("PLAYBOOK EXECUTION ID: %v\n", ex.ID)
		fmt.Printf("PLAYBOOK EXECUTION STATE: %v\n", ex.State)
	}

	//Delete PlaybookInstance
	instance, err = playbookSvc.DeletePlaybookInstance(instance.ID)
	testErr(err)

	//Delete Playbook
	playbook, err = playbookSvc.DeletePlaybook(playbook.ID)
	testErr(err)

	//instance represents a PlaybookInstance Output
	//use dot notations to access specific PlaybookInstance fields
	fmt.Printf("PLAYBOOK INSTANCE ID: %v\n", instance.ID)

	//playbook represents a Playbook Output
	//use dot notations to access specific Playbook fields
	fmt.Printf("DELETED PLAYBOOK ID: %v\n", playbook.ID)
	fmt.Printf("DELETED PLAYBOOK NAME: %v\n", playbook.Name)

}
