package events

import (
	"github.com/secureworks/tdr-sdk-go/client"
	"github.com/secureworks/tdr-sdk-go/graphql"
)

const IDs = "ids"
const ID = "id"

type EventsSvc interface {
	GetEvents(ids []string, opts ...graphql.RequestOption) ([]*Event, error)
	GetEventQuery(id string, opts ...graphql.RequestOption) (*EventQuery, error)
	GetEventQueries(opts ...graphql.RequestOption) ([]*EventQuery, error)
	Subscriptions
}

var _ EventsSvc = (*eventsSvc)(nil)

type eventsSvc struct {
	client *client.Client
	url    string
}

func New(url string, opts ...client.Option) *eventsSvc {
	client := client.NewClient(opts...)
	return &eventsSvc{
		client: client,
		url:    url,
	}
}

func (s *eventsSvc) GetEvents(ids []string, opts ...graphql.RequestOption) ([]*Event, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(getEvents, IDs), opts...)
	req.Var(IDs, ids)

	var data struct {
		Events []*Event `json:"events"`
	}
	if err := graphql.ExecuteQuery(s.client, s.url, req, &data); err != nil {
		return nil, err
	}
	return data.Events, nil
}

func (s *eventsSvc) GetEventQuery(id string, opts ...graphql.RequestOption) (*EventQuery, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(getEventQuery, ID), opts...)
	req.Var(ID, id)

	var data struct {
		EventQuery *EventQuery `json:"eventQuery"`
	}
	if err := graphql.ExecuteQuery(s.client, s.url, req, &data); err != nil {
		return nil, err
	}
	return data.EventQuery, nil
}

func (s *eventsSvc) GetEventQueries(opts ...graphql.RequestOption) ([]*EventQuery, error) {
	req := graphql.NewRequest(getEventQueries, opts...)

	var data struct {
		EventQueries []*EventQuery `json:"eventQueries"`
	}
	if err := graphql.ExecuteQuery(s.client, s.url, req, &data); err != nil {
		return nil, err
	}
	return data.EventQueries, nil
}

func (s *eventsSvc) DeleteEventQuery(id string, opts ...graphql.RequestOption) (bool, error) {
	req := graphql.NewRequest(graphql.AddVarNamesToQuery(deleteEventQuery, ID), opts...)
	req.Var(ID, id)

	var data struct {
		DeleteEventQuery bool `json:"deleteEventQuery"`
	}
	if err := graphql.ExecuteQuery(s.client, s.url, req, &data); err != nil {
		return false, err
	}
	return data.DeleteEventQuery, nil
}
