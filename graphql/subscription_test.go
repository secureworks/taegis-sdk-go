package graphql_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/99designs/gqlgen/example/chat"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/secureworks/tdr-sdk-go/graphql"
	"github.com/stretchr/testify/require"
)

func newServer() *httptest.Server {
	srv := handler.NewDefaultServer(chat.NewExecutableSchema(newChatConfig()))
	return httptest.NewServer(srv)
}

type expectedResp struct {
	MessageAdded struct {
		Text string `json:"text"`
	} `json:"messageAdded"`
}

func TestSubscription_BasicScenario(t *testing.T) {
	s := newServer()
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)
	u.Scheme = "ws"

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	sub, err := graphql.NewSubscription(ctx, u, `subscription {
	messageAdded(roomName: "test") {
		text
	}
}`,
		func() interface{} {
			return &expectedResp{}
		})
	require.NoError(t, err)

	var (
		ok bool
		m  *graphql.Message
	)
	messages := sub.Messages()
	select {
	case m, ok = <-messages:
	case <-time.After(10 * time.Second):
		require.Fail(t, "channel timed out")
	}

	require.True(t, ok)
	require.NoError(t, m.Err)
	require.Equal(t, "hello!", m.Payload.(*expectedResp).MessageAdded.Text)
	err = sub.Shutdown(context.TODO())
	require.NoError(t, err)

	_, ok = <-messages
	require.False(t, ok)
}

func TestSubscription_SendKAs(t *testing.T) {
	s := newServer()
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)
	u.Scheme = "ws"

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	sub, err := graphql.NewSubscription(ctx, u, `subscription {
			messageAdded(roomName: "test") {
				text
			}
		}`,
		func() interface{} {
			return &expectedResp{}
		},
		graphql.SubscriptionSendKAMessages)
	require.NoError(t, err)
	defer sub.Shutdown(context.TODO())

	messages := sub.Messages()
	m := <-messages
	_, ok := m.Payload.(graphql.KeepAliveMessage)
	require.True(t, ok)
}

func TestSubscription_WorksWithMock(t *testing.T) {
	q := "subscription test"
	vars := map[string]interface{}{"test": "test"}
	output := 2
	s := graphql.NewMockSubServer(t, q, vars, output)
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)
	u.Scheme = "ws"

	sub, err := graphql.NewSubscription(ctx, u, q, func() interface{} {
		return new(int)
	}, graphql.SubscriptionWithVars(vars))
	require.NoError(t, err)
	defer sub.Shutdown(context.TODO())

	messages := sub.Messages()

	m := <-messages
	require.NoError(t, m.Err)
	require.Equal(t, &output, m.Payload)
}

func TestSubscription_WorksWithMockOnError(t *testing.T) {
	q := "subscription test"
	vars := map[string]interface{}{"test": "test"}
	output := errors.New("failed")
	s := graphql.NewMockSubServer(t, q, vars, output)
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)
	u.Scheme = "ws"

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	sub, err := graphql.NewSubscription(ctx, u, q, func() interface{} {
		return new(int)
	}, graphql.SubscriptionWithVars(vars))
	require.NoError(t, err)
	defer sub.Shutdown(context.TODO())
	m := <-sub.Messages()
	require.Error(t, m.Err)
}
