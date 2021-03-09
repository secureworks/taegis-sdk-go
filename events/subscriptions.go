package events

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/graphql"
)

const Query = "query"
const Metadata = "metadata"
const Options = "options"
const PageID = "id"

const (
	eventPageSubscription = `
		subscription($%[1]s: ID!) {
		    eventPage(id: $%[1]s) {
				query {
  					id
  					query
	  				status
	  				reasons {
						id
						type
						backend
						status
						reason
						submitted
						completed
					}
	  				submitted
	  				completed
	  				expires
	  				types
	  				metadata
				}
				result {	
					id
    				type
					backend
					status
					reason
    				submitted
    				completed
    				expires
                    facets
					rows
					progress {
						totalRows
						totalRowsIsLowerBound
						resultsTruncated
					}
				}
				next
				prev
			}
		}
	`
	eventQuerySubscription = `
		subscription($%[1]s: String!, $%[2]s: JSONObject, $%[3]s: EventQueryOptions) {
		    eventQuery(query: $%[1]s, metadata: $%[2]s, options: $%[3]s){
				query {
  					id
  					query
	  				status
	  				reasons {
						id
						type
						backend
						status
						reason
						submitted
						completed
					}
	  				submitted
	  				completed
	  				expires
	  				types
	  				metadata
				}
				result {	
					id
    				type
					backend
					status
					reason
    				submitted
    				completed
    				expires
    				facets
					rows
					progress {
						totalRows
						totalRowsIsLowerBound
						resultsTruncated
					}
				}
				next
				prev
			}
		}
	`
)

type Subscriptions interface {
	EventPage(ctx context.Context, pageID string, options ...graphql.SubscriptionOption) (Subscription, error)
	EventQuery(ctx context.Context, query string, metadata common.Object, qopts *EventQueryOptions, sopts ...graphql.SubscriptionOption) (Subscription, error)
}

func (s *eventsSvc) EventPage(ctx context.Context, id string, options ...graphql.SubscriptionOption) (Subscription, error) {
	u, err := url.Parse(s.url)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "http", "ws":
		u.Scheme = "ws"
	case "https", "wss":
		u.Scheme = "wss"
	default:
		return nil, fmt.Errorf("invalid protocol: %s", u.Scheme)
	}

	gSub, err := graphql.NewSubscription(ctx, u,
		graphql.AddVarNamesToQuery(eventPageSubscription, PageID),
		func() interface{} { return &eventPageResult{} },
		append(options, graphql.SubscriptionWithVars(map[string]interface{}{
			PageID: id,
		}))...,
	)
	if err != nil {
		return nil, err
	}
	sub := &subscription{sub: gSub, msgs: gSub.Messages()}
	return sub, nil
}

func (s *eventsSvc) EventQuery(ctx context.Context, query string, metadata common.Object, qopts *EventQueryOptions, sopts ...graphql.SubscriptionOption) (Subscription, error) {
	u, err := url.Parse(s.url)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "http", "ws":
		u.Scheme = "ws"
	case "https", "wss":
		u.Scheme = "wss"
	default:
		return nil, fmt.Errorf("invalid protocol: %s", u.Scheme)
	}
	gSub, err := graphql.NewSubscription(ctx, u,
		graphql.AddVarNamesToQuery(eventQuerySubscription, Query, Metadata, Options),
		func() interface{} { return &eventQueryResult{} },
		append(sopts, graphql.SubscriptionWithVars(map[string]interface{}{
			Query:    query,
			Metadata: metadata,
			Options:  qopts,
		}))...,
	)
	if err != nil {
		return nil, err
	}
	sub := &subscription{sub: gSub, msgs: gSub.Messages()}

	return sub, nil
}

// Subscriptions
type Subscription interface {
	Next(ctx context.Context) (*EventQueryResults, error)
	io.Closer

	//GetAllEventResults is a wrapper for Next that will gather all results until the page has been exhausted.
	//It accepts a context that can be canceled to stop getting results and should return any that have already been found
	//Once all results have been retrieved for a page, it should return a list of those results and a pointer to the next page
	GetAllEventResults(ctx context.Context) (Results, *string, error)
}

type subscription struct {
	sub  *graphql.Subscription
	msgs <-chan *graphql.Message
}

type eventQueryResult struct {
	EventQueryResults *EventQueryResults `json:"eventQuery"`
}

type eventPageResult struct {
	EventQueryResults *EventQueryResults `json:"eventPage"`
}

func (s *subscription) Next(ctx context.Context) (*EventQueryResults, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg, ok := <-s.msgs:
		if !ok {
			return nil, io.EOF
		}
		if msg.Err != nil {
			return nil, msg.Err
		}
		switch t := msg.Payload.(type) {
		case graphql.KeepAliveMessage:
			return &EventQueryResults{}, nil
		case *eventQueryResult:
			return t.EventQueryResults, nil
		case *eventPageResult:
			return t.EventQueryResults, nil
		default:
			return nil, fmt.Errorf("unexpected event payload type %T", msg.Payload)
		}
	}
}

func (s *subscription) Close() error {
	return s.sub.Shutdown(context.Background())
}

//GetAllEventResults is a wrapper for Next that will gather all results until the page has been exhausted.
//It accepts a context that can be canceled to stop getting results and should return any that have already been found
//Once all results have been retrieved for a page, it should return a list of those results and a pointer to the next page
func (s *subscription) GetAllEventResults(ctx context.Context) (Results, *string, error) {
	results := make(Results, 0, 10)
	for {
		rslt, err := s.Next(ctx)
		if err != nil {
			switch err {
			case context.Canceled, context.DeadlineExceeded, io.EOF:
				return results, results.GetNextPageID(), nil
			default:
				return results, results.GetNextPageID(), err // unexpected error occurred - These errors typically indicate there is a problem with the request and are usually fatal. IE: "invalid page identifier" can cause infinite loops if we do not catch it.
			}
		}
		if rslt == nil {
			continue
		}
		if rslt.Result == nil {
			switch rslt.Query.Status { //TODO: Determine what cases need to be handled
			case "RUNNING":
				continue // If the query is still running keep going
			default: //"SUCCEEDED", "QUEUED", "CANCELED", "FAILED"
				return results, results.GetNextPageID(), nil
			}
		}
		results = append(results, rslt)
	}
}
