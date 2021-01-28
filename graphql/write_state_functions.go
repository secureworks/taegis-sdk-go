package graphql

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

type writeStateFunc func() (interface{}, writeStateFunc, error)

func (s *Subscription) writeConnectionInitMessageState() (interface{}, writeStateFunc, error) {
	s.log.Debug().Msg("writeState:writeConnectionInitMessage")
	return operationMessage{Type: connectionInitMsg}, s.writeSubscriptionQueryState, nil
}

func (s *Subscription) writeSubscriptionQueryState() (interface{}, writeStateFunc, error) {
	s.log.Debug().Msg("writeState:writeSubscriptionQuery")
	buf, err := json.Marshal(Request{
		Query:     s.query,
		Variables: s.vars,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed marshalling query")
	}

	return operationMessage{Type: startMsg, ID: "1", Payload: buf}, nil, nil
}

func (s *Subscription) connect(connectionCtx context.Context) error {
	currentState := s.writeConnectionInitMessageState
	for currentState != nil {

		select {
		case <-connectionCtx.Done():
			return connectionCtx.Err()
		case <-s.readerDone:
			return errors.New("reader down")
		default:
		}

		toWrite, next, err := currentState()
		if err != nil {

			return err
		}

		if err := s.conn.WriteJSON(toWrite); err != nil {
			return err
		}
		currentState = next
	}
	return nil
}
