package main

import (
	"fmt"
	"os"

	"github.com/secureworks/taegis-sdk-go/client"
	"github.com/secureworks/taegis-sdk-go/graphql"
	"github.com/secureworks/taegis-sdk-go/users"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	clnt := client.NewClient()
	userSvc := users.NewUserSvc(clnt)
	getUserInput := &users.GetUserInput{ID: "555"}

	out, err := userSvc.GetUser(getUserInput,
		users.DefaultFields,
		graphql.RequestWithToken(os.Getenv("ACCESS_TOKEN")),
		graphql.RequestWithTenant("1234"))

	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	fmt.Println("Retrieved user, ID: ", out.ID)
}
