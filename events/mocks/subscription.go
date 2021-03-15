package mocks

import (
	"context"
	"errors"

	"github.com/secureworks/taegis-sdk-go/events"
)

type subscription struct {
	results chan interface{}
}

func (s *subscription) Close() error {
	return nil
}

func (s *subscription) Next(_ context.Context) (*events.EventQueryResults, error) {
	val, ok := <-s.results
	if !ok {
		return nil, errors.New("event query result channel closed")
	}

	switch v := val.(type) {
	case error:
		return nil, v
	case *events.EventQueryResults:
		return v, nil
	}
	return nil, errors.New("mock subscription only supports *EventQueryResult or error types")
}

func (s *subscription) GetAllEventResults(ctx context.Context) (events.Results, *string, error) {
	results := make(events.Results, 0, len(s.results))
	for {
		if len(s.results) == 0 { //Done chan has been exhausted and all results have been found
			return results, results.GetNextPageID(), nil
		}
		rslt, err := s.Next(ctx)
		if err != nil {
			return results, results.GetNextPageID(), err // unexpected error occurred - These errors typically indicate there is a problem with the request and are usually fatal. IE: "invalid page identifier" can cause infinite loops if we do not catch it.
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
