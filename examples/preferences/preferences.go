package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/secureworks/tdr-sdk-go/client"
	"github.com/secureworks/tdr-sdk-go/preferences"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	clnt := client.NewClient()
	preferenceSvc := preferences.NewPreferencesSvc(clnt, "preferences example")

	userEmailPreferences := &preferences.UserEmail{
		Mention:             true,
		AssigneeChange:      true,
		AwaitingAction:      true,
		GroupMention:        true,
		GroupAssigneeChange: true,
		GroupAwaitingAction: true,
	}

	preferenceStorage := make(map[string]interface{})
	pStorageBytes, err := json.Marshal(userEmailPreferences)
	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	err = json.Unmarshal(pStorageBytes, &preferenceStorage)
	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	preferenceInput := &preferences.PreferencesInput{
		BearerToken: os.Getenv("ACCESS_TOKEN"),
		TenantID:    "123456789",
		Key:         preferences.EmailKey,
		Preferences: preferenceStorage,
	}

	//Create preferences example
	out, err := preferenceSvc.CreatePreferences(preferenceInput, preferences.DefaultFields)
	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	//CreatePreferences Output
	fmt.Println("")
	fmt.Println("CREATE PREFERENCES OUTPUT:")
	fmt.Println(out)

	//Get Preferences example
	out, err = preferenceSvc.GetPreferencesByKey(preferenceInput, preferences.DefaultFields)

	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	//Get Preferences Output
	fmt.Println("")
	fmt.Println("GET PREFERENCES OUTPUT:")
	fmt.Println(out)
}
