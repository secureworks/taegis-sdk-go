package main

import (
	"fmt"
	"os"

	"github.com/secureworks/taegis-sdk-go/graphql"

	"github.com/secureworks/taegis-sdk-go/client"
	"github.com/secureworks/taegis-sdk-go/investigations"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	clnt := client.NewClient()
	investigationSvc := investigations.NewInvestigationSvc(clnt, "my-app-name-or-email")

	out, err := investigationSvc.GetInvestigation(&investigations.GetInvestigationInput{
		TenantID: "123456789",
		ID:       "8569a520-63e5-4352-af0c-432c2523be46"},
		investigations.DefaultFields,
		graphql.RequestWithToken(os.Getenv("ACCESS_TOKEN")),
	)

	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	//out represents an InvestigationOutput
	//use dot notations to access specific Investigation fields
	fmt.Printf("INVESTIGATION NAME: %v\n", out.Description)
	fmt.Printf("INVESTIGATION STATUS: %v\n", out.Status)
}
