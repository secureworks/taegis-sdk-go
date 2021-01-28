package main

import (
	"fmt"
	"os"

	"github.com/secureworks/tdr-sdk-go/client"
	"github.com/secureworks/tdr-sdk-go/collectors"
	"github.com/davecgh/go-spew/spew"
)

const (
	tenantID    = "123456789"
	defaultRole = "collector"
	apiEndpoint = "https://api.ctpx.secureworks.com/graphql"
)

func main() {
	token := os.Getenv("ACCESS_TOKEN")
	c := client.NewClient(client.WithBearerToken(token))
	cl := collectors.NewWithClient(apiEndpoint, tenantID, c)
	clusters, err := cl.GetAllClusters(defaultRole)
	if err != nil {
		spew.Dump(err)
		return
	}

	for _, c := range clusters {
		fmt.Printf("collector id: %s\n", c.ID)
		if c.Name != nil {
			fmt.Printf("collector name: %v\n", *c.Name)
		}
	}
}
