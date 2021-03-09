package mocks

import (
	"context"

	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/events"
	"github.com/secureworks/taegis-sdk-go/graphql"
)

var _ events.EventsSvc = (*EventsSvc)(nil)

type EventsSvc struct {
	GetEventsError       error
	GetEventQueryError   error
	GetEventQueriesError error
	EventPageError       error
	EventQueryError      error

	GetEventsResult       []*events.Event
	GetEventQueryResult   *events.EventQuery
	GetEventQueriesResult []*events.EventQuery
	EventPageResult       chan interface{}
	EventQueryResult      chan interface{}
}

func (m *EventsSvc) GetEvents(_ []string, _ ...graphql.RequestOption) ([]*events.Event, error) {
	return m.GetEventsResult, m.GetEventsError
}

func (m *EventsSvc) GetEventQuery(_ string, _ ...graphql.RequestOption) (*events.EventQuery, error) {
	return m.GetEventQueryResult, m.GetEventQueryError
}

func (m *EventsSvc) GetEventQueries(_ ...graphql.RequestOption) ([]*events.EventQuery, error) {
	return m.GetEventQueriesResult, m.GetEventQueriesError
}

func (m *EventsSvc) EventPage(_ context.Context, _ string, _ ...graphql.SubscriptionOption) (events.Subscription, error) {
	return &subscription{results: m.EventPageResult}, m.EventPageError
}
func (m *EventsSvc) EventQuery(_ context.Context, _ string, _ common.Object, _ *events.EventQueryOptions, _ ...graphql.SubscriptionOption) (events.Subscription, error) {
	return &subscription{results: m.EventQueryResult}, m.EventQueryError
}
