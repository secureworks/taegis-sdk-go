package connectors

import (
	"github.com/secureworks/tdr-sdk-go/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestObject_Keys(t *testing.T) {
	o := make(map[string]interface{})
	o["test"] = "test"
	obj := common.Object(o)
	keys := obj.Keys()

	assert.Len(t, keys, 1)
	assert.Equal(t, "test", keys[0])
}

func TestConnector_LookupAction(t *testing.T) {
	testConnectorAction := ConnectorAction{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "testAction",
		Description: "testing actions",
		Tags:        nil,
		Interface:   nil,
		Inputs:      nil,
		Outputs:     nil,
	}

	testConnectorInterface := ConnectorInterface{
		ID:          "1234",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "test",
		Description: "test",
		Tags:        nil,
		Categories:  nil,
		Actions:     []*ConnectorAction{&testConnectorAction},
	}

	testConnector := Connector{
		ID:          "1",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Name:        "TestConnector2",
		Description: "Testing",
		Tags:        nil,
		Docs:        nil,
		Method:      nil,
		Implements:  []*ConnectorInterface{&testConnectorInterface},
		Actions:     nil,
		Parameters:  nil,
		AuthTypes:   nil,
	}

	a := testConnector.LookupAction("testAction")
	assert.Equal(t, "testAction", a.Name)

	a = testConnector.LookupAction("test")
	assert.Nil(t, a)
}
