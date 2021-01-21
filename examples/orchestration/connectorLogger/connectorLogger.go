package main

import (
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/secureworks/tdr-sdk-go/client"
	"github.com/secureworks/tdr-sdk-go/common"
	"github.com/secureworks/tdr-sdk-go/orchestration/connectorLogger"
)

func errTest(err error) {
	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}
}
func main() {

	connectorLoggerSvc := connectorLogger.New("test.com", client.WithHTTPTimeout(3*time.Second), client.WithTenant("123456789"))

	input := connectorLogger.ConnectorLogQueryInput{
		Connector: "test",
	}

	page, perPage := 1, 10
	pagination := common.Pagination{
		Page:    &page,
		PerPage: &perPage,
	}

	logs, err := connectorLoggerSvc.GetAllConnectorLogs(input, pagination)
	errTest(err)

	for _, log := range logs.Entries {
		//log represents a Log Entry
		//use dot notations to access specific Log fields
		fmt.Printf("LOG ID: %v\n", log.ID)
		fmt.Printf("LOG MESSAGE: %v\n", log.Message)
	}
}
