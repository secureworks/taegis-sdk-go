package graphql

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type readStateFunc func() (json.Unmarshaler, readStateFunc)

func (s *Subscription) startWSReader() {
	defer func() {
		close(s.ch)
		close(s.readerDone)
	}()
	s.log.Debug().Msg("started reader goroutine")

	currentState := s.readAckState
	for currentState != nil {
		var toRead json.Unmarshaler
		toRead, currentState = currentState()
		err := s.conn.ReadJSON(toRead)
		if err != nil {
			if errors.Is(err, errOmitMessage) {
				continue
			}
			s.log.WithError(err).Warn().Msg("reader ended")
			return
		}
	}
}

type readAck struct {
	s *Subscription
}

func (r *readAck) UnmarshalJSON(data []byte) error {
	var ack operationMessage
	if err := json.Unmarshal(data, &ack); err != nil {
		err = errors.Wrap(err, "invalid ack message")
		return err
	}
	if ack.Type != connectionAckMsg {
		err := fmt.Errorf("expected ack message, got %#v", ack)
		return err
	}

	return nil
}

func (s *Subscription) readAckState() (json.Unmarshaler, readStateFunc) {
	s.log.Debug().Msg("readState:readAck")
	return &readAck{s: s}, s.readSubscriptionMessagesState
}

type readMessage struct {
	s *Subscription
}

func (r *readMessage) UnmarshalJSON(data []byte) error {
	var op operationMessage
	if err := json.Unmarshal(data, &op); err != nil {
		return errors.Wrap(err, "invalid message")
	}

	r.s.log.WithFields(map[string]interface{}{
		"payload": string(op.Payload),
		"op":      op.Type,
	}).Debug().Msg("subscription got message")

	var msg *Message

	switch op.Type {
	case dataMsg:
		var resp Response
		resp.Data = r.s.responseCreator()

		if err := json.Unmarshal(op.Payload, &resp); err != nil {
			return errors.Wrap(err, "failed unmarshalling payload, make sure responseCreator is compatible with the graphql schema")
		}

		var err error
		for _, gqlErr := range resp.Error {
			err = multierror.Append(err, gqlErr)
		}

		msg = &Message{Payload: resp.Data, Err: err}
	case connectionKaMsg:
		if !r.s.sendKAMsgs {
			return errOmitMessage
		}
		msg = &Message{Payload: KeepAliveMessage{}}
	case errorMsg:
		return fmt.Errorf(string(op.Payload))
	case completedMsg:
		return errSubscriptionCompleted
	default:
		return fmt.Errorf("unexpected message type: %#v", op)
	}

	if err := r.s.informMessageReceived(msg); err != nil {
		return errors.Wrap(err, "failed sending message to channel")
	}

	return nil
}

func (s *Subscription) readSubscriptionMessagesState() (json.Unmarshaler, readStateFunc) {
	s.log.Debug().Msg("readState:readSubscriptionMessages")
	return &readMessage{s: s}, s.readSubscriptionMessagesState
}
