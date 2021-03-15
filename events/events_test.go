package events

import (
	"testing"
	"time"

	"github.com/secureworks/taegis-sdk-go/client"
	"github.com/stretchr/testify/assert"

	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/testutils"
)

var (
	testEvent = Event{
		ID:     "123",
		Values: common.Object{"value": "test"},
	}
	header = testutils.CreateHeader()
)

func TestEventsSvc_GetEvents(t *testing.T) {
	r := struct {
		Out []*Event `json:"events"`
	}{
		Out: []*Event{&testEvent},
	}

	s := testutils.NewMockGQLOutput(t, header, r)
	defer s.Close()

	c := New(s.URL, client.WithHTTPTimeout(5*time.Second))
	event, err := c.GetEvents([]string{"123"})
	assert.Nil(t, err)
	assert.Equal(t, "123", event[0].ID)
}
